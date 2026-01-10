package dto

// Role Request DTO
type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description"`
}

// Update Role Request DTO
type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

// Assign Role Request DTO
type AssignRoleRequest struct {
	UserID    int64  `json:"user_id" validate:"required"`
	RoleID    int64  `json:"role_id" validate:"required"`
	CompanyID int64  `json:"company_id" validate:"required"`
	BranchID  *int64 `json:"branch_id"`
}

// Role Permission Request DTO
type RolePermissionRequest struct {
	ModuleID  int64 `json:"module_id" validate:"required"`
	CanRead   bool  `json:"can_read"`
	CanWrite  bool  `json:"can_write"`
	CanDelete bool  `json:"can_delete"`
}

// Update Role Permissions Request DTO
type UpdateRolePermissionsRequest struct {
	Permissions []RolePermissionRequest `json:"permissions" validate:"required,dive"`
}

// Role List Request DTO
type RoleListRequest struct {
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
	Search   string `form:"search"`
	IsActive *bool  `form:"is_active"`
}

// Role Response DTO
type RoleResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Role With Permissions Response DTO
type RoleWithPermissionsResponse struct {
	RoleResponse
	Modules []RoleModulePermissionResponse `json:"modules"`
}

// Role Module Permission Response DTO
type RoleModulePermissionResponse struct {
	ModuleID   int64  `json:"module_id"`
	ModuleName string `json:"module_name"`
	ModuleURL  string `json:"module_url"`
	CanRead    bool   `json:"can_read"`
	CanWrite   bool   `json:"can_write"`
	CanDelete  bool   `json:"can_delete"`
}

// User Role Assignment Response DTO
type UserRoleAssignmentResponse struct {
	ID          int64   `json:"id"`
	UserID      int64   `json:"user_id"`
	RoleID      int64   `json:"role_id"`
	CompanyID   int64   `json:"company_id"`
	BranchID    *int64  `json:"branch_id"`
	RoleName    string  `json:"role_name"`
	CompanyName string  `json:"company_name"`
	BranchName  *string `json:"branch_name"`
	CreatedAt   string  `json:"created_at"`
}

// Role List Response DTO
type RoleListResponse struct {
	Data    []*RoleResponse `json:"data"`
	Total   int64           `json:"total"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
	HasMore bool            `json:"has_more"`
}
