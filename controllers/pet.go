package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"vet-clinic/container"
	"vet-clinic/models/dto"
	"vet-clinic/service"
	"vet-clinic/util"
)

type PetController struct {
	container container.Container
	service   *service.PetService
}

// NewPetController is constructor.
func NewPetController(container container.Container) *PetController {
	return &PetController{container: container, service: service.NewPetService(container)}
}

// Get returns one record matched pet's id.
//
// @Summary Get a pet.
// @Description Returns one record matched pet's id.
// @Tags Pets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Pet ID"
// @Success 200 {object} models.Pet "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to fetch data."
// @Failure 401 "Failed to the authentication."
// @Router /pets/{id} [get]
func (r *PetController) Get(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}

	pet, err := r.service.Get(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, pet)
}

// GetAll returns the list of pets.
//
// @Summary Get a pet list.
// @Description Returns the list of pets.
// @Tags Pets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []models.Pet "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to fetch data."
// @Failure 401 "Failed to the authentication."
// @Router /pets [get]
func (r *PetController) GetAll(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}

	pets, err := r.service.GetAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, pets)
}

// Create creates a new pet.
//
// @Summary Create a new pet.
// @Description Create a new pet.
// @Tags Pets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.PetDto true "A new pet data for creating."
// @Success 200 {object} models.Pet "Success to fetch data."
// @Failure 400 {object} dto.PetDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Router /pets [post]
func (r *PetController) Create(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}

	data := &dto.PetDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	pet, err := r.service.Create(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, pet)
}

// Update updates the existing pet.
//
// @Summary Update the existing pet.
// @Description Update the existing pet.
// @Tags Pets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Pet ID"
// @Param data body dto.PetDto true "Pet data for update."
// @Success 200 {object} models.Pet "Success to fetch data."
// @Failure 400 {object} dto.PetDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Router /pets/{id} [put]
func (r *PetController) Update(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}

	data := &dto.PetDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	pet, err := r.service.Update(data, c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, pet)
}

// Delete deletes the existing pet.
//
// @Summary Delete the existing pet. Required user's role: Superuser
// @Description Delete the existing pet.
// @Tags Pets
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Pet ID"
// @Success 200 {object} models.Pet "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /pets/{id} [delete]
func (r *PetController) Delete(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Superuser.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	pet, err := r.service.Delete(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, pet)
}
