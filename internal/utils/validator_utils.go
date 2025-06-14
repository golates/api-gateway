package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golates/api-gateway/internal/models"
	"net/http"
	"reflect"
	"strings"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		fmt.Println(name)
		if name == "-" {
			return ""
		}
		return name
	})

	return &CustomValidator{validator: validate}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}

	return nil
}

func ValidateBody(w http.ResponseWriter, r *http.Request, data interface{}) error {
	v := r.Context().Value("validator")
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		WriteJSON(w, http.StatusBadRequest, models.MessageAPIResponseError{Message: err.Error()})
		return err
	}

	switch t := v.(type) {
	case *CustomValidator:
		if err := t.Validate(data); err != nil {
			var validationErrors []string
			for _, err := range err.(validator.ValidationErrors) {
				validationErrors = append(validationErrors, parseErrTag(err))
			}

			WriteJSON(w, http.StatusBadRequest, models.ErrorsArrayAPIResponseError{
				Message: "there was a problem with data validation",
				Errors:  validationErrors,
			})
			return errors.New("there was a problem with data validation")
		}
	default:
		WriteJSON(w, http.StatusBadRequest, models.MessageAPIResponseError{Message: "there was a problem with data validation"})
		return errors.New("there was a problem with data validation")
	}

	return nil
}

func parseErrTag(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("Field '%s' cannot be blank", err.Field())
	case "email":
		return fmt.Sprintf("Field '%s' must be a valid email address", err.Field())
	case "oneof":
		return fmt.Sprintf("Field '%s' must be one of the values %v", err.Field(), err.Param())
	default:
		return fmt.Sprintf("Field '%s': '%v' must satisfy '%s' '%v' criteria", err.Field(), err.Value(), err.Tag(), err.Param())
	}
}
