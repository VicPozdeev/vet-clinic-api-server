package config

const ErrExitStatus int = 2

const (
	// AppConfigPath is the path of application.yml.
	AppConfigPath = "config/%s.yml"
)

// PasswordHashCost is hash cost for a password.
const PasswordHashCost int = 12

const (
	// Login represents the path to get the logged in account.
	Login = "/login"
	// Logout represents the path to logout.
	Logout = "/logout"
	// Profile represents a group of paths for managing the user's personal data.
	Profile = "/profile"
	// Password represents the path to change the password.
	Password = Profile + "/password"
	// Roles represents a group of role management paths.
	Roles = "/roles"
	// RolesID represents the path to get role data using the id.
	RolesID = Roles + "/:id"
	// Users represents a group of user management paths.
	Users = "/users"
	// UsersID represents the path to get user data using the id.
	UsersID = Users + "/:id"
	// Departments represents a group of department management paths.
	Departments = "/departments"
	// DepartmentsID represents the path to get department data using the id.
	DepartmentsID = Departments + "/:id"
	// Categories represents a group of category management paths.
	Categories = "/categories"
	// CategoriesID represents the path to get category data using the id.
	CategoriesID = Categories + "/:id"
	// Services represents a group of service management paths.
	Services = "/services"
	// ServicesID represents the path to get service data using the id.
	ServicesID = Services + "/:id"
	// Clients represents a group of client management paths.
	Clients = "/clients"
	// ClientsID represents the path to get client data using the id.
	ClientsID = Clients + "/:id"
	// Pets represents a group of pet management paths.
	Pets = "/pets"
	// PetsID represents the path to get pet data using the id.
	PetsID = Pets + "/:id"
	// Records represents a group of record management paths.
	Records = "/records"
	// RecordsID represents the path to get record data using the id.
	RecordsID = Records + "/:id"
	// Visits represents a group of visit management paths.
	Visits = "/visits"
	// VisitsID represents the path to get visit data using the id.
	VisitsID = Visits + "/:id"
	// Leads represents a group of lead management paths.
	Leads = "/leads"
	// LeadsID represents the path to get lead data using the id.
	LeadsID = Leads + "/:id"
	// LeadsTypes represents the path to get the list of lead types.
	LeadsTypes = Leads + "/types"
	// LeadsStatuses represents the path to get the list of lead statuses.
	LeadsStatuses = Leads + "/statuses"
)

// APIv1 represents the group of API v1.
const APIv1 = "/v1"

const (
	// APIv1Login represents the API v1 to get the logged in account.
	APIv1Login = APIv1 + Login
	// APIv1Logout represents the API v1 to logout.
	APIv1Logout = APIv1 + Logout
	// APIv1Profile represents the API group for managing user's personal data.
	APIv1Profile = APIv1 + Profile
	// APIv1Password represents the API for changing the password
	APIv1Password = APIv1 + Password
	// APIv1Roles represents the group of role management API v1.
	APIv1Roles = APIv1 + Roles
	// APIv1RolesID represents the API v1 to get role data using id.
	APIv1RolesID = APIv1 + RolesID
	// APIv1Users represents the group of user management API v1.
	APIv1Users = APIv1 + Users
	// APIv1UsersID represents the API v1 to get user data using id.
	APIv1UsersID = APIv1 + UsersID
	// APIv1Departments represents a group of department management API v1.
	APIv1Departments = APIv1 + Departments
	// APIv1DepartmentsID represents the API v1 to get department data using the id.
	APIv1DepartmentsID = APIv1 + DepartmentsID
	// APIv1Categories represents a group of category management API v1.
	APIv1Categories = APIv1 + Categories
	// APIv1CategoriesID represents the API v1 to get category data using the id.
	APIv1CategoriesID = APIv1 + CategoriesID
	// APIv1Services represents a group of service management API v1.
	APIv1Services = APIv1 + Services
	// APIv1ServicesID represents the API v1 to get service data using the id.
	APIv1ServicesID = APIv1 + ServicesID
	// APIv1Clients represents a group of client management API v1.
	APIv1Clients = APIv1 + Clients
	// APIv1ClientsID represents the API v1 to get client data using the id.
	APIv1ClientsID = APIv1 + ClientsID
	// APIv1Pets represents a group of pet management API v1.
	APIv1Pets = APIv1 + Pets
	// APIv1PetsID represents the API v1 to get pet data using the id.
	APIv1PetsID = APIv1 + PetsID
	// APIv1Records represents a group of record management API v1.
	APIv1Records = APIv1 + Records
	// APIv1RecordsID represents the API v1 to get record data using the id.
	APIv1RecordsID = APIv1 + RecordsID
	// APIv1Visits represents a group of visit management API v1.
	APIv1Visits = APIv1 + Visits
	// APIv1VisitsID represents the API v1 to get visit data using the id.
	APIv1VisitsID = APIv1 + VisitsID
	// APIv1Leads represents a group of lead management API v1.
	APIv1Leads = APIv1 + Leads
	// APIv1LeadsID represents the API v1 to get lead data using the id.
	APIv1LeadsID = APIv1 + LeadsID
	// APIv1LeadsTypes represents the API v1 to get the list of lead types.
	APIv1LeadsTypes = APIv1 + LeadsTypes
	// APIv1LeadsStatuses represents the API v1 to get the list of lead statuses.
	APIv1LeadsStatuses = APIv1 + LeadsStatuses
)

const (
	// APIv1Health represents the API v1 to get the status of this application.
	APIv1Health = APIv1 + "/health"
)
