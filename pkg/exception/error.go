package exception

import (
	"github.com/gofiber/fiber/v2"
)

func CreateErrorResponse(ctx *fiber.Ctx, statusCode int, errMessage string, err interface{}) error {
	var errDetail interface{}
	if e, ok := err.(error); ok {
		errDetail = e.Error()
	} else {
		errDetail = err
	}

	ctxError := GenerateErrorCtx(statusCode, errMessage, errDetail)
	return ctx.Status(statusCode).JSON(ctxError)
}

func GenerateErrorCtx(statusCode int, errMessage string, err interface{}) ErrResponseCtx {
	return ErrResponseCtx{
		IsError:    true,
		StatusCode: statusCode,
		Message:    errMessage,
		Error:      err,
	}
}
