package migration

import (
	"vet-clinic/container"
	"vet-clinic/models"
)

// CreateTables creates the tables used in this application.
func CreateTables(container container.Container) {
	if container.Config().Database.Migration {
		rep := container.Repository()

		_ = rep.DropTableIfExists(&models.User{})
		_ = rep.DropTableIfExists(&models.Role{})
		_ = rep.DropTableIfExists(&models.Category{})
		_ = rep.DropTableIfExists(&models.Service{})
		_ = rep.DropTableIfExists(&models.Department{})
		_ = rep.DropTableIfExists(&models.Client{})
		_ = rep.DropTableIfExists(&models.Pet{})
		_ = rep.DropTableIfExists(&models.Visit{})
		_ = rep.DropTableIfExists(&models.Lead{})
		_ = rep.DropTableIfExists("users_departments")
		_ = rep.DropTableIfExists("users_services")
		_ = rep.DropTableIfExists("departments_services")

		_ = rep.AutoMigrate(&models.User{})
		_ = rep.AutoMigrate(&models.Role{})
		_ = rep.AutoMigrate(&models.Category{})
		_ = rep.AutoMigrate(&models.Service{})
		_ = rep.AutoMigrate(&models.Department{})
		_ = rep.AutoMigrate(&models.Client{})
		_ = rep.AutoMigrate(&models.Pet{})
		_ = rep.AutoMigrate(&models.Visit{})
		_ = rep.AutoMigrate(&models.Lead{})
	}
}
