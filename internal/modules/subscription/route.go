package subscription

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

// @Summary      Get all subscription plans
// @Description  Mendapatkan daftar semua subscription plans (public endpoint)
// @Tags         Subscription Plans
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=[]subscription.SubscriptionPlanResponse}  "Subscription plans berhasil diambil"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Router       /api/v1/subscription-plans [get]
func (h *Handler) GetAllPlans(c *gin.Context) {
	result, err := h.service.GetSubscriptionPlans()
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionPlansRetrieved, result)
}

// @Summary      Get subscription plan by ID
// @Description  Mendapatkan detail subscription plan berdasarkan ID (public endpoint)
// @Tags         Subscription Plans
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Subscription Plan ID"
// @Success      200  {object}  response.Response{data=subscription.SubscriptionPlanResponse}  "Subscription plan berhasil diambil"
// @Failure      400  {object}  response.Response  "Bad request - Invalid plan ID"
// @Failure      404  {object}  response.Response  "Subscription plan tidak ditemukan"
// @Router       /api/v1/subscription-plans/{id} [get]
func (h *Handler) GetPlanByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid plan ID")
		return
	}

	result, err := h.service.GetSubscriptionPlanByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgSubscriptionPlanNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionPlanRetrieved, result)
}

// @Summary      Create subscription plan (Admin)
// @Description  Membuat subscription plan baru (admin only)
// @Tags         Subscription Plans (Admin)
// @Accept       json
// @Produce      json
// @Param        plan  body      subscription.CreateSubscriptionPlanRequest  true  "Subscription plan data"
// @Success      201   {object}  response.Response{data=subscription.SubscriptionPlanResponse}  "Subscription plan berhasil dibuat"
// @Failure      400   {object}  response.Response  "Bad request - validation failed"
// @Failure      409   {object}  response.Response  "Conflict - plan name sudah ada"
// @Failure      500   {object}  response.Response  "Internal server error"
// @Router       /api/v1/admin/subscription-plans [post]
// @Security     BearerAuth
func (h *Handler) CreateSubscriptionPlan(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*CreateSubscriptionPlanRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.CreateSubscriptionPlan(req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create subscription plan", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgSubscriptionPlanCreated, result)
}

// @Summary      Update subscription plan (Admin)
// @Description  Memperbarui subscription plan (admin only)
// @Tags         Subscription Plans (Admin)
// @Accept       json
// @Produce      json
// @Param        id    path      int                                         true  "Subscription Plan ID"
// @Param        plan  body      subscription.UpdateSubscriptionPlanRequest  true  "Subscription plan data yang akan diupdate"
// @Success      200   {object}  response.Response{data=subscription.SubscriptionPlanResponse}  "Subscription plan berhasil diupdate"
// @Failure      400   {object}  response.Response  "Bad request - Invalid plan ID atau validation failed"
// @Failure      404   {object}  response.Response  "Subscription plan tidak ditemukan"
// @Failure      500   {object}  response.Response  "Internal server error"
// @Router       /api/v1/admin/subscription-plans/{id} [put]
// @Security     BearerAuth
func (h *Handler) UpdateSubscriptionPlan(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid plan ID")
		return
	}

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*UpdateSubscriptionPlanRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.UpdateSubscriptionPlan(id, req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionPlanUpdated, result)
}

// @Summary      Delete subscription plan (Admin)
// @Description  Menghapus subscription plan (admin only)
// @Tags         Subscription Plans (Admin)
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Subscription Plan ID"
// @Success      200  {object}  response.Response  "Subscription plan berhasil dihapus"
// @Failure      400  {object}  response.Response  "Bad request - Invalid plan ID"
// @Failure      404  {object}  response.Response  "Subscription plan tidak ditemukan"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Router       /api/v1/admin/subscription-plans/{id} [delete]
// @Security     BearerAuth
func (h *Handler) DeleteSubscriptionPlan(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid plan ID")
		return
	}

	if err := h.service.DeleteSubscriptionPlan(id); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionPlanDeleted, nil)
}

// @Summary      Get all subscriptions
// @Description  Mendapatkan daftar semua subscriptions dengan filter opsional
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        limit      query     int     false  "Limit jumlah data"
// @Param        offset     query     int     false  "Offset data"
// @Param        company_id query     int     false  "Filter by company ID"
// @Param        plan_id    query     int     false  "Filter by plan ID"
// @Param        status     query     string  false  "Filter by status"
// @Success      200        {object}  response.Response{data=subscription.SubscriptionListResponse}  "Subscriptions berhasil diambil"
// @Failure      400        {object}  response.Response  "Bad request"
// @Failure      500        {object}  response.Response  "Internal server error"
// @Router       /api/v1/subscriptions [get]
// @Security     BearerAuth
func (h *Handler) GetAllSubscriptions(c *gin.Context) {
	var req SubscriptionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.service.GetSubscriptions(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionsRetrieved, result)
}

