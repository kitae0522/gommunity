package utils

import (
	"github.com/gofiber/fiber/v2"
)

func GetRequest(ctx *fiber.Ctx, param interface{}) error {
	if err := ctx.ParamsParser(param); err != nil {
		return err
	}

	if err := ctx.ReqHeaderParser(param); err != nil {
		return err
	}

	if err := ctx.QueryParser(param); err != nil {
		return err
	}

	if ctx.GetReqHeaders()["Content-Type"] != nil && ctx.GetReqHeaders()["Content-Type"][0] == "application/json" {
		err := ctx.BodyParser(param)
		if err != nil {
			return err
		}
	}

	return nil
}
