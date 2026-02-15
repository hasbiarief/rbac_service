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

// @Summary      Get all roles
// @Description  Mendapatkan daftar semua role dengan filter opsional dan pagination
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        limit      query     int     false  "Limit jumlah data"
// @Param        offset     query     int     false  "Offset data"
// @Param        search     query     string  false  "Search by name"
// @Param        is_active  query     bool    false  "Filter by active status"
// @Success      200        {object}  response.Response{data=role.RoleListResponse}  "Roles berhasil diambil"
// @Failure      400        {object}  response.Response  "Bad request"
// @Failure      500        {object}  response.Response  "Internal server error"
// @Router       /api/v1/roles [get]
// @Security     BearerAuth
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

// @Summary      Get role by ID
// @Description  Mendapatkan detail role berdasarkan ID
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Role ID"
// @Success      200  {object}  response.Response{data=role.RoleResponse}  "Role berhasil diambil"
// @Failure      400  {object}  response.Response  "Bad request - Invalid role ID"
// @Failure      404  {object}  response.Response  "Role tidak ditemukan"
// @Router       /api/v1/roles/{id} [get]
// @Security     BearerAuth
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

// @Summary      Get role with permissions
// @Description  Mendapatkan detail role dengan daftar permissions/modules
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Role ID"
// @Success      200  {object}  response.Response{data=role.RoleWithPermissionsResponse}  "Role with permissions berhasil diambil"
// @Failure      400  {object}  response.Response  "Bad request - Invalid role ID"
// @Failure      404  {object}  response.Response  "Role tidak ditemukan"
// @Router       /api/v1/roles/{id}/permissions [get]
// @Security     BearerAuth
func (h *Handler) GetRoleWithPermissions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid role ID")
		return
	}

	result, err := h.service.GetRoleWithPermissions(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgRoleNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Role with permissions successfully retrieved", result)
}

// @Summary      Create new role
// @Description  Membuat role baru
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        role  body      role.CreateRoleRequest  true  "Role data"
// @Success      201   {object}  response.Response{data=role.RoleResponse}  "Role berhasil dibuat"
// @Failure      400   {object}  response.Response  "Bad request - validation failed"
// @Failure      409   {object}  response.Response  "Conflict - role name sudah ada"
// @Failure      500   {object}  response.Response  "Internal server error"
// @Router       /api/v1/roles [post]
// @Security     BearerAuth
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

// @Summary      Update role
// @Description  Memperbarui informasi role
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        id    path      int                     true  "Role ID"
// @Param        role  body      role.UpdateRoleRequest  true  "Role data yang akan diupdate"
// @Success      200   {object}  response.Response{data=role.RoleResponse}  "Role berhasil diupdate"
// @Failure      400   {object}  response.Response  "Bad request - Invalid role ID atau validation failed"
// @Failure      404   {object}  response.Response  "Role tidak ditemukan"
// @Failure      500   {object}  response.Response  "Internal server error"
// @Router       /api/v1/roles/{id} [put]
// @Security     BearerAuth
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

// @Summary      Delete role
// @Description  Menghapus role berdasarkan ID
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Role ID"
// @Success      200  {object}  response.Response  "Role berhasil dihapus"
// @Failure      400  {object}  response.Response  "Bad request - Invalid role ID"
// @Failure      404  {object}  response.Response  "Role tidak ditemukan"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Router       /api/v1/roles/{id} [delete]
// @Security     BearerAuth
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

// @Summary      Assign role to user
// @Description  Menugaskan role ke user dengan scope company, branch, atau unit
// @Tags         Role Management
// @Accept       json
// @Produce      json
// @Param        assignment  body      role.AssignRoleRequest  true  "Role assignment data"
// @Success      200         {object}  response.Response{data=role.UserRoleAssignmentResponse}  "Role berhasil ditugaskan"
// @Failure      400         {object}  response.Response  "Bad request - validation failed"
// @Failure      404         {object}  response.Response  "User atau role tidak ditemukan"
// @Failure      500         {object}  response.Response  "Internal server error"
// @Router       /api/v1/role-management/assign-user-role [post]
// @Security     BearerAuth
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

// @Summary      Bulk assign roles to users
// @Description  Menugaskan role ke multiple users sekaligus
// @Tags         Role Management
// @Accept       json
// @Produce      json
// @Param        bulk_assignment  body      role.BulkAssignRoleRequest  true  "Bulk role assignment data"
// @Success      200              {object}  response.Response  "Roles berhasil ditugaskan ke users"
// @Failure      400              {object}  response.Response  "Bad request - validation failed"
// @Failure      404              {object}  response.Response  "User atau role tidak ditemukan"
// @Failure      500              {object}  response.Response  "Internal server error"
// @Router       /api/v1/role-management/bulk-assign-roles [post]
// @Security     BearerAuth
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

// @Summary      Update role permissions (replace all)
// @Description  Mengganti semua permissions/modules dari role (replace operation)
// @Tags         Role Management
// @Accept       json
// @Produce      json
// @Param        roleId       path      int                                  true  "Role ID"
// @Param        permissions  body      role.UpdateRolePermissionsRequest    true  "Role permissions data"
// @Success      200          {object}  response.Response  "Permissions berhasil diupdate"
// @Failure      400          {object}  response.Response  "Bad request - Invalid role ID atau validation failed"
// @Failure      404          {object}  response.Response  "Role tidak ditemukan"
// @Failure      500          {object}  response.Response  "Internal server error"
// @Router       /api/v1/role-management/role/{roleId}/modules [put]
// @Security     BearerAuth
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

