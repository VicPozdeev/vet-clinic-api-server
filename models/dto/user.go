package dto

import (
	"time"
	"vet-clinic/models"
)

// UserCreateDto defines a data transfer object for create user.
type UserCreateDto struct {
	// Username must start with an alphabetical character.
	// It can consist of ASCII alphanumeric characters and the following symbols: _.-
	Username string `json:"username" validate:"required,username" example:"Username"`
	// Password must contain at least one uppercase letter, one lowercase letter, one digit, and one special symbol.
	// It can consist of printable ASCII characters.
	Password string `json:"password" validate:"required,min=8,max=72,password" example:"8p6R*R{3"`
	RoleID   uint   `json:"roleId" validate:"required"`
}

// ToModel creates models.User from this DTO.
func (d *UserCreateDto) ToModel() *models.User {
	return models.NewUser(d.Username, d.Password, d.RoleID)
}

// UserUpdateDto defines a data transfer object for update user.
type UserUpdateDto struct {
	Email       string    `json:"email" validate:"email" example:"mail@mail.com"` // E-mail string.
	Phone       string    `json:"phone" validate:"e164" example:"+79876543210"`   // E.164 phone number string.
	Active      bool      `json:"active"`
	Surname     string    `json:"surname" validate:"omitempty,rualpha,max=255"`    // Alphabetic characters only (Russian and English).
	Name        string    `json:"name" validate:"omitempty,rualpha,max=255"`       // Alphabetic characters only (Russian and English).
	Patronymic  string    `json:"patronymic" validate:"omitempty,rualpha,max=255"` // Alphabetic characters only (Russian and English).
	Sex         string    `json:"sex" validate:"omitempty,rualpha"`                // Alphabetic characters only (Russian and English).
	BirthDate   time.Time `json:"birthDate" format:"date"`                         // Date only
	Profession  string    `json:"profession" validate:"ruprintascii"`              // Allowed characters: printable ASCII (Russian and English).
	Info        string    `json:"info" validate:"ruprintascii"`                    // Allowed characters: printable ASCII (Russian and English).
	RoleID      uint      `json:"roleId"`
	Departments []uint    `json:"departments"`
	Services    []uint    `json:"services"`
}

// ToModel creates models.User from this DTO.
// If the owner value is true, then fields are included that can only be changed by a user with the owner role.
func (d *UserUpdateDto) ToModel(owner bool) *models.User {
	if !owner {
		return &models.User{
			Email:      d.Email,
			Phone:      d.Phone,
			Surname:    d.Surname,
			Name:       d.Name,
			Patronymic: d.Patronymic,
			Sex:        d.Sex,
			BirthDate:  d.BirthDate,
		}
	}

	var departments []*models.Department
	for _, id := range d.Departments {
		department := &models.Department{BaseModel: &models.BaseModel{ID: id}}
		departments = append(departments, department)
	}
	var services []*models.Service
	for _, id := range d.Services {
		service := &models.Service{BaseModel: &models.BaseModel{ID: id}}
		services = append(services, service)
	}
	return &models.User{
		Email:       d.Email,
		Phone:       d.Phone,
		Active:      d.Active,
		Surname:     d.Surname,
		Name:        d.Name,
		Patronymic:  d.Patronymic,
		Sex:         d.Sex,
		BirthDate:   d.BirthDate,
		Profession:  d.Profession,
		Info:        d.Info,
		RoleID:      d.RoleID,
		Departments: departments,
		Services:    services,
	}
}

// UpdatePasswordDto defines the data transfer object to update the password.
//
// @Description The 'NewPassword' and 'OldPassword' should not match, while the 'ConfirmPassword'
// @Description is required for verification but must match the new password.
type UpdatePasswordDto struct {
	// Password must contain at least one uppercase letter, one lowercase letter, one digit, and one special symbol.
	// It can consist of printable ASCII characters.
	OldPassword string `json:"oldPassword" validate:"required,min=8,max=72,password" example:"8p6R*R{3"`
	// Password must contain at least one uppercase letter, one lowercase letter, one digit, and one special symbol.
	// It can consist of printable ASCII characters.
	NewPassword string `json:"newPassword" validate:"required,nefield=OldPassword,eqfield=ConfirmPassword,min=8,max=72,password" example:"6@5JG5hG"`
	// Password must contain at least one uppercase letter, one lowercase letter, one digit, and one special symbol.
	// It can consist of printable ASCII characters.
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8,max=72,password" example:"6@5JG5hG"`
}

// LoginDto defines a data transfer object for user.
type LoginDto struct {
	// Login using username, e-mail, or phone.
	Login string `json:"login" validate:"required,username|email|e164" example:"Login"`
	// Password must contain at least one uppercase letter, one lowercase letter, one digit, and one special symbol.
	// It can consist of printable ASCII characters.
	Password string `json:"password" validate:"required,min=8,max=72,password" example:"8p6R*R{3"`
}
