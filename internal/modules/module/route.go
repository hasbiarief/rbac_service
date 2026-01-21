package module

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
func (h *Handler) GetModules(c *gin.Context) {
	var req ModuleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Internal error", "Invalid user ID type")
		return
	}

	result, err := h.service.GetModulesFiltered(userIDInt64, &req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get modules", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgModulesRetrieved, result)
}

func (h *Handler) GetModuleByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid module ID")
		return
	}

	result, err := h.service.GetModuleByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgModuleNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgModuleRetrieved, result)
}

func (h *Handler) CreateModule(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	createReq, ok := validatedBody.(*CreateModuleRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.CreateModule(createReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create module", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgModuleCreated, result)
}

func (h *Handler) UpdateModule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid module ID")
		return
	}

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	updateReq, ok := validatedBody.(*UpdateModuleRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.UpdateModule(id, updateReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgModuleUpdated, result)
}

func (h *Handler) DeleteModule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid module ID")
		return
	}

	if err := h.service.DeleteModule(id); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgModuleDeleted, nil)
}

func (h *Handler) GetModuleTree(c *gin.Context) {
	category := c.Query("category")
	parentName := c.Query("parent")

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Internal error", "Invalid user ID type")
		return
	}

	var result []*ModuleTreeResponse
	var err error

	if parentName != "" {
		result, err = h.service.GetModuleTreeByParentFiltered(userIDInt64, parentName)
	} else {
		result, err = h.service.GetModuleTreeFiltered(userIDInt64, category)
	}

	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgModulesRetrieved, result)
}

func (h *Handler) GetModuleChildren(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid module ID")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Internal error", "Invalid user ID type")
		return
	}

	result, err := h.service.GetModuleChildrenFiltered(userIDInt64, id)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgModulesRetrieved, result)
}

func (h *Handler) GetModuleAncestors(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid module ID")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Internal error", "Invalid user ID type")
		return
	}

	result, err := h.service.GetModuleAncestorsFiltered(userIDInt64, id)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgModulesRetrieved, result)
}

// Route registration
func RegisterRoutes(api *gin.RouterGroup, handler *Handler) {
	modules := api.Group("/modules")
	{
		// GET /api/v1/modules - Get all modules with optional filters
		modules.GET("", handler.GetModules)

		// GET /api/v1/modules/:id - Get module by ID
		modules.GET("/:id", handler.GetModuleByID)

		// POST /api/v1/modules - Create new module
		modules.POST("",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &CreateModuleRequest{},
			}),
			handler.CreateModule,
		)

		// PUT /api/v1/modules/:id - Update module by ID
		modules.PUT("/:id",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &UpdateModuleRequest{},
			}),
			handler.UpdateModule,
		)

		// DELETE /api/v1/modules/:id - Delete module by ID
		modules.DELETE("/:id", handler.DeleteModule)

		// GET /api/v1/modules/tree - Get module tree structure
		modules.GET("/tree", handler.GetModuleTree)

		// GET /api/v1/modules/:id/children - Get module children by ID
		modules.GET("/:id/children", handler.GetModuleChildren)

		// GET /api/v1/modules/:id/ancestors - Get module ancestors by ID
		modules.GET("/:id/ancestors", handler.GetModuleAncestors)
	}
}
