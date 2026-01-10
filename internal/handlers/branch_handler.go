package handlers

import (
	"net/http"
	"strconv"

	"gin-scalable-api/internal/constants"
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/interfaces"
	"gin-scalable-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type BranchHandler struct {
	branchService interfaces.BranchServiceInterface
}

func NewBranchHandler(branchService interfaces.BranchServiceInterface) *BranchHandler {
	return &BranchHandler{
		branchService: branchService,
	}
}

func (h *BranchHandler) GetBranches(c *gin.Context) {
	var req dto.BranchListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request parameter", err.Error())
		return
	}

	// Check if nested structure is requested
	nested := c.DefaultQuery("nested", "false") == "true"

	if nested {
		result, err := h.branchService.GetBranchesNested(&req)
		if err != nil {
			response.ErrorWithAutoStatus(c, "Failed to get branches", err.Error())
			return
		}
		response.Success(c, http.StatusOK, constants.MsgBranchesRetrieved+" (nested)", result)
		return
	}

	result, err := h.branchService.GetBranches(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get branches", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgBranchesRetrieved, result)
}

func (h *BranchHandler) GetBranchByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Branch ID is invalid", "Branch ID must be a valid number")
		return
	}

	result, err := h.branchService.GetBranchByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgBranchNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgBranchRetrieved, result)
}

func (h *BranchHandler) CreateBranch(c *gin.Context) {
	// Ambil body yang sudah divalidasi dari context (diset oleh middleware validasi)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	// Type assertion ke DTO yang diharapkan
	req, ok := validatedBody.(*dto.CreateBranchRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "body structure is invalid")
		return
	}

	result, err := h.branchService.CreateBranch(req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create branch", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgBranchCreated, result)
}

func (h *BranchHandler) UpdateBranch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Branch ID is invalid", "Branch ID must be a valid number")
		return
	}

	// Ambil body yang sudah divalidasi dari context (diset oleh middleware validasi)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	// Type assertion ke DTO yang diharapkan
	req, ok := validatedBody.(*dto.UpdateBranchRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "body structure is invalid")
		return
	}

	result, err := h.branchService.UpdateBranch(id, req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to update branch", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgBranchUpdated, result)
}

func (h *BranchHandler) DeleteBranch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Branch ID is invalid", "Branch ID must be a valid number")
		return
	}

	if err := h.branchService.DeleteBranch(id); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to delete branch", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgBranchDeleted, nil)
}

func (h *BranchHandler) GetCompanyBranches(c *gin.Context) {
	companyID, err := strconv.ParseInt(c.Param("companyId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Company ID is invalid", "Company ID must be a valid number")
		return
	}

	// Ambil includeHierarchy dari parameter query, default ke true
	includeHierarchy := c.DefaultQuery("include_hierarchy", "true") == "true"

	// Periksa apakah struktur bersarang diminta
	nested := c.DefaultQuery("nested", "false") == "true"

	if nested {
		result, err := h.branchService.GetCompanyBranchesNested(companyID)
		if err != nil {
			response.ErrorWithAutoStatus(c, "Failed to get company branches", err.Error())
			return
		}
		response.Success(c, http.StatusOK, constants.MsgCompanyBranchesRetrieved+" (nested)", result)
		return
	}

	result, err := h.branchService.GetCompanyBranches(companyID, includeHierarchy)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get company branches", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgCompanyBranchesRetrieved, result)
}

func (h *BranchHandler) GetBranchChildren(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Branch ID is invalid", "Branch ID must be a valid number")
		return
	}

	result, err := h.branchService.GetBranchChildren(id)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get branch children", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgBranchChildrenRetrieved, result)
}

func (h *BranchHandler) GetBranchHierarchy(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Branch ID is invalid", "Branch ID must be a valid number")
		return
	}

	// Periksa apakah struktur bersarang diminta
	nested := c.DefaultQuery("nested", "false") == "true"

	result, err := h.branchService.GetBranchHierarchyByID(id, nested)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get branch hierarchy", err.Error())
		return
	}

	message := constants.MsgBranchHierarchyRetrieved
	if nested {
		message += " (nested)"
	}

	response.Success(c, http.StatusOK, message, result)
}
