package user

import (
	"context"

	"github.com/Adhiana46/echo-boilerplate/entity"
)

type UserRepository interface {
	Create(ctx context.Context, e *entity.User) (*entity.User, error)
	Update(ctx context.Context, e *entity.User) (*entity.User, error)
	Destroy(ctx context.Context, e *entity.User) error
	FindById(ctx context.Context, id int) (*entity.User, error)
	FindByUuid(ctx context.Context, uuid string) (*entity.User, error)
	FindByUsernameOrEmail(ctx context.Context, username string) (*entity.User, error)
	FindAll(ctx context.Context, offset int, limit int, sorts map[string]string, search string) ([]*entity.User, error)

	CountByUsername(ctx context.Context, username string) (int, error)
	CountByEmail(ctx context.Context, email string) (int, error)
	CountAll(ctx context.Context, search string) (int, error)
}

type UserDeviceRepository interface {
	Create(ctx context.Context, e *entity.UserDevice) (*entity.UserDevice, error)
	Update(ctx context.Context, e *entity.UserDevice) (*entity.UserDevice, error)
	Destroy(ctx context.Context, e *entity.UserDevice) error
	FindById(ctx context.Context, id int) (*entity.UserDevice, error)
	FindByUuid(ctx context.Context, uuid string) (*entity.UserDevice, error)
	FindByToken(ctx context.Context, userId int, token string) (*entity.UserDevice, error)
}
