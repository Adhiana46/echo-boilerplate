package seeds

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type RoleSeeder struct{}

func (s *RoleSeeder) Name() string {
	return "Role Seeder"
}

func (s *RoleSeeder) Up(db *sql.DB) error {
	roles := map[string][]string{
		"super-admin": {
			"permissions",
			"permissions.create",
			"permissions.read",
			"permissions.update",
			"permissions.delete",
			"roles",
			"roles.create",
			"roles.read",
			"roles.update",
			"roles.delete",
			"users",
			"users.create",
			"users.read",
			"users.update",
			"users.delete",
		},
		"admin": {
			"roles",
			"roles.create",
			"roles.read",
			"roles.update",
			"roles.delete",
			"users",
			"users.create",
			"users.read",
			"users.update",
			"users.delete",
		},
	}

	sqlRole := `
		INSERT INTO roles 
		(uuid, name, created_at, updated_at)
		VALUES
		($1, $2, $3, $4)
		RETURNING id
	`
	sqlRolePerm := `
		INSERT INTO role_permissions
		(role_id, permission_id)
		SELECT $1, (SELECT id FROM permissions WHERE name = $2)
	`

	roleId := 0
	for role, permissions := range roles {
		err := db.QueryRow(
			sqlRole,
			uuid.NewString(),
			role,
			time.Now(),
			time.Now(),
		).Scan(&roleId)

		if err != nil {
			return err
		}

		for _, permission := range permissions {
			err := db.QueryRow(
				sqlRolePerm,
				roleId,
				permission,
			).Err()

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *RoleSeeder) Down(db *sql.DB) error {
	var err error

	_, err = db.Exec("TRUNCATE role_permissions")
	if err != nil {
		return err
	}

	_, err = db.Exec("TRUNCATE roles CASCADE")
	if err != nil {
		return err
	}

	return nil
}
