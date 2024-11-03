package utils

import (
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

func Bind(ctx *fiber.Ctx, targetStruct interface{}) []exception.ErrValidateResult {
	errs := make([]exception.ErrValidateResult, 0)

	if err := parsePayload(targetStruct, ctx.ParamsParser, "Params"); err != nil {
		errs = append(errs, *err)
	}

	if err := parsePayload(targetStruct, ctx.ReqHeaderParser, "Headers"); err != nil {
		errs = append(errs, *err)
	}

	if err := parsePayload(targetStruct, ctx.QueryParser, "Query"); err != nil {
		errs = append(errs, *err)
	}

	if err := parsePayload(targetStruct, ctx.BodyParser, "Body"); err != nil {
		errs = append(errs, *err)
	}

	if ctx.GetReqHeaders()["Content-Type"] != nil && ctx.GetReqHeaders()["Content-Type"][0] == "application/json" {
		if err := parsePayload(targetStruct, ctx.BodyParser, "Body"); err != nil {
			errs = append(errs, *err)
		}
	}

	errs = append(errs, Validate(targetStruct)...)
	return errs
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
