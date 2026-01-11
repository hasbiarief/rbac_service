package constants

// Pesan Status HTTP
const (
	MsgSuccess          = "Success"
	MsgCreated          = "Created"
	MsgUpdated          = "Updated"
	MsgDeleted          = "Deleted"
	MsgNotFound         = "Not Found"
	MsgValidationFailed = "Validation Failed"
	MsgUnauthorized     = "Unauthorized"
	MsgForbidden        = "Forbidden"
)

// Pesan autentikasi
const (
	MsgLoginSuccess         = "Login successful"
	MsgLogoutSuccess        = "Logout successful"
	MsgRegisterSuccess      = "Registration successful"
	MsgTokenRefreshed       = "Token successfully refreshed"
	MsgPasswordChanged      = "Password successfully changed"
	MsgPasswordResetSent    = "Password reset email has been sent"
	MsgPasswordResetSuccess = "Password reset successful"
	MsgInvalidCredentials   = "Invalid credentials"
	MsgTokenExpired         = "Token has expired"
	MsgTokenInvalid         = "Token is invalid"
)

// Pesan pengguna
const (
	MsgUserRetrieved      = "User successfully retrieved"
	MsgUsersRetrieved     = "Users list successfully retrieved"
	MsgUserCreated        = "User successfully created"
	MsgUserUpdated        = "User successfully updated"
	MsgUserDeleted        = "User successfully deleted"
	MsgUserNotFound       = "User not found"
	MsgEmailAlreadyExists = "Email already exists"
)

// Pesan perusahaan
const (
	MsgCompanyRetrieved   = "Company successfully retrieved"
	MsgCompaniesRetrieved = "Companies list successfully retrieved"
	MsgCompanyCreated     = "Company successfully created"
	MsgCompanyUpdated     = "Company successfully updated"
	MsgCompanyDeleted     = "Company successfully deleted"
	MsgCompanyNotFound    = "Company not found"
	MsgCompanyCodeExists  = "Company code already exists"
)

// Pesan peran
const (
	MsgRoleRetrieved      = "Role successfully retrieved"
	MsgRolesRetrieved     = "Roles list successfully retrieved"
	MsgRoleCreated        = "Role successfully created"
	MsgRoleUpdated        = "Role successfully updated"
	MsgRoleDeleted        = "Role successfully deleted"
	MsgRoleNotFound       = "Role not found"
	MsgRoleNameExists     = "Role name already exists"
	MsgRoleAssigned       = "Role successfully assigned"
	MsgRoleUnassigned     = "Role successfully unassigned"
	MsgPermissionsUpdated = "Permissions successfully updated"
)

// Pesan modul
const (
	MsgModuleRetrieved      = "Module successfully retrieved"
	MsgModulesRetrieved     = "Modules list successfully retrieved"
	MsgModuleCreated        = "Module successfully created"
	MsgModuleUpdated        = "Module successfully updated"
	MsgModuleDeleted        = "Module successfully deleted"
	MsgModuleNotFound       = "Module not found"
	MsgModuleURLExists      = "module's URL already available"
	MsgUserModulesRetrieved = "User's module successfully retrieved"
)

// Pesan cabang
const (
	MsgBranchRetrieved          = "Branch successfully retrieved"
	MsgBranchesRetrieved        = "Branches list successfully retrieved"
	MsgBranchCreated            = "Branch successfully created"
	MsgBranchUpdated            = "Branch successfully updated"
	MsgBranchDeleted            = "Branch successfully deleted"
	MsgBranchNotFound           = "Branch not found"
	MsgBranchHasChildren        = "Cannot delete branch that has child branches"
	MsgCompanyBranchesRetrieved = "Company branches successfully retrieved"
	MsgBranchChildrenRetrieved  = "Branch children successfully retrieved"
	MsgBranchHierarchyRetrieved = "Branch hierarchy successfully retrieved"
)

