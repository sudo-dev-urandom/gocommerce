package helper

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/thedevsaddam/govalidator"
)

func GetValidationError(err interface{}) []map[string]interface{} {
	errorList := []map[string]interface{}{}
	for _, err := range err.(validator.ValidationErrors) {

		errorMessage := map[string]interface{}{
			"namespace":        err.Namespace(),
			"field":            err.Field(),
			"struct_namespace": err.StructNamespace(),
			"struct_field":     err.StructField(),
			"tag":              err.Tag(),
			"actual_tag":       err.ActualTag(),
			"kind":             err.Kind(),
			"type":             err.Type(),
			"value":            err.Value(),
			"param":            err.Param(),
		}

		errorList = append(errorList, errorMessage)
	}

	return errorList
}

func ValidateRequestPayload(c echo.Context, rules govalidator.MapData, data interface{}) (i interface{}) {
	opts := govalidator.Options{
		Request: c.Request(),
		Data:    data,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	mappedError := v.ValidateJSON()
	if len(mappedError) > 0 {
		i = mappedError
	}

	return i
}

func ValidateRequestFormData(c echo.Context, rules govalidator.MapData) (i interface{}) {
	opts := govalidator.Options{
		Request: c.Request(),
		Rules:   rules,
	}

	v := govalidator.New(opts)
	mappedError := v.Validate()
	if len(mappedError) > 0 {
		i = mappedError
	}

	return i
}
