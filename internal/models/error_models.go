package models

type MessageAPIResponseError struct {
	Message string `json:"message"`
}

type ErrorsArrayAPIResponseError struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}
