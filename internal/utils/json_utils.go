package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Add("Content-Type", "application/json")

	jsonValue, err := json.Marshal(value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(status)
		w.Write(jsonValue)
	}
}
