package application

import (
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
	return &Handler{
		service: service,
	}
}

// Handler methods

// GetApplications godoc
// @Summary      Get all applications
// @Description  Mendapatkan daftar semua applications dengan filter opsional dan pagination
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Param        limit      query     int     false  "Limit jumlah data"
// @Param        offset     query     int     false  "Offset data"
// @Param        search     query     string  false  "Search by name atau code"
// @Param        is_active  query     bool    false  "Filter by active status"
// @Success      200        {object}  response.Response{data=application.ApplicationListResponse}  "Applications berhasil diambil"
// @Failure      400        {object}  response.Response  "Bad request"
// @Failure      500        {object}  response.Response  "Internal server error"
// @Router       /api/v1/applications [get]
// @Security     BearerAuth
func (h *Handler) GetApplications(c *gin.Context) {
	var req ApplicationListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.service.GetApplications(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get applications", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Applications successfully retrieved", result)
}

// GetApplicationByID godoc
// @Summary      Get application by ID
// @Description  Mendapatkan detail application berdasarkan ID
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Application ID"
// @Success      200  {object}  response.Response{data=application.ApplicationResponse}  "Application berhasil diambil"
// @Failure      400  {object}  response.Response  "Bad request - Invalid application ID"
// @Failure      404  {object}  response.Response  "Application tidak ditemukan"
// @Router       /api/v1/applications/{id} [get]
// @Security     BearerAuth
func (h *Handler) GetApplicationByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid application ID")
		return
	}

	result, err := h.service.GetApplicationByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Application not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Application successfully retrieved", result)
}

// GetApplicationByCode godoc
// @Summary      Get application by code
// @Description  Mendapatkan detail application berdasarkan application code
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Param        code  path      string  true  "Application code"
// @Success      200   {object}  response.Response{data=application.ApplicationResponse}  "Application berhasil diambil"
// @Failure      400   {object}  response.Response  "Bad request - Application code diperlukan"
// @Failure      404   {object}  response.Response  "Application tidak ditemukan"
// @Router       /api/v1/applications/code/{code} [get]
// @Security     BearerAuth
func (h *Handler) GetApplicationByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.Error(c, http.StatusBadRequest, "Bad request", "Application code is required")
		return
	}

	result, err := h.service.GetApplicationByCode(code)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Application not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Application successfully retrieved", result)
}

