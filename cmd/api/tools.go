package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/Adhiana46/echo-boilerplate/database/seeds"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

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
	err := seeder.Up(context.Background())
	if err != nil {
		return err
	}

	log.Println("[Seeder]:", "Seeder running sucessfully")

	return nil
}

func rollbackSeeder(db *sql.DB) error {
	log.Println("[Seeder]:", "Rollback database seeder")

	seeder := seeds.NewSeeder(db)
	err := seeder.Down(context.Background())
	if err != nil {
		return err
	}

	log.Println("[Seeder]:", "Seeder rollbacked sucessfully")

	return nil
}
