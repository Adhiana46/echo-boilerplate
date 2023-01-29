package http

import (
	"net/http"
	"strings"

	"github.com/Adhiana46/echo-boilerplate/dto"
	"github.com/Adhiana46/echo-boilerplate/internal/user"
	"github.com/Adhiana46/echo-boilerplate/pkg/utils"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Store() func(echo.Context) error
	Update() func(echo.Context) error
	Delete() func(echo.Context) error
	GetByUuid() func(echo.Context) error
	GetAll() func(echo.Context) error

	SignIn() func(echo.Context) error
	SignOut() func(echo.Context) error
	RefreshToken() func(echo.Context) error
}

type handler struct {
	uc user.UserUsecase
}

func NewUserHttpHandler(uc user.UserUsecase) Handler {
	return &handler{
		uc: uc,
	}
}

func (h *handler) Store() func(echo.Context) error {
	return func(c echo.Context) error {
		input := dto.CreateUserRequest{}

		if err := c.Bind(&input); err != nil {
			return err
		}

		if err := c.Validate(input); err != nil {
			return err
		}

		res, err := h.uc.CreateUser(c.Request().Context(), &input)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, utils.JsonSuccess(http.StatusCreated, "", res, nil))
	}
}

func (h *handler) Update() func(echo.Context) error {
	return func(c echo.Context) error {
		input := dto.UpdateUserRequest{}

		if err := c.Bind(&input); err != nil {
			return err
		}

		input.Uuid = strings.Trim(c.Param("uuid"), "/")

		if err := c.Validate(input); err != nil {
			return err
		}

		res, err := h.uc.UpdateUser(c.Request().Context(), &input)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, utils.JsonSuccess(http.StatusOK, "", res, nil))
	}
}

func (h *handler) Delete() func(echo.Context) error {
	return func(c echo.Context) error {
		input := dto.DeleteUserRequest{
			Uuid: strings.Trim(c.Param("uuid"), "/"),
		}

		if err := c.Validate(input); err != nil {
			return err
		}

		res, err := h.uc.DeleteUser(c.Request().Context(), &input)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, utils.JsonSuccess(http.StatusOK, "", res, nil))
	}
}

func (h *handler) GetByUuid() func(echo.Context) error {
	return func(c echo.Context) error {
		input := dto.GetUserRequest{
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
		input := dto.GetListUserRequest{}

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

func (h *handler) SignIn() func(echo.Context) error {
	return func(c echo.Context) error {
		input := dto.SignInRequest{}

		if err := c.Bind(&input); err != nil {
			return err
		}

		if err := c.Validate(input); err != nil {
			return err
		}

		res, err := h.uc.SignIn(c.Request().Context(), &input)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, utils.JsonSuccess(http.StatusOK, "", res, nil))
	}
}

func (h *handler) SignOut() func(echo.Context) error {
	return func(c echo.Context) error {
		input := dto.SignOutRequest{}

		if err := c.Bind(&input); err != nil {
			return err
		}

		if err := c.Validate(input); err != nil {
			return err
		}

		res, err := h.uc.SignOut(c.Request().Context(), &input)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, utils.JsonSuccess(http.StatusOK, "", res, nil))
	}
}

func (h *handler) RefreshToken() func(echo.Context) error {
	return func(c echo.Context) error {
		input := dto.RefreshTokenRequest{}

		if err := c.Bind(&input); err != nil {
			return err
		}

		if err := c.Validate(input); err != nil {
			return err
		}

		res, err := h.uc.RefreshToken(c.Request().Context(), &input)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, utils.JsonSuccess(http.StatusOK, "", res, nil))
	}
}
