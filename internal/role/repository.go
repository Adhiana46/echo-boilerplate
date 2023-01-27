package role

import (
	"context"

	"github.com/Adhiana46/echo-boilerplate/entity"
)

type RoleRepository interface {
	Create(ctx context.Context, e *entity.Role) (*entity.Role, error)
	Update(ctx context.Context, e *entity.Role) (*entity.Role, error)
	Destroy(ctx context.Context, e *entity.Role) error
	FindById(ctx context.Context, id int) (*entity.Role, error)
	FindByUuid(ctx context.Context, uuid string) (*entity.Role, error)
	FindByName(ctx context.Context, name string) (*entity.Role, error)
	FindAll(ctx context.Context, offset int, limit int, sorts map[string]string, search string) ([]*entity.Role, error)

	CountByName(ctx context.Context, name string) (int, error)
	CountAll(ctx context.Context, search string) (int, error)
}
