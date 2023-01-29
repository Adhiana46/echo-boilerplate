package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/Adhiana46/echo-boilerplate/config"
	"github.com/Adhiana46/echo-boilerplate/internal/permission"
	permissionHttpHandler "github.com/Adhiana46/echo-boilerplate/internal/permission/delivery/http"
	permissionRepo "github.com/Adhiana46/echo-boilerplate/internal/permission/repository"
	permissionUsecase "github.com/Adhiana46/echo-boilerplate/internal/permission/usecase"
	"github.com/Adhiana46/echo-boilerplate/internal/role"
	roleHttpHandler "github.com/Adhiana46/echo-boilerplate/internal/role/delivery/http"
	roleRepo "github.com/Adhiana46/echo-boilerplate/internal/role/repository"
	roleUsecase "github.com/Adhiana46/echo-boilerplate/internal/role/usecase"
	"github.com/Adhiana46/echo-boilerplate/internal/user"
	userHttpHandler "github.com/Adhiana46/echo-boilerplate/internal/user/delivery/http"
	userRepo "github.com/Adhiana46/echo-boilerplate/internal/user/repository"
	userUsecase "github.com/Adhiana46/echo-boilerplate/internal/user/usecase"
	cachePkg "github.com/Adhiana46/echo-boilerplate/pkg/cache"
	"github.com/Adhiana46/echo-boilerplate/pkg/errors"
	tokenmanager "github.com/Adhiana46/echo-boilerplate/pkg/token-manager"
	"github.com/Adhiana46/echo-boilerplate/pkg/utils"
	validatorPkg "github.com/Adhiana46/echo-boilerplate/pkg/validator"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e            *echo.Echo
	cfg          *config.Config
	db           *sqlx.DB
	cache        cachePkg.Cache
	tokenManager *tokenmanager.TokenManager

	// repositories
	repoPermission permission.PermissionRepository
	repoRole       role.RoleRepository
	repoUser       user.UserRepository
	repoUserDevice user.UserDeviceRepository

	// usecases
	usecasePermission permission.PermissionUsecase
	usecaseRole       role.RoleUsecase
	usecaseUser       user.UserUsecase

	// handlers
	permissionHandler permissionHttpHandler.Handler
	roleHandler       roleHttpHandler.Handler
	userHandler       userHttpHandler.Handler
}

func NewServer(cfg *config.Config, db *sqlx.DB, cache cachePkg.Cache, tokenManager *tokenmanager.TokenManager) *Server {
	// Validator
	validate := validator.New()
	// register function to get tag name from json tags.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	e := echo.New()
	e.Validator = validatorPkg.NewEchoValidator(validate)

	// error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var statusCode int = 500
		var message string = ""
		var errorsData any = nil
		stackTraces := []string{
			err.Error(),
			"",
			"",
		}
		stackTraces = append(stackTraces, strings.Split(strings.ReplaceAll(string(debug.Stack()), "\t", "     "), "\n")...)

		if err == sql.ErrNoRows {
			statusCode = 404
		} else {
			switch err.(type) {
			case validator.ValidationErrors:
				errs := err.(validator.ValidationErrors)

				statusCode = 400
				errorsData = utils.ValidationErrors(errs, nil)
			case errors.CustomError:
				e := err.(errors.CustomError)

				statusCode = e.StatusCode()
				message = e.Message()

				if len(e.Errors()) > 0 {
					errorsData = e.Errors()
				}
			default:
				statusCode = 500

				// log unexpected error
				// TODO: log error perhaps using sentry.io
			}
		}

		if !cfg.App.Debug {
			stackTraces = nil
		}

		resp := utils.JsonError(statusCode, message, errorsData, stackTraces)

		c.Logger().Error(err)
		c.JSON(statusCode, resp)
	}

	e.Pre(middleware.AddTrailingSlash())

	// Middlewares
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	srv := &Server{
		e:            e,
		cfg:          cfg,
		db:           db,
		cache:        cache,
		tokenManager: tokenManager,
	}

	srv.setupRepo()
	srv.setupUsecase()
	srv.setupHttpHandler()
	srv.setupRoutes()

	return srv
}

func (s *Server) Run() error {
	log.Println("[Server]:", fmt.Sprintf("Running server on %s:%s", s.cfg.Http.Host, s.cfg.Http.Port))
	return s.e.Start(fmt.Sprintf("%s:%s", s.cfg.Http.Host, s.cfg.Http.Port))
}

func (s *Server) setupRepo() {
	s.repoPermission = permissionRepo.NewPgPermissionRepository(s.db)
	s.repoRole = roleRepo.NewPgRoleRepository(s.db)
	s.repoUser = userRepo.NewPgUserRepository(s.db, s.repoRole)
	s.repoUserDevice = userRepo.NewPgUserDeviceRepository(s.db)
}

func (s *Server) setupUsecase() {
	s.usecasePermission = permissionUsecase.NewPermissionUsecase(s.repoPermission)
	s.usecaseRole = roleUsecase.NewRoleUsecase(s.repoRole, s.repoPermission)
	s.usecaseUser = userUsecase.NewUserUsecase(s.repoUser, s.repoRole, s.repoUserDevice, s.tokenManager)
}

func (s *Server) setupHttpHandler() {
	s.permissionHandler = permissionHttpHandler.NewPermissionHttpHandler(s.usecasePermission)
	s.roleHandler = roleHttpHandler.NewRoleHttpHandler(s.usecaseRole)
	s.userHandler = userHttpHandler.NewUserHttpHandler(s.usecaseUser)
}

func (s *Server) setupRoutes() {
	groupPermission := s.e.Group("/api/v1/permissions")
	groupPermission.POST("/", s.permissionHandler.Store())
	groupPermission.PUT("/:uuid", s.permissionHandler.Update())
	groupPermission.DELETE("/:uuid", s.permissionHandler.Delete())
	groupPermission.GET("/:uuid", s.permissionHandler.GetByUuid())
	groupPermission.GET("/", s.permissionHandler.GetAll())

	groupRole := s.e.Group("/api/v1/roles")
	groupRole.POST("/", s.roleHandler.Store())
	groupRole.PUT("/:uuid", s.roleHandler.Update())
	groupRole.DELETE("/:uuid", s.roleHandler.Delete())
	groupRole.GET("/:uuid", s.roleHandler.GetByUuid())
	groupRole.GET("/", s.roleHandler.GetAll())

	groupUser := s.e.Group("/api/v1/users")
	groupUser.POST("/", s.userHandler.Store())
	groupUser.PUT("/:uuid", s.userHandler.Update())
	groupUser.DELETE("/:uuid", s.userHandler.Delete())
	groupUser.GET("/:uuid", s.userHandler.GetByUuid())
	groupUser.GET("/", s.userHandler.GetAll())

	groupAuth := s.e.Group("/api/v1/auth")
	groupAuth.POST("/signin/", s.userHandler.SignIn())
	groupAuth.POST("/signout/", s.userHandler.SignOut())
	groupAuth.POST("/refresh-token/", s.userHandler.RefreshToken())
}
