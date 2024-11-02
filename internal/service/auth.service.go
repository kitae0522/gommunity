package service

import (
	"github.com/kitae0522/gommunity/internal/repository"
	"github.com/kitae0522/gommunity/pkg/crypt"
	"github.com/kitae0522/gommunity/pkg/dto"
	"github.com/kitae0522/gommunity/pkg/exception"
)

type AuthService struct {
	authRepo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{authRepo: repo}
}

func (s *AuthService) Register(req dto.AuthRegisterReq) error {
	// 1. Check if password and confirmation password match
	if err := s.comprePassword(req.Password, req.PasswordConfirm); err != nil {
		return err
	}

	// 2. Create User
	if _, err := s.authRepo.CreateUser(req); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	// 1. Get PasswordInfo
	passwordInfo, err := s.authRepo.GetUserPassword(email)
	if err != nil {
		return "", err
	}

	// 2. Check if password field and password payload match
	if !crypt.VerifyPassword(passwordInfo.HashPassword, password, passwordInfo.Salt) {
		return "", exception.ErrWrongPassword
	}

	// 3. Generate JWT Token
	token, err := crypt.NewToken(string(passwordInfo.Role), passwordInfo.Email, []byte("tempSecret"))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) HandleReset(req dto.AuthHandleResetEntity) error {
	if err := s.authRepo.UpdateUserHandle(req.Email, req.Handle); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) PasswordReset(req dto.AuthPasswordResetEntity) error {
	// 1. Compare NewPassword, NewPasswordConfirm
	if err := s.comprePassword(req.PasswordPayload.NewPassword, req.PasswordPayload.NewPasswordConfirm); err != nil {
		return err
	}

	// 2. Get UserPassword
	passwordInfo, err := s.authRepo.GetUserPassword(req.Email)
	if err != nil {
		return err
	}

	// 3. Compare Password
	if !crypt.VerifyPassword(passwordInfo.HashPassword, req.PasswordPayload.OldPassword, passwordInfo.Salt) {
		return exception.ErrWrongPassword
	}

	// 4. Update UserPassword
	if err := s.authRepo.UpdateUserPassword(passwordInfo.Email, passwordInfo.Salt, req.PasswordPayload.NewPassword); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Withdraw(email string) error {
	if ok, err := s.authRepo.DeleteUser(email); err != nil {
		return err
	} else if !ok {
		return exception.ErrUnableToDeleteUser
	}
	return nil
}

func (s *AuthService) comprePassword(password, confirmPassword string) error {
	if password != confirmPassword {
		return exception.ErrIncorrectConfirmPassword
	}
	return nil
}