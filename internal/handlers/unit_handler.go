package handlers

import (
	"gin-scalable-api/internal/constants"
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/interfaces"
	"gin-scalable-api/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UnitHandler struct {
	unitService      interfaces.UnitServiceInterface
	responseHelper   *utils.ResponseHelper
	validationHelper *utils.ValidationHelper
}

func NewUnitHandler(
	unitService interfaces.UnitServiceInterface,
	responseHelper *utils.ResponseHelper,
	validationHelper *utils.ValidationHelper,
) *UnitHandler {
	return &UnitHandler{
		unitService:      unitService,
		responseHelper:   responseHelper,
		validationHelper: validationHelper,
	}
}

// GetUnits godoc
// @Summary Get units with filtering and pagination
// @Description Retrieve units with optional filtering by branch, search, and active status
// @Tags units
// @Accept json
// @Produce json
// @Param branch_id query int false "Branch ID"
// @Param search query string false "Search term"
// @Param is_active query bool false "Active status"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} utils.SuccessResponse{data=dto.UnitListResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /units [get]
func (h *UnitHandler) GetUnits(c *gin.Context) {
	var req dto.UnitListRequest
	if err := h.validationHelper.GetValidatedQuery(c, &req); err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	result, err := h.unitService.GetUnits(&req)
	if err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	h.responseHelper.Success(c, constants.MsgDataRetrieved, result)
}

// GetUnitByID godoc
// @Summary Get unit by ID
// @Description Retrieve a specific unit by its ID
// @Tags units
// @Accept json
// @Produce json
// @Param id path int true "Unit ID"
// @Success 200 {object} utils.SuccessResponse{data=dto.UnitResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /units/{id} [get]
func (h *UnitHandler) GetUnitByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.responseHelper.BadRequest(c, constants.MsgInvalidID)
		return
	}

	result, err := h.unitService.GetUnitByID(id)
	if err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	h.responseHelper.Success(c, constants.MsgDataRetrieved, result)
}

// GetUnitHierarchy godoc
// @Summary Get unit hierarchy for a branch
// @Description Retrieve hierarchical structure of units within a branch
// @Tags units
// @Accept json
// @Produce json
// @Param id path int true "Branch ID"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.UnitHierarchyResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /branches/{id}/units/hierarchy [get]
func (h *UnitHandler) GetUnitHierarchy(c *gin.Context) {
	branchID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.responseHelper.BadRequest(c, constants.MsgInvalidID)
		return
	}

	result, err := h.unitService.GetUnitHierarchy(branchID)
	if err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	h.responseHelper.Success(c, constants.MsgDataRetrieved, result)
}

// GetUnitWithStats godoc
// @Summary Get unit with statistics
// @Description Retrieve unit information along with usage statistics
// @Tags units
// @Accept json
// @Produce json
// @Param id path int true "Unit ID"
// @Success 200 {object} utils.SuccessResponse{data=dto.UnitWithStatsResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /units/{id}/stats [get]
func (h *UnitHandler) GetUnitWithStats(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.responseHelper.BadRequest(c, constants.MsgInvalidID)
		return
	}

	result, err := h.unitService.GetUnitWithStats(id)
	if err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	h.responseHelper.Success(c, constants.MsgDataRetrieved, result)
}

// CreateUnit godoc
// @Summary Create a new unit
// @Description Create a new unit within a branch
// @Tags units
// @Accept json
// @Produce json
// @Param unit body dto.CreateUnitRequest true "Unit data"
// @Success 201 {object} utils.SuccessResponse{data=dto.UnitResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /units [post]
func (h *UnitHandler) CreateUnit(c *gin.Context) {
	var req dto.CreateUnitRequest
	if err := h.validationHelper.GetValidatedBody(c, &req); err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	result, err := h.unitService.CreateUnit(&req)
	if err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	h.responseHelper.Created(c, constants.MsgDataCreated, result)
}

// UpdateUnit godoc
// @Summary Update an existing unit
// @Description Update unit information
// @Tags units
// @Accept json
// @Produce json
// @Param id path int true "Unit ID"
// @Param unit body dto.UpdateUnitRequest true "Updated unit data"
// @Success 200 {object} utils.SuccessResponse{data=dto.UnitResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /units/{id} [put]
func (h *UnitHandler) UpdateUnit(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.responseHelper.BadRequest(c, constants.MsgInvalidID)
		return
	}

	var req dto.UpdateUnitRequest
	if err := h.validationHelper.GetValidatedBody(c, &req); err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	result, err := h.unitService.UpdateUnit(id, &req)
	if err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	h.responseHelper.Success(c, constants.MsgDataUpdated, result)
}

// DeleteUnit godoc
// @Summary Delete a unit
// @Description Soft delete a unit (mark as inactive)
// @Tags units
// @Accept json
// @Produce json
// @Param id path int true "Unit ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /units/{id} [delete]
func (h *UnitHandler) DeleteUnit(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.responseHelper.BadRequest(c, constants.MsgInvalidID)
		return
	}

	if err := h.unitService.DeleteUnit(id); err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	h.responseHelper.Success(c, constants.MsgDataDeleted, nil)
}

