package middlewares

import (
	"github.com/go-chi/chi/v5"
	"github.com/golates/api-gateway/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidatorMiddleware(t *testing.T) {
	t.Run("ValidatorMiddleware test: ", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/hello", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		// Setup test router
		r := chi.NewRouter()
		r.Use(ValidatorMiddleware)
		r.Post("/hello", func(writer http.ResponseWriter, request *http.Request) {
			v := request.Context().Value("validator")
			switch v.(type) {
			case *utils.CustomValidator:
				return
			default:
				t.Errorf("Validator not found")
			}
		})

		r.ServeHTTP(rr, req)
	})
}
