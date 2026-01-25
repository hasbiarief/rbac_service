package constants

// Package constants provides shared constants across all modules
// in the module-based architecture. This is one of the few exceptions
// where cross-module sharing is allowed, as constants are immutable
// and don't create coupling between modules.

// HTTP Status Messages
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

// Auth Module Messages
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

// User Module Messages
const (
	MsgUserRetrieved      = "User successfully retrieved"
	MsgUsersRetrieved     = "Users list successfully retrieved"
	MsgUserCreated        = "User successfully created"
	MsgUserUpdated        = "User successfully updated"
	MsgUserDeleted        = "User successfully deleted"
	MsgUserNotFound       = "User not found"
	MsgEmailAlreadyExists = "Email already exists"
)

// Company Module Messages
const (
	MsgCompanyRetrieved   = "Company successfully retrieved"
	MsgCompaniesRetrieved = "Companies list successfully retrieved"
	MsgCompanyCreated     = "Company successfully created"
	MsgCompanyUpdated     = "Company successfully updated"
	MsgCompanyDeleted     = "Company successfully deleted"
	MsgCompanyNotFound    = "Company not found"
	MsgCompanyCodeExists  = "Company code already exists"
)

// Role Module Messages
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

// Module Module Messages
const (
	MsgModuleRetrieved      = "Module successfully retrieved"
	MsgModulesRetrieved     = "Modules list successfully retrieved"
	MsgModuleCreated        = "Module successfully created"
	MsgModuleUpdated        = "Module successfully updated"
	MsgModuleDeleted        = "Module successfully deleted"
	MsgModuleNotFound       = "Module not found"
	MsgModuleURLExists      = "Module URL already exists"
	MsgUserModulesRetrieved = "User modules successfully retrieved"
)

// Branch Module Messages
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

// Unit Module Messages
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

// Audit Module Messages
const (
	MsgAuditLogRetrieved   = "Audit log successfully retrieved"
	MsgAuditLogsRetrieved  = "Audit logs list successfully retrieved"
	MsgAuditLogCreated     = "Audit log successfully created"
	MsgAuditStatsRetrieved = "Audit statistics successfully retrieved"
	MsgAuditLogsCleanedUp  = "Old audit logs successfully cleaned up"
)

// Subscription Module Messages
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

// Validation Constants
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

// Database Constants
const (
	DefaultLimit  = 10
	MaxLimit      = 100
	DefaultOffset = 0
)

// Query Parameters
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

// Subscription Tiers
const (
	TierBasic        = "basic"
	TierProfessional = "professional"
	TierEnterprise   = "enterprise"
)

// Subscription Status
const (
	StatusActive    = "active"
	StatusInactive  = "inactive"
	StatusCancelled = "cancelled"
	StatusExpired   = "expired"
)

// Payment Status
const (
	PaymentPending = "pending"
	PaymentPaid    = "paid"
	PaymentFailed  = "failed"
)

// Billing Cycles
const (
	BillingMonthly = "monthly"
	BillingYearly  = "yearly"
)

// Audit Status
const (
	AuditSuccess = "success"
	AuditError   = "error"
	AuditWarning = "warning"
)

// API Documentation Module IDs
const (
	ModuleAPIDocumentation   = 139 // Main API Documentation module
	ModuleAPIDocCollections  = 140 // Collections Management sub-module
	ModuleAPIDocEndpoints    = 141 // Endpoints Management sub-module
	ModuleAPIDocEnvironments = 142 // Environments Management sub-module
	ModuleAPIDocExport       = 143 // Export Documentation sub-module
)

// API Documentation Module Messages
const (
	MsgCollectionRetrieved     = "Collection successfully retrieved"
	MsgCollectionsRetrieved    = "Collections list successfully retrieved"
	MsgCollectionCreated       = "Collection successfully created"
	MsgCollectionUpdated       = "Collection successfully updated"
	MsgCollectionDeleted       = "Collection successfully deleted"
	MsgCollectionNotFound      = "Collection not found"
	MsgFolderRetrieved         = "Folder successfully retrieved"
	MsgFoldersRetrieved        = "Folders list successfully retrieved"
	MsgFolderCreated           = "Folder successfully created"
	MsgFolderUpdated           = "Folder successfully updated"
	MsgFolderDeleted           = "Folder successfully deleted"
	MsgFolderNotFound          = "Folder not found"
	MsgEndpointRetrieved       = "Endpoint successfully retrieved"
	MsgEndpointsRetrieved      = "Endpoints list successfully retrieved"
	MsgEndpointCreated         = "Endpoint successfully created"
	MsgEndpointUpdated         = "Endpoint successfully updated"
	MsgEndpointDeleted         = "Endpoint successfully deleted"
	MsgEndpointNotFound        = "Endpoint not found"
	MsgEnvironmentRetrieved    = "Environment successfully retrieved"
	MsgEnvironmentsRetrieved   = "Environments list successfully retrieved"
	MsgEnvironmentCreated      = "Environment successfully created"
	MsgEnvironmentUpdated      = "Environment successfully updated"
	MsgEnvironmentDeleted      = "Environment successfully deleted"
	MsgEnvironmentNotFound     = "Environment not found"
	MsgVariableRetrieved       = "Environment variable successfully retrieved"
	MsgVariablesRetrieved      = "Environment variables list successfully retrieved"
	MsgVariableCreated         = "Environment variable successfully created"
	MsgVariableUpdated         = "Environment variable successfully updated"
	MsgVariableDeleted         = "Environment variable successfully deleted"
	MsgVariableNotFound        = "Environment variable not found"
	MsgExportGenerated         = "Export successfully generated"
	MsgInsufficientPermissions = "Insufficient permissions for this operation"
)