// @Summary      Create subscription
// @Description  Membuat subscription baru untuk company
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body      subscription.CreateSubscriptionRequest  true  "Subscription data"
// @Success      201           {object}  response.Response{data=subscription.SubscriptionResponse}  "Subscription berhasil dibuat"
// @Failure      400           {object}  response.Response  "Bad request - validation failed"
// @Failure      409           {object}  response.Response  "Conflict - company sudah memiliki subscription aktif"
// @Failure      500           {object}  response.Response  "Internal server error"
// @Router       /api/v1/subscriptions [post]
// @Security     BearerAuth
func (h *Handler) CreateSubscription(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*CreateSubscriptionRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.CreateSubscription(req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgSubscriptionCreated, result)
}

// @Summary      Get subscription by ID
// @Description  Mendapatkan detail subscription berdasarkan ID
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Subscription ID"
// @Success      200  {object}  response.Response{data=subscription.SubscriptionResponse}  "Subscription berhasil diambil"
// @Failure      400  {object}  response.Response  "Bad request - Invalid subscription ID"
// @Failure      404  {object}  response.Response  "Subscription tidak ditemukan"
// @Router       /api/v1/subscriptions/{id} [get]
// @Security     BearerAuth
func (h *Handler) GetSubscriptionByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid subscription ID")
		return
	}

	result, err := h.service.GetSubscriptionByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgSubscriptionNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionRetrieved, result)
}

// @Summary      Update subscription
// @Description  Memperbarui subscription
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        id            path      int                                     true  "Subscription ID"
// @Param        subscription  body      subscription.UpdateSubscriptionRequest  true  "Subscription data yang akan diupdate"
// @Success      200           {object}  response.Response{data=subscription.SubscriptionResponse}  "Subscription berhasil diupdate"
// @Failure      400           {object}  response.Response  "Bad request - Invalid subscription ID atau validation failed"
// @Failure      404           {object}  response.Response  "Subscription tidak ditemukan"
// @Failure      500           {object}  response.Response  "Internal server error"
// @Router       /api/v1/subscriptions/{id} [put]
// @Security     BearerAuth
func (h *Handler) UpdateSubscription(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid subscription ID")
		return
	}

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*UpdateSubscriptionRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.UpdateSubscription(id, req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionUpdated, result)
}

// @Summary      Get company subscription
// @Description  Mendapatkan subscription aktif dari company
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Company ID"
// @Success      200  {object}  response.Response{data=subscription.SubscriptionResponse}  "Company subscription berhasil diambil"
// @Failure      400  {object}  response.Response  "Bad request - Invalid company ID"
// @Failure      404  {object}  response.Response  "Subscription tidak ditemukan"
// @Router       /api/v1/companies/{id}/subscription [get]
// @Security     BearerAuth
func (h *Handler) GetCompanySubscription(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid company ID")
		return
	}

	result, err := h.service.GetCompanySubscription(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgSubscriptionNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionRetrieved, result)
}

// Plan Modules Management Handlers

// @Summary      Get plan modules (Admin)
// @Description  Mendapatkan daftar modules yang termasuk dalam subscription plan (admin only)
// @Tags         Subscription Plans (Admin)
// @Accept       json
// @Produce      json
// @Param        plan_id  path      int  true  "Subscription Plan ID"
// @Success      200      {object}  response.Response  "Plan modules berhasil diambil"
// @Failure      400      {object}  response.Response  "Bad request - Invalid plan ID"
// @Failure      500      {object}  response.Response  "Internal server error"
// @Router       /api/v1/admin/plan-modules/{plan_id} [get]
// @Security     BearerAuth
func (h *Handler) GetPlanModules(c *gin.Context) {
	planID, err := strconv.ParseInt(c.Param("plan_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid plan ID")
		return
	}

	result, err := h.service.GetPlanModules(planID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get plan modules", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Plan modules retrieved successfully", result)
}

