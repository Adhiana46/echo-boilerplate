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

	createdAt := ""
	updatedAt := ""

	if e.CreatedAt.Valid {
		createdAt = e.CreatedAt.Time.Format(time.RFC3339)
	}
	if e.UpdatedAt.Valid {
		updatedAt = e.UpdatedAt.Time.Format(time.RFC3339)
	}

	return &dto.PermissionResponse{
		Uuid:      e.Uuid,
		ParentId:  e.ParentId,
		Name:      e.Name,
		Type:      e.Type,
		CreatedAt: createdAt,
		CreatedBy: e.CreatedBy,
		UpdatedAt: updatedAt,
		UpdatedBy: e.UpdatedBy,
	}, err
}

func (uc *permissionUsecase) UpdatePermission(ctx context.Context, input *dto.UpdatePermissionRequest) (*dto.PermissionResponse, error) {
	// TODO: validation logic

	panic("not implemented.")
}

func (uc *permissionUsecase) DeletePermission(ctx context.Context, input *dto.DeletePermissionRequest) (*dto.PermissionResponse, error) {
	// TODO: validation logic

	panic("not implemented.")
}

func (uc *permissionUsecase) Get(ctx context.Context, input *dto.GetPermissionRequest) (*dto.PermissionResponse, error) {
	// TODO: validation logic

	panic("not implemented.")
}

func (uc *permissionUsecase) GetList(ctx context.Context, input *dto.GetListPermissionRequest) (*dto.PermissionResponse, error) {
	// TODO: validation logic

	panic("not implemented.")
}
