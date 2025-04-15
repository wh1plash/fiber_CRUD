package api

import (
	"fiber/store"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

var config = fiber.Config{
	ErrorHandler: ErrorHandler,
}

type Server struct {
	listenAddr string
	logger     *slog.Logger
}

func NewServer(addr string) *Server {
	return &Server{
		listenAddr: addr,
		logger:     slog.Default(),
	}
}

func (s *Server) Stop() {
	s.logger.Info("server stopped")
}

func (s *Server) Run() {
	db, err := store.NewPostgresStore()
	if err != nil {
		s.logger.Error("error to connect to Posgres database", "error", err.Error())
		return
	}

	if err := db.Init(); err != nil {
		s.logger.Error("error to create tables", "error", err.Error())
		return
	}

	var (
		app         = fiber.New()
		apiv1       = app.Group("/api/v1")
		userHandler = NewUserHandler(db)
	)

	apiv1.Get("/home", s.HomeHandler)

	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUserByID)

	err = app.Listen(s.listenAddr)
	if err != nil {
		s.logger.Error("error to start server", "error", err.Error())
		return
	}

}

func (s *Server) HomeHandler(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "All works fine"})

}
