package entity

type Status string

const (
	// Active status
	Active Status = "active"
	// Inactive status
	Inactive Status = "inactive"
	// Pending status
	Pending Status = "pending"
	// Deleted status
	Deleted Status = "deleted"
)

// Role represents the role of a user
type Role string

const (
	// Admin role
	AdminRole   Role = "admin"
	ManagerRole Role = "manager"
	UserRole    Role = "user"
)

type AuthUserKey string

const (
	AuthUser AuthUserKey = "user"
)
