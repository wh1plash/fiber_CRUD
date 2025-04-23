package middleware

import (
	"database/sql"
	"errors"
	"fiber/api"
	"fiber/store"
	"fiber/types"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userStore store.UserStore
}

func NewAuthHandler(userStore store.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleLoginUser(c *fiber.Ctx) error {
	var params types.UserLoginRequest
	if err := c.BodyParser(&params); err != nil {
		return api.ErrBadRequest()
	}

	if errors := params.Validate(); len(errors) > 0 {
		return api.NewValidationError(errors)
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return api.ErrNotFound(params.Email, "User")
		}
		return err
	}
	return c.JSON(user)
}
