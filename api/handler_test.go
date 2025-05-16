package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fiber/store"
	"fiber/types"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type testdb struct {
	store.UserStore
}

func (tdb *testdb) SeedUsers(t *testing.T) {
	for i := range 2 {
		params := types.CreateUserParams{
			FirstName: fmt.Sprintf("FName_%d", i),
			LastName:  "foi",
			Email:     "some@mail.com",
			Password:  "qwerty",
		}
		user, err := types.NewUserFromParams(params)
		if err != nil {
			t.Fatal(err)
		}
		_, err = tdb.UserStore.InsertUser(context.TODO(), user)
		if err != nil {
			t.Fatal(err)
		}

	}
}

func (tdb *testdb) teardown(t *testing.T) {
	err := tdb.UserStore.DropTable("users")
	if err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	connStr := "host=postgres port=5432 user=postgres password=postgres dbname=test sslmode=disable"
	db, err := store.NewPostgresStore(connStr)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Init(); err != nil {
		t.Fatal("error to create table", err)
	}
	return &testdb{
		UserStore: db,
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})
	userHandler := NewUserHandler(tdb)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "Test1",
		LastName:  "foi",
		Email:     "some@mail.com",
		Password:  "qwerty",
	}
	if errors := params.Validate(); len(errors) > 0 {
		t.Fatal("validation fail")
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if user.ID == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting the EncryptedPassword not to be included in the json response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected last name %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status code %d but got %d", http.StatusOK, resp.StatusCode)
	}

}

func TestHandleGetUsers(t *testing.T) {
	tdb := setup(t)
	tdb.SeedUsers(t)
	defer tdb.teardown(t)

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	userHandler := NewUserHandler(tdb)
	app.Get("/user", userHandler.HandleGetUsers)
	req := httptest.NewRequest("GET", "/user", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var users []types.User
	json.NewDecoder(resp.Body).Decode(&users)

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status code %d but got %d", http.StatusOK, resp.StatusCode)
	}
	if len(users) < 2 {
		t.Errorf("expected users in database = 2 but got %d", len(users))
	}
}

func TestHandleGetUsersNoRows(t *testing.T) {
	tdb := setup(t)
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})
	UserHandler := NewUserHandler(tdb)
	app.Get("/user", UserHandler.HandleGetUsers)
	req := httptest.NewRequest("GET", "/user", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("expected response status code %d but got %d", fiber.StatusNotFound, resp.StatusCode)
	}
}

func TestHandleGetUserByID(t *testing.T) {
	tdb := setup(t)
	tdb.SeedUsers(t)
	defer tdb.teardown(t)

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})
	userHandler := NewUserHandler(tdb)
	app.Get("/user/:id", userHandler.HandleGetUserByID)

	req := httptest.NewRequest("GET", "/user/1", nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status code %d but got %d", http.StatusOK, resp.StatusCode)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if user.ID == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting the EncryptedPassword not to be included in the json response")
	}
	if len(user.FirstName) == 0 {
		t.Errorf("expected to receive user FirstName len > 0 but got %s", user.FirstName)
	}
	if len(user.LastName) == 0 {
		t.Errorf("expected to receive user LastName len > 0 but got %s", user.LastName)
	}
	if len(user.Email) == 0 {
		t.Errorf("expected to receive user Email len > 0 but got %s", user.Email)
	}

}

func TestHandleGetUserByIDNotFound(t *testing.T) {
	tdb := setup(t)
	tdb.SeedUsers(t)
	defer tdb.teardown(t)

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})
	userHandler := NewUserHandler(tdb)
	app.Get("/user/:id", userHandler.HandleGetUserByID)

	req := httptest.NewRequest("GET", "/user/", nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("expected status code %d but got %d", http.StatusNotFound, resp.StatusCode)
	}
}
