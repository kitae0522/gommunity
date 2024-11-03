package dto

import "github.com/kitae0522/gommunity/internal/model"

type RegisterRequest struct {
	Handle          string `json:"handle" validate:"required"`
	Name            string `json:"name" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	IsError    bool   `json:"isError"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Token      string `json:"token"`
}

type HandleResetEntity struct {
	ID     string
	Handle string
}

type PasswordResetRequest struct {
	OldPassword        string `json:"oldPassword" validate:"required"`
	NewPassword        string `json:"newPassword" validate:"required"`
	NewPasswordConfirm string `json:"newPasswordConfirm" validate:"required"`
}

type PasswordResetEntity struct {
	ID              string                `json:"id"`
	PasswordPayload *PasswordResetRequest `json:"passwordPayload"`
}

type WithdrawRequest struct {
	ID string `json:"id" validate:"required"`
}

type PasswordEntity struct {
	ID           string
	HashPassword string
	Salt         string
	Role         model.UserRoles
	Handle       string
}
