package api

import (
	"database/sql"
	"fiber/store"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

type testdb struct {
	store.UserStore
}

type PGStore struct {
	db *sql.DB
}

func PGStoreTest() (*PGStore, error) {
	port, _ := strconv.Atoi(os.Getenv("PG_PORT"))
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", os.Getenv("PG_HOST"), port, os.Getenv("PG_USER"), os.Getenv("PG_PASS"), os.Getenv("PG_DB_NAME"))
	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
		//log.Fatal("error to connect to Posgres database", err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PGStore{
		db: db,
	}, nil
}

func TestHandleGetUserByID(t *testing.T) {
	mustLoadEnvVariables()
	b, _ := PGStoreTest()
	fmt.Println(b)

	// s := fiber.New(fiber.Config{
	// 	ErrorHandler: ErrorHandler,
	// })
	// resp, err := http.Get()

}

func mustLoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}
