package dto

import "github.com/kitae0522/gommunity/internal/model"

type DefaultRes struct {
	IsError    bool   `json:"isError"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type AuthRegisterReq struct {
	Handle          string `json:"handle" validate:"required"`
	Name            string `json:"name" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required"`
	Email           string `json:"email" validate:"required"`
}

type AuthLoginReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthLoginRes struct {
	IsError    bool   `json:"isError"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Token      string `json:"token"`
}

type AuthHandleResetEntity struct {
	Email  string
	Handle string
}

type AuthPasswordResetReq struct {
	OldPassword        string `json:"oldPassword" validate:"required"`
	NewPassword        string `json:"newPassword" validate:"required"`
	NewPasswordConfirm string `json:"newPasswordConfirm" validate:"required"`
}

type AuthPasswordResetEntity struct {
	Email           string `json:"email" validate:"requierd"`
	PasswordPayload AuthPasswordResetReq
}

type AuthWithdrawEntity struct {
	Email string `json:"email" validate:"requierd"`
}

type PasswordEntity struct {
	Email        string
	HashPassword string
	Salt         string
	Role         model.UserRoles
	Handle       string
}
