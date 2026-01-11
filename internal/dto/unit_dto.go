package dto

// Unit DTOs
type CreateUnitRequest struct {
	BranchID    int64  `json:"branch_id" validate:"required"`
	ParentID    *int64 `json:"parent_id"`
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Code        string `json:"code" validate:"required,min=2,max=50"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

type UpdateUnitRequest struct {
	ParentID    *int64 `json:"parent_id"`
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Code        string `json:"code" validate:"required,min=2,max=50"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

type UnitResponse struct {
	ID          int64  `json:"id"`
	BranchID    int64  `json:"branch_id"`
	ParentID    *int64 `json:"parent_id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Level       int    `json:"level"`
	Path        string `json:"path"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`

	// Additional fields for extended responses
	BranchName  string `json:"branch_name,omitempty"`
	BranchCode  string `json:"branch_code,omitempty"`
	CompanyName string `json:"company_name,omitempty"`
	CompanyCode string `json:"company_code,omitempty"`
}

type UnitHierarchyResponse struct {
	UnitResponse
	Children []UnitHierarchyResponse `json:"children,omitempty"`
}

type UnitWithStatsResponse struct {
	UnitResponse
	TotalUsers    int `json:"total_users"`
	TotalSubUnits int `json:"total_sub_units"`
	TotalRoles    int `json:"total_roles"`
}

type UnitListRequest struct {
	BranchID *int64 `json:"branch_id" form:"branch_id"`
	ParentID *int64 `json:"parent_id" form:"parent_id"`
	Search   string `json:"search" form:"search"`
	IsActive *bool  `json:"is_active" form:"is_active"`
	Limit    int    `json:"limit" form:"limit" validate:"min=1,max=100"`
	Offset   int    `json:"offset" form:"offset" validate:"min=0"`
}

type UnitListResponse struct {
	Data    []*UnitResponse `json:"data"`
	Total   int64           `json:"total"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
	HasMore bool            `json:"has_more"`
}

// Unit Role DTOs
type CreateUnitRoleRequest struct {
	UnitID int64 `json:"unit_id" validate:"required"`
	RoleID int64 `json:"role_id" validate:"required"`
}

type UnitRoleResponse struct {
	ID        int64  `json:"id"`
	UnitID    int64  `json:"unit_id"`
	RoleID    int64  `json:"role_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	// Additional fields for joined data
	UnitName string `json:"unit_name,omitempty"`
	RoleName string `json:"role_name,omitempty"`
}

type UnitRoleListRequest struct {
	UnitID int64 `json:"unit_id" form:"unit_id"`
	RoleID int64 `json:"role_id" form:"role_id"`
	Limit  int   `json:"limit" form:"limit" validate:"min=1,max=100"`
	Offset int   `json:"offset" form:"offset" validate:"min=0"`
}

type UnitRoleListResponse struct {
	Data    []*UnitRoleResponse `json:"data"`
	Total   int64               `json:"total"`
	Limit   int                 `json:"limit"`
	Offset  int                 `json:"offset"`
	HasMore bool                `json:"has_more"`
}

// Unit Role Module DTOs
type CreateUnitRoleModuleRequest struct {
	UnitRoleID int64 `json:"unit_role_id" validate:"required"`
	ModuleID   int64 `json:"module_id" validate:"required"`
	CanRead    bool  `json:"can_read"`
	CanWrite   bool  `json:"can_write"`
	CanDelete  bool  `json:"can_delete"`
	CanApprove bool  `json:"can_approve"`
}

type UpdateUnitRoleModuleRequest struct {
	CanRead    bool `json:"can_read"`
	CanWrite   bool `json:"can_write"`
	CanDelete  bool `json:"can_delete"`
	CanApprove bool `json:"can_approve"`
}

type UnitRoleModuleResponse struct {
	ID         int64  `json:"id"`
	UnitRoleID int64  `json:"unit_role_id"`
	ModuleID   int64  `json:"module_id"`
	CanRead    bool   `json:"can_read"`
	CanWrite   bool   `json:"can_write"`
	CanDelete  bool   `json:"can_delete"`
	CanApprove bool   `json:"can_approve"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`

	// Additional fields for joined data
	ModuleName     string `json:"module_name,omitempty"`
	ModuleCategory string `json:"module_category,omitempty"`
	ModuleURL      string `json:"module_url,omitempty"`
	UnitName       string `json:"unit_name,omitempty"`
	RoleName       string `json:"role_name,omitempty"`
	IsCustomized   bool   `json:"is_customized,omitempty"` // Different from default role permission
}

type UnitRoleWithPermissionsResponse struct {
	UnitRoleResponse
	Modules []UnitRoleModuleResponse `json:"modules"`
}

type BulkUpdateUnitRoleModulesRequest struct {
	UnitRoleID int64                            `json:"unit_role_id" validate:"required"`
	Modules    []UpdateUnitRoleModulePermission `json:"modules" validate:"required,dive"`
}

type UpdateUnitRoleModulePermission struct {
	ModuleID   int64 `json:"module_id" validate:"required"`
	CanRead    bool  `json:"can_read"`
	CanWrite   bool  `json:"can_write"`
	CanDelete  bool  `json:"can_delete"`
	CanApprove bool  `json:"can_approve"`
}

// Enhanced User Role DTOs with Unit support
type CreateUserRoleRequest struct {
	UserID    int64  `json:"user_id" validate:"required"`
	RoleID    int64  `json:"role_id" validate:"required"`
	CompanyID int64  `json:"company_id" validate:"required"`
	BranchID  *int64 `json:"branch_id"`
	UnitID    *int64 `json:"unit_id"`
}

type UpdateUserRoleRequest struct {
	RoleID   int64  `json:"role_id" validate:"required"`
	BranchID *int64 `json:"branch_id"`
	UnitID   *int64 `json:"unit_id"`
}

type UserRoleWithUnitResponse struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	RoleID    int64  `json:"role_id"`
	CompanyID int64  `json:"company_id"`
	BranchID  *int64 `json:"branch_id"`
	UnitID    *int64 `json:"unit_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	// Additional fields for joined data
	RoleName    string `json:"role_name,omitempty"`
	CompanyName string `json:"company_name,omitempty"`
	BranchName  string `json:"branch_name,omitempty"`
	UnitName    string `json:"unit_name,omitempty"`
}

// Unit Permission Summary
type UnitPermissionSummaryResponse struct {
	UnitID      int64                       `json:"unit_id"`
	UnitName    string                      `json:"unit_name"`
	UnitCode    string                      `json:"unit_code"`
	BranchName  string                      `json:"branch_name"`
	CompanyName string                      `json:"company_name"`
	Roles       []UnitRolePermissionSummary `json:"roles"`
}

type UnitRolePermissionSummary struct {
	RoleID            int64  `json:"role_id"`
	RoleName          string `json:"role_name"`
	TotalModules      int    `json:"total_modules"`
	ReadAccess        int    `json:"read_access"`
	WriteAccess       int    `json:"write_access"`
	DeleteAccess      int    `json:"delete_access"`
	ApprovalAccess    int    `json:"approval_access"`
	CustomizedModules int    `json:"customized_modules"` // Modules with custom permissions
}

// Copy Unit Permissions
type CopyUnitPermissionsRequest struct {
	SourceUnitID      int64 `json:"source_unit_id" validate:"required"`
	TargetUnitID      int64 `json:"target_unit_id" validate:"required"`
	RoleID            int64 `json:"role_id" validate:"required"`
	OverwriteExisting bool  `json:"overwrite_existing"`
}