// CreateApplication godoc
// @Summary      Create new application
// @Description  Membuat application baru
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Param        application  body      application.CreateApplicationRequest  true  "Application data"
// @Success      201          {object}  response.Response{data=application.ApplicationResponse}  "Application berhasil dibuat"
// @Failure      400          {object}  response.Response  "Bad request - validation failed"
// @Failure      409          {object}  response.Response  "Conflict - application code sudah ada"
// @Failure      500          {object}  response.Response  "Internal server error"
// @Router       /api/v1/applications [post]
// @Security     BearerAuth
func (h *Handler) CreateApplication(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	createReq, ok := validatedBody.(*CreateApplicationRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.CreateApplication(createReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create application", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Application successfully created", result)
}

// UpdateApplication godoc
// @Summary      Update application
// @Description  Memperbarui informasi application
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Param        id           path      int                                   true  "Application ID"
// @Param        application  body      application.UpdateApplicationRequest  true  "Application data yang akan diupdate"
// @Success      200          {object}  response.Response{data=application.ApplicationResponse}  "Application berhasil diupdate"
// @Failure      400          {object}  response.Response  "Bad request - Invalid application ID atau validation failed"
// @Failure      404          {object}  response.Response  "Application tidak ditemukan"
// @Failure      500          {object}  response.Response  "Internal server error"
// @Router       /api/v1/applications/{id} [put]
// @Security     BearerAuth
func (h *Handler) UpdateApplication(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid application ID")
		return
	}

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	updateReq, ok := validatedBody.(*UpdateApplicationRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.UpdateApplication(id, updateReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Application successfully updated", result)
}

// DeleteApplication godoc
// @Summary      Delete application
// @Description  Menghapus application berdasarkan ID
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Application ID"
// @Success      200  {object}  response.Response  "Application berhasil dihapus"
// @Failure      400  {object}  response.Response  "Bad request - Invalid application ID"
// @Failure      404  {object}  response.Response  "Application tidak ditemukan"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Router       /api/v1/applications/{id} [delete]
// @Security     BearerAuth
func (h *Handler) DeleteApplication(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid application ID")
		return
	}

	if err := h.service.DeleteApplication(id); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Application successfully deleted", nil)
}

// GetPlanApplications godoc
// @Summary      Get plan applications (Admin)
// @Description  Mendapatkan daftar applications yang termasuk dalam subscription plan (admin only)
// @Tags         Applications (Admin)
// @Accept       json
// @Produce      json
// @Param        planId  path      int  true  "Subscription Plan ID"
// @Success      200     {object}  response.Response{data=[]application.ApplicationResponse}  "Plan applications berhasil diambil"
// @Failure      400     {object}  response.Response  "Bad request - Invalid plan ID"
// @Failure      500     {object}  response.Response  "Internal server error"
// @Router       /api/v1/admin/plan-applications/{planId} [get]
// @Security     BearerAuth
func (h *Handler) GetPlanApplications(c *gin.Context) {
	planID, err := strconv.ParseInt(c.Param("planId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid plan ID")
		return
	}

	result, err := h.service.GetPlanApplications(planID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Plan applications successfully retrieved", result)
}

// AddApplicationsToPlan godoc
// @Summary      Add applications to plan (Admin)
// @Description  Menambahkan applications ke subscription plan (admin only)
// @Tags         Applications (Admin)
// @Accept       json
// @Produce      json
// @Param        planId        path      int                                 true  "Subscription Plan ID"
// @Param        applications  body      application.PlanApplicationRequest  true  "Applications to add"
// @Success      200           {object}  response.Response  "Applications berhasil ditambahkan ke plan"
// @Failure      400           {object}  response.Response  "Bad request - Invalid plan ID atau validation failed"
// @Failure      404           {object}  response.Response  "Subscription plan tidak ditemukan"
// @Failure      500           {object}  response.Response  "Internal server error"
// @Router       /api/v1/admin/plan-applications/{planId} [post]
// @Security     BearerAuth
func (h *Handler) AddApplicationsToPlan(c *gin.Context) {
	planID, err := strconv.ParseInt(c.Param("planId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid plan ID")
		return
	}

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	addReq, ok := validatedBody.(*PlanApplicationRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	if err := h.service.AddApplicationsToPlan(planID, addReq); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to add applications to plan", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Applications successfully added to plan", nil)
}

// RemoveApplicationFromPlan godoc
// @Summary      Remove application from plan (Admin)
// @Description  Menghapus application dari subscription plan (admin only)
// @Tags         Applications (Admin)
// @Accept       json
// @Produce      json
// @Param        planId         path      int  true  "Subscription Plan ID"
// @Param        applicationId  path      int  true  "Application ID"
// @Success      200            {object}  response.Response  "Application berhasil dihapus dari plan"
// @Failure      400            {object}  response.Response  "Bad request - Invalid ID"
// @Failure      404            {object}  response.Response  "Plan atau application tidak ditemukan"
// @Failure      500            {object}  response.Response  "Internal server error"
// @Router       /api/v1/admin/plan-applications/{planId}/{applicationId} [delete]
// @Security     BearerAuth
func (h *Handler) RemoveApplicationFromPlan(c *gin.Context) {
	planID, err := strconv.ParseInt(c.Param("planId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid plan ID")
		return
	}

	applicationID, err := strconv.ParseInt(c.Param("applicationId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid application ID")
		return
	}

	if err := h.service.RemoveApplicationFromPlan(planID, applicationID); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Application successfully removed from plan", nil)
}

// Route registration
func RegisterRoutes(api *gin.RouterGroup, handler *Handler) {
	applications := api.Group("/applications")
	{
		// GET /api/v1/applications - Get all applications
		applications.GET("", handler.GetApplications)

		// GET /api/v1/applications/:id - Get application by ID
		applications.GET("/:id", handler.GetApplicationByID)

		// GET /api/v1/applications/code/:code - Get application by code
		applications.GET("/code/:code", handler.GetApplicationByCode)

		// POST /api/v1/applications - Create new application
		applications.POST("",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &CreateApplicationRequest{},
			}),
			handler.CreateApplication,
		)

		// PUT /api/v1/applications/:id - Update application by ID
		applications.PUT("/:id",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &UpdateApplicationRequest{},
			}),
			handler.UpdateApplication,
		)

		// DELETE /api/v1/applications/:id - Delete application by ID
		applications.DELETE("/:id", handler.DeleteApplication)
	}

	// Admin routes for plan-application management
	admin := api.Group("/admin")
	{
		// GET /api/v1/admin/plan-applications/:planId - Get applications for plan
		admin.GET("/plan-applications/:planId", handler.GetPlanApplications)

		// POST /api/v1/admin/plan-applications/:planId - Add applications to plan
		admin.POST("/plan-applications/:planId",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &PlanApplicationRequest{},
			}),
			handler.AddApplicationsToPlan,
		)

		// DELETE /api/v1/admin/plan-applications/:planId/:applicationId - Remove application from plan
		admin.DELETE("/plan-applications/:planId/:applicationId", handler.RemoveApplicationFromPlan)
	}
}
