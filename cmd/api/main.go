package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Adhiana46/echo-boilerplate/config"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Panic("Error loading config: ", err)
	}

	db, err := openDb(cfg.Pg)
	if err != nil {
		log.Panic("Error connecting to DB: ", err)
	}

	// TODO: cache

	// TODO: documentstore

	// TODO: notif

	// TODO: server
	// srv := server.NewServer(cfg, db)

	// if err := srv.Run(); err != nil {
	// 	log
	// }
}

func openDb(cfg config.PgConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=false password=%s",
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
		log.Println("Database connected successfully!")
	}

	return dbConn, nil
}

func openCache() {
	//
}
