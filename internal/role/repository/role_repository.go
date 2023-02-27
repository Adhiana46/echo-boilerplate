package repository

import (
	"context"

	"github.com/Adhiana46/echo-boilerplate/entity"
	"github.com/Adhiana46/echo-boilerplate/internal/role"
)

type roleRepository struct {
	persistent role.RolePersistent
}

func NewRoleRepository(persistent role.RolePersistent) role.RoleRepository {
	return &roleRepository{
		persistent: persistent,
	}
}

func (r *roleRepository) Create(ctx context.Context, e *entity.Role) (*entity.Role, error) {
	return r.persistent.Create(ctx, e)
}

func (r *roleRepository) Update(ctx context.Context, e *entity.Role) (*entity.Role, error) {
	return r.persistent.Update(ctx, e)
}

func (r *roleRepository) Destroy(ctx context.Context, e *entity.Role) error {
	return r.persistent.Destroy(ctx, e)
}

func (r *roleRepository) FindById(ctx context.Context, id int) (*entity.Role, error) {
	return r.persistent.FindById(ctx, id)
}

func (r *roleRepository) FindByUuid(ctx context.Context, uuid string) (*entity.Role, error) {
	return r.persistent.FindByUuid(ctx, uuid)
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	return r.persistent.FindByName(ctx, name)
}

func (r *roleRepository) FindAll(ctx context.Context, offset int, limit int, sorts map[string]string, search string) ([]*entity.Role, error) {
	return r.persistent.FindAll(ctx, offset, limit, sorts, search)
}

func (r *roleRepository) CountByName(ctx context.Context, name string) (int, error) {
	return r.persistent.CountByName(ctx, name)
}

func (r *roleRepository) CountAll(ctx context.Context, search string) (int, error) {
	return r.persistent.CountAll(ctx, search)
}
