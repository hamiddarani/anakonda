package routes

import (
	"github.com/hamiddarani/anakonda/internal/api/handler"
	"github.com/hamiddarani/anakonda/internal/config"
	"github.com/labstack/echo/v4"
)

func HealthRoutes(r *echo.Group, cfg *config.Config) {
	h := handler.NewHealthHandler(cfg)

	r.GET("", h.Health)
}
