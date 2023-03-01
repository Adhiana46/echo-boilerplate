package data

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/Adhiana46/echo-boilerplate/entity"
	"github.com/Adhiana46/echo-boilerplate/internal/permission"
	"github.com/Adhiana46/echo-boilerplate/pkg/errors"
	"github.com/jmoiron/sqlx"
)

var (
	postgresPersistInstance     *postgresPermissionPersistent
	postgresPersistInstanceOnce sync.Once
)

type postgresPermissionPersistent struct {
	db *sqlx.DB
}

func NewPostgresPermissionPersistent(db *sqlx.DB) permission.PermissionPersistent {
	postgresPersistInstanceOnce.Do(func() {
		postgresPersistInstance = &postgresPermissionPersistent{
			db: db,
		}
	})

	return postgresPersistInstance
}

func (r *postgresPermissionPersistent) Create(ctx context.Context, e *entity.Permission) (*entity.Permission, error) {
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

func (r *postgresPermissionPersistent) Update(ctx context.Context, e *entity.Permission) (*entity.Permission, error) {
	sql := `
		UPDATE permissions
		SET parent_id = $1, 
			name = $2, 
			type = $3, 
			updated_at = $4, 
			updated_by = $5
		WHERE
			uuid = $6
	`

	err := r.db.QueryRowContext(
		ctx,
		sql,
		e.ParentId,
		e.Name,
		e.Type,
		e.UpdatedAt,
		e.UpdatedBy,
		e.Uuid,
	).Err()

	if err != nil {
		return nil, err
	}

	return r.FindByUuid(ctx, e.Uuid)
}

func (r *postgresPermissionPersistent) Destroy(ctx context.Context, e *entity.Permission) error {
	sql := `DELETE FROM permissions WHERE id = $1`

	_, err := r.db.ExecContext(ctx, sql, e.Id)

	return err
}

func (r *postgresPermissionPersistent) FindById(ctx context.Context, id int) (*entity.Permission, error) {
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

func (r *postgresPermissionPersistent) FindByUuid(ctx context.Context, uuid string) (*entity.Permission, error) {
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

func (r *postgresPermissionPersistent) FindAll(ctx context.Context, offset int, limit int, sorts map[string]string, search string) ([]*entity.Permission, error) {
	sql := `
		SELECT id, uuid, parent_id, name, type, created_at, created_by, updated_at, updated_by
		FROM permissions
	`
	aWheres := []string{}
	aOrders := []string{}

	// search
	if search != "" {
		aWheres = append(aWheres, fmt.Sprintf("LOWER(name) LIKE '%s'", "%"+search+"%"))
	}
	if len(aWheres) > 0 {
		sql += " WHERE " + strings.Join(aWheres, " AND ")
	}

	// orders
	if len(sorts) > 0 {
		for field, dir := range sorts {
			if strings.ToLower(dir) != "asc" && strings.ToLower(dir) != "desc" {
				return nil, errors.NewBadRequestError(fmt.Sprintf("Order direction for field '%s' should be 'asc' or 'desc'", field))
			}

			aOrders = append(aOrders, fmt.Sprintf("%s %s", field, dir))
		}
	}
	if len(aOrders) > 0 {
		sql += " ORDER BY " + strings.Join(aOrders, ", ")
	}

	// limit offset
	sql += fmt.Sprintf(" OFFSET %v LIMIT %v ", offset, limit)

	rows := []*entity.Permission{}
	err := r.db.SelectContext(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *postgresPermissionPersistent) FindAllByNames(ctx context.Context, names []string) ([]*entity.Permission, error) {
	sql, args, err := sqlx.In(`
		SELECT id, uuid, parent_id, name, type, created_at, created_by, updated_at, updated_by
		FROM permissions
		WHERE name IN (?)
	`, names)

	sql = r.db.Rebind(sql)

	if err != nil {
		return nil, err
	}

	rows := []*entity.Permission{}
	err = r.db.SelectContext(ctx, &rows, sql, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *postgresPermissionPersistent) CountByName(ctx context.Context, name string) (int, error) {
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

func (r *postgresPermissionPersistent) CountAll(ctx context.Context, search string) (int, error) {
	sql := `
		SELECT COUNT(id) AS numrows
		FROM permissions
	`
	aWheres := []string{}

	// search
	if search != "" {
		aWheres = append(aWheres, fmt.Sprintf("LOWER(name) LIKE '%s'", "%"+search+"%"))
	}
	if len(aWheres) > 0 {
		sql += " WHERE " + strings.Join(aWheres, " AND ")
	}

	numrows := 0
	err := r.db.QueryRow(sql).Scan(&numrows)
	if err != nil {
		return 0, err
	}

	return numrows, nil
}
