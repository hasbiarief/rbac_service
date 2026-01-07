package handlers

import (
	"gin-scalable-api/internal/service"
	"strconv"

	"gin-scalable-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// Basic Role Management
func (h *RoleHandler) GetRoles(c *gin.Context) {
	var req service.RoleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.roleService.GetRoles(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get roles", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid role ID")
		return
	}

	result, err := h.roleService.GetRoleByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert to the expected struct
	req, ok := validatedBody.(*struct {
		Name        string `json:"name" validate:"required,min=2,max=100"`
		Description string `json:"description"`
	})
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	// Convert to service request
	createReq := &service.CreateRoleRequest{
		Name:        req.Name,
		Description: req.Description,
	}

	result, err := h.roleService.CreateRole(createReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create role", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Created successfully", result)
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

	// Type assert to the expected struct
	req, ok := validatedBody.(*struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	})
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	// Convert to service request
	updateReq := &service.UpdateRoleRequest{
		Name:        req.Name,
		Description: req.Description,
	}

	result, err := h.roleService.UpdateRole(id, updateReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
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

	response.Success(c, http.StatusOK, "Role deleted successfully", nil)
}

// Advanced Role Management System
func (h *RoleHandler) AssignUserRole(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert to the expected struct
	req, ok := validatedBody.(*struct {
		UserID    int64  `json:"user_id" validate:"required,min=1"`
		RoleID    int64  `json:"role_id" validate:"required,min=1"`
		CompanyID int64  `json:"company_id" validate:"required,min=1"`
		BranchID  *int64 `json:"branch_id"`
	})
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	// Convert to service request
	assignReq := &service.AssignUserRoleRequest{
		UserID:    req.UserID,
		RoleID:    req.RoleID,
		CompanyID: req.CompanyID,
		BranchID:  req.BranchID,
	}

	if err := h.roleService.AssignUserRole(assignReq); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to assign user role", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User role assigned successfully", nil)
}

func (h *RoleHandler) BulkAssignRoles(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert to the expected struct
	req, ok := validatedBody.(*struct {
		UserIDs   []int64 `json:"user_ids" validate:"required,min=1"`
		RoleID    int64   `json:"role_id" validate:"required,min=1"`
		CompanyID int64   `json:"company_id" validate:"required,min=1"`
		BranchID  *int64  `json:"branch_id"`
	})
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	// Convert to service request
	var assignments []service.AssignUserRoleRequest
	for _, userID := range req.UserIDs {
		assignments = append(assignments, service.AssignUserRoleRequest{
			UserID:    userID,
			RoleID:    req.RoleID,
			CompanyID: req.CompanyID,
			BranchID:  req.BranchID,
		})
	}

	bulkReq := &service.BulkAssignRolesRequest{
		Assignments: assignments,
	}

	if err := h.roleService.BulkAssignRoles(bulkReq); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to assign roles", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Roles assigned successfully", nil)
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

	// Type assert to the expected struct
	req, ok := validatedBody.(*struct {
		Modules []struct {
			ModuleID  int64 `json:"module_id" validate:"required,min=1"`
			CanRead   bool  `json:"can_read"`
			CanWrite  bool  `json:"can_write"`
			CanDelete bool  `json:"can_delete"`
		} `json:"modules" validate:"required,min=1"`
	})
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	// Convert to service request
	var modules []service.RoleModulePermission
	for _, m := range req.Modules {
		modules = append(modules, service.RoleModulePermission{
			ModuleID:  m.ModuleID,
			CanRead:   m.CanRead,
			CanWrite:  m.CanWrite,
			CanDelete: m.CanDelete,
		})
	}

	updateReq := &service.UpdateRoleModulesRequest{
		Modules: modules,
	}

	if err := h.roleService.UpdateRoleModules(roleID, updateReq); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Role modules updated successfully", nil)
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

	if err := h.roleService.RemoveUserRole(userID, roleID); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User role removed successfully", nil)
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

	response.Success(c, http.StatusOK, "Success", result)
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

	response.Success(c, http.StatusOK, "Success", result)
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

	response.Success(c, http.StatusOK, "Success", result)
}
