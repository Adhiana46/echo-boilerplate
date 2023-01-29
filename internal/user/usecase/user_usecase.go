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
	tokenmanager "github.com/Adhiana46/echo-boilerplate/pkg/token-manager"
	"github.com/Adhiana46/echo-boilerplate/pkg/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type userUsecase struct {
	userRepo       user.UserRepository
	roleRepo       role.RoleRepository
	userDeviceRepo user.UserDeviceRepository
	tokenManager   *tokenmanager.TokenManager
}

func NewUserUsecase(
	userRepo user.UserRepository,
	roleRepo role.RoleRepository,
	userDeviceRepo user.UserDeviceRepository,
	tokenManager *tokenmanager.TokenManager,
) user.UserUsecase {
	return &userUsecase{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		userDeviceRepo: userDeviceRepo,
		tokenManager:   tokenManager,
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

func (uc *userUsecase) SignIn(ctx context.Context, input *dto.SignInRequest) (*dto.SignInResponse, error) {
	invalidCredsErr := errors.NewBadRequestError("Invalid Credentials")

	user, err := uc.userRepo.FindByUsernameOrEmail(ctx, input.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, invalidCredsErr
		}
		return nil, err
	}

	// set user last_login_at
	user.LastLoginAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	err = utils.ComparePassword(user.Password, input.Password)
	if err != nil {
		return nil, invalidCredsErr
	}

	// save device
	var userDevice *entity.UserDevice
	if input.Device.Token != "" {
		userDevice, err = uc.userDeviceRepo.FindByToken(ctx, user.Id, input.Device.Token)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		if userDevice != nil {
			userDevice.IP = input.Device.IP
			userDevice.Location = input.Device.Location
			userDevice.Platform = input.Device.Platform
			userDevice.UserAgent = input.Device.UserAgent
			userDevice.AppVersion = input.Device.AppVersion
			userDevice.Vendor = input.Device.Vendor
			userDevice.UpdatedAt = sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			}
			userDevice.UpdatedBy = sql.NullInt64{
				Int64: int64(user.Id),
				Valid: true,
			}

			_, err = uc.userDeviceRepo.Update(ctx, userDevice)
			if err != nil {
				return nil, err
			}
		} else {
			userDevice, err = uc.userDeviceRepo.Create(ctx, &entity.UserDevice{
				Uuid:       uuid.NewString(),
				UserId:     user.Id,
				Token:      input.Device.Token,
				IP:         input.Device.IP,
				Location:   input.Device.Location,
				Platform:   input.Device.Platform,
				UserAgent:  input.Device.UserAgent,
				AppVersion: input.Device.AppVersion,
				Vendor:     input.Device.Vendor,
				CreatedAt: sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				},
				CreatedBy: sql.NullInt64{
					Int64: int64(user.Id),
					Valid: true,
				},
				UpdatedAt: sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				},
				UpdatedBy: sql.NullInt64{
					Int64: int64(user.Id),
					Valid: true,
				},
			})
			if err != nil {
				return nil, err
			}
		}
	}

	// Access Token
	accessToken, err := uc.tokenManager.GenerateToken(dto.NewUserClaims(user, userDevice, jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(constants.ACCESS_TOKEN_DURATION),
		},
	}))
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshToken, err := uc.tokenManager.GenerateToken(dto.NewUserClaims(user, userDevice, jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(constants.REFRESH_TOKEN_DURATION),
		},
	}))
	if err != nil {
		return nil, err
	}

	// Update user last_login_at
	_, err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return &dto.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Device:       *dto.NewUserDevice(userDevice),
	}, nil
}

func (uc *userUsecase) SignOut(ctx context.Context, input *dto.SignOutRequest) (*dto.SignOutResponse, error) {
	_, claims, err := uc.tokenManager.ParseToken(input.AccessToken)
	if err != nil {
		if err == tokenmanager.ErrBlacklistedToken || err == tokenmanager.ErrInvalidToken {
			return nil, errors.NewBadRequestError(err.Error())
		}
		return nil, err
	}

	uc.tokenManager.BlacklistToken(input.AccessToken)
	uc.tokenManager.BlacklistToken(input.RefreshToken)

	// delete device if any
	if claims.Device.Uuid != "" {
		userDevice, err := uc.userDeviceRepo.FindByUuid(ctx, claims.Device.Uuid)
		if err != nil {
			return nil, err
		}

		err = uc.userDeviceRepo.Destroy(ctx, userDevice)
		if err != nil {
			return nil, err
		}
	}

	return &dto.SignOutResponse{}, nil
}

func (uc *userUsecase) RefreshToken(ctx context.Context, input *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	_, claims, err := uc.tokenManager.ParseToken(input.RefreshToken)
	if err != nil {
		return nil, errors.NewBadRequestError("Invalid Refresh Token")
	}

	// Access Token
	accessToken, err := uc.tokenManager.GenerateToken(&dto.UserClaims{
		User:   claims.User,
		Device: claims.Device,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(constants.ACCESS_TOKEN_DURATION),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshToken, err := uc.tokenManager.GenerateToken(&dto.UserClaims{
		User:   claims.User,
		Device: claims.Device,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(constants.REFRESH_TOKEN_DURATION),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Blacklist current refresh token
	uc.tokenManager.BlacklistToken(input.RefreshToken)

	return &dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, err
}
