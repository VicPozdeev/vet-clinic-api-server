package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"vet-clinic/container"
	"vet-clinic/models/dto"
	"vet-clinic/service"
	"vet-clinic/util"
)

type VisitController struct {
	container container.Container
	service   *service.VisitService
}

// NewVisitController is constructor.
func NewVisitController(container container.Container) *VisitController {
	return &VisitController{container: container, service: service.NewVisitService(container)}
}

// Get returns one record matched visit's id.
//
// @Summary Get a visit.
// @Description Returns one record matched visit's id.
// @Tags Visits
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Visit ID"
// @Success 200 {object} models.Visit "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to fetch data."
// @Failure 401 "Failed to the authentication."
// @Router /visits/{id} [get]
func (r *VisitController) Get(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}

	visit, err := r.service.Get(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, visit)
}

// GetAll returns the list of visits.
//
// @Summary Get a visit list.
// @Description Returns the list of visits.
// @Tags Visits
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []models.Visit "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to fetch data."
// @Failure 401 "Failed to the authentication."
// @Router /visits [get]
func (r *VisitController) GetAll(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}

	visits, err := r.service.GetAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, visits)
}

// Create creates a new visit.
//
// @Summary Create a new visit. Required user's role: Admin
// @Description Create a new visit.
// @Tags Visits
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.VisitDto true "A new visit data for creating."
// @Success 200 {object} models.Visit "Success to fetch data."
// @Failure 400 {object} dto.VisitDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /visits [post]
func (r *VisitController) Create(c echo.Context) error {
	user := getUser(c, r.container)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Administrator.AccessAllowed(util.ToAccessLevel(user.Role.Name)) {
		return c.NoContent(http.StatusForbidden)
	}

	data := &dto.VisitDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	data.LastUpdatedByID = user.ID
	visit, err := r.service.Create(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, visit)
}

// Update updates the existing visit.
//
// @Summary Update the existing visit. Required user's role: Admin
// @Description Update the existing visit.
// @Tags Visits
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Visit ID"
// @Param data body dto.VisitDto true "Visit data for update."
// @Success 200 {object} models.Visit "Success to fetch data."
// @Failure 400 {object} dto.VisitDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /visits/{id} [put]
func (r *VisitController) Update(c echo.Context) error {
	user := getUser(c, r.container)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Administrator.AccessAllowed(util.ToAccessLevel(user.Role.Name)) {
		return c.NoContent(http.StatusForbidden)
	}

	data := &dto.VisitDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	data.LastUpdatedByID = user.ID
	visit, err := r.service.Update(data, c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, visit)
}

// Delete deletes the existing visit.
//
// @Summary Delete the existing visit. Required user's role: Superuser
// @Description Delete the existing visit.
// @Tags Visits
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Visit ID"
// @Success 200 {object} models.Visit "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /visits/{id} [delete]
func (r *VisitController) Delete(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Superuser.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	visit, err := r.service.Delete(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, visit)
}
