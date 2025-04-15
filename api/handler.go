package api

import (
	"database/sql"
	"errors"
	"fiber/store"
	"fiber/types"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userStore store.UserStore
}

func NewUserHandler(userStore store.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(insertedUser)

}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	par := c.Params("id")
	id, err := strconv.Atoi(par)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}
	var params types.UpdateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	res, err := h.userStore.UpdateUser(c.Context(), id, params.FirstName, params.LastName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("user with id %d not found", id),
			})
		}
		return err
	}
	return c.JSON(res)

}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	// var params types.DeleteUserParams
	// if err := c.BodyParser(&params); err != nil {
	// 	return ErrBadRequest()
	// }
	// if errors := params.Validate(); len(errors) > 0 {
	// 	return c.JSON(errors)
	// }

	par := c.Params("id")
	id, err := strconv.Atoi(par)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	deletedID, err := h.userStore.DeleteUser(c.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("user with id %d not found", id),
			})
		}
		return err
	}
	return c.JSON(map[string]string{"deleted": fmt.Sprintf("user with id %d", deletedID)})
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleGetUserByID(c *fiber.Ctx) error {
	// var params types.GetUserParams
	// if err := c.BodyParser(&params); err != nil {
	// 	return ErrBadRequest()
	// }
	// if errors := params.Validate(); len(errors) > 0 {
	// 	return c.JSON(errors)
	// }
	par := c.Params("id")
	id, err := strconv.Atoi(par)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}
	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("user with id %d not found", id),
			})
		}
		return err
	}

	return c.JSON(user)
}
