package api

import (
	"fmt"

	"github.com/hamiddarani/anakonda/docs"
	"github.com/hamiddarani/anakonda/internal/api/middlewares"
	"github.com/hamiddarani/anakonda/internal/api/routes"
	"github.com/hamiddarani/anakonda/internal/config"
	"github.com/hamiddarani/anakonda/pkg/logging"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	cfg    *config.Config
	app    *echo.Echo
	logger logging.Logger
}

func New(cfg *config.Config, lg logging.Logger) *Server {
	e := &Server{cfg: cfg, logger: lg}

	e.app = echo.New()
	e.app.Use(middleware.Recover())
	e.app.Use(middlewares.StructuredLogger(cfg, lg))

	registerRoutes(e.app, cfg)
	registerSwagger(e.app, cfg)

	return e
}

func registerRoutes(e *echo.Echo, cfg *config.Config) {
	api := e.Group("/api")
	v1 := api.Group("/v1")
	{
		health := v1.Group("/health")
		routes.HealthRoutes(health, cfg)

		task := v1.Group("/tasks")
		routes.TaskRoutes(task, cfg)
	}
}

func (server *Server) Serve() error {
	addr := fmt.Sprintf(":%d", server.cfg.App.Port)

	if err := server.app.Start(addr); err != nil {
		server.logger.Fatal(logging.Internal, logging.Startup, err.Error(), nil)
		return err
	}

	return nil
}

func registerSwagger(e *echo.Echo, cfg *config.Config) {
	docs.SwaggerInfo.Title = "anakonda"
	docs.SwaggerInfo.Description = "anakonda"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", "8080")
	docs.SwaggerInfo.Schemes = []string{"http"}

	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
