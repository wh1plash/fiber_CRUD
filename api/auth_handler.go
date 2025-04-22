package api

import (
	"database/sql"
	"errors"
	"fiber/types"

	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) HandleLogging(c *fiber.Ctx) error {
	var params types.UserLoginRequest
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	if errors := params.Validate(); len(errors) > 0 {
		return NewValidationError(errors)
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound(params.Email, "User")
		}
		return err
	}
	return c.JSON(user)
}
