package repository

import (
	"context"

	"github.com/Adhiana46/echo-boilerplate/entity"
	"github.com/Adhiana46/echo-boilerplate/internal/permission"
)

type permissionRepository struct {
	persistent permission.PermissionPersistent
}

func NewPermissionRepository(persistent permission.PermissionPersistent) permission.PermissionRepository {
	return &permissionRepository{
		persistent: persistent,
	}
}

func (r *permissionRepository) Create(ctx context.Context, e *entity.Permission) (*entity.Permission, error) {
	return r.persistent.Create(ctx, e)
}

func (r *permissionRepository) Update(ctx context.Context, e *entity.Permission) (*entity.Permission, error) {
	return r.persistent.Update(ctx, e)
}

func (r *permissionRepository) Destroy(ctx context.Context, e *entity.Permission) error {
	return r.persistent.Destroy(ctx, e)
}

func (r *permissionRepository) FindById(ctx context.Context, id int) (*entity.Permission, error) {
	return r.persistent.FindById(ctx, id)
}

func (r *permissionRepository) FindByUuid(ctx context.Context, uuid string) (*entity.Permission, error) {
	return r.persistent.FindByUuid(ctx, uuid)
}

func (r *permissionRepository) FindAll(ctx context.Context, offset int, limit int, sorts map[string]string, search string) ([]*entity.Permission, error) {
	return r.persistent.FindAll(ctx, offset, limit, sorts, search)
}

func (r *permissionRepository) FindAllByNames(ctx context.Context, names []string) ([]*entity.Permission, error) {
	return r.persistent.FindAllByNames(ctx, names)
}

func (r *permissionRepository) CountByName(ctx context.Context, name string) (int, error) {
	return r.persistent.CountByName(ctx, name)
}

func (r *permissionRepository) CountAll(ctx context.Context, search string) (int, error) {
	return r.persistent.CountAll(ctx, search)
}
