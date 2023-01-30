package middlewares

import (
	"context"

	"github.com/Adhiana46/echo-boilerplate/dto"
	"github.com/Adhiana46/echo-boilerplate/pkg/errors"
	tokenmanager "github.com/Adhiana46/echo-boilerplate/pkg/token-manager"
	"github.com/Adhiana46/echo-boilerplate/pkg/utils"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Authenticate(tokenManager *tokenmanager.TokenManager) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ErrorHandler: func(c echo.Context, err error) error {
			if err != nil {
				return errors.NewUnauthorizedError(err.Error())
			}

			return nil
		},
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			token, _, err := tokenManager.ParseToken(auth)
			if err != nil {
				return nil, err
			}

			return token, nil
		},
		SuccessHandler: func(c echo.Context) {
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(*dto.UserClaims)

			// Set echo context (useless btw)
			c.Set("user", claims.User)
			c.Set("device", claims.Device)

			// Set request context
			ctx := context.WithValue(c.Request().Context(), "user", claims.User)
			ctx = context.WithValue(ctx, "device", claims.Device)

			c.SetRequest(c.Request().WithContext(ctx))
		},
	})
}

func Permissions(permissions ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			user := utils.GetUserFromContext(ctx)

			if user == nil {
				return errors.NewUnauthorizedError("")
			}

			// check if user does not have role
			if user.Role == nil {
				return errors.NewForbiddenError("")
			}

			userPermissionsMap := map[string]string{}
			for _, userPermission := range user.Role.Permissions {
				userPermissionsMap[userPermission] = userPermission
			}

			hasPermission := true
			for _, permission := range permissions {
				_, isExists := userPermissionsMap[permission]

				hasPermission = hasPermission && isExists
			}

			if !hasPermission {
				return errors.NewForbiddenError("")
			}

			return next(c)
		}
	}
}
