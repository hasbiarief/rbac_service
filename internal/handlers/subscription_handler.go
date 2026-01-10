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

type SubscriptionHandler struct {
	subscriptionService interfaces.SubscriptionServiceInterface
}

func NewSubscriptionHandler(subscriptionService interfaces.SubscriptionServiceInterface) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// Public endpoints (no auth required)
func (h *SubscriptionHandler) GetAllPlans(c *gin.Context) {
	result, err := h.subscriptionService.GetSubscriptionPlans()
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionPlansRetrieved, result)
}

func (h *SubscriptionHandler) GetPlanByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid plan ID")
		return
	}

	result, err := h.subscriptionService.GetSubscriptionPlanByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgSubscriptionPlanNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionPlanRetrieved, result)
}

// Admin endpoints for subscription plan management
func (h *SubscriptionHandler) CreateSubscriptionPlan(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert ke DTO yang diharapkan
	createReq, ok := validatedBody.(*dto.CreateSubscriptionPlanRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.subscriptionService.CreateSubscriptionPlan(createReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create subscription plan", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgSubscriptionPlanCreated, result)
}

func (h *SubscriptionHandler) UpdateSubscriptionPlan(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid plan ID")
		return
	}

	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert ke DTO yang diharapkan
	updateReq, ok := validatedBody.(*dto.UpdateSubscriptionPlanRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.subscriptionService.UpdateSubscriptionPlan(id, updateReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionPlanUpdated, result)
}

func (h *SubscriptionHandler) DeleteSubscriptionPlan(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid plan ID")
		return
	}

	if err := h.subscriptionService.DeleteSubscriptionPlan(id); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionPlanDeleted, nil)
}

// Protected endpoints
func (h *SubscriptionHandler) GetAllSubscriptions(c *gin.Context) {
	var req dto.SubscriptionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.subscriptionService.GetSubscriptions(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionsRetrieved, result)
}

func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert to the expected DTO struct
	req, ok := validatedBody.(*dto.CreateSubscriptionBasicRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	// Convert to full DTO request with default values
	createReq := &dto.CreateSubscriptionRequest{
		CompanyID:    req.CompanyID,
		PlanID:       req.PlanID,
		BillingCycle: req.BillingCycle,
		// Default values will be set by service
	}

	result, err := h.subscriptionService.CreateSubscription(createReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgSubscriptionCreated, result)
}

func (h *SubscriptionHandler) GetSubscriptionByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid subscription ID")
		return
	}

	result, err := h.subscriptionService.GetSubscriptionByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgSubscriptionNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionRetrieved, result)
}

func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid subscription ID")
		return
	}

	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert to the expected DTO struct
	req, ok := validatedBody.(*dto.UpdateSubscriptionRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.subscriptionService.UpdateSubscription(id, req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionUpdated, result)
}

func (h *SubscriptionHandler) RenewSubscription(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid subscription ID")
		return
	}

	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert ke DTO yang diharapkan
	req, ok := validatedBody.(*dto.RenewSubscriptionRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.subscriptionService.RenewSubscription(id, req.PlanID, req.BillingCycle)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionRenewed, result)
}

func (h *SubscriptionHandler) CancelSubscription(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid subscription ID")
		return
	}

	// The interface method only accepts ID, additional parameters would need interface update
	if err := h.subscriptionService.CancelSubscription(id); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionCancelled, nil)
}

func (h *SubscriptionHandler) GetCompanySubscription(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid company ID")
		return
	}

	result, err := h.subscriptionService.GetCompanySubscription(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgSubscriptionNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionRetrieved, result)
}

func (h *SubscriptionHandler) GetCompanySubscriptionStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid company ID")
		return
	}

	result, err := h.subscriptionService.GetCompanySubscription(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgSubscriptionNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionRetrieved, result)
}

func (h *SubscriptionHandler) CheckModuleAccess(c *gin.Context) {
	companyID, err := strconv.ParseInt(c.Param("companyId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid company ID")
		return
	}

	moduleID, err := strconv.ParseInt(c.Param("moduleId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid module ID")
		return
	}

	hasAccess, err := h.subscriptionService.CheckModuleAccess(companyID, moduleID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Module access successfully checked", gin.H{
		"has_access": hasAccess,
	})
}

func (h *SubscriptionHandler) GetExpiringSubscriptions(c *gin.Context) {
	days := 30
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil {
			days = parsed
		}
	}

	result, err := h.subscriptionService.GetExpiringSubscriptions(days)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Expiring subscriptions successfully retrieved", result)
}

func (h *SubscriptionHandler) UpdateExpiredSubscriptions(c *gin.Context) {
	if err := h.subscriptionService.UpdateExpiredSubscriptions(); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Expired subscriptions successfully updated", nil)
}

func (h *SubscriptionHandler) GetSubscriptionStats(c *gin.Context) {
	result, err := h.subscriptionService.GetSubscriptionStats()
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Subscription statistics successfully retrieved", result)
}

func (h *SubscriptionHandler) MarkPaymentAsPaid(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid subscription ID")
		return
	}

	if err := h.subscriptionService.MarkPaymentAsPaid(id); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Payment successfully marked as paid", nil)
}
