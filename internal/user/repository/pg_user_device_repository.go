package repository

import (
	"context"

	"github.com/Adhiana46/echo-boilerplate/entity"
	"github.com/Adhiana46/echo-boilerplate/internal/user"
	"github.com/jmoiron/sqlx"
)

type pgUserDeviceRepository struct {
	db *sqlx.DB
}

func NewPgUserDeviceRepository(db *sqlx.DB) user.UserDeviceRepository {
	return &pgUserDeviceRepository{
		db: db,
	}
}

func (r *pgUserDeviceRepository) Create(ctx context.Context, e *entity.UserDevice) (*entity.UserDevice, error) {
	sql := `
		INSERT INTO user_devices
		(uuid, user_id, token, ip, location, platform, user_agent, app_version, vendor, created_at, created_by, updated_at, updated_by)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`

	insertId := 0
	err := r.db.QueryRowContext(
		ctx,
		sql,
		e.Uuid,
		e.UserId,
		e.Token,
		e.IP,
		e.Location,
		e.Platform,
		e.UserAgent,
		e.AppVersion,
		e.Vendor,
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

func (r *pgUserDeviceRepository) Update(ctx context.Context, e *entity.UserDevice) (*entity.UserDevice, error) {
	sql := `
		UPDATE user_devices
			SET user_id = $1,
				token = $2,
				ip = $3,
				location = $4,
				platform = $5,
				user_agent = $6,
				app_version = $7,
				vendor = $8,
				updated_at = $9,
				updated_by = $10
		WHERE id = $11
	`

	_, err := r.db.ExecContext(
		ctx,
		sql,
		e.UserId,
		e.Token,
		e.IP,
		e.Location,
		e.Platform,
		e.UserAgent,
		e.AppVersion,
		e.Vendor,
		e.UpdatedAt,
		e.UpdatedBy,
		e.Id,
	)
	if err != nil {
		return nil, err
	}

	return r.FindById(ctx, e.Id)
}

func (r *pgUserDeviceRepository) Destroy(ctx context.Context, e *entity.UserDevice) error {
	sql := `DELETE FROM user_devices WHERE id = $1`

	_, err := r.db.ExecContext(ctx, sql, e.Id)

	return err
}

func (r *pgUserDeviceRepository) FindById(ctx context.Context, id int) (*entity.UserDevice, error) {
	sql := `
		SELECT id, uuid, user_id, token, ip, location, platform, user_agent, app_version, vendor, created_at, created_by, updated_at, updated_by
		FROM user_devices
		WHERE id = $1
	`

	row := &entity.UserDevice{}
	err := r.db.GetContext(ctx, row, sql, id)
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (r *pgUserDeviceRepository) FindByUuid(ctx context.Context, uuid string) (*entity.UserDevice, error) {
	sql := `
		SELECT id, uuid, user_id, token, ip, location, platform, user_agent, app_version, vendor, created_at, created_by, updated_at, updated_by
		FROM user_devices
		WHERE uuid = $1
	`

	row := &entity.UserDevice{}
	err := r.db.GetContext(ctx, row, sql, uuid)
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (r *pgUserDeviceRepository) FindByToken(ctx context.Context, userId int, token string) (*entity.UserDevice, error) {
	sql := `
		SELECT id, uuid, user_id, token, ip, location, platform, user_agent, app_version, vendor, created_at, created_by, updated_at, updated_by
		FROM user_devices
		WHERE user_id = $1 AND token = $2
	`

	row := &entity.UserDevice{}
	err := r.db.GetContext(ctx, row, sql, userId, token)
	if err != nil {
		return nil, err
	}

	return row, nil
}
