package dto

import (
	"time"

	"github.com/Adhiana46/echo-boilerplate/entity"
)

type UserResponse struct {
	Uuid        string        `json:"uuid"`
	Username    string        `json:"username"`
	Email       string        `json:"email"`
	Name        string        `json:"name"`
	Status      int           `json:"status"`
	LastLoginAt string        `json:"last_login_at"`
	CreatedAt   string        `json:"created_at"`
	CreatedBy   int           `json:"created_by"`
	UpdatedAt   string        `json:"updated_at"`
	UpdatedBy   int           `json:"updated_by"`
	Role        *RoleResponse `json:"role,omitempty"`
}

type UserResponseWithID struct {
	ID          int           `json:"id"`
	Uuid        string        `json:"uuid"`
	Username    string        `json:"username"`
	Email       string        `json:"email"`
	Name        string        `json:"name"`
	Status      int           `json:"status"`
	LastLoginAt string        `json:"last_login_at"`
	CreatedAt   string        `json:"created_at"`
	CreatedBy   int           `json:"created_by"`
	UpdatedAt   string        `json:"updated_at"`
	UpdatedBy   int           `json:"updated_by"`
	Role        *RoleResponse `json:"role,omitempty"`
}

type UserCollectionResponse struct {
	Data       []*UserResponse    `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

func NewUserResponse(e *entity.User) *UserResponse {
	lastLoginAt := ""
	createdAt := ""
	updatedAt := ""
	var roleResponse *RoleResponse = nil

	if e.LastLoginAt.Valid {
		lastLoginAt = e.LastLoginAt.Time.Format(time.RFC3339)
	}
	if e.CreatedAt.Valid {
		createdAt = e.CreatedAt.Time.Format(time.RFC3339)
	}
	if e.UpdatedAt.Valid {
		updatedAt = e.UpdatedAt.Time.Format(time.RFC3339)
	}

	if e.Role != nil {
		roleResponse = NewRoleResponse(e.Role)
	}

	return &UserResponse{
		Uuid:        e.Uuid,
		Username:    e.Username,
		Email:       e.Email,
		Name:        e.Name,
		Status:      e.Status,
		LastLoginAt: lastLoginAt,
		CreatedAt:   createdAt,
		CreatedBy:   int(e.CreatedBy.Int64),
		UpdatedAt:   updatedAt,
		UpdatedBy:   int(e.UpdatedBy.Int64),
		Role:        roleResponse,
	}
}

func NewUserResponseWithID(e *entity.User) *UserResponseWithID {
	lastLoginAt := ""
	createdAt := ""
	updatedAt := ""
	var roleResponse *RoleResponse = nil

	if e.LastLoginAt.Valid {
		lastLoginAt = e.LastLoginAt.Time.Format(time.RFC3339)
	}
	if e.CreatedAt.Valid {
		createdAt = e.CreatedAt.Time.Format(time.RFC3339)
	}
	if e.UpdatedAt.Valid {
		updatedAt = e.UpdatedAt.Time.Format(time.RFC3339)
	}

	if e.Role != nil {
		roleResponse = NewRoleResponse(e.Role)
	}

	return &UserResponseWithID{
		ID:          e.Id,
		Uuid:        e.Uuid,
		Username:    e.Username,
		Email:       e.Email,
		Name:        e.Name,
		Status:      e.Status,
		LastLoginAt: lastLoginAt,
		CreatedAt:   createdAt,
		CreatedBy:   int(e.CreatedBy.Int64),
		UpdatedAt:   updatedAt,
		UpdatedBy:   int(e.UpdatedBy.Int64),
		Role:        roleResponse,
	}
}

func NewUserCollectionResponse(rows []*entity.User, pagination PaginationResponse) *UserCollectionResponse {
	response := &UserCollectionResponse{
		Data:       []*UserResponse{},
		Pagination: pagination,
	}

	for _, row := range rows {
		response.Data = append(response.Data, NewUserResponse(row))
	}

	return response
}

type CreateUserRequest struct {
	Username             string `json:"username" validate:"required,min=3,max=30"`
	Email                string `json:"email" validate:"required,email"`
	Name                 string `json:"name" validate:"required,min=3,max=100"`
	Status               int    `json:"status" validate:"numeric,min=0"`
	Role                 string `json:"role" validate:"required"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type UpdateUserRequest struct {
	Uuid                 string `json:"uuid" validate:"required"`
	Username             string `json:"username" validate:"required,min=3,max=30"`
	Email                string `json:"email" validate:"required,email"`
	Name                 string `json:"name" validate:"required,min=3,max=100"`
	Status               int    `json:"status" validate:"required,numeric"`
	Role                 string `json:"role" validate:"required"`
	Password             string `json:"password" validate:"omitempty,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"eqfield=Password"`
}

type DeleteUserRequest struct {
	Uuid string `json:"uuid" validate:"required"`
}

type GetUserRequest struct {
	Uuid string `json:"uuid" validate:"required"`
}

type GetListUserRequest struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	SortBy string `query:"sortBy"`
	Filter string `query:"filter"`
	// TODO: more filters
}
