package models

import (
	"fmt"
	"github.com/gosimple/slug"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
	"vet-clinic/config"
	"vet-clinic/repository"
)

// User defines struct of user data.
type User struct {
	*BaseModel
	Username    string        `json:"username" gorm:"unique;not null;size:255"`
	Email       string        `json:"email" gorm:"unique;size:255"`
	Phone       string        `json:"phone" gorm:"unique;size:255"`
	Active      bool          `json:"active"`
	Password    string        `json:"-"`
	Surname     string        `json:"surname" gorm:"size:255"`
	Name        string        `json:"name" gorm:"size:255"`
	Patronymic  string        `json:"patronymic" gorm:"size:255"`
	Sex         string        `json:"sex"`
	BirthDate   time.Time     `json:"birthDate"`
	Profession  string        `json:"profession"`
	Info        string        `json:"info"`
	Slug        string        `json:"slug"`
	RoleID      uint          `json:"roleId"`
	Role        *Role         `json:"role" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Departments []*Department `json:"departments" gorm:"many2many:users_departments;"`
	Services    []*Service    `json:"services" gorm:"many2many:users_services;"`
}

// TableName returns the table name of user struct and it is used by gorm.
func (*User) TableName() string {
	return "user_master"
}

// NewUser is constructor.
func NewUser(username, password string, roleID uint) *User {
	return &User{Username: username, Password: password, RoleID: roleID}
}

// Exist returns true if a given user exits.
func (m *User) Exist(rep repository.Repository, id uint) (bool, error) {
	if err := rep.First(&User{}, id).Error; err != nil {
		return false, err
	}
	return true, nil
}

// Login find user by using username, e-mail or phone.
func (m *User) Login(rep repository.Repository, login, password string) (*User, error) {
	user := &User{}

	if err := rep.Preload("Role").Preload("Departments").Preload("Services").
		Where("username = ?", login).Or("email = ?", login).Or("phone = ?", login).
		First(user).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}

// Get returns user full matched given user ID.
func (m *User) Get(rep repository.Repository, id uint) (*User, error) {
	user := &User{}
	if err := rep.Preload("Role").Preload("Departments").
		Preload("Services").First(user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetBySlug returns user full matched given user slug.
func (m *User) GetBySlug(rep repository.Repository, slug string) (*User, error) {
	user := &User{}
	if err := rep.Preload("Role").Preload("Departments").
		Preload("Services").First(user, "slug = ?", slug).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetAll returns a slice of all users.
func (m *User) GetAll(rep repository.Repository) ([]*User, error) {
	var users []*User
	if err := rep.Preload("Role").Preload("Departments").
		Preload("Services").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// Create persists this user data.
func (m *User) Create(rep repository.Repository) (*User, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		return txCreateUser(tx, m)
	}); err != nil {
		return nil, err
	}

	return m.Get(rep, m.ID)
}

func txCreateUser(tx repository.Repository, m *User) error {
	role := &Role{}
	if _, err := role.Exist(tx, m.RoleID); err != nil {
		return err
	}

	m.Active = true
	if hashed, err := bcrypt.GenerateFromPassword([]byte(m.Password), config.PasswordHashCost); err != nil {
		return err
	} else {
		m.Password = string(hashed)
	}

	if err := tx.Select("username", "active", "password", "role_id").Create(m).Error; err != nil {
		return err
	}
	return tx.Model(m).Where("id = ?", m.ID).Update("slug", strconv.Itoa(int(m.ID))).Error
}

// Update updates this user data.
// If the owner value is true, then fields are included that can only be changed by a user with the owner role.
func (m *User) Update(rep repository.Repository, id uint, owner bool) (*User, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		var err error
		if owner {
			err = txUpdateUserOwner(tx, m, id)
		} else {
			err = txUpdateUser(tx, m, id)
		}
		return err
	}); err != nil {
		return nil, err
	}

	return m.Get(rep, id)
}

func txUpdateUser(tx repository.Repository, m *User, id uint) error {
	if _, err := m.Exist(tx, id); err != nil {
		return err
	}

	makeUserSlug(m)

	return tx.Model(&User{}).Where("id = ?", id).
		Select("email", "phone", "surname", "name", "patronymic", "sex", "birth_date", "slug").
		Updates(m).Error
}

func txUpdateUserOwner(tx repository.Repository, m *User, id uint) error {
	if _, err := m.Exist(tx, id); err != nil {
		return err
	}

	role := &Role{}
	if _, err := role.Exist(tx, m.RoleID); err != nil {
		return err
	}

	i := 0
	for _, department := range m.Departments {
		if ok, _ := department.Exist(tx, department.ID); ok {
			m.Departments[i] = department
			i++
		}
	}
	for j := i; j < len(m.Departments); j++ {
		m.Departments[j] = nil
	}
	m.Departments = m.Departments[:i]

	if err := tx.Model(&User{BaseModel: &BaseModel{ID: id}}).
		Association("Departments").Replace(m.Departments); err != nil {
		return err
	}

	i = 0
	for _, service := range m.Services {
		if ok, _ := service.Exist(tx, service.ID); ok {
			m.Services[i] = service
			i++
		}
	}
	for j := i; j < len(m.Services); j++ {
		m.Services[j] = nil
	}
	m.Services = m.Services[:i]

	if err := tx.Model(&User{BaseModel: &BaseModel{ID: id}}).
		Association("Services").Replace(m.Services); err != nil {
		return err
	}

	makeUserSlug(m)

	return tx.Model(&User{}).Where("id = ?", id).
		Select("email", "phone", "active", "surname", "name", "patronymic",
			"sex", "birth_date", "profession", "info", "slug", "role_id").
		Updates(m).Error
}

// UpdatePassword updates this user password.
func (m *User) UpdatePassword(rep repository.Repository, id uint, old, new string) error {
	if err := rep.Transaction(func(tx repository.Repository) error {
		user := &User{}
		var err error
		if user, err = m.Get(tx, id); err != nil {
			return err
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(old)); err != nil {
			return err
		}

		var hashed []byte
		if hashed, err = bcrypt.GenerateFromPassword([]byte(new), config.PasswordHashCost); err != nil {
			return err
		}

		return tx.Model(m).Where("id = ?", id).Update("password", string(hashed)).Error
	}); err != nil {
		return err
	}

	return nil
}

// Delete deletes this user data.
func (m *User) Delete(rep repository.Repository, id uint) (*User, error) {
	user := &User{}

	if err := rep.Transaction(func(tx repository.Repository) error {
		var err error
		if user, err = m.Get(tx, id); err != nil {
			return err
		}
		return tx.Delete(&User{}, id).Error
	}); err != nil {
		return nil, err
	}

	return user, nil
}

func makeUserSlug(user *User) {
	slug.MaxLength = 40
	slug.EnableSmartTruncate = false
	slug.CustomSub = map[string]string{
		"ь": "",
		"Ь": "",
		"ъ": "",
		"Ъ": "",
	}
	if user.Surname != "" || user.Name != "" || user.Patronymic != "" {
		user.Slug = slug.Make(fmt.Sprintf("%s %s %s", user.Surname, user.Name, user.Patronymic))
	} else {
		user.Slug = strconv.Itoa(int(user.ID))
	}
}
