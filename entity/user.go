package entity

import (
	"database/sql"
)

type User struct {
	Id          int           `db:"id" json:"id"`
	Uuid        string        `db:"uuid" json:"uuid"`
	Username    string        `db:"username" json:"username"`
	Email       string        `db:"email" json:"email"`
	Password    string        `db:"password" json:"password"`
	Name        string        `db:"name" json:"name"`
	RoleId      int           `db:"role_id" json:"role_id"`
	Status      int           `db:"status" json:"status"`
	LastLoginAt sql.NullTime  `db:"last_login_at" json:"last_login_at"`
	CreatedAt   sql.NullTime  `db:"created_at" json:"created_at"`
	CreatedBy   sql.NullInt64 `db:"created_by" json:"created_by"`
	UpdatedAt   sql.NullTime  `db:"updated_at" json:"updated_at"`
	UpdatedBy   sql.NullInt64 `db:"updated_by" json:"updated_by"`
	Role        *Role         `json:"role"`
}
