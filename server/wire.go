//go:build wireinject
// +build wireinject

package server

import (
	permissionData "github.com/Adhiana46/echo-boilerplate/internal/permission/data"
	permissionHttpHandler "github.com/Adhiana46/echo-boilerplate/internal/permission/delivery/http"
	permissionRepo "github.com/Adhiana46/echo-boilerplate/internal/permission/repository"
	permissionUsecase "github.com/Adhiana46/echo-boilerplate/internal/permission/usecase"
	roleData "github.com/Adhiana46/echo-boilerplate/internal/role/data"
	roleHttpHandler "github.com/Adhiana46/echo-boilerplate/internal/role/delivery/http"
	roleRepo "github.com/Adhiana46/echo-boilerplate/internal/role/repository"
	roleUsecase "github.com/Adhiana46/echo-boilerplate/internal/role/usecase"
	userData "github.com/Adhiana46/echo-boilerplate/internal/user/data"
	userHttpHandler "github.com/Adhiana46/echo-boilerplate/internal/user/delivery/http"
	userRepo "github.com/Adhiana46/echo-boilerplate/internal/user/repository"
	userUsecase "github.com/Adhiana46/echo-boilerplate/internal/user/usecase"
	cachePkg "github.com/Adhiana46/echo-boilerplate/pkg/cache"
	tokenmanager "github.com/Adhiana46/echo-boilerplate/pkg/token-manager"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

var ProviderSet = wire.NewSet(
	// Http Handler
	permissionHttpHandler.NewPermissionHttpHandler,
	roleHttpHandler.NewRoleHttpHandler,
	userHttpHandler.NewUserHttpHandler,

	// Usecase
	permissionUsecase.NewPermissionUsecase,
	roleUsecase.NewRoleUsecase,
	userUsecase.NewUserUsecase,

	// Repository
	permissionRepo.NewPermissionRepository,
	roleRepo.NewRoleRepository,
	userRepo.NewUserRepository,
	userRepo.NewUserDeviceRepository,

	// Data Source
	permissionData.NewPostgresPermissionPersistent,
	roleData.NewPostgresRolePersistent,
	userData.NewPostgresUserPersistent,
	userData.NewPostgresUserDevicePersistent,
)

func InitializedPermissionHandler(db *sqlx.DB, cache cachePkg.Cache, tokenManager *tokenmanager.TokenManager) permissionHttpHandler.Handler {
	panic(wire.Build(
		ProviderSet,
	))
}

func InitializedRoleHandler(db *sqlx.DB, cache cachePkg.Cache, tokenManager *tokenmanager.TokenManager) roleHttpHandler.Handler {
	panic(wire.Build(
		ProviderSet,
	))
}

func InitializedUserHandler(db *sqlx.DB, cache cachePkg.Cache, tokenManager *tokenmanager.TokenManager) userHttpHandler.Handler {
	panic(wire.Build(
		ProviderSet,
	))
}
