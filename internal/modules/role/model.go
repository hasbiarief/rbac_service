package role

import "time"

type Role struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type UserRole struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	RoleID    int64     `json:"role_id" db:"role_id"`
	CompanyID int64     `json:"company_id" db:"company_id"`
	BranchID  *int64    `json:"branch_id" db:"branch_id"`
	UnitID    *int64    `json:"unit_id" db:"unit_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	RoleName    string `json:"role_name,omitempty" db:"-"`
	CompanyName string `json:"company_name,omitempty" db:"-"`
	BranchName  string `json:"branch_name,omitempty" db:"-"`
	UnitName    string `json:"unit_name,omitempty" db:"-"`
}

type RoleModule struct {
	ID        int64     `json:"id" db:"id"`
	RoleID    int64     `json:"role_id" db:"role_id"`
	ModuleID  int64     `json:"module_id" db:"module_id"`
	CanRead   bool      `json:"can_read" db:"can_read"`
	CanWrite  bool      `json:"can_write" db:"can_write"`
	CanDelete bool      `json:"can_delete" db:"can_delete"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (Role) TableName() string {
	return "roles"
}

func (UserRole) TableName() string {
	return "user_roles"
}

func (RoleModule) TableName() string {
	return "role_modules"
}

// RoleWithPermissions represents a role with its module permissions
type RoleWithPermissions struct {
	Role    Role
	Modules []RoleModulePermission
}

// RoleModulePermission represents module permission details
type RoleModulePermission struct {
	ModuleID   int64
	ModuleName string
	ModuleURL  string
	CanRead    bool
	CanWrite   bool
	CanDelete  bool
}

// User model for role module - minimal fields needed
type User struct {
	ID           int64
	Name         string
	Email        string
	UserIdentity *string
	IsActive     bool
}
