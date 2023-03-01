package repository

import (
	"context"
	"sync"

	"github.com/Adhiana46/echo-boilerplate/entity"
	"github.com/Adhiana46/echo-boilerplate/internal/role"
	"github.com/Adhiana46/echo-boilerplate/internal/user"
)

var (
	userRepoInstance     *userRepository
	userRepoInstanceOnce sync.Once
)

type userRepository struct {
	userPersistent user.UserPersistent
	rolePersistent role.RolePersistent
}

func NewUserRepository(userPersistent user.UserPersistent, rolePersistent role.RolePersistent) user.UserRepository {
	userRepoInstanceOnce.Do(func() {
		userRepoInstance = &userRepository{
			userPersistent: userPersistent,
			rolePersistent: rolePersistent,
		}
	})

	return userRepoInstance
}

func (r *userRepository) Create(ctx context.Context, e *entity.User) (*entity.User, error) {
	return r.userPersistent.Create(ctx, e)
}

func (r *userRepository) Update(ctx context.Context, e *entity.User) (*entity.User, error) {
	return r.userPersistent.Update(ctx, e)
}

func (r *userRepository) Destroy(ctx context.Context, e *entity.User) error {
	return r.userPersistent.Destroy(ctx, e)
}

func (r *userRepository) FindById(ctx context.Context, id int) (*entity.User, error) {
	userEntity, err := r.userPersistent.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	roleEntity, err := r.rolePersistent.FindById(ctx, userEntity.RoleId)
	if err != nil {
		return nil, err
	}

	userEntity.Role = roleEntity

	return userEntity, nil
}

func (r *userRepository) FindByUuid(ctx context.Context, uuid string) (*entity.User, error) {
	userEntity, err := r.userPersistent.FindByUuid(ctx, uuid)
	if err != nil {
		return nil, err
	}

	roleEntity, err := r.rolePersistent.FindById(ctx, userEntity.RoleId)
	if err != nil {
		return nil, err
	}

	userEntity.Role = roleEntity

	return userEntity, nil
}

func (r *userRepository) FindByUsernameOrEmail(ctx context.Context, username string) (*entity.User, error) {
	userEntity, err := r.userPersistent.FindByUsernameOrEmail(ctx, username)
	if err != nil {
		return nil, err
	}

	roleEntity, err := r.rolePersistent.FindById(ctx, userEntity.RoleId)
	if err != nil {
		return nil, err
	}

	userEntity.Role = roleEntity

	return userEntity, nil
}

func (r *userRepository) FindAll(ctx context.Context, offset int, limit int, sorts map[string]string, search string) ([]*entity.User, error) {
	return r.userPersistent.FindAll(ctx, offset, limit, sorts, search)
}

func (r *userRepository) CountByUsername(ctx context.Context, username string) (int, error) {
	return r.userPersistent.CountByUsername(ctx, username)
}

func (r *userRepository) CountByEmail(ctx context.Context, email string) (int, error) {
	return r.userPersistent.CountByEmail(ctx, email)
}

func (r *userRepository) CountAll(ctx context.Context, search string) (int, error) {
	return r.userPersistent.CountAll(ctx, search)
}
