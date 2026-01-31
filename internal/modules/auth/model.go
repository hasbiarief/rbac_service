package auth

import (
	"time"
)

// User model for auth module - only fields needed for authentication
type User struct {
	ID           int64     `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	UserIdentity *string   `json:"user_identity" db:"user_identity"`
	PasswordHash string    `json:"-" db:"password_hash"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserRole model for auth module
type UserRole struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	RoleID    int64     `json:"role_id" db:"role_id"`
	CompanyID int64     `json:"company_id" db:"company_id"`
	BranchID  *int64    `json:"branch_id" db:"branch_id"`
	UnitID    *int64    `json:"unit_id" db:"unit_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

func (UserRole) TableName() string {
	return "user_roles"
}
