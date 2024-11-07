package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kitae0522/gommunity/internal/config"
	"github.com/kitae0522/gommunity/internal/dto"
	"github.com/kitae0522/gommunity/internal/model"
	"github.com/kitae0522/gommunity/internal/repository"
	"github.com/kitae0522/gommunity/pkg/crypt"
	"github.com/kitae0522/gommunity/pkg/exception"
)

type AuthService struct {
	authRepo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{authRepo: repo}
}

func (s *AuthService) Register(req dto.RegisterRequest) *exception.ErrResponseCtx {
	if err := s.comparePassword(req.Password, req.PasswordConfirm); err != nil {
		return exception.GenerateErrorCtx(fiber.StatusBadRequest, "❌ 회원가입 실패. 패스워드가 일치하지 않습니다.", err)
	}

	if _, err := s.authRepo.CreateUser(req); err != nil {
		if _, uniqueErr := model.IsErrUniqueConstraint(err); uniqueErr {
			return exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 회원가입 실패. 중복된 유저가 존재합니다.", err)
		}
		return exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 회원가입 실패. Repository에서 문제 발생", err)
	}

	return nil
}

func (s *AuthService) Login(email, password string) (string, *exception.ErrResponseCtx) {
	passwordInfo, err := s.authRepo.GetUserPasswordByEmail(email)
	if err != nil {
		switch err {
		case model.ErrNotFound:
			return "", exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 로그인 실패. 존재하지 않는 사용자입니다.", err)
		default:
			return "", exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 로그인 실패. Repository에서 문제 발생", err)
		}
	}

	if !crypt.VerifyPassword(passwordInfo.HashPassword, password, passwordInfo.Salt) {
		return "", exception.GenerateErrorCtx(fiber.StatusBadRequest, "❌ 로그인 실패. 패스워드가 일치하지 않습니다.", err)
	}

	token, err := crypt.NewToken(string(passwordInfo.Role), passwordInfo.ID, []byte(config.Envs.JWTSecret))
	if err != nil {
		return "", exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 로그인 실패. 토큰 생성 중 문제가 발생했습니다.", err)
	}

	return token, nil
}

func (s *AuthService) HandleReset(req dto.HandleResetEntity) error {
	if err := s.authRepo.UpdateUserHandle(req.ID, req.Handle); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) PasswordReset(req dto.PasswordResetEntity) *exception.ErrResponseCtx {
	if err := s.comparePassword(req.PasswordPayload.NewPassword, req.PasswordPayload.NewPasswordConfirm); err != nil {
		return exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 비밀번호 초기화 실패. 패스워드가 일치하지 않습니다.", err)
	}

	passwordInfo, err := s.authRepo.GetUserPasswordByID(req.ID)
	if err != nil {
		switch err {
		case model.ErrNotFound:
			return exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 비밀번호 초기화 실패. 존재하지 않는 사용자입니다.", err)
		default:
			return exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 비밀번호 초기화 실패. Repository에서 문제 발생", err)
		}
	}

	if !crypt.VerifyPassword(passwordInfo.HashPassword, req.PasswordPayload.OldPassword, passwordInfo.Salt) {
		return exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 비밀번호 초기화 실패. 패스워드가 일치하지 않습니다.", err)
	}

	if err := s.authRepo.UpdateUserPassword(passwordInfo.ID, passwordInfo.Salt, req.PasswordPayload.NewPassword); err != nil {
		return exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 비밀번호 초기화 실패. Repository에서 문제 발생", err)
	}

	return nil
}

func (s *AuthService) Withdraw(ID string) *exception.ErrResponseCtx {
	if ok, err := s.authRepo.DeleteUser(ID); err != nil {
		switch err {
		case model.ErrNotFound:
			return exception.GenerateErrorCtx(fiber.StatusNotFound, "❌ 유저 탈퇴 실패. 존재하지 않는 사용자입니다.", err)
		default:
			return exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 유저 탈퇴 실패. Repository에서 문제 발생", err)
		}
	} else if !ok {
		return exception.GenerateErrorCtx(fiber.StatusInternalServerError, "❌ 유저 탈퇴 실패. 유저를 삭제할 수 없습니다.", err)
	}

	return nil
}

func (s *AuthService) comparePassword(password, confirmPassword string) error {
	if password != confirmPassword {
		return exception.ErrIncorrectConfirmPassword
	}
	return nil
}
