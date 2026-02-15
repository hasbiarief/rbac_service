package company

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

// @Summary      Get all companies
// @Description  Mendapatkan daftar semua company dengan filter opsional dan pagination
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        limit      query     int     false  "Limit jumlah data"
// @Param        offset     query     int     false  "Offset data"
// @Param        search     query     string  false  "Search by name atau code"
// @Param        is_active  query     bool    false  "Filter by active status"
// @Success      200        {object}  response.Response{data=company.CompanyListResponse}  "Companies berhasil diambil"
// @Failure      400        {object}  response.Response  "Bad request"
// @Failure      500        {object}  response.Response  "Internal server error"
// @Router       /api/v1/companies [get]
// @Security     BearerAuth
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

// @Summary      Get company by ID
// @Description  Mendapatkan detail company berdasarkan ID
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Company ID"
// @Success      200  {object}  response.Response{data=company.CompanyResponse}  "Company berhasil diambil"
// @Failure      400  {object}  response.Response  "Bad request - Invalid company ID"
// @Failure      404  {object}  response.Response  "Company tidak ditemukan"
// @Router       /api/v1/companies/{id} [get]
// @Security     BearerAuth
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

// @Summary      Create new company
// @Description  Membuat company baru
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        company  body      company.CreateCompanyRequest  true  "Company data"
// @Success      201      {object}  response.Response{data=company.CompanyResponse}  "Company berhasil dibuat"
// @Failure      400      {object}  response.Response  "Bad request - validation failed"
// @Failure      409      {object}  response.Response  "Conflict - company code sudah ada"
// @Failure      500      {object}  response.Response  "Internal server error"
// @Router       /api/v1/companies [post]
// @Security     BearerAuth
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

// @Summary      Update company
// @Description  Memperbarui informasi company
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        id       path      int                           true  "Company ID"
// @Param        company  body      company.UpdateCompanyRequest  true  "Company data yang akan diupdate"
// @Success      200      {object}  response.Response{data=company.CompanyResponse}  "Company berhasil diupdate"
// @Failure      400      {object}  response.Response  "Bad request - Invalid company ID atau validation failed"
// @Failure      404      {object}  response.Response  "Company tidak ditemukan"
// @Failure      500      {object}  response.Response  "Internal server error"
// @Router       /api/v1/companies/{id} [put]
// @Security     BearerAuth
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

// @Summary      Delete company
// @Description  Menghapus company berdasarkan ID
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Company ID"
// @Success      200  {object}  response.Response  "Company berhasil dihapus"
// @Failure      400  {object}  response.Response  "Bad request - Invalid company ID"
// @Failure      404  {object}  response.Response  "Company tidak ditemukan"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Router       /api/v1/companies/{id} [delete]
// @Security     BearerAuth
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

// Route registration
func RegisterRoutes(api *gin.RouterGroup, handler *Handler) {
	companies := api.Group("/companies")
	{
		// GET /api/v1/companies - Get all companies with optional filters
		companies.GET("", handler.GetCompanies)

		// GET /api/v1/companies/:id - Get company by ID
		companies.GET("/:id", handler.GetCompanyByID)

		// POST /api/v1/companies - Create new company
		companies.POST("",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &CreateCompanyRequest{},
			}),
			handler.CreateCompany,
		)

		// PUT /api/v1/companies/:id - Update company by ID
		companies.PUT("/:id",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &UpdateCompanyRequest{},
			}),
			handler.UpdateCompany,
		)

		// DELETE /api/v1/companies/:id - Delete company by ID
		companies.DELETE("/:id", handler.DeleteCompany)
	}
}
