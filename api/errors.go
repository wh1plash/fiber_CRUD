package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if ApiError, ok := err.(Error); ok {
		return c.Status(ApiError.Code).JSON(ApiError)
	}
	ApiError := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(ApiError.Code).JSON(ApiError)
}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

// Error implements the Error interface
func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid JSON request",
	}
}

func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid id given",
	}
}

func ErrUnAuthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "unauthorized",
	}
}

func ErrNotFound(id int, resource string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  fmt.Sprintf("%s with %d not found", resource, id),
	}
}

func ErrNoRecords(resource string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  fmt.Sprintf("%s not found", resource),
	}
}
