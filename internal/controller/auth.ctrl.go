package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kitae0522/gommunity/internal/dto"
	"github.com/kitae0522/gommunity/internal/middleware"
	"github.com/kitae0522/gommunity/internal/model"
	"github.com/kitae0522/gommunity/internal/repository"
	"github.com/kitae0522/gommunity/internal/service"
	"github.com/kitae0522/gommunity/pkg/exception"
	"github.com/kitae0522/gommunity/pkg/utils"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(service *service.AuthService) *AuthController {
	return &AuthController{authService: service}
}

func initAuthDI(dbconn *model.PrismaClient) *AuthController {
	repository := repository.NewAuthRepository(dbconn)
	service := service.NewAuthService(repository)
	handler := NewAuthController(service)
	return handler
}

func initAuthRouter(router fiber.Router, handler *AuthController) {
	authRouter := router.Group("/auth")
	handler.Accessible(authRouter)
	handler.Restricted(authRouter)
}

func (c *AuthController) Accessible(router fiber.Router) {
	router.Post("/register", c.Register)
	router.Post("/login", c.Login)
}

func (c *AuthController) Restricted(router fiber.Router) {
	router.Use(middleware.JWTMiddleware)
	router.Patch("/reset", c.PasswordReset)
	router.Delete("/withdraw", c.Withdraw)
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var createUserPayload dto.RegisterRequest
	if errs := utils.Bind(ctx, &createUserPayload); len(errs) > 0 {
		return exception.CreateErrorRes(ctx, fiber.StatusBadRequest, "❌ 회원가입 실패. Body Binding 과정에서 문제 발생", errs)
	}

	err := c.authService.Register(createUserPayload)
	if err != nil {
		if _, uniqueErr := model.IsErrUniqueConstraint(err); uniqueErr {
			return exception.CreateErrorRes(ctx, fiber.StatusInternalServerError, "❌ 회원가입 실패. 중복된 유저가 존재합니다.", err)
		}
		if err == exception.ErrIncorrectConfirmPassword {
			return exception.CreateErrorRes(ctx, fiber.StatusInternalServerError, "❌ 회원가입 실패. 패스워드가 일치하지 않습니다.", err)
		}
		return exception.CreateErrorRes(ctx, fiber.StatusInternalServerError, "❌ 회원가입 실패. Repository에서 문제 발생", err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.DefaultResponse{
		IsError:    false,
		StatusCode: fiber.StatusCreated,
		Message:    "✅ 회원가입 완료",
	})
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var loginPayload dto.LoginRequest
	if errs := utils.Bind(ctx, &loginPayload); len(errs) > 0 {
		return exception.CreateErrorRes(ctx, fiber.StatusBadRequest, "❌ 회원가입 실패. Body Binding 과정에서 문제 발생", errs)
	}

	token, err := c.authService.Login(loginPayload.Email, loginPayload.Password)
	if err != nil {
		switch err {
		case model.ErrNotFound:
			return exception.CreateErrorRes(ctx, fiber.StatusInternalServerError, "❌ 로그인 실패. 존재하지 않는 사용자입니다.", err)
		case exception.ErrWrongPassword:
			return exception.CreateErrorRes(ctx, fiber.StatusInternalServerError, "❌ 로그인 실패. 패스워드가 일치하지 않습니다.", err)
		default:
			return exception.CreateErrorRes(ctx, fiber.StatusInternalServerError, "❌ 로그인 실패. Repository에서 문제 발생", err)
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.LoginResponse{
		IsError:    false,
		StatusCode: fiber.StatusOK,
		Message:    "✅ 로그인 완료",
		Token:      token,
	})
}

func (c *AuthController) PasswordReset(ctx *fiber.Ctx) error {
	var passwordResetPayload dto.PasswordResetRequest
	if errs := utils.Bind(ctx, passwordResetPayload); len(errs) > 0 {
		return exception.CreateErrorRes(ctx, fiber.StatusBadRequest, "❌ 비밀번호 초기화 실패. Body Binding 과정에서 문제 발생", errs)
	}

	resetEntity := dto.PasswordResetEntity{
		Email:           middleware.GetEmailFromMiddleware(ctx),
		PasswordPayload: &passwordResetPayload,
	}

	if err := c.authService.PasswordReset(resetEntity); err != nil {
		switch err {
		case model.ErrNotFound:
			return exception.CreateErrorRes(ctx, fiber.StatusInternalServerError, "❌ 비밀번호 초기화 실패. 존재하지 않는 사용자입니다.", err)
		case exception.ErrWrongPassword:
			return exception.CreateErrorRes(ctx, fiber.StatusInternalServerError, "❌ 비밀번호 초기화 실패. 패스워드가 일치하지 않습니다.", err)
		default:
			return exception.CreateErrorRes(ctx, fiber.StatusInternalServerError, "❌ 비밀번호 초기화 실패. Repository에서 문제 발생", err)
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.DefaultResponse{
		IsError:    false,
		StatusCode: fiber.StatusOK,
		Message:    "✅ 비밀번호 초기화 완료",
	})
}

func (c *AuthController) Withdraw(ctx *fiber.Ctx) error {
	var withdrawPayload dto.WithdrawRequest
	withdrawPayload.Email = middleware.GetEmailFromMiddleware(ctx)

	if err := c.authService.Withdraw(withdrawPayload.Email); err != nil {
		switch err {
		case model.ErrNotFound:
			return exception.CreateErrorRes(ctx, fiber.StatusInternalServerError, "❌ 유저 탈퇴 실패. 존재하지 않는 사용자입니다.", err)
		case exception.ErrUnableToDeleteUser:
			return exception.CreateErrorRes(ctx, fiber.StatusInternalServerError, "❌ 유저 탈퇴 실패. 유저를 삭제할 수 없습니다.", err)
		default:
			return exception.CreateErrorRes(ctx, fiber.StatusInternalServerError, "❌ 유저 탈퇴 실패. Repository에서 문제 발생", err)
		}
	}

	return ctx.Status(fiber.StatusNoContent).JSON(dto.DefaultResponse{
		IsError:    false,
		StatusCode: fiber.StatusNoContent,
		Message:    "✅ 유저 탈퇴 완료",
	})
}
