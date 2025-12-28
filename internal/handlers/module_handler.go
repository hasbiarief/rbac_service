package handlers

import (
	"gin-scalable-api/internal/service"
	"strconv"

	"gin-scalable-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ModuleHandler struct {
	moduleService *service.ModuleService
}

func NewModuleHandler(moduleService *service.ModuleService) *ModuleHandler {
	return &ModuleHandler{
		moduleService: moduleService,
	}
}

func (h *ModuleHandler) GetModules(c *gin.Context) {
	var req service.ModuleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.moduleService.GetModules(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *ModuleHandler) GetModuleByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid module ID")
		return
	}

	result, err := h.moduleService.GetModuleByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *ModuleHandler) CreateModule(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert to the expected struct
	req, ok := validatedBody.(*struct {
		Category         string `json:"category" validate:"required,min=2,max=50"`
		Name             string `json:"name" validate:"required,min=2,max=100"`
		URL              string `json:"url" validate:"required,min=1,max=255"`
		Icon             string `json:"icon" validate:"max=50"`
		Description      string `json:"description" validate:"max=500"`
		ParentID         *int64 `json:"parent_id"`
		SubscriptionTier string `json:"subscription_tier" validate:"required,oneof=basic pro enterprise"`
	})
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	// Convert to service request
	createReq := &service.CreateModuleRequest{
		Category:         req.Category,
		Name:             req.Name,
		URL:              req.URL,
		Icon:             req.Icon,
		Description:      req.Description,
		ParentID:         req.ParentID,
		SubscriptionTier: req.SubscriptionTier,
	}

	result, err := h.moduleService.CreateModule(createReq)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Created successfully", result)
}

func (h *ModuleHandler) UpdateModule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid module ID")
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
		Category         string `json:"category"`
		Name             string `json:"name"`
		URL              string `json:"url"`
		Icon             string `json:"icon"`
		Description      string `json:"description"`
		ParentID         *int64 `json:"parent_id"`
		SubscriptionTier string `json:"subscription_tier"`
	})
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	// Convert to service request
	updateReq := &service.UpdateModuleRequest{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    nil, // This field is not provided in the request struct
	}

	result, err := h.moduleService.UpdateModule(id, updateReq)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *ModuleHandler) DeleteModule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid module ID")
		return
	}

	if err := h.moduleService.DeleteModule(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Module deleted successfully", nil)
}

func (h *ModuleHandler) GetModuleTree(c *gin.Context) {
	category := c.Query("category")
	result, err := h.moduleService.GetModuleTree(category)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *ModuleHandler) GetModuleChildren(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid module ID")
		return
	}

	result, err := h.moduleService.GetModuleChildren(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *ModuleHandler) GetModuleAncestors(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid module ID")
		return
	}

	result, err := h.moduleService.GetModuleAncestors(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}
