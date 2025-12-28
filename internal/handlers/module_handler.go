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
	var req service.CreateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.moduleService.CreateModule(&req)
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

	var req service.UpdateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.moduleService.UpdateModule(id, &req)
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
