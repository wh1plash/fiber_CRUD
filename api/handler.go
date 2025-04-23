package api

import (
	"database/sql"
	"errors"
	"fiber/store"
	"fiber/types"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	UserStore store.UserStore
}

func NewUserHandler(userStore store.UserStore) *UserHandler {
	return &UserHandler{
		UserStore: userStore,
	}
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return NewValidationError(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	insertedUser, err := h.UserStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(insertedUser)

}

func (h *UserHandler) HandleGetUserByID(c *fiber.Ctx) error {
	//time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	par := c.Params("id")
	id, err := strconv.Atoi(par)
	if err != nil {
		return ErrInvalidID()
	}
	user, err := h.UserStore.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound(id, "User")
		}
		return err
	}

	return c.JSON(user)
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
		return NewValidationError(errors)
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
	if len(querySet) == 0 {
		return ErrBadRequest()
	}

	res, err := h.UserStore.UpdateUser(c.Context(), id, querySet)
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

	deletedID, err := h.UserStore.DeleteUser(c.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound(id, "User")
		}
		return err
	}
	return c.JSON(map[string]string{"deleted": fmt.Sprintf("user with id %d", deletedID)})
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.UserStore.GetUsers(c.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound("Users", "no condition")
		}
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params types.AuthParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	if errors := params.Validate(); len(errors) > 0 {
		return NewValidationError(errors)
	}

	user, err := h.UserStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound(params.Email, "User")
		}
		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return ErrInvalidCredentials()
	}

	token, err := CreateTokenFromUser(user)
	if err != nil {
		return err
	}

	resp := types.AuthResponse{
		User:  user,
		Token: token,
	}

	//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImJhc3JAZm9vLmNvbSIsImV4cGlyZXMiOjE3NDU0ODY3ODQsImlkIjoyfQ.gBN5DSkrmscxUbakFdqEFozRhEzkqJwYyFH_j42UZjg
	return c.JSON(resp)
}

func CreateTokenFromUser(u *types.User) (string, error) {
	now := time.Now()
	expires := now.Add(time.Hour * 24).Unix()
	claims := jwt.MapClaims{
		"id":      u.ID,
		"email":   u.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
