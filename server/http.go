package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	"github.com/Adhiana46/echo-boilerplate/config"
	permissionHttpHandler "github.com/Adhiana46/echo-boilerplate/internal/permission/delivery/http"
	roleHttpHandler "github.com/Adhiana46/echo-boilerplate/internal/role/delivery/http"
	userHttpHandler "github.com/Adhiana46/echo-boilerplate/internal/user/delivery/http"
	cachePkg "github.com/Adhiana46/echo-boilerplate/pkg/cache"
	"github.com/Adhiana46/echo-boilerplate/pkg/errors"
	"github.com/Adhiana46/echo-boilerplate/pkg/logger"
	m "github.com/Adhiana46/echo-boilerplate/pkg/middlewares"
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

		logger.WithFields(logger.Fields{
			"at":     time.Now().Format("2006-01-02 15:04:05"),
			"method": c.Request().Method,
			"uri":    c.Request().URL.String(),
			"ip":     c.Request().RemoteAddr,
		}).Error(err)
		c.JSON(statusCode, resp)
	}

	e.Pre(middleware.AddTrailingSlash())

	// Middlewares
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.WithFields(logger.Fields{
				"at":     time.Now().Format("2006-01-02 15:04:05"),
				"method": c.Request().Method,
				"uri":    c.Request().URL.String(),
				"ip":     c.Request().RemoteAddr,
			}).Info("incoming request")

			return next(c)
		}
	})

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

	srv.setupHttpHandler()
	srv.setupRoutes()

	return srv
}

func (s *Server) Run() error {
	log.Println("[Server]:", fmt.Sprintf("Running server on %s:%s", s.cfg.Http.Host, s.cfg.Http.Port))
	return s.e.Start(fmt.Sprintf("%s:%s", s.cfg.Http.Host, s.cfg.Http.Port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func (s *Server) setupHttpHandler() {
	s.permissionHandler = InitializedPermissionHandler(s.db, s.cache, s.tokenManager)
	s.roleHandler = InitializedRoleHandler(s.db, s.cache, s.tokenManager)
	s.userHandler = InitializedUserHandler(s.db, s.cache, s.tokenManager)
}

func (s *Server) setupRoutes() {
	groupPermission := s.e.Group("/api/v1/permissions", m.Authenticate(s.tokenManager))
	groupPermission.POST("/", s.permissionHandler.Store(), m.Permissions("permissions.create"))
	groupPermission.PUT("/:uuid", (s.permissionHandler.Update()), m.Permissions("permissions.update"))
	groupPermission.DELETE("/:uuid", s.permissionHandler.Delete(), m.Permissions("permissions.delete"))
	groupPermission.GET("/:uuid", s.permissionHandler.GetByUuid(), m.Permissions("permissions.read"))
	groupPermission.GET("/", s.permissionHandler.GetAll(), m.Permissions("permissions.read"))

	groupRole := s.e.Group("/api/v1/roles", m.Authenticate(s.tokenManager))
	groupRole.POST("/", s.roleHandler.Store(), m.Permissions("roles.create"))
	groupRole.PUT("/:uuid", s.roleHandler.Update(), m.Permissions("roles.update"))
	groupRole.DELETE("/:uuid", s.roleHandler.Delete(), m.Permissions("roles.delete"))
	groupRole.GET("/:uuid", s.roleHandler.GetByUuid(), m.Permissions("roles.read"))
	groupRole.GET("/", s.roleHandler.GetAll(), m.Permissions("roles.read"))

	groupUser := s.e.Group("/api/v1/users", m.Authenticate(s.tokenManager))
	groupUser.POST("/", s.userHandler.Store(), m.Permissions("users.create"))
	groupUser.PUT("/:uuid", s.userHandler.Update(), m.Permissions("users.update"))
	groupUser.DELETE("/:uuid", s.userHandler.Delete(), m.Permissions("users.delete"))
	groupUser.GET("/:uuid", s.userHandler.GetByUuid(), m.Permissions("users.read"))
	groupUser.GET("/", s.userHandler.GetAll(), m.Permissions("users.read"))

	groupAuth := s.e.Group("/api/v1/auth")
	groupAuth.POST("/signin/", s.userHandler.SignIn())
	groupAuth.POST("/signout/", s.userHandler.SignOut())
	groupAuth.POST("/refresh-token/", s.userHandler.RefreshToken())
}
