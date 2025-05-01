package middleware

import (
	"fiber/api"
	"fiber/store"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(h fiber.Handler, userStore store.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return api.ErrUnAuthorized("unauthorized")
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			return api.ErrUnAuthorized("unauthorized")
		}

		claims, err := validateToken(token)
		if err != nil {
			return err
		}

		// Check token expiration
		if float64(time.Now().Unix()) > claims["expires"].(float64) {
			return api.ErrUnAuthorized("token is expired")
		}

		userID := claims["id"].(float64)
		user, err := userStore.GetUserByID(c.Context(), int(userID))
		if err != nil {
			return api.ErrUnAuthorized("unauthorized")
		}
		_ = user
		// Set the current authenticated user to the context.
		//c.Context().SetUserValue("user", user)

		err = h(c)
		if err != nil {
			return err
		}
		return nil
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, api.ErrUnAuthorized("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		return nil, api.ErrUnAuthorized("unauthorized")
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return nil, api.ErrUnAuthorized("unauthorized")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, api.ErrUnAuthorized("unauthorized")
	}
	return claims, nil
}
