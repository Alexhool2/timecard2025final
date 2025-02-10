package types

const (
	admin Role = "admin"
	user  Role = "user"
)

func IsValidRole(r Role) bool {
	switch r {
	case admin, user:
		return true
	default:
		return false
	}
}
