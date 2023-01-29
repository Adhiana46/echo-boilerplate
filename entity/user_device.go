package entity

import (
	"database/sql"
)

type UserDevice struct {
	Id         int           `db:"id" json:"id"`
	Uuid       string        `db:"uuid" json:"uuid"`
	UserId     int           `db:"user_id" json:"user_id"`
	Token      string        `db:"token" json:"token"`
	IP         string        `db:"ip" json:"ip"`
	Location   string        `db:"location" json:"location"`
	Platform   string        `db:"platform" json:"platform"`
	UserAgent  string        `db:"user_agent" json:"user_agent"`
	AppVersion string        `db:"app_version" json:"app_version"`
	Vendor     string        `db:"vendor" json:"vendor"`
	CreatedAt  sql.NullTime  `db:"created_at" json:"created_at"`
	CreatedBy  sql.NullInt64 `db:"created_by" json:"created_by"`
	UpdatedAt  sql.NullTime  `db:"updated_at" json:"updated_at"`
	UpdatedBy  sql.NullInt64 `db:"updated_by" json:"updated_by"`
}
