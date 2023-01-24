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
		CreatedBy: e.CreatedBy,
		UpdatedAt: updatedAt,
		UpdatedBy: e.UpdatedBy,
	}
}

type PermissionCollectionResponse struct {
	Data       []PermissionResponse `json:"data"`
	Pagination PaginationResponse   `json:"pagination"`
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
	Page   int `query:"page"`
	Limit  int `query:"limit"`
	SortBy int `query:"sortBy"`
	Filter int `query:"filter"`
}
