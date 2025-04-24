package api

import (
	"database/sql"
	"errors"
	"fiber/store"
	"fiber/types"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

func (p AuthParams) Validate() map[string]string {
	errors := map[string]string{}
	if p.Email == "" {
		errors["Email"] = "email is required"
	}
	if p.Password == "" {
		errors["Password"] = "password is required"
	}
	return errors
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
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

	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return ErrInvalidCredentials()
	}

	token, err := CreateTokenFromUser(user)
	if err != nil {
		return err
	}

	resp := AuthResponse{
		User:  user,
		Token: token,
	}

	//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImJhc3JAZm9vLmNvbSIsImV4cGlyZXMiOjE3NDU1NTk1MjMsImlkIjoyfQ.rmDAQCd8iWME7nGIOCJJW64mKPWmyHoGZEFfZKa8gQw
	return c.JSON(resp)
}

func CreateTokenFromUser(u *types.User) (string, error) {
	now := time.Now()
	expires := now.Add(time.Minute * 1).Unix()
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
