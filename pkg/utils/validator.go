package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ValidateErrRes struct {
	IsError bool
	Field   string
	Tag     string
	Value   interface{}
}

var validate *validator.Validate = validator.New()

func Validate(i interface{}) []ValidateErrRes {
	errlog := make([]ValidateErrRes, 0)
	if errs := validate.Struct(i); errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			errlog = append(errlog, ValidateErrRes{
				IsError: true,
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   err.Value(),
			})
		}
	}
	return errlog
}

func Bind(ctx *fiber.Ctx, targetStruct interface{}) []ValidateErrRes {
	// BodyParse
	if err := ctx.BodyParser(targetStruct); err != nil {
		return []ValidateErrRes{{
			IsError: true,
			Field:   "Body",
			Value:   err.Error(),
		}}
	}

	// Validate
	if err := Validate(targetStruct); len(err) != 0 {
		return err
	}

	return nil
}
