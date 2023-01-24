package http

import (
	"net/http"
	"strings"

	"github.com/Adhiana46/echo-boilerplate/dto"
	"github.com/Adhiana46/echo-boilerplate/internal/permission"
	"github.com/Adhiana46/echo-boilerplate/pkg/utils"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Store() func(echo.Context) error
	Update() func(echo.Context) error
	Delete() func(echo.Context) error
	GetByUuid() func(echo.Context) error
}

type handler struct {
	uc permission.PermissionUsecase
}

func NewPermissionHttpHandler(uc permission.PermissionUsecase) Handler {
	return &handler{
		uc: uc,
	}
}

func (h *handler) Store() func(echo.Context) error {
	return func(c echo.Context) error {
		input := dto.CreatePermissionRequest{}

		if err := c.Bind(&input); err != nil {
			panic(err)
		}

		if err := c.Validate(input); err != nil {
			panic(err)
		}

		res, err := h.uc.CreatePermission(c.Request().Context(), &input)
		if err != nil {
			panic(err)
		}

		return c.JSON(http.StatusCreated, utils.JsonSuccess(http.StatusCreated, "", res, nil))
	}
}

func (h *handler) Update() func(echo.Context) error {
	return func(c echo.Context) error {
		input := dto.UpdatePermissionRequest{}

		if err := c.Bind(&input); err != nil {
			panic(err)
		}

		input.Uuid = strings.Trim(c.Param("uuid"), "/")

		if err := c.Validate(input); err != nil {
			panic(err)
		}

		res, err := h.uc.UpdatePermission(c.Request().Context(), &input)
		if err != nil {
			panic(err)
		}

		return c.JSON(http.StatusOK, utils.JsonSuccess(http.StatusOK, "", res, nil))
	}
}

func (h *handler) Delete() func(echo.Context) error {
	return func(c echo.Context) error {
		input := dto.DeletePermissionRequest{
			Uuid: strings.Trim(c.Param("uuid"), "/"),
		}

		if err := c.Validate(input); err != nil {
			panic(err)
		}

		res, err := h.uc.DeletePermission(c.Request().Context(), &input)
		if err != nil {
			panic(err)
		}

		return c.JSON(http.StatusOK, utils.JsonSuccess(http.StatusOK, "", res, nil))
	}
}

func (h *handler) GetByUuid() func(echo.Context) error {
	return func(c echo.Context) error {
		input := dto.GetPermissionRequest{
			Uuid: strings.Trim(c.Param("uuid"), "/"),
		}

		if err := c.Validate(input); err != nil {
			panic(err)
		}

		res, err := h.uc.Get(c.Request().Context(), &input)
		if err != nil {
			panic(err)
		}

		return c.JSON(http.StatusOK, utils.JsonSuccess(http.StatusOK, "", res, nil))
	}
}
