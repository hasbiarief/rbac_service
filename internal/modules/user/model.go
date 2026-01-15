package user

import (
	"gin-scalable-api/pkg/model"
	"time"
)

type User struct {
	model.BaseModel
	Name         string     `json:"name" db:"name"`
	Email        string     `json:"email" db:"email"`
	UserIdentity *string    `json:"user_identity" db:"user_identity"`
	PasswordHash string     `json:"-" db:"password_hash"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

func (User) TableName() string {
	return "users"
}

// UserModule represents user's access to a module
type UserModule struct {
	UserID     int64
	ModuleID   int64
	ModuleName string
	ModuleURL  string
	Category   string
	CanRead    bool
	CanWrite   bool
	CanDelete  bool
	CanApprove bool
	CompanyID  int64
	BranchID   *int64
	UnitID     *int64
}

// UserRole represents user's role assignment
type UserRole struct {
	ID          int64
	UserID      int64
	RoleID      int64
	CompanyID   int64
	BranchID    *int64
	UnitID      *int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	RoleName    string
	CompanyName string
}
