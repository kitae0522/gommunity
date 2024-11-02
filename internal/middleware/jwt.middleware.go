package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/kitae0522/gommunity/pkg/crypt"
	"github.com/kitae0522/gommunity/pkg/exception"
)

func JWTMiddleware(ctx *fiber.Ctx) error {
	authHeader := strings.Split(ctx.Get("Authorization"), " ")
	if len(authHeader) != 2 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	token := authHeader[1]

	email, err := crypt.ParseJWT(token)
	if err != nil {
		return exception.CreateErrorRes(ctx, fiber.StatusUnauthorized, "❌ 유효하지 않는 토큰 값입니다.", err)
	}
	ctx.Locals("email", email)
	return ctx.Next()
}

func GetEmailFromMiddleware(ctx *fiber.Ctx) string {
	return ctx.Locals("email").(string)
}
