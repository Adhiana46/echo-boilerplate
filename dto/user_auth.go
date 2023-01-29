package dto

import (
	"github.com/Adhiana46/echo-boilerplate/entity"
	"github.com/golang-jwt/jwt/v4"
)

func NewUserClaims(eUser *entity.User, eDevice *entity.UserDevice, regClaims jwt.RegisteredClaims) *UserClaims {
	return &UserClaims{
		User:             NewUserResponse(eUser),
		Device:           NewUserDevice(eDevice),
		RegisteredClaims: regClaims,
	}
}

type UserClaims struct {
	User   *UserResponse `json:"user,omitempty"`
	Device *UserDevice   `json:"device,omitempty"`
	jwt.RegisteredClaims
}

func NewUserDevice(e *entity.UserDevice) *UserDevice {
	return &UserDevice{
		Uuid:       e.Uuid,
		Token:      e.Token,
		IP:         e.IP,
		Location:   e.Location,
		Platform:   e.Platform,
		UserAgent:  e.UserAgent,
		AppVersion: e.AppVersion,
		Vendor:     e.Vendor,
	}
}

type UserDevice struct {
	Uuid       string `json:"uuid" validate:"omitempty,required"`
	Token      string `json:"token" validate:"omitempty,required"`
	IP         string `json:"ip" validate:"omitempty,required"`
	Location   string `json:"location" validate:"omitempty,required"`
	Platform   string `json:"platform" validate:"omitempty,required"`
	UserAgent  string `json:"user_agent" validate:"omitempty,required"`
	AppVersion string `json:"app_version" validate:"omitempty,required"`
	Vendor     string `json:"vendor" validate:"omitempty,required"`
}

type SignInRequest struct {
	Username string      `json:"username" validate:"required"`
	Password string      `json:"password" validate:"required"`
	Device   *UserDevice `json:"device" validate:"omitempty"`
}

type SignInResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	Device       UserDevice `json:"device"`
}

type SignOutRequest struct {
	AccessToken  string     `json:"access_token" validate:"required"`
	RefreshToken string     `json:"refresh_token"`
	Device       UserDevice `json:"device" validate:"omitempty"`
}

type SignOutResponse struct {
	// Empty
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
