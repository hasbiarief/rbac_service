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

// CopyUnitRolePermissions - More flexible version using unit_role_id directly
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
