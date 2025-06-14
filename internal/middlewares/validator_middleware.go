package middlewares

import (
	"context"
	"github.com/golates/api-gateway/internal/utils"
	"net/http"
)

func ValidatorMiddleware(next http.Handler) http.Handler {
	v := utils.NewCustomValidator()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := context.WithValue(r.Context(), "validator", v)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
