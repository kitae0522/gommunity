package controller

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"

	"github.com/kitae0522/gommunity/internal/dto"
	"github.com/kitae0522/gommunity/internal/middleware"
	"github.com/kitae0522/gommunity/internal/model"
	"github.com/kitae0522/gommunity/internal/repository"
	"github.com/kitae0522/gommunity/internal/service"
	"github.com/kitae0522/gommunity/pkg/utils"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(service *service.AuthService) *AuthController {
	return &AuthController{authService: service}
}

func initAuthDI(dbconn *model.PrismaClient, rdconn *redis.Client) *AuthController {
	repository := repository.NewAuthRepository(dbconn)
	service := service.NewAuthService(repository, rdconn)
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
	if err := utils.Bind(ctx, &createUserPayload, "회원가입"); err != nil {
		return ctx.Status(err.StatusCode).JSON(err)
	}

	err := c.authService.Register(ctx.Context(), createUserPayload)
	if err != nil {
		return ctx.Status(err.StatusCode).JSON(err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.DefaultResponse{
		IsError:    false,
		StatusCode: fiber.StatusCreated,
		Message:    "✅ 회원가입 완료",
	})
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var loginPayload dto.LoginRequest
	if err := utils.Bind(ctx, &loginPayload, "로그인"); err != nil {
		return ctx.Status(err.StatusCode).JSON(err)
	}

	token, err := c.authService.Login(ctx.Context(), loginPayload.Email, loginPayload.Password)
	if err != nil {
		return ctx.Status(err.StatusCode).JSON(err)
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
	if err := utils.Bind(ctx, &passwordResetPayload, "비밀번호 초기화"); err != nil {
		return ctx.Status(err.StatusCode).JSON(err)
	}

	resetEntity := dto.PasswordResetEntity{
		ID:              middleware.GetIdFromMiddleware(ctx),
		PasswordPayload: &passwordResetPayload,
	}

	if err := c.authService.PasswordReset(ctx.Context(), resetEntity); err != nil {
		return ctx.Status(err.StatusCode).JSON(err)
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.DefaultResponse{
		IsError:    false,
		StatusCode: fiber.StatusOK,
		Message:    "✅ 비밀번호 초기화 완료",
	})
}

func (c *AuthController) Withdraw(ctx *fiber.Ctx) error {
	var withdrawPayload dto.WithdrawRequest
	withdrawPayload.ID = middleware.GetIdFromMiddleware(ctx)

	if err := c.authService.Withdraw(ctx.Context(), withdrawPayload.ID); err != nil {
		return ctx.Status(err.StatusCode).JSON(err)
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.DefaultResponse{
		IsError:    false,
		StatusCode: fiber.StatusOK,
		Message:    "✅ 유저 탈퇴 완료",
	})
}
