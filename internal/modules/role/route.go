package role

import (
	"gin-scalable-api/internal/constants"
	"gin-scalable-api/middleware"
	"gin-scalable-api/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler struct
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Handler methods
func (h *Handler) GetRoles(c *gin.Context) {
	var req RoleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.service.GetRoles(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get roles", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgRolesRetrieved, result)
}

func (h *Handler) GetRoleByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid role ID")
		return
	}

	result, err := h.service.GetRoleByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgRoleNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgRoleRetrieved, result)
}

func (h *Handler) CreateRole(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	createReq, ok := validatedBody.(*CreateRoleRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.CreateRole(createReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create role", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgRoleCreated, result)
}

func (h *Handler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid role ID")
		return
	}

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	updateReq, ok := validatedBody.(*UpdateRoleRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.UpdateRole(id, updateReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgRoleUpdated, result)
}

func (h *Handler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid role ID")
		return
	}

	if err := h.service.DeleteRole(id); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgRoleDeleted, nil)
}

func (h *Handler) AssignUserRole(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	assignReq, ok := validatedBody.(*AssignRoleRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.AssignRoleToUser(assignReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to assign user role", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgRoleAssigned, result)
}

func (h *Handler) BulkAssignUserRole(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	bulkAssignReq, ok := validatedBody.(*BulkAssignRoleRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	results, err := h.service.BulkAssignRoleToUsers(bulkAssignReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to bulk assign user roles", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Roles successfully assigned to users", map[string]interface{}{
		"assignments": results,
		"total":       len(results),
	})
}

func (h *Handler) UpdateRoleModules(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("roleId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid role ID")
		return
	}

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	updateReq, ok := validatedBody.(*UpdateRolePermissionsRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	if err := h.service.UpdateRolePermissions(roleID, updateReq); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgPermissionsUpdated, nil)
}

func (h *Handler) RemoveUserRole(c *gin.Context) {
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

	if err := h.service.RemoveRoleFromUser(userID, roleID, companyID); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgRoleUnassigned, nil)
}

func (h *Handler) GetUsersByRole(c *gin.Context) {
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

	result, err := h.service.GetUsersByRole(roleID, limit)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgUsersRetrieved, result)
}

func (h *Handler) GetUserRoles(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	result, err := h.service.GetUserRoles(userID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgRolesRetrieved, result)
}

func (h *Handler) GetUserAccessSummary(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	result, err := h.service.GetUserAccessSummary(userID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User access summary successfully retrieved", result)
}

// Route registration
func RegisterRoutes(api *gin.RouterGroup, handler *Handler) {
	roleManagement := api.Group("/role-management")
	{
		// POST /api/v1/role-management/assign-user-role - Assign role to user
		roleManagement.POST("/assign-user-role",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &AssignRoleRequest{},
			}),
			handler.AssignUserRole,
		)

		// POST /api/v1/role-management/bulk-assign-roles - Bulk assign roles to users
		roleManagement.POST("/bulk-assign-roles",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &BulkAssignRoleRequest{},
			}),
			handler.BulkAssignUserRole,
		)

		// PUT /api/v1/role-management/role/:roleId/modules - Update role permissions
		roleManagement.PUT("/role/:roleId/modules",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &UpdateRolePermissionsRequest{},
			}),
			handler.UpdateRoleModules,
		)

		// DELETE /api/v1/role-management/user/:userId/role/:roleId - Remove role from user
		roleManagement.DELETE("/user/:userId/role/:roleId", handler.RemoveUserRole)

		// GET /api/v1/role-management/role/:roleId/users - Get users by role
		roleManagement.GET("/role/:roleId/users", handler.GetUsersByRole)

		// GET /api/v1/role-management/user/:userId/roles - Get user roles
		roleManagement.GET("/user/:userId/roles", handler.GetUserRoles)

		// GET /api/v1/role-management/user/:userId/access-summary - Get user access summary
		roleManagement.GET("/user/:userId/access-summary", handler.GetUserAccessSummary)
	}

	roles := api.Group("/roles")
	{
		// GET /api/v1/roles - Get all roles with optional filters
		roles.GET("", handler.GetRoles)

		// GET /api/v1/roles/:id - Get role by ID
		roles.GET("/:id", handler.GetRoleByID)

		// POST /api/v1/roles - Create new role
		roles.POST("",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &CreateRoleRequest{},
			}),
			handler.CreateRole,
		)

		// PUT /api/v1/roles/:id - Update role by ID
		roles.PUT("/:id",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &UpdateRoleRequest{},
			}),
			handler.UpdateRole,
		)

		// DELETE /api/v1/roles/:id - Delete role by ID
		roles.DELETE("/:id", handler.DeleteRole)
	}
}
