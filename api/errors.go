package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if ApiError, ok := err.(Error); ok {
		return c.Status(ApiError.Code).JSON(ApiError)
	}

	ApiError := NewError(err.(*fiber.Error).Code, err.Error())
	curTime := time.Now()
	fmt.Printf("%s Request failed with code %d and message: %s\n", &curTime, ApiError.Code, ApiError.Message)
	return c.Status(ApiError.Code).JSON(ApiError)

}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

// Error implements the Error interface
func (e Error) Error() string {
	return e.Message
}

func NewError(code int, err string) Error {
	return Error{
		Code:    code,
		Message: err,
	}
}

func ErrBadRequest() Error {
	return Error{
		Code:    http.StatusBadRequest,
		Message: "invalid JSON request",
	}
}

func ErrInvalidID() Error {
	return Error{
		Code:    http.StatusBadRequest,
		Message: "invalid id given",
	}
}

func ErrUnAuthorized() Error {
	return Error{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized",
	}
}

func ErrNotFound(id int, resource string) Error {
	return Error{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s with %d not found", resource, id),
	}
}

func ErrNoRecords(resource string) Error {
	return Error{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}
