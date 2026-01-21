package subscription

import (
	"net/http"
	"strconv"

	"gin-scalable-api/internal/constants"
	"gin-scalable-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAllPlans(c *gin.Context) {
	result, err := h.service.GetSubscriptionPlans()
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionPlansRetrieved, result)
}

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

func (h *Handler) CreateSubscriptionPlan(c *gin.Context) {
	var req CreateSubscriptionPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.service.CreateSubscriptionPlan(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create subscription plan", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgSubscriptionPlanCreated, result)
}

func (h *Handler) UpdateSubscriptionPlan(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid plan ID")
		return
	}

	var req UpdateSubscriptionPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.service.UpdateSubscriptionPlan(id, &req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionPlanUpdated, result)
}

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

func (h *Handler) CreateSubscription(c *gin.Context) {
	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.service.CreateSubscription(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgSubscriptionCreated, result)
}

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

func (h *Handler) UpdateSubscription(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid subscription ID")
		return
	}

	var req UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.service.UpdateSubscription(id, &req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgSubscriptionUpdated, result)
}

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

func (h *Handler) AddModulesToPlan(c *gin.Context) {
	planID, err := strconv.ParseInt(c.Param("plan_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid plan ID")
		return
	}

	var req AddModulesToPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	if err := h.service.AddModulesToPlan(planID, &req); err != nil {
		response.ErrorWithAutoStatus(c, "Failed to add modules to plan", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Modules successfully added to plan", nil)
}

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
