package util

type AccessLevel int

const (
	Unauthorized AccessLevel = iota
	Staff
	Administrator
	Owner
	Superuser
)

func (l AccessLevel) ToString() string {
	switch l {
	case Superuser:
		return "Superuser"
	case Owner:
		return "Owner"
	case Administrator:
		return "Admin"
	case Staff:
		return "Staff"
	default:
		return "Unauthorized"
	}
}

func ToAccessLevel(name string) AccessLevel {
	switch name {
	case "Superuser":
		return Superuser
	case "Owner":
		return Owner
	case "Admin":
		return Administrator
	case "Staff":
		return Staff
	default:
		return Unauthorized
	}
}

func (l AccessLevel) AccessAllowed(level AccessLevel) bool {
	if level >= l {
		return true
	}
	return false
}