// AssignRoleToUnit godoc
// @Summary Assign a role to a unit
// @Description Create a role assignment for a unit
// @Tags units
// @Accept json
// @Produce json
// @Param id path int true "Unit ID"
// @Param role_id path int true "Role ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /units/{id}/roles/{role_id} [post]
func (h *UnitHandler) AssignRoleToUnit(c *gin.Context) {
	unitID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.responseHelper.BadRequest(c, constants.MsgInvalidID)
		return
	}

	roleID, err := strconv.ParseInt(c.Param("role_id"), 10, 64)
	if err != nil {
		h.responseHelper.BadRequest(c, constants.MsgInvalidID)
		return
	}

	if err := h.unitService.AssignRoleToUnit(unitID, roleID); err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	h.responseHelper.Success(c, "Role berhasil ditugaskan ke unit", nil)
}

// RemoveRoleFromUnit godoc
// @Summary Remove a role from a unit
// @Description Remove a role assignment from a unit
// @Tags units
// @Accept json
// @Produce json
// @Param id path int true "Unit ID"
// @Param role_id path int true "Role ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /units/{id}/roles/{role_id} [delete]
func (h *UnitHandler) RemoveRoleFromUnit(c *gin.Context) {
	unitID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.responseHelper.BadRequest(c, constants.MsgInvalidID)
		return
	}

	roleID, err := strconv.ParseInt(c.Param("role_id"), 10, 64)
	if err != nil {
		h.responseHelper.BadRequest(c, constants.MsgInvalidID)
		return
	}

	if err := h.unitService.RemoveRoleFromUnit(unitID, roleID); err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	h.responseHelper.Success(c, "Role berhasil dihapus dari unit", nil)
}

// GetUnitRoles godoc
// @Summary Get roles assigned to a unit
// @Description Retrieve all roles assigned to a specific unit
// @Tags units
// @Accept json
// @Produce json
// @Param id path int true "Unit ID"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.UnitRoleResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /units/{id}/roles [get]
func (h *UnitHandler) GetUnitRoles(c *gin.Context) {
	unitID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.responseHelper.BadRequest(c, constants.MsgInvalidID)
		return
	}

	result, err := h.unitService.GetUnitRoles(unitID)
	if err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	h.responseHelper.Success(c, constants.MsgDataRetrieved, result)
}

// GetUnitPermissions godoc
// @Summary Get permissions for a unit-role combination
// @Description Retrieve permissions for a specific unit and role combination
// @Tags units
// @Accept json
// @Produce json
// @Param id path int true "Unit ID"
// @Param role_id path int true "Role ID"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.UnitRoleModuleResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /units/{id}/roles/{role_id}/permissions [get]
func (h *UnitHandler) GetUnitPermissions(c *gin.Context) {
	unitID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.responseHelper.BadRequest(c, constants.MsgInvalidID)
		return
	}

	roleID, err := strconv.ParseInt(c.Param("role_id"), 10, 64)
	if err != nil {
		h.responseHelper.BadRequest(c, constants.MsgInvalidID)
		return
	}

	result, err := h.unitService.GetUnitPermissions(unitID, roleID)
	if err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	h.responseHelper.Success(c, constants.MsgDataRetrieved, result)
}

// UpdateUnitPermissions godoc
// @Summary Update permissions for a unit role
// @Description Bulk update permissions for a unit-role combination
// @Tags units
// @Accept json
// @Produce json
// @Param unit_role_id path int true "Unit Role ID"
// @Param permissions body dto.BulkUpdateUnitRoleModulesRequest true "Permission updates"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /unit-roles/{unit_role_id}/permissions [put]
func (h *UnitHandler) UpdateUnitPermissions(c *gin.Context) {
	unitRoleID, err := strconv.ParseInt(c.Param("unit_role_id"), 10, 64)
	if err != nil {
		h.responseHelper.BadRequest(c, constants.MsgInvalidID)
		return
	}

	var req dto.BulkUpdateUnitRoleModulesRequest
	if err := h.validationHelper.GetValidatedBody(c, &req); err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	if err := h.unitService.UpdateUnitPermissions(unitRoleID, &req); err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	h.responseHelper.Success(c, "Permissions berhasil diperbarui", nil)
}

// CopyPermissions godoc
// @Summary Copy permissions between units
// @Description Copy permissions from one unit to another for a specific role
// @Tags units
// @Accept json
// @Produce json
// @Param copy_request body dto.CopyUnitPermissionsRequest true "Copy request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /units/copy-permissions [post]
func (h *UnitHandler) CopyPermissions(c *gin.Context) {
	var req dto.CopyUnitPermissionsRequest
	if err := h.validationHelper.GetValidatedBody(c, &req); err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	if err := h.unitService.CopyPermissions(&req); err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	h.responseHelper.Success(c, "Permissions berhasil disalin", nil)
}

// GetUserEffectivePermissions godoc
// @Summary Get effective permissions for a user
// @Description Retrieve all effective permissions for a user across all their unit assignments
// @Tags units
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.UnitRoleModuleResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /users/{id}/effective-permissions [get]
func (h *UnitHandler) GetUserEffectivePermissions(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.responseHelper.BadRequest(c, constants.MsgInvalidID)
		return
	}

	result, err := h.unitService.GetUserEffectivePermissions(userID)
	if err != nil {
		h.responseHelper.HandleError(c, err)
		return
	}

	h.responseHelper.Success(c, constants.MsgDataRetrieved, result)
}
