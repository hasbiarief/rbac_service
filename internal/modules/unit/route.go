package unit

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
	return &Handler{service: service}
}

// Handler methods

// GetUnits godoc
// @Summary      Get all units
// @Description  Mendapatkan daftar semua unit dengan filter opsional dan pagination
// @Tags         Units
// @Accept       json
// @Produce      json
// @Param        limit      query     int     false  "Limit jumlah data"
// @Param        offset     query     int     false  "Offset data"
// @Param        search     query     string  false  "Search by name atau code"
// @Param        branch_id  query     int     false  "Filter by branch ID"
// @Param        is_active  query     bool    false  "Filter by active status"
// @Success      200        {object}  response.Response{data=unit.UnitListResponse}  "Units berhasil diambil"
// @Failure      400        {object}  response.Response  "Bad request"
// @Failure      500        {object}  response.Response  "Internal server error"
// @Router       /api/v1/units [get]
// @Security     BearerAuth
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

// GetUnitByID godoc
// @Summary      Get unit by ID
// @Description  Mendapatkan detail unit berdasarkan ID
// @Tags         Units
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Unit ID"
// @Success      200  {object}  response.Response{data=unit.UnitResponse}  "Unit berhasil diambil"
// @Failure      400  {object}  response.Response  "Bad request - Invalid unit ID"
// @Failure      404  {object}  response.Response  "Unit tidak ditemukan"
// @Router       /api/v1/units/{id} [get]
// @Security     BearerAuth
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

// GetUnitHierarchy godoc
// @Summary      Get unit hierarchy for branch
// @Description  Mendapatkan unit hierarchy (tree structure) untuk branch tertentu
// @Tags         Units
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Branch ID"
// @Success      200  {object}  response.Response{data=[]unit.UnitHierarchyResponse}  "Unit hierarchy berhasil diambil"
// @Failure      400  {object}  response.Response  "Bad request - Invalid branch ID"
// @Failure      404  {object}  response.Response  "Branch tidak ditemukan"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Router       /api/v1/branches/{id}/units/hierarchy [get]
// @Security     BearerAuth
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

// GetUnitWithStats godoc
// @Summary      Get unit with statistics
// @Description  Mendapatkan detail unit dengan statistik (jumlah users, roles, dll)
// @Tags         Units
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Unit ID"
// @Success      200  {object}  response.Response{data=unit.UnitWithStatsResponse}  "Unit stats berhasil diambil"
// @Failure      400  {object}  response.Response  "Bad request - Invalid unit ID"
// @Failure      404  {object}  response.Response  "Unit tidak ditemukan"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Router       /api/v1/units/{id}/stats [get]
// @Security     BearerAuth
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

// CreateUnit godoc
// @Summary      Create new unit
// @Description  Membuat unit baru dengan support untuk hierarchical structure
// @Tags         Units
// @Accept       json
// @Produce      json
// @Param        unit  body      unit.CreateUnitRequest  true  "Unit data"
// @Success      201   {object}  response.Response{data=unit.UnitResponse}  "Unit berhasil dibuat"
// @Failure      400   {object}  response.Response  "Bad request - validation failed"
// @Failure      409   {object}  response.Response  "Conflict - unit code sudah ada"
// @Failure      500   {object}  response.Response  "Internal server error"
// @Router       /api/v1/units [post]
// @Security     BearerAuth
func (h *Handler) CreateUnit(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	req, ok := validatedBody.(*CreateUnitRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "invalid body structure")
		return
	}

	result, err := h.service.CreateUnit(req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create unit", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgDataCreated, result)
}

