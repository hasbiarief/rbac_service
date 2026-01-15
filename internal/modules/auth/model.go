package auth

import (
	"time"
)

// User model for auth module - only fields needed for authentication
type User struct {
	ID           int64
	Name         string
	Email        string
	UserIdentity *string
	PasswordHash string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserRole model for auth module
type UserRole struct {
	ID        int64
	UserID    int64
	RoleID    int64
	CompanyID int64
	BranchID  *int64
	UnitID    *int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
