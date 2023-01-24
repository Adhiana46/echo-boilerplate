package seeds

import (
	"database/sql"
	"time"

	"github.com/Adhiana46/echo-boilerplate/pkg/utils"
	"github.com/google/uuid"
)

func seedUsers(db *sql.DB) error {
	password, err := utils.HashPassword("pass1234")
	if err != nil {
		return err
	}

	users := []map[string]string{
		map[string]string{
			"username": "root",
			"email":    "root@example.com",
			"password": password,
			"name":     "Super Admin",
			"role":     "super-admin",
			"status":   "1",
		},
		map[string]string{
			"username": "admin",
			"email":    "admin@example.com",
			"password": password,
			"name":     "Admin",
			"role":     "admin",
			"status":   "1",
		},
	}

	sql := `
		INSERT INTO users
		(uuid, username, email, password, name, role_id, status, created_at, updated_at)
		SELECT $1, $2, $3, $4, $5, (SELECT id FROM roles WHERE name = $6), $7, $8, $9
	`

	for _, user := range users {
		err := db.QueryRow(
			sql,
			uuid.NewString(),
			user["username"],
			user["email"],
			user["password"],
			user["name"],
			user["role"],
			user["status"],
			time.Now(),
			time.Now(),
		).Err()

		if err != nil {
			return err
		}
	}

	return nil
}

func unseedUsers(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE users")
	if err != nil {
		return err
	}

	return nil
}
