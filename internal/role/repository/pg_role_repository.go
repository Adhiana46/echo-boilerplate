package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Adhiana46/echo-boilerplate/entity"
	"github.com/Adhiana46/echo-boilerplate/internal/role"
	"github.com/Adhiana46/echo-boilerplate/pkg/errors"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type pgRoleRepository struct {
	db *sqlx.DB
}

func NewPgRoleRepository(db *sqlx.DB) role.RoleRepository {
	return &pgRoleRepository{
		db: db,
	}
}

func (r *pgRoleRepository) Create(ctx context.Context, e *entity.Role) (*entity.Role, error) {
	// set squirrel
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	sqlInsertRole := `
		INSERT INTO roles
		(uuid, name, created_at, created_by, updated_at, updated_by)
		VALUES
		($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	// Insert Role
	roleId := 0
	err = tx.QueryRowContext(
		ctx,
		sqlInsertRole,
		e.Uuid,
		e.Name,
		e.CreatedAt,
		e.CreatedBy,
		e.UpdatedAt,
		e.UpdatedBy,
	).Scan(&roleId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert role_permissions
	qInsertRolePerms := sq.Insert("role_permissions").Columns("role_id", "permission_id")
	for _, perm := range e.Permissions {
		qInsertRolePerms = qInsertRolePerms.Values(roleId, perm.Id)
	}

	sqlInsertRolePerms, args, err := qInsertRolePerms.ToSql()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.ExecContext(ctx, sqlInsertRolePerms, args...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return r.FindById(ctx, roleId)
}

func (r *pgRoleRepository) Update(ctx context.Context, e *entity.Role) (*entity.Role, error) {
	// set squirrel
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sqlUpdateRole := `
		UPDATE roles
		SET name = $1,
			updated_at = $2,
			updated_by = $3
		WHERE id = $4
	`

	sqlDeleteRolePermissions := `
		DELETE FROM role_permissions
		WHERE role_id = $1
	`

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, sqlUpdateRole, e.Name, e.UpdatedAt, e.UpdatedBy, e.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.ExecContext(ctx, sqlDeleteRolePermissions, e.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert role_permissions
	qInsertRolePerms := sq.Insert("role_permissions").Columns("role_id", "permission_id")
	for _, perm := range e.Permissions {
		qInsertRolePerms = qInsertRolePerms.Values(e.Id, perm.Id)
	}

	sqlInsertRolePerms, args, err := qInsertRolePerms.ToSql()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.ExecContext(ctx, sqlInsertRolePerms, args...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return r.FindById(ctx, e.Id)
}

func (r *pgRoleRepository) Destroy(ctx context.Context, e *entity.Role) error {
	sql := `DELETE FROM roles WHERE id = $1`

	_, err := r.db.ExecContext(ctx, sql, e.Id)

	return err
}

func (r *pgRoleRepository) FindById(ctx context.Context, id int) (*entity.Role, error) {
	sql := `
		SELECT id, uuid, name, created_at, created_by, updated_at, updated_by
		FROM roles
		WHERE id = $1
	`
	sqlPerms := `
		SELECT id, uuid, parent_id, name, type, created_at, created_by, updated_at, updated_by
		FROM permissions
		WHERE id IN (SELECT permission_id FROM role_permissions WHERE role_id = $1)
	`

	e := &entity.Role{}
	err := r.db.GetContext(ctx, e, sql, id)
	if err != nil {
		return nil, err
	}

	perms := []*entity.Permission{}
	err = r.db.SelectContext(ctx, &perms, sqlPerms, e.Id)
	if err != nil {
		return nil, err
	}

	e.Permissions = perms

	return e, nil
}

func (r *pgRoleRepository) FindByUuid(ctx context.Context, uuid string) (*entity.Role, error) {
	sql := `
		SELECT id, uuid, name, created_at, created_by, updated_at, updated_by
		FROM roles
		WHERE uuid = $1
	`
	sqlPerms := `
		SELECT id, uuid, parent_id, name, type, created_at, created_by, updated_at, updated_by
		FROM permissions
		WHERE id IN (SELECT permission_id FROM role_permissions WHERE role_id = $1)
	`

	e := &entity.Role{}
	err := r.db.GetContext(ctx, e, sql, uuid)
	if err != nil {
		return nil, err
	}

	perms := []*entity.Permission{}
	err = r.db.SelectContext(ctx, &perms, sqlPerms, e.Id)
	if err != nil {
		return nil, err
	}

	e.Permissions = perms

	return e, nil
}

func (r *pgRoleRepository) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	sql := `
		SELECT id, uuid, name, created_at, created_by, updated_at, updated_by
		FROM roles
		WHERE name = $1
	`
	sqlPerms := `
		SELECT id, uuid, parent_id, name, type, created_at, created_by, updated_at, updated_by
		FROM permissions
		WHERE id IN (SELECT permission_id FROM role_permissions WHERE role_id = $1)
	`

	e := &entity.Role{}
	err := r.db.GetContext(ctx, e, sql, name)
	if err != nil {
		return nil, err
	}

	perms := []*entity.Permission{}
	err = r.db.SelectContext(ctx, &perms, sqlPerms, e.Id)
	if err != nil {
		return nil, err
	}

	e.Permissions = perms

	return e, nil
}

func (r *pgRoleRepository) FindAll(ctx context.Context, offset int, limit int, sorts map[string]string, search string) ([]*entity.Role, error) {
	sql := `
		SELECT id, uuid, name, created_at, created_by, updated_at, updated_by
		FROM roles
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

	rows := []*entity.Role{}
	err := r.db.SelectContext(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *pgRoleRepository) CountByName(ctx context.Context, name string) (int, error) {
	sql := `
		SELECT COUNT(id) AS numrows
		FROM roles
		WHERE name = $1
	`

	numrows := 0
	err := r.db.QueryRow(sql, name).Scan(&numrows)
	if err != nil {
		return 0, err
	}

	return numrows, nil
}

func (r *pgRoleRepository) CountAll(ctx context.Context, search string) (int, error) {
	sql := `
		SELECT COUNT(id) AS numrows
		FROM roles
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
