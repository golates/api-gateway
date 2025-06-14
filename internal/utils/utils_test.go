package utils

import (
	"context"
	"errors"
	"fmt"
	validator2 "github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSONTableDriven(t *testing.T) {
	var tests = []struct {
		status         int
		value          any
		expectedValue  string
		expectedStatus int
	}{
		{http.StatusOK, nil, "null", http.StatusOK},
		{http.StatusOK, struct {
			Test string `json:"test"`
		}{
			"Test",
		}, "{\"test\":\"Test\"}",
			http.StatusOK},
		{http.StatusOK, struct {
			ErrorMessage string `json:"error_message"`
		}{
			"Test",
		}, "{\"error_message\":\"Test\"}",
			http.StatusOK},
		{http.StatusOK, make(chan int), "", http.StatusInternalServerError},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("WriteJSON test: expected %v status", test.status)
		t.Run(testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteJSON(w, test.status, test.value)
			if test.expectedStatus != w.Result().StatusCode {
				t.Errorf("status %v not equal to %v", w.Result().StatusCode, test.expectedStatus)
			}

			if test.expectedValue != strings.ReplaceAll(w.Body.String(), "\n", "") {
				t.Errorf("value %v not equal to %v", w.Body.String(), test.expectedValue)
			}
		})
	}
}

func TestParseGRPCError(t *testing.T) {
	var tests = []struct {
		err             error
		expectedMessage string
		expectedStatus  int
	}{
		{nil, "", http.StatusBadRequest},
		{errors.New("test"), "Internal server error", http.StatusInternalServerError},
		{status.Error(codes.PermissionDenied, "Permission denied"), "Permission denied", http.StatusBadRequest},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("ParseGRPCError test: expected %v error", test.expectedMessage)
		t.Run(testName, func(t *testing.T) {
			s, message := ParseGRPCError(test.err)
			if test.expectedMessage != message {
				t.Errorf("error message %v not equal to %v", message, test.expectedMessage)
			}

			if test.expectedStatus != s {
				t.Errorf("error status %v not equal to %v", s, test.expectedStatus)
			}
		})
	}
}

func TestNewCustomValidator(t *testing.T) {
	t.Run("NewCustomValidator test", func(t *testing.T) {
		NewCustomValidator()
	})
}

func TestParseErrTagTableDriven(t *testing.T) {
	var validator = NewCustomValidator()
	err := validator.Validate(struct {
		EmptyTest   string `json:"empty_test" validate:"required"`
		EmailTest   string `json:"email_test" validate:"email"`
		OneOfTest   string `json:"one_of_test" validate:"oneof=test1 test2"`
		DefaultTest int    `json:"default_test" validate:"min=2"`
		EmptyName   string `validate:"required"`
		Dash        string `json:"-"`
	}{
		EmptyTest:   "",
		EmailTest:   "test",
		OneOfTest:   "wrong",
		DefaultTest: 1,
		EmptyName:   "",
		Dash:        "",
	})
	var tests = []string{
		"Field 'empty_test' cannot be blank",
		"Field 'email_test' must be a valid email address",
		"Field 'one_of_test' must be one of the values test1 test2",
		"Field 'default_test': '1' must satisfy 'min' '2' criteria",
		"Field 'EmptyName' cannot be blank",
	}

	for index, err := range err.(validator2.ValidationErrors) {
		testName := fmt.Sprintf("ParseErrTag test: expected %v status", tests[index])
		t.Run(testName, func(t *testing.T) {
			if tests[index] != parseErrTag(err) {
				t.Errorf("error %v not equal to %v", tests[index], parseErrTag(err))
			}
		})
	}
}

func TestValidateBodyTableDriven(t *testing.T) {
	var validator = NewCustomValidator()

	var tests = []struct {
		data           interface{}
		body           string
		expectedError  error
		expectedReturn string
	}{
		{new(struct {
			Test string `json:"test" validate:"required"`
		}), "{\"test\": \"\"}", errors.New("there was a problem with data validation"),
			"{\"message\":\"there was a problem with data validation\",\"errors\":[\"Field 'test' cannot be blank\"]}"},
		{new(struct {
			Test string `json:"test" validate:"required"`
		}), "{\"test\": \"test\"}", nil,
			""},
		{new(struct {
			Test string `json:"test" validate:"required"`
		}), "{\"test\": \"test\"}", errors.New("there was a problem with data validation"),
			"{\"message\":\"there was a problem with data validation\"}"},
		{new(struct {
			Test string `json:"test" validate:"required"`
		}), "}", errors.New("invalid character '}' looking for beginning of value"),
			"{\"message\":\"invalid character '}' looking for beginning of value\"}"},
	}

	for index, test := range tests {
		testName := fmt.Sprintf("ValidateBody test: expected %v error", test.expectedError)
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			if index != 2 {
				ctx = context.WithValue(context.Background(), "validator", validator)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequestWithContext(ctx, "POST", "/", strings.NewReader(test.body))

			err := ValidateBody(w, r, test.data)
			if (err != nil && test.expectedError != nil && err.Error() != test.expectedError.Error()) ||
				(err == nil && test.expectedError != nil) || (err != nil && test.expectedError == nil) {
				t.Errorf("error %v not equal to %v", err.Error(), test.expectedError.Error())
			}
			if strings.ReplaceAll(w.Body.String(), "\n", "") != test.expectedReturn {
				t.Errorf("body error %v not equal to %v", w.Body.String(), test.expectedReturn)
			}
		})
	}
}
