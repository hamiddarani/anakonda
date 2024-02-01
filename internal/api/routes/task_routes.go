package routes

import (
	"github.com/hamiddarani/anakonda/internal/api/handler"
	"github.com/hamiddarani/anakonda/internal/config"
	"github.com/labstack/echo/v4"
)

func TaskRoutes(r *echo.Group, cfg *config.Config) {
	h := handler.NewTaskHandler(cfg)

	r.POST("", h.CreateTask)
	r.GET("/:id", h.GetById)
}
