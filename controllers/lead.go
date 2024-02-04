package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"vet-clinic/container"
	"vet-clinic/models/dto"
	"vet-clinic/service"
	"vet-clinic/util"
)

type LeadController struct {
	container container.Container
	service   *service.LeadService
}

// NewLeadController is constructor.
func NewLeadController(container container.Container) *LeadController {
	return &LeadController{container: container, service: service.NewLeadService(container)}
}

// Get returns one record matched lead's id.
//
// @Summary Get a lead.
// @Description Returns one record matched lead's id.
// @Tags Leads
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Lead ID"
// @Success 200 {object} models.Lead "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to fetch data."
// @Failure 401 "Failed to the authentication."
// @Router /leads/{id} [get]
func (r *LeadController) Get(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}

	lead, err := r.service.Get(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, lead)
}

// GetAll returns the list of leads.
//
// @Summary Get a lead list.
// @Description Returns the list of leads.
// @Tags Leads
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []models.Lead "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to fetch data."
// @Failure 401 "Failed to the authentication."
// @Router /leads [get]
func (r *LeadController) GetAll(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}

	leads, err := r.service.GetAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, leads)
}

// Create creates a new lead.
//
// @Summary Create a new lead.
// @Description Create a new lead.
// @Tags Leads
// @Accept json
// @Produce json
// @Param data body dto.LeadDto true "A new lead data for creating."
// @Success 200 {object} models.Lead "Success to fetch data."
// @Failure 400 {object} dto.LeadDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Router /leads [post]
func (r *LeadController) Create(c echo.Context) error {
	data := &dto.LeadDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	lead, err := r.service.Create(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, lead)
}

// Update updates the existing lead.
//
// @Summary Update the existing lead.
// @Description Update the existing lead.
// @Tags Leads
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Lead ID"
// @Param data body dto.LeadDto true "Lead data for update."
// @Success 200 {object} models.Lead "Success to fetch data."
// @Failure 400 {object} dto.LeadDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Router /leads/{id} [put]
func (r *LeadController) Update(c echo.Context) error {
	user := getUser(c, r.container)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	data := &dto.LeadDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	data.LastUpdatedByID = user.ID
	lead, err := r.service.Update(data, c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, lead)
}

// Delete deletes the existing lead.
//
// @Summary Delete the existing lead. Required user's role: Superuser
// @Description Delete the existing lead.
// @Tags Leads
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Lead ID"
// @Success 200 {object} models.Lead "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /leads/{id} [delete]
func (r *LeadController) Delete(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Superuser.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	lead, err := r.service.Delete(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, lead)
}
