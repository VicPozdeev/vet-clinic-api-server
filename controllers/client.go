package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"vet-clinic/container"
	"vet-clinic/models/dto"
	"vet-clinic/service"
	"vet-clinic/util"
)

type ClientController struct {
	container container.Container
	service   *service.ClientService
}

// NewClientController is constructor.
func NewClientController(container container.Container) *ClientController {
	return &ClientController{container: container, service: service.NewClientService(container)}
}

// Get returns one record matched client's id.
//
// @Summary Get a client.
// @Description Returns one record matched client's id.
// @Tags Clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Client ID"
// @Success 200 {object} models.Client "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to fetch data."
// @Failure 401 "Failed to the authentication."
// @Router /clients/{id} [get]
func (r *ClientController) Get(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}

	client, err := r.service.Get(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, client)
}

// GetAll returns the list of clients.
//
// @Summary Get a client list.
// @Description Returns the list of clients.
// @Tags Clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []models.Client "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to fetch data."
// @Failure 401 "Failed to the authentication."
// @Router /clients [get]
func (r *ClientController) GetAll(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}

	clients, err := r.service.GetAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, clients)
}

// Create creates a new client.
//
// @Summary Create a new client. Required user's role: Admin
// @Description Create a new client.
// @Tags Clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.ClientDto true "A new client data for creating."
// @Success 200 {object} models.Client "Success to fetch data."
// @Failure 400 {object} dto.ClientDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /clients [post]
func (r *ClientController) Create(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Administrator.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	data := &dto.ClientDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	client, err := r.service.Create(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, client)
}

// Update updates the existing client.
//
// @Summary Update the existing client. Required user's role: Admin
// @Description Update the existing client.
// @Tags Clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Client ID"
// @Param data body dto.ClientDto true "Client data for update."
// @Success 200 {object} models.Client "Success to fetch data."
// @Failure 400 {object} dto.ClientDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /clients/{id} [put]
func (r *ClientController) Update(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Administrator.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	data := &dto.ClientDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	client, err := r.service.Update(data, c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, client)
}

// Delete deletes the existing client.
//
// @Summary Delete the existing client. Required user's role: Superuser
// @Description Delete the existing client.
// @Tags Clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Client ID"
// @Success 200 {object} models.Client "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /clients/{id} [delete]
func (r *ClientController) Delete(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Superuser.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	client, err := r.service.Delete(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, client)
}
