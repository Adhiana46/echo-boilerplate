package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Adhiana46/echo-boilerplate/constants"
	"github.com/Adhiana46/echo-boilerplate/dto"
	"github.com/Adhiana46/echo-boilerplate/entity"
	"github.com/Adhiana46/echo-boilerplate/internal/permission"
	"github.com/Adhiana46/echo-boilerplate/pkg/errors"
	"github.com/Adhiana46/echo-boilerplate/pkg/utils"
	"github.com/google/uuid"
)

var (
	permissionUcInstance     *permissionUsecase
	permissionUcInstanceOnce sync.Once
)

type permissionUsecase struct {
	repo permission.PermissionRepository
}

func NewPermissionUsecase(repo permission.PermissionRepository) permission.PermissionUsecase {
	permissionUcInstanceOnce.Do(func() {
		permissionUcInstance = &permissionUsecase{
			repo: repo,
		}
	})
	return permissionUcInstance
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

	createdBy := sql.NullInt64{}
	updatedBy := sql.NullInt64{}
	user := utils.GetUserFromContext(ctx)
	if user != nil {
		createdBy = sql.NullInt64{
			Int64: int64(user.ID),
			Valid: true,
		}

		updatedBy = sql.NullInt64{
			Int64: int64(user.ID),
			Valid: true,
		}
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
		CreatedBy: createdBy,
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedBy: updatedBy,
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

	updatedBy := sql.NullInt64{}
	user := utils.GetUserFromContext(ctx)
	if user != nil {
		updatedBy = sql.NullInt64{
			Int64: int64(user.ID),
			Valid: true,
		}
	}

	e.ParentId = input.ParentId
	e.Name = input.Name
	e.Type = input.Type
	e.UpdatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	e.UpdatedBy = updatedBy

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

func (uc *permissionUsecase) GetList(ctx context.Context, input *dto.GetListPermissionRequest) (*dto.PermissionCollectionResponse, error) {
	var err error
	offset := 0                                 // default
	limit := constants.DEFAULT_PAGINATION_LIMIT // default
	sorts := map[string]string{}
	filter := input.Filter

	if input.Limit > 0 {
		limit = input.Limit
	}

	if input.Page > 0 {
		offset = (input.Page - 1) * limit
	}

	if input.SortBy != "" {
		sorts, err = utils.QuerySortToMap(input.SortBy)
		if err != nil {
			return nil, err
		}
	}

	rows, err := uc.repo.FindAll(ctx, offset, limit, sorts, filter)
	if err != nil {
		return nil, err
	}

	numrows, err := uc.repo.CountAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	return dto.NewPermissionCollectionResponse(rows, dto.PaginationResponse{
		Size:        len(rows),
		Total:       numrows,
		TotalPages:  int(math.Ceil(float64(numrows) / float64(limit))),
		CurrentPage: input.Page,
	}), nil
}
