package handlers

import (
	"gin-scalable-api/internal/service"
	"strconv"

	"gin-scalable-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CompanyHandler struct {
	companyService *service.CompanyService
}

func NewCompanyHandler(companyService *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
	}
}

func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	var req service.CompanyListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.companyService.GetCompanies(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get companies", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *CompanyHandler) GetCompanyByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid company ID")
		return
	}

	result, err := h.companyService.GetCompanyByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert to the expected struct
	req, ok := validatedBody.(*struct {
		Name string `json:"name" validate:"required,min=2,max=100"`
		Code string `json:"code" validate:"required,min=2,max=20"`
	})
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	// Convert to service request
	createReq := &service.CreateCompanyRequest{
		Name: req.Name,
		Code: req.Code,
	}

	result, err := h.companyService.CreateCompany(createReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create company", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Created successfully", result)
}

func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid company ID")
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
		Name string `json:"name"`
		Code string `json:"code"`
	})
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	// Convert to service request
	updateReq := &service.UpdateCompanyRequest{
		Name: req.Name,
		Code: req.Code,
	}

	result, err := h.companyService.UpdateCompany(id, updateReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid company ID")
		return
	}

	if err := h.companyService.DeleteCompany(id); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Company deleted successfully", nil)
}
