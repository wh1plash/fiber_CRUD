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
		userHandler = api.NewUserHandler(db)
		authHandler = api.NewAuthHandler(db)
		promMetrics = middleware.NewPromMetrics()
		auth        = app.Group("/api")
		apiv1       = app.Group("/api/v1")
	)
	RegisterMetrics(app)

	auth.Post("/auth", WrapHandler(promMetrics, authHandler.HandleAuthenticate, "HandleAuthenticate"))

	apiv1.Post("/user", WrapHandler(promMetrics, WithAuth(userHandler.HandlePostUser, db), "HandlePostUser"))
	apiv1.Put("/user/:id", WrapHandler(promMetrics, WithAuth(userHandler.HandlePutUser, db), "HandlePutUser"))
	apiv1.Delete("/user/:id", WrapHandler(promMetrics, WithAuth(userHandler.HandleDeleteUser, db), "HandleDeleteUser"))
	apiv1.Get("/user/:id", WrapHandler(promMetrics, WithAuth(userHandler.HandleGetUserByID, db), "HandleGetUserByID"))

	//apiv1.Get("/auth", WrapHandler(promMetrics, authHandler.HandleAuthenticate, "HandleAuthenticate"))
	//apiv1.Use(authHandler)
	apiv1.Get("/users", WrapHandler(promMetrics, WithAuth(userHandler.HandleGetUsers, db), "HandleGetUsers"))

	err = app.Listen(s.listenAddr)
	if err != nil {
		s.logger.Error("error to start server", "error", err.Error())
		return
	}
}

func WithAuth(handler fiber.Handler, db store.UserStore) fiber.Handler {
	return middleware.JWTAuthentication(handler, db)
}

func WithLogging(handler fiber.Handler) fiber.Handler {
	return middleware.LoggingHandlerDecorator(handler)
}

func WrapHandler(p *middleware.PromMetrics, handler fiber.Handler, handlerName string) fiber.Handler {
	return p.WithMetrics(WithLogging(handler), handlerName)
}
