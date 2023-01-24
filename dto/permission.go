package dto

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
