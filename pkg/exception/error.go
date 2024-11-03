package exception

import (
	"github.com/gofiber/fiber/v2"
)

func CreateErrorRes(ctx *fiber.Ctx, statusCode int, errMessage string, err interface{}) error {
	var errorDetail interface{}
	if e, ok := err.(error); ok {
		errorDetail = e.Error()
	} else {
		errorDetail = err
	}

	return ctx.Status(statusCode).JSON(ErrResponseCtx{
		IsError:    true,
		StatusCode: statusCode,
		Message:    errMessage,
		Error:      errorDetail,
	})
}
