package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/Adhiana46/echo-boilerplate/constants"
	"github.com/Adhiana46/echo-boilerplate/dto"
	"github.com/Adhiana46/echo-boilerplate/entity"
	"github.com/Adhiana46/echo-boilerplate/internal/role"
	"github.com/Adhiana46/echo-boilerplate/internal/user"
	"github.com/Adhiana46/echo-boilerplate/pkg/errors"
	"github.com/Adhiana46/echo-boilerplate/pkg/utils"
	"github.com/google/uuid"
)

type userUsecase struct {
	userRepo user.UserRepository
	roleRepo role.RoleRepository
}

func NewUserUsecase(userRepo user.UserRepository, roleRepo role.RoleRepository) user.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (uc *userUsecase) CreateUser(ctx context.Context, input *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Validation
	numrows, err := uc.userRepo.CountByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if numrows > 0 {
		return nil, errors.NewBadRequestError(fmt.Sprintf("User with email '%s' already exists", input.Email))
	}

	numrows, err = uc.userRepo.CountByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	if numrows > 0 {
		return nil, errors.NewBadRequestError(fmt.Sprintf("User with username '%s' already exists", input.Username))
	}

	role, err := uc.roleRepo.FindByName(ctx, input.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewBadRequestError(fmt.Sprintf("Role '%s' is not exists", input.Role))
		}
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	row, err := uc.userRepo.Create(ctx, &entity.User{
		Uuid:        uuid.NewString(),
		Username:    input.Username,
		Email:       input.Email,
		Password:    hashedPassword,
		Name:        input.Name,
		RoleId:      role.Id,
		Status:      input.Status,
		LastLoginAt: sql.NullTime{},
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		CreatedBy: sql.NullInt64{}, // TODO
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedBy: sql.NullInt64{}, // TODO
		Role:      role,
	})

	if err != nil {
		return nil, err
	}

	return dto.NewUserResponse(row), nil
}

func (uc *userUsecase) UpdateUser(ctx context.Context, input *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	e, err := uc.userRepo.FindByUuid(ctx, input.Uuid)
	if err != nil {
		return nil, err
	}

	// validation logic
	if e.Email != input.Email {
		numrows, err := uc.userRepo.CountByEmail(ctx, input.Email)
		if err != nil {
			return nil, err
		}
		if numrows > 0 {
			return nil, errors.NewBadRequestError(fmt.Sprintf("User with email '%s' already exists", input.Email))
		}
	}

	if e.Username != input.Username {
		numrows, err := uc.userRepo.CountByUsername(ctx, input.Username)
		if err != nil {
			return nil, err
		}
		if numrows > 0 {
			return nil, errors.NewBadRequestError(fmt.Sprintf("User with username '%s' already exists", input.Username))
		}
	}

	role, err := uc.roleRepo.FindByName(ctx, input.Role)
	if err != nil {
		return nil, err
	}
	if role.Id == 0 {
		return nil, errors.NewBadRequestError(fmt.Sprintf("Role '%s' is not exists", input.Role))
	}

	// Update e
	e.Name = input.Name
	e.Username = input.Username
	e.Email = input.Email
	e.Name = input.Name
	e.Status = input.Status
	e.RoleId = role.Id
	e.Role = role
	e.UpdatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	e.UpdatedBy = sql.NullInt64{} // TODO: user
	if input.Password != "" {
		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			return nil, err
		}

		e.Password = hashedPassword
	}

	// Update
	updatedE, err := uc.userRepo.Update(ctx, e)
	if err != nil {
		return nil, err
	}

	return dto.NewUserResponse(updatedE), nil
}

func (uc *userUsecase) DeleteUser(ctx context.Context, input *dto.DeleteUserRequest) (*dto.UserResponse, error) {
	e, err := uc.userRepo.FindByUuid(ctx, input.Uuid)
	if err != nil {
		return nil, err
	}

	err = uc.userRepo.Destroy(ctx, e)

	return nil, err
}

func (uc *userUsecase) Get(ctx context.Context, input *dto.GetUserRequest) (*dto.UserResponse, error) {
	e, err := uc.userRepo.FindByUuid(ctx, input.Uuid)
	if err != nil {
		return nil, err
	}

	return dto.NewUserResponse(e), nil
}

func (uc *userUsecase) GetList(ctx context.Context, input *dto.GetListUserRequest) (*dto.UserCollectionResponse, error) {
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

	rows, err := uc.userRepo.FindAll(ctx, offset, limit, sorts, filter)
	if err != nil {
		return nil, err
	}

	numrows, err := uc.userRepo.CountAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	return dto.NewUserCollectionResponse(rows, dto.PaginationResponse{
		Size:        len(rows),
		Total:       numrows,
		TotalPages:  int(math.Ceil(float64(numrows) / float64(limit))),
		CurrentPage: input.Page,
	}), nil
}