// Pesan unit
const (
	MsgUnitRetrieved          = "Unit successfully retrieved"
	MsgUnitsRetrieved         = "Units list successfully retrieved"
	MsgUnitCreated            = "Unit successfully created"
	MsgUnitUpdated            = "Unit successfully updated"
	MsgUnitDeleted            = "Unit successfully deleted"
	MsgUnitNotFound           = "Unit not found"
	MsgUnitCodeExists         = "Unit code already exists in this branch"
	MsgUnitHasChildren        = "Cannot delete unit that has child units"
	MsgUnitHasUsers           = "Cannot delete unit with assigned users"
	MsgUnitHierarchyRetrieved = "Unit hierarchy successfully retrieved"
	MsgUnitStatsRetrieved     = "Unit statistics successfully retrieved"
	MsgUnitRoleAssigned       = "Role successfully assigned to unit"
	MsgUnitRoleRemoved        = "Role successfully removed from unit"
	MsgUnitRolesRetrieved     = "Unit roles successfully retrieved"
	MsgUnitPermissionsUpdated = "Unit permissions successfully updated"
	MsgPermissionsCopied      = "Permissions successfully copied between units"
	MsgEffectivePermissions   = "Effective permissions successfully retrieved"
)

// Pesan audit
const (
	MsgAuditLogRetrieved   = "Audit log successfully retrieved"
	MsgAuditLogsRetrieved  = "Audit logs list successfully retrieved"
	MsgAuditLogCreated     = "Audit log successfully created"
	MsgAuditStatsRetrieved = "Audit statistics successfully retrieved"
	MsgAuditLogsCleanedUp  = "Old audit logs successfully cleaned up"
)

// Pesan langganan
const (
	MsgSubscriptionPlanRetrieved  = "Subscription plan successfully retrieved"
	MsgSubscriptionPlansRetrieved = "Subscription plans list successfully retrieved"
	MsgSubscriptionPlanCreated    = "Subscription plan successfully created"
	MsgSubscriptionPlanUpdated    = "Subscription plan successfully updated"
	MsgSubscriptionPlanDeleted    = "Subscription plan successfully deleted"
	MsgSubscriptionRetrieved      = "Subscription successfully retrieved"
	MsgSubscriptionsRetrieved     = "Subscriptions list successfully retrieved"
	MsgSubscriptionCreated        = "Subscription successfully created"
	MsgSubscriptionUpdated        = "Subscription successfully updated"
	MsgSubscriptionRenewed        = "Subscription successfully renewed"
	MsgSubscriptionCancelled      = "Subscription successfully cancelled"
	MsgSubscriptionNotFound       = "Subscription not found"
	MsgSubscriptionPlanNotFound   = "Subscription plan not found"
)

// Konstanta validasi
const (
	MinNameLength     = 2
	MaxNameLength     = 100
	MinCodeLength     = 2
	MaxCodeLength     = 20
	MinPasswordLength = 6
	MsgInvalidID      = "Invalid ID format"
	MsgDataRetrieved  = "Data successfully retrieved"
	MsgDataCreated    = "Data successfully created"
	MsgDataUpdated    = "Data successfully updated"
	MsgDataDeleted    = "Data successfully deleted"
)

// Konstanta database
const (
	DefaultLimit  = 10
	MaxLimit      = 100
	DefaultOffset = 0
)

// Parameter query
const (
	QueryNested           = "nested"
	QueryIncludeHierarchy = "include_hierarchy"
	QueryLimit            = "limit"
	QueryOffset           = "offset"
	QuerySearch           = "search"
	QueryIsActive         = "is_active"
	QueryCompanyID        = "company_id"
	QueryUserID           = "user_id"
	QueryCategory         = "category"
	QueryStatus           = "status"
	QueryDateFrom         = "date_from"
	QueryDateTo           = "date_to"
)

// Tier langganan
const (
	TierBasic        = "basic"
	TierProfessional = "professional"
	TierEnterprise   = "enterprise"
)

// Status langganan
const (
	StatusActive    = "active"
	StatusInactive  = "inactive"
	StatusCancelled = "cancelled"
	StatusExpired   = "expired"
)

// Status pembayaran
const (
	PaymentPending = "pending"
	PaymentPaid    = "paid"
	PaymentFailed  = "failed"
)

// Siklus penagihan
const (
	BillingMonthly = "monthly"
	BillingYearly  = "yearly"
)

// Status audit
const (
	AuditSuccess = "success"
	AuditError   = "error"
	AuditWarning = "warning"
)
