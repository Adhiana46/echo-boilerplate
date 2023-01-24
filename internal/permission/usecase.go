package permission

import (
	"context"

	"github.com/Adhiana46/echo-boilerplate/dto"
)

type PermissionUsecase interface {
	CreatePermission(ctx context.Context, input *dto.CreatePermissionRequest) (*dto.PermissionResponse, error)
	UpdatePermission(ctx context.Context, input *dto.UpdatePermissionRequest) (*dto.PermissionResponse, error)
	DeletePermission(ctx context.Context, input *dto.DeletePermissionRequest) (*dto.PermissionResponse, error)
	Get(ctx context.Context, input *dto.GetPermissionRequest) (*dto.PermissionResponse, error)
	GetList(ctx context.Context, input *dto.GetListPermissionRequest) (*dto.PermissionResponse, error)
}
