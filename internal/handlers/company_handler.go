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

type CompanyHandler struct {
	companyService interfaces.CompanyServiceInterface
}

func NewCompanyHandler(companyService interfaces.CompanyServiceInterface) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
	}
}

func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	var req dto.CompanyListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.companyService.GetCompanies(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get companies", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgCompaniesRetrieved, result)
}

func (h *CompanyHandler) GetCompanyByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid company ID")
		return
	}

	result, err := h.companyService.GetCompanyByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgCompanyNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgCompanyRetrieved, result)
}

func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	// Ambil body yang sudah divalidasi dari context (diset oleh middleware validasi)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assertion ke DTO yang diharapkan
	req, ok := validatedBody.(*dto.CreateCompanyRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.companyService.CreateCompany(req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create company", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgCompanyCreated, result)
}

func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid company ID")
		return
	}

	// Ambil body yang sudah divalidasi dari context (diset oleh middleware validasi)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assertion ke DTO yang diharapkan
	req, ok := validatedBody.(*dto.UpdateCompanyRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.companyService.UpdateCompany(id, req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgCompanyUpdated, result)
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

	response.Success(c, http.StatusOK, constants.MsgCompanyDeleted, nil)
}
