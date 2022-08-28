package errs

import "net/http"

type AppErr struct {
	Message string `json:"message"`
	Code    int    `json:",omitempty"`
}

func (aErr *AppErr) GetMessage() *AppErr {
	return &AppErr{
		Message: aErr.Message,
	}
}

func NotFoundError(message string) *AppErr {
	return &AppErr{
		Message: message,
		Code:    http.StatusNotFound,
	}
}

func InternalServerError(message string) *AppErr {
	return &AppErr{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}
