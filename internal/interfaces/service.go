package interfaces

import "gin-scalable-api/internal/dto"

// AuthServiceInterface mendefinisikan kontrak untuk service autentikasi
type AuthServiceInterface interface {
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	RefreshToken(req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
	Logout(token string) error
	LogoutByUserID(userID int64) error
	ForgotPassword(req *dto.ForgotPasswordRequest) error
	ResetPassword(req *dto.ResetPasswordRequest) error
	CheckUserTokens(userID int64) (interface{}, error)
	GetUserSessionCount(userID int64) (int64, error)
	GetUserRefreshTokenCount(userID int64) (int64, error)
	CleanupExpiredTokens() error
}

// UserServiceInterface mendefinisikan kontrak untuk service pengguna
type UserServiceInterface interface {
	GetUsers(req *dto.UserListRequest) (*dto.UserListResponse, error)
	GetUsersFiltered(requestingUserID int64, req *dto.UserListRequest) (*dto.UserListResponse, error)
	GetUserByID(id int64) (*dto.UserResponse, error)
	GetUserWithRoles(id int64) (*dto.UserWithRolesResponse, error)
	CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error)
	UpdateUser(id int64, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(id int64) error
	ChangePassword(userID int64, req *dto.ChangePasswordRequest) error
	ChangeUserPassword(userID int64, req *dto.ChangePasswordRequest) error
}

// CompanyServiceInterface mendefinisikan kontrak untuk service perusahaan
type CompanyServiceInterface interface {
	GetCompanies(req *dto.CompanyListRequest) (*dto.CompanyListResponse, error)
	GetCompanyByID(id int64) (*dto.CompanyResponse, error)
	GetCompanyWithStats(id int64) (*dto.CompanyWithStatsResponse, error)
	CreateCompany(req *dto.CreateCompanyRequest) (*dto.CompanyResponse, error)
	UpdateCompany(id int64, req *dto.UpdateCompanyRequest) (*dto.CompanyResponse, error)
	DeleteCompany(id int64) error
}

// RoleServiceInterface mendefinisikan kontrak untuk service peran
type RoleServiceInterface interface {
	GetRoles(req *dto.RoleListRequest) (*dto.RoleListResponse, error)
	GetRoleByID(id int64) (*dto.RoleResponse, error)
	GetRoleWithPermissions(id int64) (*dto.RoleWithPermissionsResponse, error)
	CreateRole(req *dto.CreateRoleRequest) (*dto.RoleResponse, error)
	UpdateRole(id int64, req *dto.UpdateRoleRequest) (*dto.RoleResponse, error)
	DeleteRole(id int64) error
	UpdateRolePermissions(roleID int64, req *dto.UpdateRolePermissionsRequest) error
	AssignRoleToUser(req *dto.AssignRoleRequest) (*dto.UserRoleAssignmentResponse, error)
	RemoveRoleFromUser(userID, roleID, companyID int64) error
	GetUsersByRole(roleID int64, limit int) (interface{}, error)
	GetUserRoles(userID int64) (interface{}, error)
	GetUserAccessSummary(userID int64) (interface{}, error)
}

// ModuleServiceInterface mendefinisikan kontrak untuk service modul
type ModuleServiceInterface interface {
	GetModules(req *dto.ModuleListRequest) (*dto.ModuleListResponse, error)
	GetModulesFiltered(requestingUserID int64, req *dto.ModuleListRequest) (*dto.ModuleListResponse, error)
	GetModulesNested(req *dto.ModuleListRequest) ([]*dto.NestedModuleResponse, error)
	GetModuleByID(id int64) (*dto.ModuleResponse, error)
	GetUserModules(userID int64, category string, limit int) ([]*dto.ModuleResponse, error)
	GetModuleTreeByParentFiltered(userID int64, parentName string) ([]*dto.ModuleTreeResponse, error)
	GetModuleTreeFiltered(userID int64, category string) ([]*dto.ModuleTreeResponse, error)
	GetModuleChildrenFiltered(userID int64, id int64) ([]*dto.ModuleResponse, error)
	GetModuleAncestorsFiltered(userID int64, id int64) ([]*dto.ModuleResponse, error)
	CreateModule(req *dto.CreateModuleRequest) (*dto.ModuleResponse, error)
	UpdateModule(id int64, req *dto.UpdateModuleRequest) (*dto.ModuleResponse, error)
	DeleteModule(id int64) error
}

// AuditServiceInterface mendefinisikan kontrak untuk service audit
type AuditServiceInterface interface {
	GetAuditLogs(req *dto.AuditListRequest) (*dto.AuditListResponse, error)
	GetAuditLogByID(id int64) (*dto.AuditLogWithUserResponse, error)
	GetAuditStats() (*dto.AuditStatsResponse, error)
	CreateAuditLog(req *dto.CreateAuditLogRequest) (*dto.AuditLogResponse, error)
	GetUserAuditLogs(userID int64, limit int) (*dto.AuditListResponse, error)
	GetUserAuditLogsByIdentity(identity string, limit int) (*dto.AuditListResponse, error)
	CleanupOldLogs(daysToKeep int) error
}

// SubscriptionServiceInterface mendefinisikan kontrak untuk service langganan
type SubscriptionServiceInterface interface {
	// Paket Langganan
	GetSubscriptionPlans() ([]*dto.SubscriptionPlanResponse, error)
	GetSubscriptionPlanByID(id int64) (*dto.SubscriptionPlanResponse, error)
	CreateSubscriptionPlan(req *dto.CreateSubscriptionPlanRequest) (*dto.SubscriptionPlanResponse, error)
	UpdateSubscriptionPlan(id int64, req *dto.UpdateSubscriptionPlanRequest) (*dto.SubscriptionPlanResponse, error)
	DeleteSubscriptionPlan(id int64) error

	// Langganan
	GetSubscriptions(req *dto.SubscriptionListRequest) (*dto.SubscriptionListResponse, error)
	GetSubscriptionByID(id int64) (*dto.SubscriptionResponse, error)
	GetCompanySubscription(companyID int64) (*dto.SubscriptionResponse, error)
	CreateSubscription(req *dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, error)
	UpdateSubscription(id int64, req *dto.UpdateSubscriptionRequest) (*dto.SubscriptionResponse, error)
	RenewSubscription(subscriptionID int64, planID *int64, billingCycle string) (*dto.SubscriptionResponse, error)
	CancelSubscription(id int64) error

	// Additional methods used by handler
	CheckModuleAccess(companyID, moduleID int64) (bool, error)
	GetExpiringSubscriptions(days int) (interface{}, error)
	UpdateExpiredSubscriptions() error
	GetSubscriptionStats() (interface{}, error)
	MarkPaymentAsPaid(subscriptionID int64) error
}

// BranchServiceInterface mendefinisikan kontrak untuk service cabang
type BranchServiceInterface interface {
	GetBranches(req *dto.BranchListRequest) (*dto.BranchListResponse, error)
	GetBranchesNested(req *dto.BranchListRequest) ([]*dto.NestedBranchResponse, error)
	GetBranchByID(id int64) (*dto.BranchResponse, error)
	CreateBranch(req *dto.CreateBranchRequest) (*dto.BranchResponse, error)
	UpdateBranch(id int64, req *dto.UpdateBranchRequest) (*dto.BranchResponse, error)
	DeleteBranch(id int64) error
	GetCompanyBranches(companyID int64, includeHierarchy bool) ([]*dto.BranchResponse, error)
	GetCompanyBranchesNested(companyID int64) ([]*dto.NestedBranchResponse, error)
	GetBranchChildren(parentID int64) ([]*dto.BranchResponse, error)
	GetBranchHierarchyByID(id int64, nested bool) (interface{}, error)
}
