package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Adhiana46/echo-boilerplate/dto"
	"github.com/Adhiana46/echo-boilerplate/entity"
	"github.com/Adhiana46/echo-boilerplate/internal/permission"
	"github.com/Adhiana46/echo-boilerplate/pkg/errors"
	"github.com/google/uuid"
)

type permissionUsecase struct {
	repo permission.PermissionRepository
}

func NewPermissionUsecase(repo permission.PermissionRepository) permission.PermissionUsecase {
	return &permissionUsecase{
		repo: repo,
	}
}

func (uc *permissionUsecase) CreatePermission(ctx context.Context, input *dto.CreatePermissionRequest) (*dto.PermissionResponse, error) {
	// validation logic
	numrows, err := uc.repo.CountByName(ctx, input.Name)
	if err != nil {
		return nil, err
	}

	if numrows > 0 {
		return nil, errors.NewBadRequestError(fmt.Sprintf("Permission with name '%s' already exists", input.Name))
	}

	e, err := uc.repo.Create(ctx, &entity.Permission{
		Uuid:     uuid.NewString(),
		ParentId: input.ParentId,
		Name:     input.Name,
		Type:     input.Type,
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		CreatedBy: 0, // todo get from ctx
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedBy: 0, // todo get from ctx
	})

	if err != nil {
		return nil, err
	}

	return dto.NewPermissionResponse(e), nil
}

func (uc *permissionUsecase) UpdatePermission(ctx context.Context, input *dto.UpdatePermissionRequest) (*dto.PermissionResponse, error) {
	e, err := uc.repo.FindByUuid(ctx, input.Uuid)
	if err != nil {
		return nil, err
	}

	// validation logic
	if e.Name != input.Name {
		numrows, err := uc.repo.CountByName(ctx, input.Name)
		if err != nil {
			return nil, err
		}

		if numrows > 0 {
			return nil, errors.NewBadRequestError(fmt.Sprintf("Permission with name '%s' already exists", input.Name))
		}
	}

	e.ParentId = input.ParentId
	e.Name = input.Name
	e.Type = input.Type
	e.UpdatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	e.UpdatedBy = 0 // TODO: user

	updatedE, err := uc.repo.Update(ctx, e)
	if err != nil {
		return nil, err
	}

	return dto.NewPermissionResponse(updatedE), nil
}

func (uc *permissionUsecase) DeletePermission(ctx context.Context, input *dto.DeletePermissionRequest) (*dto.PermissionResponse, error) {
	e, err := uc.repo.FindByUuid(ctx, input.Uuid)
	if err != nil {
		return nil, err
	}

	err = uc.repo.Destroy(ctx, e)

	return nil, err
}

func (uc *permissionUsecase) Get(ctx context.Context, input *dto.GetPermissionRequest) (*dto.PermissionResponse, error) {
	e, err := uc.repo.FindByUuid(ctx, input.Uuid)
	if err != nil {
		return nil, err
	}

	return dto.NewPermissionResponse(e), nil
}

func (uc *permissionUsecase) GetList(ctx context.Context, input *dto.GetListPermissionRequest) (*dto.PermissionResponse, error) {
	// TODO: validation logic

	panic("not implemented.")
}
