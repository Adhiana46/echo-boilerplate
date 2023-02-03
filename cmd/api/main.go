package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Adhiana46/echo-boilerplate/config"
	cachePkg "github.com/Adhiana46/echo-boilerplate/pkg/cache"
	tokenmanager "github.com/Adhiana46/echo-boilerplate/pkg/token-manager"
	"github.com/Adhiana46/echo-boilerplate/server"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	cfg          *config.Config
	db           *sqlx.DB
	cache        cachePkg.Cache
	tokenManager *tokenmanager.TokenManager
)

func main() {
	var err error

	cfg, err = config.LoadConfig()
	if err != nil {
		log.Panic("[Error][Config]", err)
	}
	log.Println("[Boot]:", "Config loaded")

	db, err = openDb(cfg.Pg)
	if err != nil {
		log.Panic("[Error][DB]:", err)
	}
	defer db.Close()

	// cache
	cache, err = openCache(cfg)
	if err != nil {
		log.Panic("[Error][Cache]:", err)
	}

	if cache != nil {
		defer cache.Close()
	}

	tokenManager = tokenmanager.NewTokenManager(&cfg.JWT, cache)

	// TODO: documentstore

	// TODO: notif

	run()
}

func run() {
	flag.Parse()
	args := flag.Args()

	cmd := "serve"

	if len(args) > 0 {
		cmd = args[0]
	}

	switch cmd {
	case "migrate":
		if err := runMigration(db.DB); err != nil {
			panic(err)
		}
	case "migrate:rollback":
		if err := rollbackMigration(db.DB); err != nil {
			panic(err)
		}
	case "migrate:fresh":
		if err := rollbackMigration(db.DB); err != nil {
			panic(err)
		}
		if err := runMigration(db.DB); err != nil {
			panic(err)
		}
		if err := runSeeder(db.DB); err != nil {
			panic(err)
		}
	case "migrate:seed":
		if err := runSeeder(db.DB); err != nil {
			panic(err)
		}
	case "migrate:unseed":
		if err := rollbackSeeder(db.DB); err != nil {
			panic(err)
		}
	default:
		// run server
		srv := server.NewServer(cfg, db, cache, tokenManager)

		go func() {
			if err := srv.Run(); err != nil {
				log.Fatal("shutting down the server ", err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(0)
}

func openDb(cfg config.PgConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DbName,
		cfg.Pass,
	)

	dbConn, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	dbConn.SetMaxOpenConns(60)
	dbConn.SetConnMaxLifetime(120 * time.Second)
	dbConn.SetMaxIdleConns(30)
	dbConn.SetConnMaxIdleTime(20 * time.Second)
	if err = dbConn.Ping(); err != nil {
		return nil, err
	} else {
		log.Println("[Boot]:", "Database connected successfully!")
	}

	return dbConn, nil
}

func openCache(cfg *config.Config) (cachePkg.Cache, error) {
	var instance cachePkg.Cache
	var err error

	switch cfg.Cache.Driver {
	case "redis":
		instance, err = cachePkg.NewRedisCache(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, 0)
	case "memcached":
		instance, err = cachePkg.NewMcCache(cfg.Memcached.Hosts...)
	}

	return instance, err
}
