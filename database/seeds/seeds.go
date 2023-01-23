package seeds

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Seed struct {
	db *sql.DB
}

func NewSeeder(db *sql.DB) *Seed {
	return &Seed{
		db: db,
	}
}

func (s *Seed) Run(ctx context.Context) error {
	lists := map[string]func(*sql.DB) error{
		"Permissions Seeder": seedPermissions,
		"Roles Seeder":       seedRoles,
		"Users Seeder":       seedUsers,
	}

	for lbl, fn := range lists {
		log.Println("[Seeder]:\t-", fmt.Sprintf("Running %s", lbl))
		if err := fn(s.db); err != nil {
			return err
		}
	}

	return nil
}

func (s *Seed) Rollback(ctx context.Context) error {
	lists := map[string]func(*sql.DB) error{
		"Users Seeder":       unseedUsers,
		"Roles Seeder":       unseedRoles,
		"Permissions Seeder": unseedPermissions,
	}

	for lbl, fn := range lists {
		log.Println("[Seeder]:\t-", fmt.Sprintf("Rollback %s", lbl))
		if err := fn(s.db); err != nil {
			return err
		}
	}

	return nil
}
