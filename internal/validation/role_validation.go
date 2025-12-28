package validation

import "gin-scalable-api/middleware"

// Role validation rules
var CreateRoleValidation = middleware.ValidationRules{
	Body: &struct {
		Name        string `json:"name" validate:"required,min=2,max=100"`
		Description string `json:"description"`
	}{},
}

var UpdateRoleValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body: &struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{},
}

var AssignUserRoleValidation = middleware.ValidationRules{
	Body: &struct {
		UserID    int64  `json:"user_id" validate:"required,min=1"`
		RoleID    int64  `json:"role_id" validate:"required,min=1"`
		CompanyID int64  `json:"company_id" validate:"required,min=1"`
		BranchID  *int64 `json:"branch_id"`
	}{},
}

var BulkAssignRolesValidation = middleware.ValidationRules{
	Body: &struct {
		UserIDs   []int64 `json:"user_ids" validate:"required,min=1"`
		RoleID    int64   `json:"role_id" validate:"required,min=1"`
		CompanyID int64   `json:"company_id" validate:"required,min=1"`
		BranchID  *int64  `json:"branch_id"`
	}{},
}

var UpdateRoleModulesValidation = middleware.ValidationRules{
	Params: []middleware.ParamValidation{
		{Name: "roleId", Type: "int", Required: true, Min: IntPtr(1)},
	},
	Body: &struct {
		Modules []struct {
			ModuleID  int64 `json:"module_id" validate:"required,min=1"`
			CanRead   bool  `json:"can_read"`
			CanWrite  bool  `json:"can_write"`
			CanDelete bool  `json:"can_delete"`
		} `json:"modules" validate:"required,min=1"`
	}{},
}

var RoleUserValidation = middleware.ValidationRules{
	Params: []middleware.ParamValidation{
		{Name: "userId", Type: "int", Required: true, Min: IntPtr(1)},
		{Name: "roleId", Type: "int", Required: true, Min: IntPtr(1)},
	},
}
