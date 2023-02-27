package permission

import (
	"context"

	"github.com/Adhiana46/echo-boilerplate/entity"
)

type PermissionPersistent interface {
	Create(ctx context.Context, e *entity.Permission) (*entity.Permission, error)
	Update(ctx context.Context, e *entity.Permission) (*entity.Permission, error)
	Destroy(ctx context.Context, e *entity.Permission) error
	FindById(ctx context.Context, id int) (*entity.Permission, error)
	FindByUuid(ctx context.Context, uuid string) (*entity.Permission, error)
	FindAll(ctx context.Context, offset int, limit int, sorts map[string]string, search string) ([]*entity.Permission, error)
	FindAllByNames(ctx context.Context, names []string) ([]*entity.Permission, error)

	CountByName(ctx context.Context, name string) (int, error)
	CountAll(ctx context.Context, search string) (int, error)
}
