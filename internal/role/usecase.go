package role

import (
	"context"

	"github.com/Adhiana46/echo-boilerplate/dto"
)

type RoleUsecase interface {
	CreateRole(ctx context.Context, input *dto.CreateRoleRequest) (*dto.RoleResponse, error)
	UpdateRole(ctx context.Context, input *dto.UpdateRoleRequest) (*dto.RoleResponse, error)
	DeleteRole(ctx context.Context, input *dto.DeleteRoleRequest) (*dto.RoleResponse, error)
	Get(ctx context.Context, input *dto.GetRoleRequest) (*dto.RoleResponse, error)
	GetList(ctx context.Context, input *dto.GetListRoleRequest) (*dto.RoleCollectionResponse, error)
}