// UpdateUnit godoc
// @Summary      Update unit
// @Description  Memperbarui informasi unit
// @Tags         Units
// @Accept       json
// @Produce      json
// @Param        id    path      int                     true  "Unit ID"
// @Param        unit  body      unit.UpdateUnitRequest  true  "Unit data yang akan diupdate"
// @Success      200   {object}  response.Response{data=unit.UnitResponse}  "Unit berhasil diupdate"
// @Failure      400   {object}  response.Response  "Bad request - Invalid unit ID atau validation failed"
// @Failure      404   {object}  response.Response  "Unit tidak ditemukan"
// @Failure      500   {object}  response.Response  "Internal server error"
// @Router       /api/v1/units/{id} [put]
// @Security     BearerAuth
func (h *Handler) UpdateUnit(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, constants.MsgInvalidID, "Invalid unit ID")
		return
	}

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	req, ok := validatedBody.(*UpdateUnitRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "invalid body structure")
		return
	}

	result, err := h.service.UpdateUnit(id, req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to update unit", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgDataUpdated, result)
}

// DeleteUnit godoc
// @Summary      Delete unit
// @Description  Menghapus unit berdasarkan ID
// @Tags         Units
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Unit ID"
// @Success      200  {object}  response.Response  "Unit berhasil dihapus"
// @Failure      400  {object}  response.Response  "Bad request - Invalid unit ID"
// @Failure      404  {object}  response.Response  "Unit tidak ditemukan"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Router       /api/v1/units/{id} [delete]
// @Security     BearerAuth
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

// AssignRoleToUnit godoc
// @Summary      Assign role to unit
// @Description  Menugaskan role ke unit
// @Tags         Units
// @Accept       json
// @Produce      json
// @Param        id       path      int  true  "Unit ID"
// @Param        role_id  path      int  true  "Role ID"
// @Success      200      {object}  response.Response  "Role berhasil ditugaskan ke unit"
// @Failure      400      {object}  response.Response  "Bad request - Invalid ID"
// @Failure      404      {object}  response.Response  "Unit atau role tidak ditemukan"
// @Failure      500      {object}  response.Response  "Internal server error"
// @Router       /api/v1/units/{id}/roles/{role_id} [post]
// @Security     BearerAuth
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

// RemoveRoleFromUnit godoc
// @Summary      Remove role from unit
// @Description  Menghapus role dari unit
// @Tags         Units
// @Accept       json
// @Produce      json
// @Param        id       path      int  true  "Unit ID"
// @Param        role_id  path      int  true  "Role ID"
// @Success      200      {object}  response.Response  "Role berhasil dihapus dari unit"
// @Failure      400      {object}  response.Response  "Bad request - Invalid ID"
// @Failure      404      {object}  response.Response  "Unit role assignment tidak ditemukan"
// @Failure      500      {object}  response.Response  "Internal server error"
// @Router       /api/v1/units/{id}/roles/{role_id} [delete]
// @Security     BearerAuth
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

// GetUnitRoles godoc
// @Summary      Get unit roles
// @Description  Mendapatkan daftar roles yang ditugaskan ke unit
// @Tags         Units
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Unit ID"
// @Success      200  {object}  response.Response{data=[]unit.UnitRoleResponse}  "Unit roles berhasil diambil"
// @Failure      400  {object}  response.Response  "Bad request - Invalid unit ID"
// @Failure      404  {object}  response.Response  "Unit tidak ditemukan"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Router       /api/v1/units/{id}/roles [get]
// @Security     BearerAuth
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

// GetUnitPermissions godoc
// @Summary      Get unit permissions for specific role
// @Description  Mendapatkan permissions/modules yang dapat diakses unit untuk role tertentu
// @Tags         Units
// @Accept       json
// @Produce      json
// @Param        id       path      int  true  "Unit ID"
// @Param        role_id  path      int  true  "Role ID"
// @Success      200      {object}  response.Response  "Unit permissions berhasil diambil"
// @Failure      400      {object}  response.Response  "Bad request - Invalid ID"
// @Failure      404      {object}  response.Response  "Unit atau role tidak ditemukan"
// @Failure      500      {object}  response.Response  "Internal server error"
// @Router       /api/v1/units/{id}/roles/{role_id}/permissions [get]
// @Security     BearerAuth
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

