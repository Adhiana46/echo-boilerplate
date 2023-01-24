package repository

import (
	"context"

	"github.com/Adhiana46/echo-boilerplate/entity"
	"github.com/Adhiana46/echo-boilerplate/internal/permission"
	"github.com/jmoiron/sqlx"
)

type pgPermissionRepository struct {
	db *sqlx.DB
}

func NewPgPermissionRepository(db *sqlx.DB) permission.PermissionRepository {
	return &pgPermissionRepository{
		db: db,
	}
}

func (r *pgPermissionRepository) Create(ctx context.Context, e *entity.Permission) (*entity.Permission, error) {
	sql := `
		INSERT INTO permissions
		(uuid, parent_id, name, type, created_at, created_by, updated_at, updated_by)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	insertedId := 0
	err := r.db.QueryRowContext(
		ctx,
		sql,
		e.Uuid,
		e.ParentId,
		e.Name,
		e.Type,
		e.CreatedAt,
		e.CreatedBy,
		e.UpdatedAt,
		e.UpdatedBy,
	).Scan(&insertedId)

	if err != nil {
		return nil, err
	}

	return r.FindById(ctx, insertedId)
}

func (r *pgPermissionRepository) Update(ctx context.Context, e *entity.Permission) (*entity.Permission, error) {
	panic("Not implemented")
}

func (r *pgPermissionRepository) Destroy(ctx context.Context, e *entity.Permission) error {
	panic("Not implemented")
}

func (r *pgPermissionRepository) FindById(ctx context.Context, id int) (*entity.Permission, error) {
	sql := `
		SELECT id, uuid, parent_id, name, type, created_at, created_by, updated_at, updated_by
		FROM permissions
		WHERE id = $1
	`

	e := &entity.Permission{}
	err := r.db.GetContext(ctx, e, sql, id)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (r *pgPermissionRepository) FindByUuid(ctx context.Context, uuid string) (*entity.Permission, error) {
	sql := `
		SELECT id, uuid, parent_id, name, type, created_at, created_by, updated_at, updated_by
		FROM permissions
		WHERE uuid = $1
	`

	e := &entity.Permission{}
	err := r.db.GetContext(ctx, e, sql, uuid)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (r *pgPermissionRepository) FindAll(ctx context.Context, offset int, limit int) ([]*entity.Permission, error) {
	panic("Not implemented")
}

func (r *pgPermissionRepository) CountByName(ctx context.Context, name string) (int, error) {
	sql := `
		SELECT COUNT(id) AS numrows
		FROM permissions
		WHERE name = $1
	`

	numrows := 0
	err := r.db.QueryRow(sql, name).Scan(&numrows)
	if err != nil {
		return 0, err
	}

	return numrows, nil
}
