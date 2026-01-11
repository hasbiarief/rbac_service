package dto

// CreateUser Request DTO
type CreateUserRequest struct {
	Name         string  `json:"name" validate:"required,min=2,max=100"`
	Email        string  `json:"email" validate:"required,email"`
	UserIdentity *string `json:"user_identity"`
	Password     string  `json:"password" validate:"required,min=6"`
}

// UpdateUser Request DTO
type UpdateUserRequest struct {
	Name         string  `json:"name"`
	Email        string  `json:"email" validate:"omitempty,email"`
	UserIdentity *string `json:"user_identity"`
	IsActive     *bool   `json:"is_active"`
}

// ChangePassword Request DTO
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// AccessCheck Request DTO
type AccessCheckRequest struct {
	UserIdentity string `json:"user_identity" validate:"required"`
	ModuleURL    string `json:"module_url" validate:"required"`
}

// UserList Request DTO
type UserListRequest struct {
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
	Search   string `form:"search"`
	IsActive *bool  `form:"is_active"`
}

// User Response DTO
type UserResponse struct {
	ID              int64                    `json:"id"`
	Name            string                   `json:"name"`
	Email           string                   `json:"email"`
	UserIdentity    *string                  `json:"user_identity"`
	IsActive        bool                     `json:"is_active"`
	CreatedAt       string                   `json:"created_at"`
	UpdatedAt       string                   `json:"updated_at"`
	Roles           []string                 `json:"roles,omitempty"`   // For backward compatibility
	Modules         map[string][][]string    `json:"modules,omitempty"` // For login response
	RoleAssignments []map[string]interface{} `json:"role_assignments"`  // Enhanced role assignments
	TotalRoles      int                      `json:"total_roles"`       // Total role count
}

// UserWithRoles Response DTO
type UserWithRolesResponse struct {
	UserResponse
	Roles []UserRoleResponse `json:"roles"`
}

// UserRole Response DTO
type UserRoleResponse struct {
	ID          int64   `json:"id"`
	RoleID      int64   `json:"role_id"`
	RoleName    string  `json:"role_name"`
	CompanyID   int64   `json:"company_id"`
	CompanyName string  `json:"company_name"`
	BranchID    *int64  `json:"branch_id"`
	BranchName  *string `json:"branch_name"`
}

// UserList Response DTO
type UserListResponse struct {
	Data    []*UserResponse `json:"data"`
	Total   int64           `json:"total"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
	HasMore bool            `json:"has_more"`
}
