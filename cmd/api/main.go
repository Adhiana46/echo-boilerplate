package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Adhiana46/echo-boilerplate/config"
	"github.com/Adhiana46/echo-boilerplate/database/seeds"
	cachePkg "github.com/Adhiana46/echo-boilerplate/pkg/cache"
	"github.com/Adhiana46/echo-boilerplate/server"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	cfg *config.Config
	db  *sqlx.DB
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
	cache, err := openCache(cfg)
	if err != nil {
		log.Panic("[Error][Cache]:", err)
	}

	if cache != nil {
		defer cache.Close()
	}

	// TODO: documentstore

	// TODO: notif

	handleArgs()

	// server
	srv := server.NewServer(cfg, db, cache)

	if err := srv.Run(); err != nil {
		log.Panic("[Error][Server]", err)
	}
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

func handleArgs() {
	flag.Parse()
	args := flag.Args()

	if len(args) >= 1 {
		switch args[0] {
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
		}

		os.Exit(0)
	}
}

func runMigration(db *sql.DB) error {
	log.Println("[Migration]:", "Running database migrations...")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Println("[Migration]: ", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./database/migrations",
		"postgres", driver)
	if err != nil {
		log.Println("[Migration]:", err)
		return err
	}

	err = m.Up()
	if err != nil {
		log.Println("[Migration]:", err)
		return err
	}

	log.Println("[Migration]:", "Database migration running successfully")

	return nil
}

func rollbackMigration(db *sql.DB) error {
	log.Println("[Migration]:", "Rollback database migrations...")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Println("[Migration]: ", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./database/migrations",
		"postgres", driver)
	if err != nil {
		log.Println("[Migration]:", err)
		return err
	}

	err = m.Down()
	if err != nil {
		log.Println("[Migration]:", err)
		return err
	}

	log.Println("[Migration]:", "Database migration rollback successfully")

	return nil
}

func runSeeder(db *sql.DB) error {
	log.Println("[Seeder]:", "Running database seeder")

	seeder := seeds.NewSeeder(db)
	err := seeder.Run(context.Background())
	if err != nil {
		return err
	}

	log.Println("[Seeder]:", "Seeder running sucessfully")

	return nil
}

func rollbackSeeder(db *sql.DB) error {
	log.Println("[Seeder]:", "Rollback database seeder")

	seeder := seeds.NewSeeder(db)
	err := seeder.Rollback(context.Background())
	if err != nil {
		return err
	}

	log.Println("[Seeder]:", "Seeder rollbacked sucessfully")

	return nil
}
