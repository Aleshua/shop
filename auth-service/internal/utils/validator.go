package utils

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(ctx context.Context, s interface{}) error {
	return validate.StructCtx(ctx, s)
}

// Выдаёт текст для каждой найденной ошибки валидации (параметр - текст ошибки)
func TranslateValidatorError(err error) map[string]string {
	var verrs validator.ValidationErrors
	errs := map[string]string{}
	if errors.As(err, &verrs) {
		for _, verr := range verrs {
			field := CamelToSnake(verr.Field())
			switch verr.Tag() {
			case "required":
				errs[field] = "обязательный параметр"
			case "email":
				errs[field] = "параметр должен быть в формате email"
			case "url":
				errs[field] = "параметр должен быть в формате url"
			case "datetime":
				errs[field] = "параметр должен быть в формате datetime"
			case "min":
				if verr.Kind() == reflect.String {
					errs[field] = fmt.Sprintf("минимальная длина строки: %v", verr.Param())
				} else if verr.Kind() == reflect.Int || verr.Kind() == reflect.Int16 ||
					verr.Kind() == reflect.Int32 || verr.Kind() == reflect.Int64 ||
					verr.Kind() == reflect.Int8 {
					errs[field] = fmt.Sprintf("минимальное значение: %v", verr.Param())
				} else {
					errs[field] = "ошибка валидации параметра"
				}
			case "max":
				if verr.Kind() == reflect.String {
					errs[field] = fmt.Sprintf("максимальная длина строки: %v", verr.Param())
				} else if verr.Kind() == reflect.Int || verr.Kind() == reflect.Int16 ||
					verr.Kind() == reflect.Int32 || verr.Kind() == reflect.Int64 ||
					verr.Kind() == reflect.Int8 {
					errs[field] = fmt.Sprintf("максимальное значение: %v", verr.Param())
				} else {
					errs[field] = "ошибка валидации параметра"
				}
			default:
				errs[field] = "ошибка валидации параметра"
			}
		}
	}

	return errs
}
