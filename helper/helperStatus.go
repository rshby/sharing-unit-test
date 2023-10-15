package helper

import "net/http"

func GenerateStatusFromCode(statusCode int) string {
	switch statusCode {
	case http.StatusOK:
		return "ok"
	case http.StatusNotFound:
		return "not found"
	case http.StatusBadRequest:
		return "bad request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusMethodNotAllowed:
		return "method not allowed"
	default:
		return "internal server error"
	}
}