// @Summary      Remove role from user
// @Description  Menghapus role assignment dari user
// @Tags         Role Management
// @Accept       json
// @Produce      json
// @Param        userId      path      int    true  "User ID"
// @Param        roleId      path      int    true  "Role ID"
// @Param        company_id  query     int    true  "Company ID"
// @Success      200         {object}  response.Response  "Role berhasil dihapus dari user"
// @Failure      400         {object}  response.Response  "Bad request - Invalid ID atau company_id diperlukan"
// @Failure      404         {object}  response.Response  "User role assignment tidak ditemukan"
// @Failure      500         {object}  response.Response  "Internal server error"
// @Router       /api/v1/role-management/user/{userId}/role/{roleId} [delete]
// @Security     BearerAuth
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

// @Summary      Get users by role
// @Description  Mendapatkan daftar users yang memiliki role tertentu
// @Tags         Role Management
// @Accept       json
// @Produce      json
// @Param        roleId  path      int  true   "Role ID"
// @Param        limit   query     int  false  "Limit jumlah data (default: 10)"
// @Success      200     {object}  response.Response  "Users berhasil diambil"
// @Failure      400     {object}  response.Response  "Bad request - Invalid role ID"
// @Failure      500     {object}  response.Response  "Internal server error"
// @Router       /api/v1/role-management/role/{roleId}/users [get]
// @Security     BearerAuth
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

// @Summary      Get user roles
// @Description  Mendapatkan semua role assignments dari user
// @Tags         Role Management
// @Accept       json
// @Produce      json
// @Param        userId  path      int  true  "User ID"
// @Success      200     {object}  response.Response{data=[]role.UserRoleAssignmentResponse}  "User roles berhasil diambil"
// @Failure      400     {object}  response.Response  "Bad request - Invalid user ID"
// @Failure      500     {object}  response.Response  "Internal server error"
// @Router       /api/v1/role-management/user/{userId}/roles [get]
// @Security     BearerAuth
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

// @Summary      Get user access summary
// @Description  Mendapatkan ringkasan akses user termasuk semua roles dan permissions
// @Tags         Role Management
// @Accept       json
// @Produce      json
// @Param        userId  path      int  true  "User ID"
// @Success      200     {object}  response.Response  "User access summary berhasil diambil"
// @Failure      400     {object}  response.Response  "Bad request - Invalid user ID"
// @Failure      500     {object}  response.Response  "Internal server error"
// @Router       /api/v1/role-management/user/{userId}/access-summary [get]
// @Security     BearerAuth
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

		// PUT /api/v1/role-management/role/:roleId/modules - Update role permissions (REPLACE ALL)
		roleManagement.PUT("/role/:roleId/modules",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &UpdateRolePermissionsRequest{},
			}),
			handler.UpdateRoleModules,
		)

		// POST /api/v1/role-management/role/:roleId/modules - Add modules to role (APPEND)
		roleManagement.POST("/role/:roleId/modules",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &AddRoleModulesRequest{},
			}),
			handler.AddRoleModules,
		)

		// DELETE /api/v1/role-management/role/:roleId/modules - Remove modules from role
		roleManagement.DELETE("/role/:roleId/modules",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &RemoveRoleModulesRequest{},
			}),
			handler.RemoveRoleModules,
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

		// GET /api/v1/roles/:id/permissions - Get role with permissions
		roles.GET("/:id/permissions", handler.GetRoleWithPermissions)

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

// @Summary      Add modules to role
// @Description  Menambahkan modules ke role (append operation)
// @Tags         Role Management
// @Accept       json
// @Produce      json
// @Param        roleId   path      int                           true  "Role ID"
// @Param        modules  body      role.AddRoleModulesRequest    true  "Modules to add"
// @Success      200      {object}  response.Response  "Modules berhasil ditambahkan"
// @Failure      400      {object}  response.Response  "Bad request - Invalid role ID atau validation failed"
// @Failure      404      {object}  response.Response  "Role tidak ditemukan"
// @Failure      500      {object}  response.Response  "Internal server error"
// @Router       /api/v1/role-management/role/{roleId}/modules [post]
// @Security     BearerAuth
func (h *Handler) AddRoleModules(c *gin.Context) {
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

	addReq, ok := validatedBody.(*AddRoleModulesRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	if err := h.service.AddRoleModules(roleID, addReq); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to add modules to role", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Modules successfully added to role", nil)
}

// @Summary      Remove modules from role
// @Description  Menghapus modules dari role
// @Tags         Role Management
// @Accept       json
// @Produce      json
// @Param        roleId   path      int                              true  "Role ID"
// @Param        modules  body      role.RemoveRoleModulesRequest    true  "Module IDs to remove"
// @Success      200      {object}  response.Response  "Modules berhasil dihapus"
// @Failure      400      {object}  response.Response  "Bad request - Invalid role ID atau validation failed"
// @Failure      404      {object}  response.Response  "Role tidak ditemukan"
// @Failure      500      {object}  response.Response  "Internal server error"
// @Router       /api/v1/role-management/role/{roleId}/modules [delete]
// @Security     BearerAuth
func (h *Handler) RemoveRoleModules(c *gin.Context) {
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

	removeReq, ok := validatedBody.(*RemoveRoleModulesRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	if err := h.service.RemoveRoleModules(roleID, removeReq); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to remove modules from role", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Modules successfully removed from role", nil)
}
