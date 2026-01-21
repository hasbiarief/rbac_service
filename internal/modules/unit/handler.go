package unit

import (
	"net/http"
	"strconv"

	"gin-scalable-api/internal/constants"
	"gin-scalable-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetUnits(c *gin.Context) {
	var req UnitListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request parameters", err.Error())
		return
	}

	result, err := h.service.GetUnits(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get units", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgDataRetrieved, result)
}

func (h *Handler) GetUnitByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, constants.MsgInvalidID, "Invalid unit ID")
		return
	}

	result, err := h.service.GetUnitByID(id)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get unit", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgDataRetrieved, result)
}

func (h *Handler) GetUnitHierarchy(c *gin.Context) {
	branchID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, constants.MsgInvalidID, "Invalid branch ID")
		return
	}

	result, err := h.service.GetUnitHierarchy(branchID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get unit hierarchy", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgDataRetrieved, result)
}

func (h *Handler) GetUnitWithStats(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, constants.MsgInvalidID, "Invalid unit ID")
		return
	}

	result, err := h.service.GetUnitWithStats(id)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get unit stats", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgDataRetrieved, result)
}

func (h *Handler) CreateUnit(c *gin.Context) {
	var req CreateUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	result, err := h.service.CreateUnit(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create unit", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgDataCreated, result)
}

func (h *Handler) UpdateUnit(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, constants.MsgInvalidID, "Invalid unit ID")
		return
	}

	var req UpdateUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	result, err := h.service.UpdateUnit(id, &req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to update unit", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgDataUpdated, result)
}

func (h *Handler) DeleteUnit(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, constants.MsgInvalidID, "Invalid unit ID")
		return
	}

	if err := h.service.DeleteUnit(id); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to delete unit", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgDataDeleted, nil)
}

func (h *Handler) AssignRoleToUnit(c *gin.Context) {
	unitID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, constants.MsgInvalidID, "Invalid unit ID")
		return
	}

	roleID, err := strconv.ParseInt(c.Param("role_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, constants.MsgInvalidID, "Invalid role ID")
		return
	}

	if err := h.service.AssignRoleToUnit(unitID, roleID); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to assign role", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Role successfully assigned to unit", nil)
}

func (h *Handler) RemoveRoleFromUnit(c *gin.Context) {
	unitID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, constants.MsgInvalidID, "Invalid unit ID")
		return
	}

	roleID, err := strconv.ParseInt(c.Param("role_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, constants.MsgInvalidID, "Invalid role ID")
		return
	}

	if err := h.service.RemoveRoleFromUnit(unitID, roleID); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to remove role", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Role successfully removed from unit", nil)
}

func (h *Handler) GetUnitRoles(c *gin.Context) {
	unitID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, constants.MsgInvalidID, "Invalid unit ID")
		return
	}

	result, err := h.service.GetUnitRoles(unitID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get unit roles", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgDataRetrieved, result)
}

func (h *Handler) GetUnitPermissions(c *gin.Context) {
	unitID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, constants.MsgInvalidID, "Invalid unit ID")
		return
	}

	roleID, err := strconv.ParseInt(c.Param("role_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, constants.MsgInvalidID, "Invalid role ID")
		return
	}

	result, err := h.service.GetUnitPermissions(unitID, roleID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get permissions", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgDataRetrieved, result)
}

func (h *Handler) UpdateUnitPermissions(c *gin.Context) {
	unitRoleID, err := strconv.ParseInt(c.Param("unit_role_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, constants.MsgInvalidID, "Invalid unit role ID")
		return
	}

	var req BulkUpdateUnitRoleModulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	if err := h.service.UpdateUnitPermissions(unitRoleID, &req); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to update permissions", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Permissions successfully updated", nil)
}

func (h *Handler) CopyPermissions(c *gin.Context) {
	var req CopyUnitPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	if err := h.service.CopyPermissions(&req); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to copy permissions", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Permissions successfully copied", nil)
}

// CopyUnitRolePermissions - More flexible version using unit_role_id directly
func (h *Handler) CopyUnitRolePermissions(c *gin.Context) {
	var req CopyUnitRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	if err := h.service.CopyUnitRolePermissions(&req); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to copy permissions", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Permissions successfully copied", nil)
}

func (h *Handler) GetUserEffectivePermissions(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, constants.MsgInvalidID, "Invalid user ID")
		return
	}

	result, err := h.service.GetUserEffectivePermissions(userID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get effective permissions", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgDataRetrieved, result)
}

// GetUnitRoleInfo - Helper endpoint to get unit_role_id information
func (h *Handler) GetUnitRoleInfo(c *gin.Context) {
	unitID, err := strconv.ParseInt(c.Query("unit_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid unit_id parameter", err.Error())
		return
	}

	result, err := h.service.GetUnitRoleInfo(unitID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get unit role info", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Unit role information retrieved", result)
}
