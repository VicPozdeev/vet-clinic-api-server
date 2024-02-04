package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"vet-clinic/container"
	"vet-clinic/models/dto"
	"vet-clinic/service"
	"vet-clinic/util"
)

type ServiceController struct {
	container container.Container
	service   *service.ServiceService
}

// NewServiceController is constructor.
func NewServiceController(container container.Container) *ServiceController {
	return &ServiceController{container: container, service: service.NewServiceService(container)}
}

// Get returns one record matched service's id.
//
// @Summary Get a service.
// @Description Returns one record matched service's id.
// @Tags Services
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Service ID"
// @Success 200 {object} models.Service "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to fetch data."
// @Failure 401 "Failed to the authentication."
// @Router /services/{id} [get]
func (r *ServiceController) Get(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}

	serv, err := r.service.Get(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, serv)
}

// GetAll returns the list of services.
//
// @Summary Get a service list.
// @Description Returns the list of services.
// @Tags Services
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []models.Service "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to fetch data."
// @Failure 401 "Failed to the authentication."
// @Router /services [get]
func (r *ServiceController) GetAll(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}

	serv, err := r.service.GetAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, serv)
}

// Create creates a new service.
//
// @Summary Create a new service. Required user's role: Owner
// @Description Create a new service.
// @Tags Services
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.ServiceDto true "A new service data for creating."
// @Success 200 {object} models.Service "Success to fetch data."
// @Failure 400 {object} dto.ServiceDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /services [post]
func (r *ServiceController) Create(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Owner.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	data := &dto.ServiceDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	serv, err := r.service.Create(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, serv)
}

// Update updates the existing service.
//
// @Summary Update the existing service. Required user's role: Owner
// @Description Update the existing service.
// @Tags Services
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Service ID"
// @Param data body dto.ServiceDto true "Service data for update."
// @Success 200 {object} models.Service "Success to fetch data."
// @Failure 400 {object} dto.ServiceDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /services/{id} [put]
func (r *ServiceController) Update(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Owner.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	data := &dto.ServiceDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	serv, err := r.service.Update(data, c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, serv)
}

// Delete deletes the existing service.
//
// @Summary Delete the existing service. Required user's role: Owner
// @Description Delete the existing service.
// @Tags Services
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Service ID"
// @Success 200 {object} models.Service "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /services/{id} [delete]
func (r *ServiceController) Delete(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Owner.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	serv, err := r.service.Delete(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, serv)
}
