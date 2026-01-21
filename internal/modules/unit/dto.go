package unit

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

type UnitListRequest struct {
	BranchID *int64 `json:"branch_id" form:"branch_id"`
	ParentID *int64 `json:"parent_id" form:"parent_id"`
	Search   string `json:"search" form:"search"`
	IsActive *bool  `json:"is_active" form:"is_active"`
	Limit    int    `json:"limit" form:"limit" validate:"min=1,max=100"`
	Offset   int    `json:"offset" form:"offset" validate:"min=0"`
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

type UnitListResponse struct {
	Data    []*UnitResponse `json:"data"`
	Total   int64           `json:"total"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
	HasMore bool            `json:"has_more"`
}

type UnitRoleResponse struct {
	ID        int64  `json:"id"`
	UnitID    int64  `json:"unit_id"`
	RoleID    int64  `json:"role_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	UnitName  string `json:"unit_name,omitempty"`
	RoleName  string `json:"role_name,omitempty"`
}

type UnitRoleModuleResponse struct {
	ID             int64  `json:"id"`
	UnitRoleID     int64  `json:"unit_role_id"`
	ModuleID       int64  `json:"module_id"`
	CanRead        bool   `json:"can_read"`
	CanWrite       bool   `json:"can_write"`
	CanDelete      bool   `json:"can_delete"`
	CanApprove     bool   `json:"can_approve"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	ModuleName     string `json:"module_name,omitempty"`
	ModuleCategory string `json:"module_category,omitempty"`
	ModuleURL      string `json:"module_url,omitempty"`
	UnitName       string `json:"unit_name,omitempty"`
	RoleName       string `json:"role_name,omitempty"`
	IsCustomized   bool   `json:"is_customized,omitempty"`
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

type CopyUnitPermissionsRequest struct {
	SourceUnitID      int64 `json:"source_unit_id" validate:"required"`
	TargetUnitID      int64 `json:"target_unit_id" validate:"required"`
	RoleID            int64 `json:"role_id" validate:"required"`
	OverwriteExisting bool  `json:"overwrite_existing"`
}

// Alternative request for more flexible copying between different unit roles
type CopyUnitRolePermissionsRequest struct {
	SourceUnitRoleID  int64 `json:"source_unit_role_id" validate:"required"`
	TargetUnitRoleID  int64 `json:"target_unit_role_id" validate:"required"`
	OverwriteExisting bool  `json:"overwrite_existing"`
}
