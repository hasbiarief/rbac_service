package company

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

func (h *Handler) GetCompanies(c *gin.Context) {
	var req CompanyListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.service.GetCompanies(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get companies", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgCompaniesRetrieved, result)
}

func (h *Handler) GetCompanyByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid company ID")
		return
	}

	result, err := h.service.GetCompanyByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgCompanyNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgCompanyRetrieved, result)
}

func (h *Handler) CreateCompany(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*CreateCompanyRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.CreateCompany(req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create company", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgCompanyCreated, result)
}

func (h *Handler) UpdateCompany(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid company ID")
		return
	}

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*UpdateCompanyRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.UpdateCompany(id, req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgCompanyUpdated, result)
}

func (h *Handler) DeleteCompany(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid company ID")
		return
	}

	if err := h.service.DeleteCompany(id); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgCompanyDeleted, nil)
}
