package handlers

import (
	"gin-scalable-api/internal/constants"
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/interfaces"
	"strconv"

	"gin-scalable-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService interfaces.RoleServiceInterface
}

func NewRoleHandler(roleService interfaces.RoleServiceInterface) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// Basic Role Management
func (h *RoleHandler) GetRoles(c *gin.Context) {
	var req dto.RoleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.roleService.GetRoles(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get roles", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgRolesRetrieved, result)
}

func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid role ID")
		return
	}

	result, err := h.roleService.GetRoleByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgRoleNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgRoleRetrieved, result)
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert ke DTO yang diharapkan
	createReq, ok := validatedBody.(*dto.CreateRoleRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.roleService.CreateRole(createReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create role", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgRoleCreated, result)
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid role ID")
		return
	}

	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert ke DTO yang diharapkan
	updateReq, ok := validatedBody.(*dto.UpdateRoleRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.roleService.UpdateRole(id, updateReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgRoleUpdated, result)
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid role ID")
		return
	}

	if err := h.roleService.DeleteRole(id); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgRoleDeleted, nil)
}

// Advanced Role Management System
func (h *RoleHandler) AssignUserRole(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert ke DTO yang diharapkan
	assignReq, ok := validatedBody.(*dto.AssignRoleRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.roleService.AssignRoleToUser(assignReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to assign user role", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgRoleAssigned, result)
}

func (h *RoleHandler) UpdateRoleModules(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("roleId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid role ID")
		return
	}

	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert ke DTO yang diharapkan
	updateReq, ok := validatedBody.(*dto.UpdateRolePermissionsRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	if err := h.roleService.UpdateRolePermissions(roleID, updateReq); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgPermissionsUpdated, nil)
}

func (h *RoleHandler) RemoveUserRole(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	roleID, err := strconv.ParseInt(c.Param("roleId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid role ID")
		return
	}

	// Get company ID from query parameter or request body
	companyIDStr := c.Query("company_id")
	if companyIDStr == "" {
		response.Error(c, http.StatusBadRequest, "Bad request", "Company ID is required")
		return
	}

	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid company ID")
		return
	}

	if err := h.roleService.RemoveRoleFromUser(userID, roleID, companyID); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgRoleUnassigned, nil)
}

func (h *RoleHandler) GetUsersByRole(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("roleId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid role ID")
		return
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	result, err := h.roleService.GetUsersByRole(roleID, limit)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgUsersRetrieved, result)
}

func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	result, err := h.roleService.GetUserRoles(userID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgRolesRetrieved, result)
}

func (h *RoleHandler) GetUserAccessSummary(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	result, err := h.roleService.GetUserAccessSummary(userID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User access summary successfully retrieved", result)
}
