package http

import (
	"net/http"
	"strings"
	"sync"

	"github.com/Adhiana46/echo-boilerplate/dto"
	"github.com/Adhiana46/echo-boilerplate/internal/role"
	"github.com/Adhiana46/echo-boilerplate/pkg/utils"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Store() func(echo.Context) error
	Update() func(echo.Context) error
	Delete() func(echo.Context) error
	GetByUuid() func(echo.Context) error
	GetAll() func(echo.Context) error
}

var (
	handlerInstance     *handler
	handlerInstanceOnce sync.Once
)

type handler struct {
	uc role.RoleUsecase
}

func NewRoleHttpHandler(uc role.RoleUsecase) Handler {
	handlerInstanceOnce.Do(func() {
		handlerInstance = &handler{
			uc: uc,
		}
	})

	return handlerInstance
}

func (h *handler) Store() func(echo.Context) error {
	return func(c echo.Context) error {
		input := dto.CreateRoleRequest{}

		if err := c.Bind(&input); err != nil {
			return err
		}

		if err := c.Validate(input); err != nil {
			return err
		}

		res, err := h.uc.CreateRole(c.Request().Context(), &input)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, utils.JsonSuccess(http.StatusCreated, "", res, nil))
	}
}

func (h *handler) Update() func(echo.Context) error {
	return func(c echo.Context) error {
		input := dto.UpdateRoleRequest{}

		if err := c.Bind(&input); err != nil {
			return err
		}

		input.Uuid = strings.Trim(c.Param("uuid"), "/")

		if err := c.Validate(input); err != nil {
			return err
		}

		res, err := h.uc.UpdateRole(c.Request().Context(), &input)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, utils.JsonSuccess(http.StatusOK, "", res, nil))
	}
}

func (h *handler) Delete() func(echo.Context) error {
	return func(c echo.Context) error {
		input := dto.DeleteRoleRequest{
			Uuid: strings.Trim(c.Param("uuid"), "/"),
		}

		if err := c.Validate(input); err != nil {
			return err
		}

		res, err := h.uc.DeleteRole(c.Request().Context(), &input)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, utils.JsonSuccess(http.StatusOK, "", res, nil))
	}
}

func (h *handler) GetByUuid() func(echo.Context) error {
	return func(c echo.Context) error {
		input := dto.GetRoleRequest{
			Uuid: strings.Trim(c.Param("uuid"), "/"),
		}

		if err := c.Validate(input); err != nil {
			return err
		}

		res, err := h.uc.Get(c.Request().Context(), &input)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, utils.JsonSuccess(http.StatusOK, "", res, nil))
	}
}

func (h *handler) GetAll() func(echo.Context) error {
	return func(c echo.Context) error {
		input := dto.GetListRoleRequest{}

		if err := c.Bind(&input); err != nil {
			return err
		}

		if err := c.Validate(input); err != nil {
			return err
		}

		res, err := h.uc.GetList(c.Request().Context(), &input)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, utils.JsonSuccess(http.StatusOK, "", res.Data, res.Pagination))
	}
}
