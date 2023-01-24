package entity

import (
	"database/sql"
)

type Role struct {
	Id          int          `db:"id" json:"id"`
	Uuid        string       `db:"uuid" json:"uuid"`
	Name        string       `db:"name" json:"name"`
	CreatedAt   sql.NullTime `db:"created_at" json:"created_at"`
	CreatedBy   int          `db:"created_by" json:"created_by"`
	UpdatedAt   sql.NullTime `db:"updated_at" json:"updated_at"`
	UpdatedBy   int          `db:"updated_by" json:"updated_by"`
	Permissions []Permission `json:"permissions"`
}
