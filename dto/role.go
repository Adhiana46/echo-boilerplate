package dto

import (
	"time"

	"github.com/Adhiana46/echo-boilerplate/entity"
)

type RoleResponse struct {
	Uuid      string `json:"uuid"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	CreatedBy int    `json:"created_by"`
	UpdatedAt string `json:"updated_at"`
	UpdatedBy int    `json:"updated_by"`
}

type RoleResponseWithPermissions struct {
	Uuid        string   `json:"uuid"`
	Name        string   `json:"name"`
	CreatedAt   string   `json:"created_at"`
	CreatedBy   int      `json:"created_by"`
	UpdatedAt   string   `json:"updated_at"`
	UpdatedBy   int      `json:"updated_by"`
	Permissions []string `json:"permissions"`
}

type RoleCollectionResponse struct {
	Data       []*RoleResponse    `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

func NewRoleResponse(e *entity.Role) *RoleResponseWithPermissions {
	createdAt := ""
	updatedAt := ""
	permissions := []string{}

	if e.CreatedAt.Valid {
		createdAt = e.CreatedAt.Time.Format(time.RFC3339)
	}
	if e.UpdatedAt.Valid {
		updatedAt = e.UpdatedAt.Time.Format(time.RFC3339)
	}

	for _, perm := range e.Permissions {
		permissions = append(permissions, perm.Name)
	}

	return &RoleResponseWithPermissions{
		Uuid:        e.Uuid,
		Name:        e.Name,
		CreatedAt:   createdAt,
		CreatedBy:   int(e.CreatedBy.Int64),
		UpdatedAt:   updatedAt,
		UpdatedBy:   int(e.UpdatedBy.Int64),
		Permissions: permissions,
	}
}

func NewRoleCollectionResponse(rows []*entity.Role, pagination PaginationResponse) *RoleCollectionResponse {
	response := &RoleCollectionResponse{
		Data:       []*RoleResponse{},
		Pagination: pagination,
	}

	for _, row := range rows {
		createdAt := ""
		updatedAt := ""

		if row.CreatedAt.Valid {
			createdAt = row.CreatedAt.Time.Format(time.RFC3339)
		}
		if row.UpdatedAt.Valid {
			updatedAt = row.UpdatedAt.Time.Format(time.RFC3339)
		}

		response.Data = append(response.Data, &RoleResponse{
			Uuid:      row.Uuid,
			Name:      row.Name,
			CreatedAt: createdAt,
			CreatedBy: int(row.CreatedBy.Int64),
			UpdatedAt: updatedAt,
			UpdatedBy: int(row.UpdatedBy.Int64),
		})
	}

	return response
}

type CreateRoleRequest struct {
	Name        string   `json:"name" validate:"required"`
	Permissions []string `json:"permissions" validate:""`
}

type UpdateRoleRequest struct {
	Uuid        string   `json:"uuid" validate:"required"`
	Name        string   `json:"name" validate:"required"`
	Permissions []string `json:"permissions" validate:""`
}

type DeleteRoleRequest struct {
	Uuid string `json:"uuid" validate:"required"`
}

type GetRoleRequest struct {
	Uuid string `json:"uuid" validate:"required"`
}

type GetListRoleRequest struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	SortBy string `query:"sortBy"`
	Filter string `query:"filter"`
}
