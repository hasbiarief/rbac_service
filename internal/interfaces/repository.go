package interfaces

import "gin-scalable-api/internal/models"

// UserRepositoryInterface mendefinisikan kontrak untuk repository pengguna
type UserRepositoryInterface interface {
	GetAll(limit, offset int, search string, isActive *bool) ([]*models.User, error)
	GetByID(id int64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUserIdentity(userIdentity string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int64) error
	GetWithRoles(id int64) (*models.User, error)
	Count(search string, isActive *bool) (int64, error)
	GetUserModulesGroupedWithSubscription(userID int64) (map[string][][]string, error)
	GetUserModulesWithSubscription(userID int64) ([]string, error)
	GetUserRoles(userID int64) ([]*models.UserRole, error)
}

// CompanyRepositoryInterface mendefinisikan kontrak untuk repository perusahaan
type CompanyRepositoryInterface interface {
	GetAll(limit, offset int, search string, isActive *bool) ([]*models.Company, error)
	GetByID(id int64) (*models.Company, error)
	GetByCode(code string) (*models.Company, error)
	GetWithStats(id int64) (*models.CompanyWithStats, error)
	Create(company *models.Company) error
	Update(company *models.Company) error
	Delete(id int64) error
	Count(search string, isActive *bool) (int64, error)
}

// RoleRepositoryInterface mendefinisikan kontrak untuk repository peran
type RoleRepositoryInterface interface {
	GetAll(limit, offset int, search string, isActive *bool) ([]*models.Role, error)
	GetByID(id int64) (*models.Role, error)
	GetByName(name string) (*models.Role, error)
	GetWithPermissions(id int64) (*models.RoleWithPermissions, error)
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(id int64) error
	Count(search string, isActive *bool) (int64, error)

	// Izin peran
	GetRoleModules(roleID int64) ([]*models.RoleModule, error)
	UpdateRoleModules(roleID int64, modules []*models.RoleModule) error

	// Penugasan peran pengguna
	AssignUserRole(userRole *models.UserRole) error
	RemoveUserRole(userID, roleID, companyID int64) error
	GetUserRoles(userID int64) ([]*models.UserRole, error)
	GetUsersByRole(roleID int64, limit int) ([]*models.User, error)
}

// ModuleRepositoryInterface mendefinisikan kontrak untuk repository modul
type ModuleRepositoryInterface interface {
	GetAll(limit, offset int, search string, category, subscriptionTier string, parentID *int64, isActive *bool) ([]*models.Module, error)
	GetByID(id int64) (*models.Module, error)
	GetByURL(url string) (*models.Module, error)
	GetUserModules(userID, companyID int64) ([]*models.UserModule, error)
	GetChildren(parentID int64) ([]*models.Module, error)
	Create(module *models.Module) error
	Update(module *models.Module) error
	Delete(id int64) error
	Count(search string, category, subscriptionTier string, parentID *int64, isActive *bool) (int64, error)
	GetTreeStructure(category string, userID int64) ([]*models.Module, error)
	GetTreeByParentName(parentName string, userID int64) ([]*models.Module, error)
	GetAncestors(moduleID int64, userID int64) ([]*models.Module, error)
	CheckUserAccess(userID int64, moduleURL string) (bool, error)
}

// AuditRepositoryInterface mendefinisikan kontrak untuk repository audit
type AuditRepositoryInterface interface {
	GetAll(limit, offset int, filters map[string]interface{}) ([]*models.AuditLogWithUser, error)
	GetByID(id int64) (*models.AuditLogWithUser, error)
	Create(auditLog *models.AuditLog) error
	GetStats() (*models.AuditStats, error)
	CleanupOldLogs(daysToKeep int) error
	Count(filters map[string]interface{}) (int64, error)
	GetUserLogs(userID int64, limit int) ([]*models.AuditLogWithUser, error)
	GetUserLogsByIdentity(userIdentity string, limit int) ([]*models.AuditLogWithUser, error)
}

// SubscriptionRepositoryInterface mendefinisikan kontrak untuk repository langganan
type SubscriptionRepositoryInterface interface {
	// Paket Langganan
	GetAllPlans() ([]*models.SubscriptionPlan, error)
	GetPlanByID(id int64) (*models.SubscriptionPlan, error)
	GetPlanByName(name string) (*models.SubscriptionPlan, error)
	CreatePlan(plan *models.SubscriptionPlan) error
	UpdatePlan(plan *models.SubscriptionPlan) error
	DeletePlan(id int64) error

	// Langganan
	GetAll(limit, offset int, filters map[string]interface{}) ([]*models.Subscription, error)
	GetByID(id int64) (*models.Subscription, error)
	GetByCompanyID(companyID int64) (*models.Subscription, error)
	Create(subscription *models.Subscription) error
	Update(subscription *models.Subscription) error
	RenewSubscription(subscriptionID int64, planID *int64, billingCycle string) error
	Delete(id int64) error
	Count(filters map[string]interface{}) (int64, error)

	// Additional methods
	CheckModuleAccess(companyID, moduleID int64) (bool, error)
	GetExpiringSubscriptions(days int) ([]*models.Subscription, error)
	UpdateExpiredSubscriptions() error
	MarkPaymentAsPaid(subscriptionID int64) error
}

// BranchRepositoryInterface mendefinisikan kontrak untuk repository cabang
type BranchRepositoryInterface interface {
	GetAll(limit, offset int, search string, companyID *int64, isActive *bool) ([]*models.Branch, error)
	GetByID(id int64) (*models.Branch, error)
	Create(branch *models.Branch) error
	Update(branch *models.Branch) error
	Delete(id int64) error
	GetByCompany(companyID int64, includeHierarchy bool) ([]*models.Branch, error)
	GetChildren(parentID int64) ([]*models.Branch, error)
}
