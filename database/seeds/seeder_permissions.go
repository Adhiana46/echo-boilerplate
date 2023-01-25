package seeds

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PermissionSeeder struct{}

func (s *PermissionSeeder) Name() string {
	return "Permission Seeder"
}

func (s *PermissionSeeder) Up(db *sql.DB) error {
	menus := []string{
		"permissions",
		"roles",
		"users",
	}

	actions := []string{
		"create",
		"read",
		"update",
		"delete",
	}

	sql := `
		INSERT INTO permissions 
		(uuid, parent_id, name, type, created_at, updated_at) 
		VALUES 
		($1, $2, $3, $4, $5, $6)
		RETURNING id`

	parentId := 0
	for _, menu := range menus {
		err := db.QueryRow(
			sql,
			uuid.NewString(),
			0,
			menu,
			"menu",
			time.Now(),
			time.Now(),
		).Scan(&parentId)

		if err != nil {
			return err
		}

		for _, action := range actions {
			err = db.QueryRow(
				sql,
				uuid.NewString(),
				parentId,
				fmt.Sprintf("%s.%s", menu, action),
				"action",
				time.Now(),
				time.Now(),
			).Err()

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *PermissionSeeder) Down(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE permissions CASCADE")
	if err != nil {
		return err
	}

	return nil
}
