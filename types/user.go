package types

import (
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	minFirstNameLen = 3
	minLastNameLen  = 3
	minPasswordLen  = 4
	bcryptCost      = 12
)

type User struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Email             string    `json:"email"`
	EncryptedPassword string    `json:"-"`
	IsAdmin           bool      `json:"isAdmin"`
	CreatedAt         time.Time `json:"createdAt"`
}

type GetUserParams struct {
	ID int `json:"id"`
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type DeleteUserParams struct {
	ID int `json:"id"`
}

type UpdateUserParams struct {
	FirstName string `db:"first_name" json:"firstName,omitempty"`
	LastName  string `db:"last_name" json:"lastName,omitempty"`
	Email     string `json:"email,omitempty" db:"email"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (params UserLoginRequest) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.Email) == 0 {
		errors["email"] = "email is required"
	}
	if len(params.Password) == 0 {
		errors["password"] = "password is required"
	}
	return errors
}

func (params UpdateUserParams) Validate() map[string]string {
	errors := map[string]string{}

	if params.FirstName != "" && len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen)
	}
	if params.LastName != "" && len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen)
	}
	if params.Email != "" && !isEmailValid(params.Email) {
		errors["email"] = "invalid email present"
	}

	return errors
}

func (params GetUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if params.ID == 0 {
		errors["id"] = "id is required"
	}
	return errors
}

func (params DeleteUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if params.ID == 0 {
		errors["id"] = "id is required"
	}
	return errors
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("email %s is invalid", params.Email)
	}
	return errors
}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
