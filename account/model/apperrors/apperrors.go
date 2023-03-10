package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

type Type string

const (
	Authorization        Type = "AUTHORIZATION"
	BadRequest           Type = "BAD_REQUEST"
	Conflict             Type = "CONFLICT"
	Internal             Type = "INTERNAL"
	NotFound             Type = "NOT_FOUND"
	UnsupportedMediaType Type = "UNSUPPORTED_MEDIA_TYPE" // for http 415
	PayloadTooLarge      Type = "PAYLOAD_TOO_LARGE"
	SeriveUnavailable    Type = "SERVICE_UNAVAILABLE"
)

type Error struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message

}

func (e *Error) StatusCode() int {
	switch e.Type {
	case Authorization:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusConflict
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case PayloadTooLarge:
		return http.StatusRequestEntityTooLarge
	case SeriveUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError

	}
}

func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.StatusCode()
	}
	return http.StatusInternalServerError
}

func NewAuthorization(reason string) *Error {
	return &Error{
		Type:    Authorization,
		Message: reason,
	}
}

func NewBadRequest(reason string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: fmt.Sprintf("Bad request. Reason %v", reason),
	}
}

func NewConflict(name, value string) *Error {
	return &Error{
		Type:    Conflict,
		Message: fmt.Sprintf("resource: %v with value: %v already exists", name, value),
	}
}

func NewInternalServerError() *Error {
	return &Error{
		Type:    Internal,
		Message: fmt.Sprintf("Internal server eror"),
	}
}

func NewNotFound(name, value string) *Error {
	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("resource: %v with value: %v, not found", name, value),
	}
}

func NewUnsupportedMediaType(reason string) *Error {
	return &Error{
		Type:    UnsupportedMediaType,
		Message: reason,
	}
}

func NewPayloadTooLarge(maxBodySize, contentLength int64) *Error {
	return &Error{
		Type:    PayloadTooLarge,
		Message: fmt.Sprintf("Max payload size of %v exceeded. Actual payload size: %v", maxBodySize, contentLength),
	}
}

func NewServiceUnavailable() *Error {
	return &Error{
		Type:    SeriveUnavailable,
		Message: fmt.Sprintf("Service unavailalbe"),
	}

}
