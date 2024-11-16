package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/kitae0522/gommunity/pkg/exception"
)

var validate *validator.Validate = validator.New()

func Validate(i interface{}) []exception.ErrValidateResult {
	errlog := make([]exception.ErrValidateResult, 0)
	if errs := validate.Struct(i); errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			errlog = append(errlog, exception.ErrValidateResult{
				IsError: true,
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   err.Value(),
			})
		}
	}
	return errlog
}

func Bind(ctx *fiber.Ctx, targetStruct interface{}, actionMessage string) *exception.ErrResponseCtx {
	var errs []exception.ErrValidateResult
	checkList := []struct {
		field  string
		parser func(interface{}) error
	}{
		{"Params", ctx.ParamsParser},
		{"Headers", ctx.ReqHeaderParser},
		{"Query", ctx.QueryParser},
	}

	for _, checkItem := range checkList {
		if err := parsePayload(targetStruct, checkItem.parser, checkItem.field); err != nil {
			errs = append(errs, *err)
		}
	}

	if ctx.GetReqHeaders()["Content-Type"] != nil && ctx.GetReqHeaders()["Content-Type"][0] == "application/json" {
		if err := parsePayload(targetStruct, ctx.BodyParser, "Body"); err != nil {
			errs = append(errs, *err)
		}
	}

	errs = append(errs, Validate(targetStruct)...)

	if len(errs) > 0 {
		errMessage := fmt.Sprintf("❌ %s 실패. Body Binding 과정에서 문제 발생", actionMessage)
		return exception.GenerateErrorCtx(fiber.StatusBadRequest, errMessage, errs)
	}

	return nil
}

func parsePayload(target interface{}, parseFunc func(interface{}) error, fieldName string) *exception.ErrValidateResult {
	if err := parseFunc(target); err != nil {
		return &exception.ErrValidateResult{
			IsError: true,
			Field:   fieldName,
			Value:   err.Error(),
		}
	}
	return nil
}
