package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Adhiana46/echo-boilerplate/entity"
	"github.com/Adhiana46/echo-boilerplate/internal/role"
	"github.com/Adhiana46/echo-boilerplate/internal/user"
	"github.com/Adhiana46/echo-boilerplate/pkg/errors"
	"github.com/jmoiron/sqlx"
)

type pgUserRepository struct {
	db       *sqlx.DB
	roleRepo role.RoleRepository
}

func NewPgRoleRepository(db *sqlx.DB, roleRepo role.RoleRepository) user.UserRepository {
	return &pgUserRepository{
		db:       db,
		roleRepo: roleRepo,
	}
}

func (r *pgUserRepository) Create(ctx context.Context, e *entity.User) (*entity.User, error) {
	sql := `
		INSERT INTO users
		(uuid, username, email, password, name, role_id, status, last_login_at, created_at, created_by, updated_at, updated_by)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`

	insertId := 0
	err := r.db.QueryRowContext(
		ctx,
		sql,
		e.Uuid,
		e.Username,
		e.Email,
		e.Password,
		e.Name,
		e.RoleId,
		e.Status,
		e.LastLoginAt,
		e.CreatedAt,
		e.CreatedBy,
		e.UpdatedAt,
		e.UpdatedBy,
	).Scan(&insertId)
	if err != nil {
		return nil, err
	}

	return r.FindById(ctx, insertId)
}

func (r *pgUserRepository) Update(ctx context.Context, e *entity.User) (*entity.User, error) {
	sql := `
		UPDATE users
			SET username = $1, 
				email = $2, 
				password = $3, 
				name = $4, 
				role_id = $5, 
				status = $6, 
				last_login_at = $7, 
				updated_at = $8, 
				updated_by = $9
		WHERE id = $10
	`

	_, err := r.db.ExecContext(
		ctx,
		sql,
		e.Username,
		e.Email,
		e.Password,
		e.Name,
		e.RoleId,
		e.Status,
		e.LastLoginAt,
		e.UpdatedAt,
		e.UpdatedBy,
		e.Id,
	)
	if err != nil {
		return nil, err
	}

	return r.FindById(ctx, e.Id)
}

func (r *pgUserRepository) Destroy(ctx context.Context, e *entity.User) error {
	sql := `DELETE FROM users WHERE id = $1`

	_, err := r.db.ExecContext(ctx, sql, e.Id)

	return err
}

func (r *pgUserRepository) FindById(ctx context.Context, id int) (*entity.User, error) {
	sql := `
		SELECT id, uuid, username, email, password, name, role_id, status, last_login_at, created_at, created_by, updated_at, updated_by
		FROM users
		WHERE id = $1
	`

	row := &entity.User{}
	err := r.db.GetContext(ctx, row, sql, id)
	if err != nil {
		return nil, err
	}

	role, err := r.roleRepo.FindById(ctx, row.RoleId)
	if err != nil {
		return nil, err
	}

	row.Role = role

	return row, nil
}

func (r *pgUserRepository) FindByUuid(ctx context.Context, uuid string) (*entity.User, error) {
	sql := `
		SELECT id, uuid, username, email, password, name, role_id, status, last_login_at, created_at, created_by, updated_at, updated_by
		FROM users
		WHERE uuid = $1
	`

	row := &entity.User{}
	err := r.db.GetContext(ctx, row, sql, uuid)
	if err != nil {
		return nil, err
	}

	role, err := r.roleRepo.FindById(ctx, row.RoleId)
	if err != nil {
		return nil, err
	}

	row.Role = role

	return row, nil
}

func (r *pgUserRepository) FindAll(ctx context.Context, offset int, limit int, sorts map[string]string, search string) ([]*entity.User, error) {
	sql := `
		SELECT id, uuid, username, email, password, name, role_id, status, last_login_at, created_at, created_by, updated_at, updated_by
		FROM users
	`
	aWheres := []string{}
	aOrders := []string{}

	// search
	if search != "" {
		aWheres = append(aWheres, fmt.Sprintf("LOWER(username) LIKE '%s'", "%"+search+"%"))
		aWheres = append(aWheres, fmt.Sprintf("LOWER(email) LIKE '%s'", "%"+search+"%"))
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

	rows := []*entity.User{}
	err := r.db.SelectContext(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *pgUserRepository) CountByUsername(ctx context.Context, username string) (int, error) {
	sql := `
		SELECT COUNT(id) AS numrows
		FROM users
		WHERE username = $1
	`

	numrows := 0
	err := r.db.QueryRow(sql, username).Scan(&numrows)
	if err != nil {
		return 0, err
	}

	return numrows, nil
}

func (r *pgUserRepository) CountByEmail(ctx context.Context, email string) (int, error) {
	sql := `
		SELECT COUNT(id) AS numrows
		FROM users
		WHERE email = $1
	`

	numrows := 0
	err := r.db.QueryRow(sql, email).Scan(&numrows)
	if err != nil {
		return 0, err
	}

	return numrows, nil
}

func (r *pgUserRepository) CountAll(ctx context.Context, search string) (int, error) {
	sql := `
		SELECT COUNT(id) AS numrows
		FROM users
	`
	aWheres := []string{}

	// search
	if search != "" {
		aWheres = append(aWheres, fmt.Sprintf("LOWER(username) LIKE '%s'", "%"+search+"%"))
		aWheres = append(aWheres, fmt.Sprintf("LOWER(email) LIKE '%s'", "%"+search+"%"))
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
