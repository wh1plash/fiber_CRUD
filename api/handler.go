package api

import (
	"database/sql"
	"errors"
	"fiber/store"
	"fiber/types"
	"fmt"
	"reflect"
	"strconv"
	"strings"

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
		return ErrInvalidID()
	}
	var params types.UpdateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}

	v := reflect.ValueOf(params)
	t := reflect.TypeOf(params)
	querySet := make(map[string]any)
	for i := range v.NumField() {
		//fieldName := t.Field(i).Name
		jsonTag := t.Field(i).Tag.Get("db")
		fieldValue := v.Field(i).Interface()
		//fmt.Printf("Field: %s (json: %s), Value: %v\n", fieldName, jsonTag, fieldValue)

		key := strings.Split(jsonTag, ",")[0]
		if value, ok := fieldValue.(string); ok && value != "" {
			querySet[key] = value
		}

	}

	res, err := h.userStore.UpdateUser(c.Context(), id, querySet)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound(id, "User")
		}
		return err
	}
	return c.JSON(res)

}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	par := c.Params("id")
	id, err := strconv.Atoi(par)
	if err != nil {
		return ErrInvalidID()
	}

	deletedID, err := h.userStore.DeleteUser(c.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound(id, "User")
		}
		return err
	}
	return c.JSON(map[string]string{"deleted": fmt.Sprintf("user with id %d", deletedID)})
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecords("Users")
		}
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleGetUserByID(c *fiber.Ctx) error {
	par := c.Params("id")
	id, err := strconv.Atoi(par)
	if err != nil {
		return ErrInvalidID()
	}
	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound(id, "User")
		}
		return err
	}

	return c.JSON(user)
}
