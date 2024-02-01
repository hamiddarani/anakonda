package handler

import (
	"net/http"

	"github.com/hamiddarani/anakonda/internal/config"
	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	cfg *config.Config
}

func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		cfg: cfg,
	}
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept Application/json
// @Produce json
// @Success 200 {object} helper.BaseHttpResponse "Success"
// @Router /v1/health [get]
func (h *HealthHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"app": h.cfg.App.Name,
	})
}
