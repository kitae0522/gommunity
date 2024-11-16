package controller

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/kitae0522/gommunity/internal/model"
)

func EnrollRouter(app *fiber.App, dbconn *model.PrismaClient, rdconn *redis.Client) {
	apiRouter := app.Group("/api")
	initAuthRouter(apiRouter, initAuthDI(dbconn, rdconn))
	initThreadRouter(apiRouter, initThreadDI(dbconn, rdconn))

	apiRouter.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "pong",
		})
	})
}
