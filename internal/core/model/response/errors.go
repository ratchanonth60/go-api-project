package response

import "github.com/gofiber/fiber/v2"

var (
	ErrAuth = &ErrorResponse{
		Code: fiber.ErrUnauthorized.Code,
		Msg:  fiber.ErrUnauthorized.Error(),
	}
	ErrNotFound = &ErrorResponse{
		Code: fiber.ErrNotFound.Code,
		Msg:  fiber.ErrNotFound.Error(),
	}
	ErrParser = &ErrorResponse{
		Code: fiber.ErrBadRequest.Code,
		Msg:  fiber.ErrBadRequest.Error(),
	}
	ErrCompareHashAndPassword = &ErrorResponse{
		Code: fiber.ErrForbidden.Code,
		Msg:  fiber.ErrForbidden.Error(),
	}
	ErrCofirmPassword = &ErrorResponse{
		Code: fiber.ErrConflict.Code,
		Msg:  "password and password confirm not match",
	}
)

type WrapError struct {
	error
	errorResponse *ErrorResponse
}

func (e *WrapError) Is(err error) bool {
	return e.errorResponse == err
}

func (e *WrapError) Unwrap() error {
	return e.error
}
