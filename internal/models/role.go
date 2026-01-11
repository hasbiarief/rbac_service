package models

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

	// Additional fields for joined data
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

type RoleWithPermissions struct {
	Role
	Modules []RoleModulePermission `json:"modules"`
}

type RoleModulePermission struct {
	ModuleID   int64  `json:"module_id" db:"module_id"`
	ModuleName string `json:"module_name" db:"module_name"`
	ModuleURL  string `json:"module_url" db:"module_url"`
	CanRead    bool   `json:"can_read" db:"can_read"`
	CanWrite   bool   `json:"can_write" db:"can_write"`
	CanDelete  bool   `json:"can_delete" db:"can_delete"`
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
