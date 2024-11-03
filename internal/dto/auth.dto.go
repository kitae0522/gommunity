package dto

import "github.com/kitae0522/gommunity/internal/model"

type DefaultResponse struct {
	IsError    bool   `json:"isError"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type AuthRegisterRequest struct {
	Handle          string `json:"handle" validate:"required"`
	Name            string `json:"name" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required"`
	Email           string `json:"email" validate:"required"`
}

type AuthLoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthLoginResponse struct {
	IsError    bool   `json:"isError"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Token      string `json:"token"`
}

type AuthHandleResetEntity struct {
	Email  string
	Handle string
}

type AuthPasswordResetRequest struct {
	OldPassword        string `json:"oldPassword" validate:"required"`
	NewPassword        string `json:"newPassword" validate:"required"`
	NewPasswordConfirm string `json:"newPasswordConfirm" validate:"required"`
}

type AuthPasswordResetEntity struct {
	Email           string                    `json:"email" validate:"required"`
	PasswordPayload *AuthPasswordResetRequest `json:"passwordPayload"`
}

type AuthWithdrawRequest struct {
	Email string `json:"email" validate:"required"`
}

type PasswordEntity struct {
	Email        string
	HashPassword string
	Salt         string
	Role         model.UserRoles
	Handle       string
}
