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

type ModuleHandler struct {
	moduleService interfaces.ModuleServiceInterface
}

func NewModuleHandler(moduleService interfaces.ModuleServiceInterface) *ModuleHandler {
	return &ModuleHandler{
		moduleService: moduleService,
	}
}

func (h *ModuleHandler) GetModules(c *gin.Context) {
	var req dto.ModuleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	// Ambil ID pengguna dari context (diset oleh middleware auth)
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

	// Gunakan metode terfilter untuk mendapatkan modul berdasarkan izin pengguna
	result, err := h.moduleService.GetModulesFiltered(userIDInt64, &req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get modules", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgModulesRetrieved, result)
}

func (h *ModuleHandler) GetModuleByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid module ID")
		return
	}

	result, err := h.moduleService.GetModuleByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgModuleNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgModuleRetrieved, result)
}

func (h *ModuleHandler) CreateModule(c *gin.Context) {
	// Ambil body yang sudah divalidasi dari context (diset oleh middleware validasi)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assertion ke DTO yang diharapkan
	createReq, ok := validatedBody.(*dto.CreateModuleRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.moduleService.CreateModule(createReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create module", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgModuleCreated, result)
}

func (h *ModuleHandler) UpdateModule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid module ID")
		return
	}

	// Ambil body yang sudah divalidasi dari context (diset oleh middleware validasi)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assertion ke DTO yang diharapkan
	updateReq, ok := validatedBody.(*dto.UpdateModuleRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.moduleService.UpdateModule(id, updateReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgModuleUpdated, result)
}

func (h *ModuleHandler) DeleteModule(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid module ID")
		return
	}

	if err := h.moduleService.DeleteModule(id); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgModuleDeleted, nil)
}

func (h *ModuleHandler) GetModuleTree(c *gin.Context) {
	category := c.Query("category")
	parentName := c.Query("parent")

	// Ambil ID pengguna dari context (diset oleh middleware auth)
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

	var result []*dto.ModuleTreeResponse
	var err error

	if parentName != "" {
		// Dapatkan tree berdasarkan nama parent dengan filtering
		result, err = h.moduleService.GetModuleTreeByParentFiltered(userIDInt64, parentName)
	} else {
		// Dapatkan tree berdasarkan kategori dengan filtering
		result, err = h.moduleService.GetModuleTreeFiltered(userIDInt64, category)
	}

	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgModulesRetrieved, result)
}

func (h *ModuleHandler) GetModuleChildren(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid module ID")
		return
	}

	// Ambil ID pengguna dari context (diset oleh middleware auth)
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

	result, err := h.moduleService.GetModuleChildrenFiltered(userIDInt64, id)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgModulesRetrieved, result)
}

func (h *ModuleHandler) GetModuleAncestors(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid module ID")
		return
	}

	// Ambil ID pengguna dari context (diset oleh middleware auth)
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

	result, err := h.moduleService.GetModuleAncestorsFiltered(userIDInt64, id)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgModulesRetrieved, result)
}
