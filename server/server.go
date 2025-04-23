package server

import (
	"fiber/api"
	"fiber/middleware"
	"fiber/store"
	"log/slog"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
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

func RegisterMetrics(app *fiber.App) {
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
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
		userHandler = api.NewUserHandler(db)
		promMetrics = middleware.NewPromMetrics()
	)
	RegisterMetrics(app)

	apiv1.Post("/user", WrapHandler(promMetrics, userHandler.HandlePostUser, "HandlePostUser"))
	apiv1.Put("/user/:id", WrapHandler(promMetrics, userHandler.HandlePutUser, "HandlePutUser"))
	apiv1.Delete("/user/:id", WrapHandler(promMetrics, userHandler.HandleDeleteUser, "HandleDeleteUser"))
	apiv1.Get("/user/:id", WrapHandler(promMetrics, userHandler.HandleGetUserByID, "HandleGetUserByID"))
	apiv1.Get("/auth", WrapHandler(promMetrics, userHandler.HandleAuthenticate, "HandleAuthUser"))
	//apiv1.Use(authHandler)
	apiv1.Get("/users", WrapHandler(promMetrics, userHandler.HandleGetUsers, "HandleGetUsers"))

	err = app.Listen(s.listenAddr)
	if err != nil {
		s.logger.Error("error to start server", "error", err.Error())
		return
	}
}

func WrapHandler(p *middleware.PromMetrics, handler fiber.Handler, handlerName string) fiber.Handler {
	return p.WithMetrics(middleware.LoggingHandlerDecorator(handler), handlerName)
}
