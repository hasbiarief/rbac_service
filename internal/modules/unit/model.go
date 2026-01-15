package unit

import "time"

type Unit struct {
	ID          int64     `json:"id" db:"id"`
	BranchID    int64     `json:"branch_id" db:"branch_id"`
	ParentID    *int64    `json:"parent_id" db:"parent_id"`
	Name        string    `json:"name" db:"name"`
	Code        string    `json:"code" db:"code"`
	Description string    `json:"description" db:"description"`
	Level       int       `json:"level" db:"level"`
	Path        string    `json:"path" db:"path"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (Unit) TableName() string {
	return "units"
}

type UnitWithBranch struct {
	Unit
	BranchName  string `json:"branch_name" db:"branch_name"`
	BranchCode  string `json:"branch_code" db:"branch_code"`
	CompanyName string `json:"company_name" db:"company_name"`
	CompanyCode string `json:"company_code" db:"company_code"`
}

type UnitRole struct {
	ID        int64     `json:"id" db:"id"`
	UnitID    int64     `json:"unit_id" db:"unit_id"`
	RoleID    int64     `json:"role_id" db:"role_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UnitName  string    `json:"unit_name,omitempty" db:"unit_name"`
	RoleName  string    `json:"role_name,omitempty" db:"role_name"`
}

func (UnitRole) TableName() string {
	return "unit_roles"
}

type UnitRoleModule struct {
	ID         int64     `json:"id" db:"id"`
	UnitRoleID int64     `json:"unit_role_id" db:"unit_role_id"`
	ModuleID   int64     `json:"module_id" db:"module_id"`
	CanRead    bool      `json:"can_read" db:"can_read"`
	CanWrite   bool      `json:"can_write" db:"can_write"`
	CanDelete  bool      `json:"can_delete" db:"can_delete"`
	CanApprove bool      `json:"can_approve" db:"can_approve"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

func (UnitRoleModule) TableName() string {
	return "unit_role_modules"
}
