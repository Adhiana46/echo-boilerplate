package seeds

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sort"
)

type Seed struct {
	db *sql.DB
}

type SeederInterface interface {
	Name() string
	Up(db *sql.DB) error
	Down(db *sql.DB) error
}

func NewSeeder(db *sql.DB) *Seed {
	return &Seed{
		db: db,
	}
}

var seeds []SeederInterface = []SeederInterface{
	&PermissionSeeder{},
	&RoleSeeder{},
	&UserSeeder{},
}

func (s *Seed) Up(ctx context.Context) error {
	for _, seed := range seeds {
		log.Println("[Seeder]:\t-", fmt.Sprintf("Running %s", seed.Name()))
		if err := seed.Up(s.db); err != nil {
			return err
		}
	}

	return nil
}

func (s *Seed) Down(ctx context.Context) error {
	sort.Slice(seeds, func(i, j int) bool {
		return i > j
	})

	for _, seed := range seeds {
		log.Println("[Seeder]:\t-", fmt.Sprintf("Rollback %s", seed.Name()))
		if err := seed.Down(s.db); err != nil {
			return err
		}
	}

	return nil
}
