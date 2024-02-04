package migration

import (
	"time"
	"vet-clinic/container"
	"vet-clinic/models"
	"vet-clinic/util"
)

// InitMasterData creates the master data used in this application.
func InitMasterData(container container.Container) {
	if container.Config().Extension.MasterGenerator {
		rep := container.Repository()

		rep.Create(models.NewRole(util.Staff.ToString()))
		rep.Create(models.NewRole(util.Administrator.ToString()))
		rep.Create(models.NewRole(util.Owner.ToString()))
		rep.Create(models.NewRole(util.Superuser.ToString()))

		rep.Create(models.NewCategory("Консультация"))
		rep.Create(models.NewCategory("Процедуры"))
		rep.Create(models.NewCategory("Кардиология"))
		rep.Create(models.NewCategory("Инструментальная диагностика"))

		rep.Create(&models.Service{Name: "Консультация", Price: 1000, CategoryID: 1})
		rep.Create(&models.Service{Name: "Прием врача терапевта", Price: 3000, CategoryID: 1})
		rep.Create(&models.Service{Name: "Стрижка когтей", Price: 800, CategoryID: 2})
		rep.Create(&models.Service{Name: "Глюкометрия", Price: 400, CategoryID: 2})
		rep.Create(&models.Service{Name: "Вакцинация", Price: 2500, CategoryID: 2})
		rep.Create(&models.Service{Name: "Залог за прибор для телеметрии", Price: 30000, CategoryID: 3})
		rep.Create(&models.Service{Name: "ЭхоКГ скрининг", Price: 3500, CategoryID: 4})
		rep.Create(&models.Service{Name: "Холтеровское мониторирование", Price: 9500, CategoryID: 4})

		dep1 := &models.Department{
			Name: "Терапия",
			Services: []*models.Service{
				{BaseModel: &models.BaseModel{ID: 1}},
				{BaseModel: &models.BaseModel{ID: 2}},
				{BaseModel: &models.BaseModel{ID: 3}},
				{BaseModel: &models.BaseModel{ID: 4}},
				{BaseModel: &models.BaseModel{ID: 5}},
			},
		}
		_, _ = dep1.Create(rep)
		dep2 := &models.Department{
			Name: "Кардиология",
			Services: []*models.Service{
				{BaseModel: &models.BaseModel{ID: 1}},
				{BaseModel: &models.BaseModel{ID: 6}},
				{BaseModel: &models.BaseModel{ID: 7}},
				{BaseModel: &models.BaseModel{ID: 8}},
			},
		}
		_, _ = dep2.Create(rep)

		user1 := &models.User{
			Username:   "Test1",
			Email:      "test1@test.com",
			Phone:      "+71111111111",
			Password:   "Password1!",
			Surname:    "Фамилия",
			Name:       "Имя",
			Patronymic: "Отчество",
			Sex:        "Женский",
			BirthDate:  time.Date(1995, time.January, 1, 0, 0, 0, 0, time.Local),
			Profession: "Терапевт",
			Info:       "Информация",
			RoleID:     4,
			Departments: []*models.Department{
				{BaseModel: &models.BaseModel{ID: 1}},
			},
			Services: []*models.Service{
				{BaseModel: &models.BaseModel{ID: 1}},
				{BaseModel: &models.BaseModel{ID: 3}},
				{BaseModel: &models.BaseModel{ID: 4}},
				{BaseModel: &models.BaseModel{ID: 5}},
			},
		}
		_, _ = user1.Create(rep)
		_, _ = user1.Update(rep, 1, true)

		client1 := &models.Client{
			Surname:    "Фамилия",
			Name:       "Имя",
			Patronymic: "Отчество",
			Sex:        "Мужской",
			BirthDate:  time.Date(1991, time.January, 1, 0, 0, 0, 0, time.Local),
			Phone:      "+78888888888",
			Email:      "mail@mail.su",
			Info:       "Информация",
		}
		_, _ = client1.Create(rep)

		pet1 := &models.Pet{
			Name:     "Китти",
			Type:     "Кошка",
			Breed:    "Дворняга",
			Colour:   "Серый полосатый",
			Sex:      "Самка",
			ClientID: 1,
		}
		_, _ = pet1.Create(rep)

		visit1 := &models.Visit{
			DateTime:        time.Date(2024, time.January, 1, 0, 0, 0, 0, time.Local),
			Info:            "Вакцинация",
			ClientID:        1,
			PetID:           1,
			DoctorID:        1,
			LastUpdatedByID: 1,
			ServiceID:       5,
		}
		_, _ = visit1.Create(rep)

		lead1 := &models.Lead{
			Name:            "Александр",
			Phone:           "+79992225566",
			Email:           "alex@test.com",
			Comment:         "Комментарий клиента",
			Type:            "callback",
			Status:          "rejected",
			DoctorID:        1,
			LastUpdatedByID: 1,
		}

		_, _ = lead1.Create(rep)
	}
}
