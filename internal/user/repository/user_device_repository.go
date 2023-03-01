package repository

import (
	"context"
	"sync"

	"github.com/Adhiana46/echo-boilerplate/entity"
	"github.com/Adhiana46/echo-boilerplate/internal/user"
)

var (
	userDeviceRepoInstance     *userDeviceRepository
	userDeviceRepoInstanceOnce sync.Once
)

type userDeviceRepository struct {
	persistent user.UserDevicePersistent
}

func NewUserDeviceRepository(persistent user.UserDevicePersistent) user.UserDeviceRepository {
	userDeviceRepoInstanceOnce.Do(func() {
		userDeviceRepoInstance = &userDeviceRepository{
			persistent: persistent,
		}
	})

	return userDeviceRepoInstance
}

func (r *userDeviceRepository) Create(ctx context.Context, e *entity.UserDevice) (*entity.UserDevice, error) {
	return r.persistent.Create(ctx, e)
}

func (r *userDeviceRepository) Update(ctx context.Context, e *entity.UserDevice) (*entity.UserDevice, error) {
	return r.persistent.Update(ctx, e)
}

func (r *userDeviceRepository) Destroy(ctx context.Context, e *entity.UserDevice) error {
	return r.persistent.Destroy(ctx, e)
}

func (r *userDeviceRepository) FindById(ctx context.Context, id int) (*entity.UserDevice, error) {
	return r.persistent.FindById(ctx, id)
}

func (r *userDeviceRepository) FindByUuid(ctx context.Context, uuid string) (*entity.UserDevice, error) {
	return r.persistent.FindByUuid(ctx, uuid)
}

func (r *userDeviceRepository) FindByToken(ctx context.Context, userId int, token string) (*entity.UserDevice, error) {
	return r.persistent.FindByToken(ctx, userId, token)
}
