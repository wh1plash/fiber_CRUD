package main

import (
	"fiber/api"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

const listenAddr = "localhost:3000"

func init() {
	mustLoadEnvVariables()
}

func main() {
	s := api.NewServer(listenAddr)

	go s.Run()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	<-sigch
	log.Println("Received shutdown signal, shutting down server...")
	s.Stop()
}

func mustLoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}
