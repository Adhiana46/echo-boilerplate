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

type job struct {
	lbl string
	f   func(*sql.DB) error
}

func (s *Seed) Run(ctx context.Context) error {
	lists := []job{
		job{
			lbl: "Permissions Seeder",
			f:   seedPermissions,
		},
		job{
			lbl: "Roles Seeder",
			f:   seedRoles,
		},
		job{
			lbl: "Users Seeder",
			f:   seedUsers,
		},
	}

	for _, j := range lists {
		log.Println("[Seeder]:\t-", fmt.Sprintf("Running %s", j.lbl))
		if err := j.f(s.db); err != nil {
			return err
		}
	}

	return nil
}

func (s *Seed) Rollback(ctx context.Context) error {
	lists := []job{
		job{
			lbl: "Permissions Seeder",
			f:   unseedPermissions,
		},
		job{
			lbl: "Roles Seeder",
			f:   unseedRoles,
		},
		job{
			lbl: "Users Seeder",
			f:   unseedUsers,
		},
	}

	for _, j := range lists {
		log.Println("[Seeder]:\t-", fmt.Sprintf("Rollback %s", j.lbl))
		if err := j.f(s.db); err != nil {
			return err
		}
	}

	return nil
}
