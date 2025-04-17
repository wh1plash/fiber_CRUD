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
		app         = fiber.New(config)
		apiv1       = app.Group("/api/v1")
		userHandler = NewUserHandler(db)
	)
	apiv1.Get("/home", s.HomeHandler)

	LoggedRoute(apiv1, "POST", "/user", userHandler.HandlePostUser)
	LoggedRoute(apiv1, "PUT", "/user/:id", userHandler.HandlePutUser)
	LoggedRoute(apiv1, "DELETE", "/user/:id", userHandler.HandleDeleteUser)
	LoggedRoute(apiv1, "GET", "/users", userHandler.HandleGetUsers)
	LoggedRoute(apiv1, "GET", "/user/:id", userHandler.HandleGetUserByID)

	err = app.Listen(s.listenAddr)
	if err != nil {
		s.logger.Error("error to start server", "error", err.Error())
		return
	}
}

func LoggedRoute(r fiber.Router, method, path string, handler fiber.Handler) {
	wrapped := LoggingHandlerDecorator(handler)
	switch method {
	case "GET":
		r.Get(path, wrapped)
	case "POST":
		r.Post(path, wrapped)
	case "PUT":
		r.Put(path, wrapped)
	case "DELETE":
		r.Delete(path, wrapped)

	}
}

func (s *Server) HomeHandler(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "All works fine"})

}
