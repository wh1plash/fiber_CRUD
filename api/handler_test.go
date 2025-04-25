package api

import (
	"bytes"
	"encoding/json"
	"fiber/store"
	"fiber/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type testdb struct {
	store.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	err := tdb.UserStore.DropTable("users")
	if err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	connStr := "host=localhost port=5444 user=postgres password=postgres dbname=test sslmode=disable"
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

	app := fiber.New()
	userHandler := NewUserHandler(tdb)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "Test1",
		LastName:  "fooo",
		Email:     "some@mail.com",
		Password:  "qwerty",
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

func TestHandleGetUserByID(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	// s := fiber.New(fiber.Config{
	// 	ErrorHandler: ErrorHandler,
	// })
	// resp, err := http.Get()

}
