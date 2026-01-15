package user

// CreateUserRequest DTO
type CreateUserRequest struct {
	Name         string  `json:"name" validate:"required,min=2,max=100"`
	Email        string  `json:"email" validate:"required,email"`
	UserIdentity *string `json:"user_identity"`
	Password     string  `json:"password" validate:"required,min=6"`
}

// UpdateUserRequest DTO
type UpdateUserRequest struct {
	Name         string  `json:"name"`
	Email        string  `json:"email" validate:"omitempty,email"`
	UserIdentity *string `json:"user_identity"`
	IsActive     *bool   `json:"is_active"`
}

// ChangePasswordRequest DTO
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// AccessCheckRequest DTO
type AccessCheckRequest struct {
	UserIdentity string `json:"user_identity" validate:"required"`
	ModuleURL    string `json:"module_url" validate:"required"`
}

// UserListRequest DTO
type UserListRequest struct {
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
	Search   string `form:"search"`
	IsActive *bool  `form:"is_active"`
}

// UserResponse DTO
type UserResponse struct {
	ID              int64                    `json:"id"`
	Name            string                   `json:"name"`
	Email           string                   `json:"email"`
	UserIdentity    *string                  `json:"user_identity"`
	IsActive        bool                     `json:"is_active"`
	CreatedAt       string                   `json:"created_at"`
	UpdatedAt       string                   `json:"updated_at"`
	Roles           []string                 `json:"roles,omitempty"`
	Modules         map[string][][]string    `json:"modules,omitempty"`
	RoleAssignments []map[string]interface{} `json:"role_assignments"`
	TotalRoles      int                      `json:"total_roles"`
}

// UserWithRolesResponse DTO
type UserWithRolesResponse struct {
	UserResponse
	Roles []UserRoleResponse `json:"roles"`
}

// UserRoleResponse DTO
type UserRoleResponse struct {
	ID          int64   `json:"id"`
	RoleID      int64   `json:"role_id"`
	RoleName    string  `json:"role_name"`
	CompanyID   int64   `json:"company_id"`
	CompanyName string  `json:"company_name"`
	BranchID    *int64  `json:"branch_id"`
	BranchName  *string `json:"branch_name"`
}

// UserListResponse DTO
type UserListResponse struct {
	Data    []*UserResponse `json:"data"`
	Total   int64           `json:"total"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
	HasMore bool            `json:"has_more"`
}