// @Summary      Add modules to plan (Admin)
// @Description  Menambahkan modules ke subscription plan (admin only)
// @Tags         Subscription Plans (Admin)
// @Accept       json
// @Produce      json
// @Param        plan_id  path      int                                     true  "Subscription Plan ID"
// @Param        modules  body      subscription.AddModulesToPlanRequest    true  "Modules to add"
// @Success      200      {object}  response.Response  "Modules berhasil ditambahkan ke plan"
// @Failure      400      {object}  response.Response  "Bad request - Invalid plan ID atau validation failed"
// @Failure      404      {object}  response.Response  "Subscription plan tidak ditemukan"
// @Failure      500      {object}  response.Response  "Internal server error"
// @Router       /api/v1/admin/plan-modules/{plan_id} [post]
// @Security     BearerAuth
func (h *Handler) AddModulesToPlan(c *gin.Context) {
	planID, err := strconv.ParseInt(c.Param("plan_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid plan ID")
		return
	}

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*AddModulesToPlanRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	if err := h.service.AddModulesToPlan(planID, req); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to add modules to plan", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Modules successfully added to plan", nil)
}

// @Summary      Remove module from plan (Admin)
// @Description  Menghapus module dari subscription plan (admin only)
// @Tags         Subscription Plans (Admin)
// @Accept       json
// @Produce      json
// @Param        plan_id    path      int  true  "Subscription Plan ID"
// @Param        module_id  path      int  true  "Module ID"
// @Success      200        {object}  response.Response  "Module berhasil dihapus dari plan"
// @Failure      400        {object}  response.Response  "Bad request - Invalid ID"
// @Failure      404        {object}  response.Response  "Plan atau module tidak ditemukan"
// @Failure      500        {object}  response.Response  "Internal server error"
// @Router       /api/v1/admin/plan-modules/{plan_id}/{module_id} [delete]
// @Security     BearerAuth
func (h *Handler) RemoveModuleFromPlan(c *gin.Context) {
	planID, err := strconv.ParseInt(c.Param("plan_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid plan ID")
		return
	}

	moduleID, err := strconv.ParseInt(c.Param("module_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid module ID")
		return
	}

	if err := h.service.RemoveModuleFromPlan(planID, moduleID); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to remove module from plan", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Module successfully removed from plan", nil)
}

// Route registration
func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	// Public plan routes (read-only)
	plans := router.Group("/subscription-plans")
	{
		// GET /api/v1/subscription-plans - Get all subscription plans
		plans.GET("", handler.GetAllPlans)

		// GET /api/v1/subscription-plans/:id - Get subscription plan by ID
		plans.GET("/:id", handler.GetPlanByID)
	}
}

// RegisterProtectedRoutes registers protected subscription routes
func RegisterProtectedRoutes(router *gin.RouterGroup, handler *Handler) {
	// Admin/Protected plan routes
	adminPlans := router.Group("/admin/subscription-plans")
	{
		// POST /api/v1/admin/subscription-plans - Create new subscription plan
		adminPlans.POST("",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &CreateSubscriptionPlanRequest{},
			}),
			handler.CreateSubscriptionPlan,
		)

		// PUT /api/v1/admin/subscription-plans/:id - Update subscription plan by ID
		adminPlans.PUT("/:id",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &UpdateSubscriptionPlanRequest{},
			}),
			handler.UpdateSubscriptionPlan,
		)

		// DELETE /api/v1/admin/subscription-plans/:id - Delete subscription plan by ID
		adminPlans.DELETE("/:id", handler.DeleteSubscriptionPlan)
	}

	// Plan modules management (separate group to avoid conflicts)
	planModules := router.Group("/admin/plan-modules")
	{
		// GET /api/v1/admin/plan-modules/:plan_id - Get modules for a plan
		planModules.GET("/:plan_id", handler.GetPlanModules)

		// POST /api/v1/admin/plan-modules/:plan_id - Add modules to a plan
		planModules.POST("/:plan_id",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &AddModulesToPlanRequest{},
			}),
			handler.AddModulesToPlan,
		)

		// DELETE /api/v1/admin/plan-modules/:plan_id/:module_id - Remove module from plan
		planModules.DELETE("/:plan_id/:module_id", handler.RemoveModuleFromPlan)
	}

	// Subscription routes (all protected)
	subscriptions := router.Group("/subscriptions")
	{
		// GET /api/v1/subscriptions - Get all subscriptions with optional filters
		subscriptions.GET("", handler.GetAllSubscriptions)

		// POST /api/v1/subscriptions - Create new subscription
		subscriptions.POST("",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &CreateSubscriptionRequest{},
			}),
			handler.CreateSubscription,
		)

		// GET /api/v1/subscriptions/:id - Get subscription by ID
		subscriptions.GET("/:id", handler.GetSubscriptionByID)

		// PUT /api/v1/subscriptions/:id - Update subscription by ID
		subscriptions.PUT("/:id",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &UpdateSubscriptionRequest{},
			}),
			handler.UpdateSubscription,
		)
	}

	// Company subscription routes
	companies := router.Group("/companies")
	{
		// GET /api/v1/companies/:id/subscription - Get company subscription
		companies.GET("/:id/subscription", handler.GetCompanySubscription)
	}
}
