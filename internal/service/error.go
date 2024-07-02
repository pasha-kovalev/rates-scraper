package service

import "net/http"

var (
	ErrBadGateway = &ApiError{
		Code:    http.StatusBadGateway,
		Message: "Bad Gateway",
	}
)

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (a ApiError) Error() string {
	return a.Message
}
