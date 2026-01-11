package handlers

import (
	"gin-scalable-api/internal/constants"
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/pkg/rbac"
	"gin-scalable-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

// UnitContextHandler handles unit context related endpoints
type UnitContextHandler struct {
	responseHelper   *utils.ResponseHelper
	validationHelper *utils.ValidationHelper
}

// NewUnitContextHandler creates a new unit context handler
func NewUnitContextHandler(
	responseHelper *utils.ResponseHelper,
	validationHelper *utils.ValidationHelper,
) *UnitContextHandler {
	return &UnitContextHandler{
		responseHelper:   responseHelper,
		validationHelper: validationHelper,
	}
}

// GetMyUnitContext godoc
// @Summary Get my unit context
// @Description Retrieve current user's unit context and permissions
// @Tags unit-context
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponse{data=dto.UnitContextInfo}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/my-unit-context [get]
func (h *UnitContextHandler) GetMyUnitContext(c *gin.Context) {
	// Get unit permissions from context (set by middleware)
	unitPermissions, exists := c.Get("unit_permissions")
	if !exists {
		h.responseHelper.Unauthorized(c, "Konteks unit tidak tersedia")
		return
	}

	permissions, ok := unitPermissions.(*rbac.UnitUserPermissions)
	if !ok {
		h.responseHelper.Unauthorized(c, "Konteks unit tidak valid")
		return
	}

	// Build unit context response
	unitContext := &dto.UnitContextInfo{
		CompanyID:      permissions.CompanyID,
		BranchID:       permissions.BranchID,
		UnitID:         permissions.UnitID,
		EffectiveUnits: permissions.EffectiveUnits,
		UnitRoles:      permissions.UnitRoles,
		AdminLevels: dto.AdminLevels{
			IsUnitAdmin:    permissions.IsUnitAdmin,
			IsBranchAdmin:  permissions.IsBranchAdmin,
			IsCompanyAdmin: permissions.IsCompanyAdmin,
		},
		Permissions: make(map[int64]dto.PermissionInfo),
	}

	// Add company/branch/unit names
	if len(permissions.UnitRoles) > 0 {
		unitContext.CompanyName = permissions.UnitRoles[0].CompanyName
		unitContext.BranchName = permissions.UnitRoles[0].BranchName
		unitContext.UnitName = permissions.UnitRoles[0].UnitName
	}

	// Convert module permissions
	for moduleID, modulePerm := range permissions.Modules {
		unitContext.Permissions[moduleID] = dto.PermissionInfo{
			ModuleID:     modulePerm.ModuleID,
			CanRead:      modulePerm.CanRead,
			CanWrite:     modulePerm.CanWrite,
			CanDelete:    modulePerm.CanDelete,
			CanApprove:   modulePerm.CanApprove,
			GrantedBy:    modulePerm.GrantedBy,
			HighestLevel: modulePerm.HighestLevel,
		}
	}

	h.responseHelper.Success(c, constants.MsgDataRetrieved, unitContext)
}

// GetMyUnitPermissions godoc
// @Summary Get my unit permissions
// @Description Retrieve current user's effective unit permissions
// @Tags unit-context
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponse{data=dto.MyUnitPermissionsResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/my-unit-permissions [get]
func (h *UnitContextHandler) GetMyUnitPermissions(c *gin.Context) {
	// Get unit permissions from context (set by middleware)
	unitPermissions, exists := c.Get("unit_permissions")
	if !exists {
		h.responseHelper.Unauthorized(c, "Konteks unit tidak tersedia")
		return
	}

	// Get additional context information
	companyID, _ := c.Get("company_id")
	branchID, _ := c.Get("branch_id")
	unitID, _ := c.Get("unit_id")
	effectiveUnits, _ := c.Get("effective_units")
	unitRoles, _ := c.Get("unit_roles")

	result := &dto.MyUnitPermissionsResponse{
		Permissions:    unitPermissions,
		CompanyID:      companyID.(int64),
		EffectiveUnits: effectiveUnits.([]int64),
		UnitRoles:      unitRoles,
	}

	if branchID != nil {
		result.BranchID = branchID.(*int64)
	}
	if unitID != nil {
		result.UnitID = unitID.(*int64)
	}

	h.responseHelper.Success(c, constants.MsgDataRetrieved, result)
}
