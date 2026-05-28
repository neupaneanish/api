package enum

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusArchived  UserStatus = "archived"
	UserStatusDeleted   UserStatus = "deleted"
	UserStatusDisabled  UserStatus = "disabled"
	UserStatusLocked    UserStatus = "locked"
	UserStatusPending   UserStatus = "pending"
	UserStatusSuspended UserStatus = "suspended"
)

func (s UserStatus) Valid() bool {
	switch s {
	case
		UserStatusActive,
		UserStatusArchived,
		UserStatusDeleted,
		UserStatusDisabled,
		UserStatusLocked,
		UserStatusPending,
		UserStatusSuspended:
		return true
	default:
		return false
	}
}
