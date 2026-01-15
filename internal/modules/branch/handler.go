package branch

import (
	"gin-scalable-api/internal/constants"
	"gin-scalable-api/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetBranches(c *gin.Context) {
	var req BranchListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request parameter", err.Error())
		return
	}

	nested := c.DefaultQuery("nested", "false") == "true"

	if nested {
		result, err := h.service.GetBranchesNested(&req)
		if err != nil {
			response.ErrorWithAutoStatus(c, "Failed to get branches", err.Error())
			return
		}
		response.Success(c, http.StatusOK, constants.MsgBranchesRetrieved+" (nested)", result)
		return
	}

	result, err := h.service.GetBranches(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get branches", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgBranchesRetrieved, result)
}

func (h *Handler) GetBranchByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Branch ID is invalid", "Branch ID must be a valid number")
		return
	}

	result, err := h.service.GetBranchByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgBranchNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgBranchRetrieved, result)
}

func (h *Handler) CreateBranch(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	req, ok := validatedBody.(*CreateBranchRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "body structure is invalid")
		return
	}

	result, err := h.service.CreateBranch(req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create branch", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgBranchCreated, result)
}

func (h *Handler) UpdateBranch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Branch ID is invalid", "Branch ID must be a valid number")
		return
	}

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	req, ok := validatedBody.(*UpdateBranchRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "body structure is invalid")
		return
	}

	result, err := h.service.UpdateBranch(id, req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to update branch", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgBranchUpdated, result)
}

func (h *Handler) DeleteBranch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Branch ID is invalid", "Branch ID must be a valid number")
		return
	}

	if err := h.service.DeleteBranch(id); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to delete branch", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgBranchDeleted, nil)
}

func (h *Handler) GetCompanyBranches(c *gin.Context) {
	companyID, err := strconv.ParseInt(c.Param("companyId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Company ID is invalid", "Company ID must be a valid number")
		return
	}

	includeHierarchy := c.DefaultQuery("include_hierarchy", "true") == "true"
	nested := c.DefaultQuery("nested", "false") == "true"

	if nested {
		result, err := h.service.GetCompanyBranchesNested(companyID)
		if err != nil {
			response.ErrorWithAutoStatus(c, "Failed to get company branches", err.Error())
			return
		}
		response.Success(c, http.StatusOK, constants.MsgCompanyBranchesRetrieved+" (nested)", result)
		return
	}

	result, err := h.service.GetCompanyBranches(companyID, includeHierarchy)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get company branches", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgCompanyBranchesRetrieved, result)
}

func (h *Handler) GetBranchChildren(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Branch ID is invalid", "Branch ID must be a valid number")
		return
	}

	result, err := h.service.GetBranchChildren(id)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get branch children", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgBranchChildrenRetrieved, result)
}

func (h *Handler) GetBranchHierarchy(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Branch ID is invalid", "Branch ID must be a valid number")
		return
	}

	nested := c.DefaultQuery("nested", "false") == "true"

	result, err := h.service.GetBranchHierarchyByID(id, nested)
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
