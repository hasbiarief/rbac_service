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

// @Summary      Get all modules
// @Description  Mendapatkan daftar semua module dengan filter opsional dan pagination. Filtered berdasarkan user access
// @Tags         Modules
// @Accept       json
// @Produce      json
// @Param        limit             query     int     false  "Limit jumlah data"
// @Param        offset            query     int     false  "Offset data"
// @Param        search            query     string  false  "Search by name"
// @Param        category          query     string  false  "Filter by category"
// @Param        subscription_tier query     string  false  "Filter by subscription tier"
// @Param        parent_id         query     int     false  "Filter by parent ID"
// @Param        is_active         query     bool    false  "Filter by active status"
// @Success      200               {object}  response.Response{data=module.ModuleListResponse}  "Modules berhasil diambil"
// @Failure      400               {object}  response.Response  "Bad request"
// @Failure      401               {object}  response.Response  "Unauthorized"
// @Failure      500               {object}  response.Response  "Internal server error"
// @Router       /api/v1/modules [get]
// @Security     BearerAuth
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

// @Summary      Get module by ID
// @Description  Mendapatkan detail module berdasarkan ID
// @Tags         Modules
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Module ID"
// @Success      200  {object}  response.Response{data=module.ModuleResponse}  "Module berhasil diambil"
// @Failure      400  {object}  response.Response  "Bad request - Invalid module ID"
// @Failure      404  {object}  response.Response  "Module tidak ditemukan"
// @Router       /api/v1/modules/{id} [get]
// @Security     BearerAuth
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

// @Summary      Create new module
// @Description  Membuat module baru dengan support untuk hierarchical structure
// @Tags         Modules
// @Accept       json
// @Produce      json
// @Param        module  body      module.CreateModuleRequest  true  "Module data"
// @Success      201     {object}  response.Response{data=module.ModuleResponse}  "Module berhasil dibuat"
// @Failure      400     {object}  response.Response  "Bad request - validation failed"
// @Failure      409     {object}  response.Response  "Conflict - module URL sudah ada"
// @Failure      500     {object}  response.Response  "Internal server error"
// @Router       /api/v1/modules [post]
// @Security     BearerAuth
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

// @Summary      Update module
// @Description  Memperbarui informasi module
// @Tags         Modules
// @Accept       json
// @Produce      json
// @Param        id      path      int                         true  "Module ID"
// @Param        module  body      module.UpdateModuleRequest  true  "Module data yang akan diupdate"
// @Success      200     {object}  response.Response{data=module.ModuleResponse}  "Module berhasil diupdate"
// @Failure      400     {object}  response.Response  "Bad request - Invalid module ID atau validation failed"
// @Failure      404     {object}  response.Response  "Module tidak ditemukan"
// @Failure      500     {object}  response.Response  "Internal server error"
// @Router       /api/v1/modules/{id} [put]
// @Security     BearerAuth
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

// @Summary      Delete module
// @Description  Menghapus module berdasarkan ID
// @Tags         Modules
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Module ID"
// @Success      200  {object}  response.Response  "Module berhasil dihapus"
// @Failure      400  {object}  response.Response  "Bad request - Invalid module ID"
// @Failure      404  {object}  response.Response  "Module tidak ditemukan"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Router       /api/v1/modules/{id} [delete]
// @Security     BearerAuth
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

// @Summary      Get module tree
// @Description  Mendapatkan module tree structure dengan filter category atau parent. Filtered berdasarkan user access
// @Tags         Modules
// @Accept       json
// @Produce      json
// @Param        category  query     string  false  "Filter by category"
// @Param        parent    query     string  false  "Filter by parent name"
// @Success      200       {object}  response.Response{data=[]module.ModuleTreeResponse}  "Module tree berhasil diambil"
// @Failure      401       {object}  response.Response  "Unauthorized"
// @Failure      500       {object}  response.Response  "Internal server error"
// @Router       /api/v1/modules/tree [get]
// @Security     BearerAuth
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

// @Summary      Get module children
// @Description  Mendapatkan daftar child modules dari module tertentu. Filtered berdasarkan user access
// @Tags         Modules
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Module ID"
// @Success      200  {object}  response.Response{data=[]module.ModuleResponse}  "Module children berhasil diambil"
// @Failure      400  {object}  response.Response  "Bad request - Invalid module ID"
// @Failure      401  {object}  response.Response  "Unauthorized"
// @Failure      404  {object}  response.Response  "Module tidak ditemukan"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Router       /api/v1/modules/{id}/children [get]
// @Security     BearerAuth
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

// @Summary      Get module ancestors
// @Description  Mendapatkan daftar ancestor modules (parent hierarchy) dari module tertentu. Filtered berdasarkan user access
// @Tags         Modules
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Module ID"
// @Success      200  {object}  response.Response{data=[]module.ModuleResponse}  "Module ancestors berhasil diambil"
// @Failure      400  {object}  response.Response  "Bad request - Invalid module ID"
// @Failure      401  {object}  response.Response  "Unauthorized"
// @Failure      404  {object}  response.Response  "Module tidak ditemukan"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Router       /api/v1/modules/{id}/ancestors [get]
// @Security     BearerAuth
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
