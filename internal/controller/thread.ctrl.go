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

type ThreadController struct {
	threadService *service.ThreadService
}

func NewThreadController(service *service.ThreadService) *ThreadController {
	return &ThreadController{threadService: service}
}

func initThreadDI(dbconn *model.PrismaClient, rdconn *redis.Client) *ThreadController {
	repository := repository.NewThreadRepository(dbconn)
	service := service.NewThreadService(repository, rdconn)
	handler := NewThreadController(service)
	return handler
}

func initThreadRouter(router fiber.Router, handler *ThreadController) {
	threadRouter := router.Group("/thread")
	handler.Accessible(threadRouter)
	handler.Restricted(threadRouter)
}

func (c *ThreadController) Accessible(router fiber.Router) {
	router.Get("/", c.ListThread)
	router.Get("/user/:handle", c.ListThreadByHandle)
	router.Get("/:threadID", c.GetThreadByID)
}

func (c *ThreadController) Restricted(router fiber.Router) {
	router.Use(middleware.JWTMiddleware)
	router.Post("/", c.CreateThread)
}

func (c *ThreadController) CreateThread(ctx *fiber.Ctx) error {
	var createThreadPayload dto.CreateThreadRequest
	createThreadPayload.UserID = middleware.GetIdFromMiddleware(ctx)

	if err := utils.Bind(ctx, &createThreadPayload, "쓰레드 생성"); err != nil {
		return ctx.Status(err.StatusCode).JSON(err)
	}

	thread, err := c.threadService.CreateThread(ctx.Context(), &createThreadPayload)
	if err != nil {
		return ctx.Status(err.StatusCode).JSON(err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.CreateThreadReponse{
		IsError:    false,
		StatusCode: fiber.StatusOK,
		Message:    "✅ 쓰레드 생성 완료",
		Thread:     *thread,
	})
}

func (c *ThreadController) ListThread(ctx *fiber.Ctx) error {
	threads, err := c.threadService.ListThread(ctx.Context())
	if err != nil {
		return ctx.Status(err.StatusCode).JSON(err)
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ListThreadResponse{
		IsError:    false,
		StatusCode: fiber.StatusOK,
		Message:    "✅ 모든 쓰레드 조회 완료",
		Threads:    threads,
	})
}

func (c *ThreadController) ListThreadByHandle(ctx *fiber.Ctx) error {
	var listThreadPayload dto.ListThreadByHandleRequest
	if err := utils.Bind(ctx, &listThreadPayload, "모든 쓰레드 조회"); err != nil {
		return ctx.Status(err.StatusCode).JSON(err)
	}

	threads, err := c.threadService.ListThreadByHandle(ctx.Context(), listThreadPayload.Handle)
	if err != nil {
		return ctx.Status(err.StatusCode).JSON(err)
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.ListThreadByHandleResponse{
		IsError:    false,
		StatusCode: fiber.StatusOK,
		Message:    "✅ 모든 쓰레드 조회 완료",
		Handle:     listThreadPayload.Handle,
		Threads:    threads,
	})
}

func (c *ThreadController) GetThreadByID(ctx *fiber.Ctx) error {
	var getThreadPayload dto.GetThreadByIDRequest
	if err := utils.Bind(ctx, &getThreadPayload, "쓰레드 조회"); err != nil {
		return ctx.Status(err.StatusCode).JSON(err)
	}

	thread, err := c.threadService.GetThreadByID(ctx.Context(), getThreadPayload.ThreadID)
	if err != nil {
		return ctx.Status(err.StatusCode).JSON(err)
	}

	comments, err := c.threadService.CommentsByID(ctx.Context(), getThreadPayload.ThreadID)
	if err != nil {
		return ctx.Status(err.StatusCode).JSON(err)
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.GetThreadByIDResponse{
		IsError:    false,
		StatusCode: fiber.StatusOK,
		Message:    "✅ 모든 쓰레드 조회 완료",
		Thread:     thread,
		SubThread:  comments,
	})
}