// UpdateUnitPermissions godoc
// @Summary      Update unit role permissions
// @Description  Memperbarui permissions/modules untuk unit role
// @Tags         Unit Roles
// @Accept       json
// @Produce      json
// @Param        unit_role_id  path      int                                     true  "Unit Role ID"
// @Param        permissions   body      unit.BulkUpdateUnitRoleModulesRequest   true  "Permissions data"
// @Success      200           {object}  response.Response  "Permissions berhasil diupdate"
// @Failure      400           {object}  response.Response  "Bad request - Invalid unit role ID atau validation failed"
// @Failure      404           {object}  response.Response  "Unit role tidak ditemukan"
// @Failure      500           {object}  response.Response  "Internal server error"
// @Router       /api/v1/unit-roles/{unit_role_id}/permissions [put]
// @Security     BearerAuth
func (h *Handler) UpdateUnitPermissions(c *gin.Context) {
	unitRoleID, err := strconv.ParseInt(c.Param("unit_role_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, constants.MsgInvalidID, "Invalid unit role ID")
		return
	}

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	req, ok := validatedBody.(*BulkUpdateUnitRoleModulesRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "invalid body structure")
		return
	}

	if err := h.service.UpdateUnitPermissions(unitRoleID, req); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to update permissions", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Permissions successfully updated", nil)
}

// CopyPermissions godoc
// @Summary      Copy permissions between units
// @Description  Menyalin permissions dari satu unit ke unit lain
// @Tags         Units
// @Accept       json
// @Produce      json
// @Param        copy  body      unit.CopyUnitPermissionsRequest  true  "Copy permissions request"
// @Success      200   {object}  response.Response  "Permissions berhasil disalin"
// @Failure      400   {object}  response.Response  "Bad request - validation failed"
// @Failure      404   {object}  response.Response  "Unit tidak ditemukan"
// @Failure      500   {object}  response.Response  "Internal server error"
// @Router       /api/v1/units/copy-permissions [post]
// @Security     BearerAuth
func (h *Handler) CopyPermissions(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	req, ok := validatedBody.(*CopyUnitPermissionsRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "invalid body structure")
		return
	}

	if err := h.service.CopyPermissions(req); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to copy permissions", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Permissions successfully copied", nil)
}

// CopyUnitRolePermissions godoc
// @Summary      Copy permissions between unit roles
// @Description  Menyalin permissions dari satu unit role ke unit role lain
// @Tags         Units
// @Accept       json
// @Produce      json
// @Param        copy  body      unit.CopyUnitRolePermissionsRequest  true  "Copy unit role permissions request"
// @Success      200   {object}  response.Response  "Permissions berhasil disalin"
// @Failure      400   {object}  response.Response  "Bad request - validation failed"
// @Failure      404   {object}  response.Response  "Unit role tidak ditemukan"
// @Failure      500   {object}  response.Response  "Internal server error"
// @Router       /api/v1/units/copy-unit-role-permissions [post]
// @Security     BearerAuth
func (h *Handler) CopyUnitRolePermissions(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	req, ok := validatedBody.(*CopyUnitRolePermissionsRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "invalid body structure")
		return
	}

	if err := h.service.CopyUnitRolePermissions(req); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to copy permissions", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Permissions successfully copied", nil)
}

// GetUserEffectivePermissions godoc
// @Summary      Get user effective permissions
// @Description  Mendapatkan effective permissions user (gabungan dari semua roles dan units)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  response.Response  "User effective permissions berhasil diambil"
// @Failure      400  {object}  response.Response  "Bad request - Invalid user ID"
// @Failure      404  {object}  response.Response  "User tidak ditemukan"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Router       /api/v1/users/{id}/effective-permissions [get]
// @Security     BearerAuth
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

