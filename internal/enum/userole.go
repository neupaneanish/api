package enum

type UserRole string

const (
	UserRoleRoot UserRole = "root"
	UserRoleUser UserRole = "user"
)

func (r UserRole) Valid() bool {
	switch r {
	case UserRoleRoot, UserRoleUser:
		return true
	default:
		return false
	}
}
