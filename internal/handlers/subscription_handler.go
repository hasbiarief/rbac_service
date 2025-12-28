package handlers

import (
	"gin-scalable-api/internal/service"
	"strconv"

	"gin-scalable-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	subscriptionService *service.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// Public endpoints (no auth required)
func (h *SubscriptionHandler) GetAllPlans(c *gin.Context) {
	result, err := h.subscriptionService.GetAllPlans()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *SubscriptionHandler) GetPlanByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid plan ID")
		return
	}

	result, err := h.subscriptionService.GetPlanByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

// Protected endpoints
func (h *SubscriptionHandler) GetAllSubscriptions(c *gin.Context) {
	var req service.SubscriptionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.subscriptionService.GetAllSubscriptions(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req service.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.subscriptionService.CreateSubscription(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Created successfully", result)
}

func (h *SubscriptionHandler) GetSubscriptionByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid subscription ID")
		return
	}

	result, err := h.subscriptionService.GetSubscriptionByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid subscription ID")
		return
	}

	var req service.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.subscriptionService.UpdateSubscription(id, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *SubscriptionHandler) RenewSubscription(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid subscription ID")
		return
	}

	var req service.RenewSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	if err := h.subscriptionService.RenewSubscription(id, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Subscription renewed successfully", nil)
}

func (h *SubscriptionHandler) CancelSubscription(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid subscription ID")
		return
	}

	var req service.CancelSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	if err := h.subscriptionService.CancelSubscription(id, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Subscription cancelled successfully", nil)
}

func (h *SubscriptionHandler) GetCompanySubscription(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid company ID")
		return
	}

	result, err := h.subscriptionService.GetCompanySubscription(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *SubscriptionHandler) GetCompanySubscriptionStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid company ID")
		return
	}

	result, err := h.subscriptionService.GetCompanySubscription(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
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
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Module access checked successfully", gin.H{
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
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *SubscriptionHandler) UpdateExpiredSubscriptions(c *gin.Context) {
	if err := h.subscriptionService.UpdateExpiredSubscriptions(); err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Expired subscriptions updated successfully", nil)
}

func (h *SubscriptionHandler) GetSubscriptionStats(c *gin.Context) {
	result, err := h.subscriptionService.GetSubscriptionStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}
