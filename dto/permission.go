package dto

import (
	"time"

	"github.com/Adhiana46/echo-boilerplate/entity"
)

type PermissionResponse struct {
	Uuid      string `json:"uuid"`
	ParentId  int    `json:"parent_id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	CreatedBy int    `json:"created_by"`
	UpdatedAt string `json:"updated_at"`
	UpdatedBy int    `json:"updated_by"`
}

func NewPermissionResponse(e *entity.Permission) *PermissionResponse {
	createdAt := ""
	updatedAt := ""

	if e.CreatedAt.Valid {
		createdAt = e.CreatedAt.Time.Format(time.RFC3339)
	}
	if e.UpdatedAt.Valid {
		updatedAt = e.UpdatedAt.Time.Format(time.RFC3339)
	}

	return &PermissionResponse{
		Uuid:      e.Uuid,
		ParentId:  e.ParentId,
		Name:      e.Name,
		Type:      e.Type,
		CreatedAt: createdAt,
		CreatedBy: int(e.CreatedBy.Int64),
		UpdatedAt: updatedAt,
		UpdatedBy: int(e.UpdatedBy.Int64),
	}
}

type PermissionCollectionResponse struct {
	Data       []*PermissionResponse `json:"data"`
	Pagination PaginationResponse    `json:"pagination"`
}

func NewPermissionCollectionResponse(rows []*entity.Permission, pagination PaginationResponse) *PermissionCollectionResponse {
	response := &PermissionCollectionResponse{
		Data:       []*PermissionResponse{},
		Pagination: pagination,
	}

	for _, row := range rows {
		response.Data = append(response.Data, NewPermissionResponse(row))
	}

	return response
}

type CreatePermissionRequest struct {
	ParentId int    `json:"parent_id" validate:"numeric"`
	Name     string `json:"name" validate:"required"`
	Type     string `json:"type" validate:"required"`
}

type UpdatePermissionRequest struct {
	Uuid     string `json:"uuid" validate:"required"`
	ParentId int    `json:"parent_id" validate:"numeric"`
	Name     string `json:"name" validate:"required"`
	Type     string `json:"type" validate:"required"`
}

type DeletePermissionRequest struct {
	Uuid string `json:"uuid"`
}

type GetPermissionRequest struct {
	Uuid string `json:"uuid"`
}

type GetListPermissionRequest struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	SortBy string `query:"sortBy"`
	Filter string `query:"filter"`
}
