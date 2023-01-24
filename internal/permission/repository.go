package permission

import (
	"context"

	"github.com/Adhiana46/echo-boilerplate/entity"
)

type PermissionRepository interface {
	Create(ctx context.Context, e *entity.Permission) (*entity.Permission, error)
	Update(ctx context.Context, e *entity.Permission) (*entity.Permission, error)
	Destroy(ctx context.Context, e *entity.Permission) error
	FindById(ctx context.Context, id int) (*entity.Permission, error)
	FindByUuid(ctx context.Context, uuid string) (*entity.Permission, error)
	FindAll(ctx context.Context, offset int, limit int) ([]*entity.Permission, error)

	CountByName(ctx context.Context, name string) (int, error)
}