// GetUnitRoleInfo godoc
// @Summary      Get unit role information
// @Description  Mendapatkan informasi unit role berdasarkan unit_id
// @Tags         Units
// @Accept       json
// @Produce      json
// @Param        unit_id  query     int  true  "Unit ID"
// @Success      200      {object}  response.Response  "Unit role info berhasil diambil"
// @Failure      400      {object}  response.Response  "Bad request - unit_id diperlukan"
// @Failure      404      {object}  response.Response  "Unit tidak ditemukan"
// @Failure      500      {object}  response.Response  "Internal server error"
// @Router       /api/v1/units/unit-role-info [get]
// @Security     BearerAuth
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

// Route registration
func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	// Unit CRUD routes
	units := router.Group("/units")
	{
		// GET /api/v1/units - Get all units with optional filters
		units.GET("", handler.GetUnits)

		// POST /api/v1/units - Create new unit
		units.POST("",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &CreateUnitRequest{},
			}),
			handler.CreateUnit,
		)

		// GET /api/v1/units/:id - Get unit by ID
		units.GET("/:id", handler.GetUnitByID)

		// PUT /api/v1/units/:id - Update unit by ID
		units.PUT("/:id",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &UpdateUnitRequest{},
			}),
			handler.UpdateUnit,
		)

		// DELETE /api/v1/units/:id - Delete unit by ID
		units.DELETE("/:id", handler.DeleteUnit)

		// GET /api/v1/units/:id/stats - Get unit with statistics
		units.GET("/:id/stats", handler.GetUnitWithStats)

		// Unit role management
		// POST /api/v1/units/:id/roles/:role_id - Assign role to unit
		units.POST("/:id/roles/:role_id", handler.AssignRoleToUnit)

		// DELETE /api/v1/units/:id/roles/:role_id - Remove role from unit
		units.DELETE("/:id/roles/:role_id", handler.RemoveRoleFromUnit)

		// GET /api/v1/units/:id/roles - Get unit roles
		units.GET("/:id/roles", handler.GetUnitRoles)

		// GET /api/v1/units/:id/roles/:role_id/permissions - Get unit permissions for specific role
		units.GET("/:id/roles/:role_id/permissions", handler.GetUnitPermissions)
	}

	// Branch unit hierarchy
	branches := router.Group("/branches")
	{
		// GET /api/v1/branches/:id/units/hierarchy - Get unit hierarchy for branch
		branches.GET("/:id/units/hierarchy", handler.GetUnitHierarchy)
	}

	// Unit role permissions
	unitRoles := router.Group("/unit-roles")
	{
		// PUT /api/v1/unit-roles/:unit_role_id/permissions - Update unit role permissions
		unitRoles.PUT("/:unit_role_id/permissions",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &BulkUpdateUnitRoleModulesRequest{},
			}),
			handler.UpdateUnitPermissions,
		)
	}

	// Permission management
	// POST /api/v1/units/copy-permissions - Copy permissions between units
	router.POST("/units/copy-permissions",
		middleware.ValidateRequest(middleware.ValidationRules{
			Body: &CopyUnitPermissionsRequest{},
		}),
		handler.CopyPermissions,
	)

	// POST /api/v1/units/copy-unit-role-permissions - Copy permissions between unit roles
	router.POST("/units/copy-unit-role-permissions",
		middleware.ValidateRequest(middleware.ValidationRules{
			Body: &CopyUnitRolePermissionsRequest{},
		}),
		handler.CopyUnitRolePermissions,
	)

	// GET /api/v1/units/unit-role-info - Get unit role information
	router.GET("/units/unit-role-info", handler.GetUnitRoleInfo)

	// User effective permissions
	users := router.Group("/users")
	{
		// GET /api/v1/users/:id/effective-permissions - Get user effective permissions
		users.GET("/:id/effective-permissions", handler.GetUserEffectivePermissions)
	}
}
