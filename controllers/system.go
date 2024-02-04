package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"vet-clinic/container"
)

type SystemController struct {
	container container.Container
}

type HealthResult struct {
	Status string `json:"status" example:"available"`
}

// NewSystemController is constructor.
func NewSystemController(container container.Container) *SystemController {
	return &SystemController{container: container}
}

// GetHealthCheck returns one record matched client's id.
//
// @Summary Get the health status.
// @Description Returns the health status of the system.
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} HealthResult "Success to fetch health status."
// @Router /health [get]
func (r *SystemController) GetHealthCheck(c echo.Context) error {
	health := &HealthResult{Status: "available"}
	return c.JSON(http.StatusOK, health)
}
