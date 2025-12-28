package handlers

import (
	"net/http"
	"strconv"

	"gin-scalable-api/internal/service"
	"gin-scalable-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type BranchHandler struct {
	branchService *service.BranchService
}

func NewBranchHandler(branchService *service.BranchService) *BranchHandler {
	return &BranchHandler{
		branchService: branchService,
	}
}

func (h *BranchHandler) GetBranches(c *gin.Context) {
	var req service.BranchListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request parameters", err.Error())
		return
	}

	result, err := h.branchService.GetBranches(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get branches", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Branches retrieved successfully", result)
}

func (h *BranchHandler) GetBranchByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid branch ID", "Branch ID must be a valid number")
		return
	}

	result, err := h.branchService.GetBranchByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Branch not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Branch retrieved successfully", result)
}

func (h *BranchHandler) CreateBranch(c *gin.Context) {
	var req service.CreateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	result, err := h.branchService.CreateBranch(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create branch", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Branch created successfully", result)
}

func (h *BranchHandler) UpdateBranch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid branch ID", "Branch ID must be a valid number")
		return
	}

	var req service.UpdateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	result, err := h.branchService.UpdateBranch(id, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update branch", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Branch updated successfully", result)
}

func (h *BranchHandler) DeleteBranch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid branch ID", "Branch ID must be a valid number")
		return
	}

	if err := h.branchService.DeleteBranch(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete branch", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Branch deleted successfully", nil)
}

func (h *BranchHandler) GetCompanyBranches(c *gin.Context) {
	companyID, err := strconv.ParseInt(c.Param("companyId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid company ID", "Company ID must be a valid number")
		return
	}

	// Get includeHierarchy from query parameter, default to true
	includeHierarchy := c.DefaultQuery("include_hierarchy", "true") == "true"

	result, err := h.branchService.GetCompanyBranches(companyID, includeHierarchy)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get company branches", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Company branches retrieved successfully", result)
}

func (h *BranchHandler) GetBranchChildren(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid branch ID", "Branch ID must be a valid number")
		return
	}

	result, err := h.branchService.GetBranchChildren(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get branch children", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Branch children retrieved successfully", result)
}
