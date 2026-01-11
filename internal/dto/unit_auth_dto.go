package dto

import "gin-scalable-api/pkg/rbac"

// UnitLoginResponse extends LoginResponse with unit context
type UnitLoginResponse struct {
	LoginResponse
	UnitContext *UnitContextInfo `json:"unit_context,omitempty"`
}

// UnitContextInfo provides unit context information
type UnitContextInfo struct {
	CompanyID      int64                    `json:"company_id"`
	CompanyName    string                   `json:"company_name"`
	BranchID       *int64                   `json:"branch_id,omitempty"`
	BranchName     string                   `json:"branch_name,omitempty"`
	UnitID         *int64                   `json:"unit_id,omitempty"`
	UnitName       string                   `json:"unit_name,omitempty"`
	EffectiveUnits []int64                  `json:"effective_units"`
	UnitRoles      []rbac.UnitRoleInfo      `json:"unit_roles"`
	AdminLevels    AdminLevels              `json:"admin_levels"`
	Permissions    map[int64]PermissionInfo `json:"permissions"`
}

// AdminLevels represents user's administrative capabilities
type AdminLevels struct {
	IsUnitAdmin    bool `json:"is_unit_admin"`
	IsBranchAdmin  bool `json:"is_branch_admin"`
	IsCompanyAdmin bool `json:"is_company_admin"`
}

// PermissionInfo represents module permission information
type PermissionInfo struct {
	ModuleID     int64                   `json:"module_id"`
	CanRead      bool                    `json:"can_read"`
	CanWrite     bool                    `json:"can_write"`
	CanDelete    bool                    `json:"can_delete"`
	CanApprove   bool                    `json:"can_approve"`
	GrantedBy    []rbac.PermissionSource `json:"granted_by"`
	HighestLevel string                  `json:"highest_level"`
}

// UnitAccessValidationRequest for validating unit access
type UnitAccessValidationRequest struct {
	UnitID int64 `json:"unit_id" validate:"required,min=1"`
}

// UnitAccessValidationResponse for unit access validation result
type UnitAccessValidationResponse struct {
	HasAccess bool   `json:"has_access"`
	UnitID    int64  `json:"unit_id"`
	Message   string `json:"message,omitempty"`
}

// UpdateUnitAssignmentRequest for updating user unit assignment
type UpdateUnitAssignmentRequest struct {
	UserID int64 `json:"user_id" validate:"required,min=1"`
}

// UnitTokenStatsResponse for unit token statistics
type UnitTokenStatsResponse struct {
	TotalTokens   int `json:"total_tokens"`
	UnitTokens    int `json:"unit_tokens"`
	RegularTokens int `json:"regular_tokens"`
	ExpiredTokens int `json:"expired_tokens"`
	CompanyAdmins int `json:"company_admins"`
	BranchAdmins  int `json:"branch_admins"`
	UnitAdmins    int `json:"unit_admins"`
}

// MyUnitPermissionsResponse for current user's unit permissions
type MyUnitPermissionsResponse struct {
	Permissions    interface{} `json:"permissions"`
	CompanyID      int64       `json:"company_id"`
	BranchID       *int64      `json:"branch_id,omitempty"`
	UnitID         *int64      `json:"unit_id,omitempty"`
	EffectiveUnits []int64     `json:"effective_units"`
	UnitRoles      interface{} `json:"unit_roles"`
}

// UnitAssignRoleRequest for assigning role to user in unit context
type UnitAssignRoleRequest struct {
	UserID int64 `json:"user_id" validate:"required,min=1"`
	RoleID int64 `json:"role_id" validate:"required,min=1"`
	UnitID int64 `json:"unit_id" validate:"required,min=1"`
}

// UnitAssignRoleResponse for unit role assignment result
type UnitAssignRoleResponse struct {
	ID      int64  `json:"id"`
	UserID  int64  `json:"user_id"`
	RoleID  int64  `json:"role_id"`
	UnitID  int64  `json:"unit_id"`
	Message string `json:"message"`
}

// UserUnitRolesResponse for user's unit role assignments
type UserUnitRolesResponse struct {
	UserID    int64                `json:"user_id"`
	UserName  string               `json:"user_name"`
	UnitRoles []UserUnitRoleDetail `json:"unit_roles"`
}

// UserUnitRoleDetail represents detailed unit role assignment
type UserUnitRoleDetail struct {
	ID          int64  `json:"id"`
	RoleID      int64  `json:"role_id"`
	RoleName    string `json:"role_name"`
	UnitID      int64  `json:"unit_id"`
	UnitName    string `json:"unit_name"`
	BranchID    int64  `json:"branch_id"`
	BranchName  string `json:"branch_name"`
	CompanyID   int64  `json:"company_id"`
	CompanyName string `json:"company_name"`
	Level       string `json:"level"` // "company", "branch", "unit"
	AssignedAt  string `json:"assigned_at"`
}

// UnitFilterRequest for filtering data by unit access
type UnitFilterRequest struct {
	IncludeSubUnits bool   `json:"include_sub_units" form:"include_sub_units"`
	UnitID          *int64 `json:"unit_id" form:"unit_id"`
	BranchID        *int64 `json:"branch_id" form:"branch_id"`
	CompanyID       *int64 `json:"company_id" form:"company_id"`
}

// UnitPermissionCheckRequest for checking specific unit permissions
type UnitPermissionCheckRequest struct {
	ModuleID   int64  `json:"module_id" validate:"required,min=1"`
	Permission string `json:"permission" validate:"required,oneof=read write delete approve"`
	UnitID     *int64 `json:"unit_id,omitempty"`
}

// UnitPermissionCheckResponse for unit permission check result
type UnitPermissionCheckResponse struct {
	HasPermission bool   `json:"has_permission"`
	ModuleID      int64  `json:"module_id"`
	Permission    string `json:"permission"`
	UnitID        *int64 `json:"unit_id,omitempty"`
	GrantedBy     string `json:"granted_by,omitempty"`
	Level         string `json:"level,omitempty"`
}

// ModuleInfo represents module information
type ModuleInfo struct {
	ID          int64  `json:"id"`
	Category    string `json:"category"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}
