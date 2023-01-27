package user

import (
	"context"

	"github.com/Adhiana46/echo-boilerplate/dto"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, input *dto.CreateUserRequest) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, input *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, input *dto.DeleteUserRequest) (*dto.UserResponse, error)
	Get(ctx context.Context, input *dto.GetUserRequest) (*dto.UserResponse, error)
	GetList(ctx context.Context, input *dto.GetListUserRequest) (*dto.UserCollectionResponse, error)
}
