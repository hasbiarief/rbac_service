package role

import "github.com/go-playground/validator/v10"

// CreateRoleRequest DTO
type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description"`
}

// UpdateRoleRequest DTO
type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=100"`
	Description string `json:"description" validate:"omitempty"`
	IsActive    *bool  `json:"is_active" validate:"omitempty"`
}

// AssignRoleRequest DTO
type AssignRoleRequest struct {
	UserID    int64  `json:"user_id" validate:"required"`
	RoleID    int64  `json:"role_id" validate:"required"`
	CompanyID int64  `json:"company_id" validate:"required"`
	BranchID  *int64 `json:"branch_id"`
	UnitID    *int64 `json:"unit_id"`
}

// BulkAssignRoleRequest DTO
type BulkAssignRoleRequest struct {
	UserIDs   []int64 `json:"user_ids" validate:"required,min=1"`
	RoleID    int64   `json:"role_id" validate:"required"`
	CompanyID int64   `json:"company_id" validate:"required"`
	BranchID  *int64  `json:"branch_id"`
	UnitID    *int64  `json:"unit_id"`
}

// RolePermissionRequest DTO
type RolePermissionRequest struct {
	ModuleID  int64 `json:"module_id" validate:"required"`
	CanRead   bool  `json:"can_read"`
	CanWrite  bool  `json:"can_write"`
	CanDelete bool  `json:"can_delete"`
}

// UpdateRolePermissionsRequest DTO
type UpdateRolePermissionsRequest struct {
	Modules []RolePermissionRequest `json:"modules" validate:"required,dive"`
}

// AddRoleModulesRequest DTO - untuk menambahkan module ke role
type AddRoleModulesRequest struct {
	Modules []RolePermissionRequest `json:"modules" validate:"required,dive"`
}

// RemoveRoleModulesRequest DTO - untuk menghapus module dari role
type RemoveRoleModulesRequest struct {
	ModuleIDs []int64 `json:"module_ids" validate:"required,min=1"`
}

// RoleListRequest DTO
type RoleListRequest struct {
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
	Search   string `form:"search"`
	IsActive *bool  `form:"is_active"`
}

// RoleResponse DTO
type RoleResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// RoleWithPermissionsResponse DTO
type RoleWithPermissionsResponse struct {
	RoleResponse
	Modules []RoleModulePermissionResponse `json:"modules"`
}

// RoleModulePermissionResponse DTO
type RoleModulePermissionResponse struct {
	ModuleID   int64  `json:"module_id"`
	ModuleName string `json:"module_name"`
	ModuleURL  string `json:"module_url"`
	CanRead    bool   `json:"can_read"`
	CanWrite   bool   `json:"can_write"`
	CanDelete  bool   `json:"can_delete"`
}

// UserRoleAssignmentResponse DTO
type UserRoleAssignmentResponse struct {
	ID          int64   `json:"id"`
	UserID      int64   `json:"user_id"`
	RoleID      int64   `json:"role_id"`
	CompanyID   int64   `json:"company_id"`
	BranchID    *int64  `json:"branch_id"`
	UnitID      *int64  `json:"unit_id"`
	RoleName    string  `json:"role_name"`
	CompanyName string  `json:"company_name"`
	BranchName  *string `json:"branch_name"`
	UnitName    *string `json:"unit_name"`
	CreatedAt   string  `json:"created_at"`
}

// RoleListResponse DTO
type RoleListResponse struct {
	Data    []*RoleResponse `json:"data"`
	Total   int64           `json:"total"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
	HasMore bool            `json:"has_more"`
}

// Validation functions
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateCreateRoleRequest validates create role request
func ValidateCreateRoleRequest(req *CreateRoleRequest) error {
	return validate.Struct(req)
}

// ValidateUpdateRoleRequest validates update role request
func ValidateUpdateRoleRequest(req *UpdateRoleRequest) error {
	return validate.Struct(req)
}

// ValidateAssignRoleRequest validates assign role request
func ValidateAssignRoleRequest(req *AssignRoleRequest) error {
	return validate.Struct(req)
}

// ValidateUpdateRolePermissionsRequest validates update role permissions request
func ValidateUpdateRolePermissionsRequest(req *UpdateRolePermissionsRequest) error {
	return validate.Struct(req)
}

// ValidateAddRoleModulesRequest validates add role modules request
func ValidateAddRoleModulesRequest(req *AddRoleModulesRequest) error {
	return validate.Struct(req)
}

// ValidateRemoveRoleModulesRequest validates remove role modules request
func ValidateRemoveRoleModulesRequest(req *RemoveRoleModulesRequest) error {
	return validate.Struct(req)
}

// ValidateRoleListRequest validates role list request
func ValidateRoleListRequest(req *RoleListRequest) error {
	return validate.Struct(req)
}
