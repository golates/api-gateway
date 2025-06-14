package utils

import (
	"google.golang.org/grpc/status"
	"net/http"
)

func ParseGRPCError(err error) (int, string) {
	st, ok := status.FromError(err)

	if !ok {
		return http.StatusInternalServerError, "Internal server error"
	}

	return http.StatusBadRequest, st.Message()
}
