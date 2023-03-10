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
	"github.com/Adhiana46/echo-boilerplate/internal/role"
	"github.com/Adhiana46/echo-boilerplate/pkg/errors"
	"github.com/Adhiana46/echo-boilerplate/pkg/utils"
	"github.com/google/uuid"
)

var (
	roleUcInstance     *roleUsecase
	roleUcInstanceOnce sync.Once
)

type roleUsecase struct {
	roleRepo role.RoleRepository
	permRepo permission.PermissionRepository
}

func NewRoleUsecase(roleRepo role.RoleRepository, permRepo permission.PermissionRepository) role.RoleUsecase {
	roleUcInstanceOnce.Do(func() {
		roleUcInstance = &roleUsecase{
			roleRepo: roleRepo,
			permRepo: permRepo,
		}
	})

	return roleUcInstance
}

func (uc *roleUsecase) CreateRole(ctx context.Context, input *dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	// Validation
	numrows, err := uc.roleRepo.CountByName(ctx, input.Name)
	if err != nil {
		return nil, err
	}

	if numrows > 0 {
		return nil, errors.NewBadRequestError(fmt.Sprintf("Role with name '%s' already exists", input.Name))
	}

	permissions, err := uc.permRepo.FindAllByNames(ctx, input.Permissions)
	if err != nil {
		return nil, err
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

	row, err := uc.roleRepo.Create(ctx, &entity.Role{
		Uuid: uuid.NewString(),
		Name: input.Name,
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		CreatedBy: createdBy,
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedBy:   updatedBy,
		Permissions: permissions,
	})

	if err != nil {
		return nil, err
	}

	return dto.NewRoleResponse(row), nil
}

func (uc *roleUsecase) UpdateRole(ctx context.Context, input *dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	e, err := uc.roleRepo.FindByUuid(ctx, input.Uuid)
	if err != nil {
		return nil, err
	}

	// validation logic
	if e.Name != input.Name {
		numrows, err := uc.roleRepo.CountByName(ctx, input.Name)
		if err != nil {
			return nil, err
		}

		if numrows > 0 {
			return nil, errors.NewBadRequestError(fmt.Sprintf("Role with name '%s' already exists", input.Name))
		}
	}

	permissions, err := uc.permRepo.FindAllByNames(ctx, input.Permissions)
	if err != nil {
		return nil, err
	}

	updatedBy := sql.NullInt64{}
	user := utils.GetUserFromContext(ctx)
	if user != nil {
		updatedBy = sql.NullInt64{
			Int64: int64(user.ID),
			Valid: true,
		}
	}

	e.Name = input.Name
	e.Permissions = permissions
	e.UpdatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	e.UpdatedBy = updatedBy

	// Update
	updatedE, err := uc.roleRepo.Update(ctx, e)
	if err != nil {
		return nil, err
	}

	return dto.NewRoleResponse(updatedE), nil
}

func (uc *roleUsecase) DeleteRole(ctx context.Context, input *dto.DeleteRoleRequest) (*dto.RoleResponse, error) {
	e, err := uc.roleRepo.FindByUuid(ctx, input.Uuid)
	if err != nil {
		return nil, err
	}

	err = uc.roleRepo.Destroy(ctx, e)

	return nil, err
}

func (uc *roleUsecase) Get(ctx context.Context, input *dto.GetRoleRequest) (*dto.RoleResponse, error) {
	e, err := uc.roleRepo.FindByUuid(ctx, input.Uuid)
	if err != nil {
		return nil, err
	}

	return dto.NewRoleResponse(e), nil
}

func (uc *roleUsecase) GetList(ctx context.Context, input *dto.GetListRoleRequest) (*dto.RoleCollectionResponse, error) {
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

	rows, err := uc.roleRepo.FindAll(ctx, offset, limit, sorts, filter)
	if err != nil {
		return nil, err
	}

	numrows, err := uc.roleRepo.CountAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	return dto.NewRoleCollectionResponse(rows, dto.PaginationResponse{
		Size:        len(rows),
		Total:       numrows,
		TotalPages:  int(math.Ceil(float64(numrows) / float64(limit))),
		CurrentPage: input.Page,
	}), nil
}
