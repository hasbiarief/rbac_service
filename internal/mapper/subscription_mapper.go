package mapper

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/models"
	"time"
)

// SubscriptionMapper handles conversion between subscription models and DTOs
type SubscriptionMapper struct{}

// NewSubscriptionMapper creates a new subscription mapper
func NewSubscriptionMapper() *SubscriptionMapper {
	return &SubscriptionMapper{}
}

// ToPlanResponse converts plan model to response DTO
func (m *SubscriptionMapper) ToPlanResponse(plan *models.SubscriptionPlan) *dto.SubscriptionPlanResponse {
	if plan == nil {
		return nil
	}

	return &dto.SubscriptionPlanResponse{
		ID:           plan.ID,
		Name:         plan.Name,
		DisplayName:  plan.DisplayName,
		Description:  plan.Description,
		PriceMonthly: plan.PriceMonthly,
		PriceYearly:  plan.PriceYearly,
		MaxUsers:     plan.MaxUsers,
		MaxBranches:  plan.MaxBranches,
		Features:     map[string]interface{}(plan.Features),
		IsActive:     plan.IsActive,
		CreatedAt:    plan.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    plan.UpdatedAt.Format(time.RFC3339),
	}
}

// ToPlanResponseList converts plan model slice to response DTO slice
func (m *SubscriptionMapper) ToPlanResponseList(plans []*models.SubscriptionPlan) []*dto.SubscriptionPlanResponse {
	if plans == nil {
		return nil
	}

	responses := make([]*dto.SubscriptionPlanResponse, len(plans))
	for i, plan := range plans {
		responses[i] = m.ToPlanResponse(plan)
	}
	return responses
}

// ToSubscriptionResponse converts subscription model to response DTO
func (m *SubscriptionMapper) ToSubscriptionResponse(subscription *models.Subscription) *dto.SubscriptionResponse {
	if subscription == nil {
		return nil
	}

	var lastPaymentDate, nextPaymentDate *string
	if subscription.LastPaymentDate != nil {
		formatted := subscription.LastPaymentDate.Format("2006-01-02")
		lastPaymentDate = &formatted
	}
	if subscription.NextPaymentDate != nil {
		formatted := subscription.NextPaymentDate.Format("2006-01-02")
		nextPaymentDate = &formatted
	}

	return &dto.SubscriptionResponse{
		ID:              subscription.ID,
		CompanyID:       subscription.CompanyID,
		PlanID:          subscription.PlanID,
		Status:          subscription.Status,
		BillingCycle:    subscription.BillingCycle,
		StartDate:       subscription.StartDate.Format("2006-01-02"),
		EndDate:         subscription.EndDate.Format("2006-01-02"),
		Price:           subscription.Price,
		Currency:        subscription.Currency,
		PaymentStatus:   subscription.PaymentStatus,
		LastPaymentDate: lastPaymentDate,
		NextPaymentDate: nextPaymentDate,
		AutoRenew:       subscription.AutoRenew,
		CompanyName:     subscription.CompanyName,
		PlanDisplayName: subscription.PlanDisplayName,
		CreatedAt:       subscription.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       subscription.UpdatedAt.Format(time.RFC3339),
	}
}

// ToSubscriptionResponseList converts subscription model slice to response DTO slice
func (m *SubscriptionMapper) ToSubscriptionResponseList(subscriptions []*models.Subscription) []*dto.SubscriptionResponse {
	if subscriptions == nil {
		return nil
	}

	responses := make([]*dto.SubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		responses[i] = m.ToSubscriptionResponse(subscription)
	}
	return responses
}

// ToPlanModel converts create plan request DTO to model
func (m *SubscriptionMapper) ToPlanModel(req *dto.CreateSubscriptionPlanRequest) *models.SubscriptionPlan {
	if req == nil {
		return nil
	}

	return &models.SubscriptionPlan{
		Name:         req.Name,
		DisplayName:  req.DisplayName,
		Description:  req.Description,
		PriceMonthly: req.PriceMonthly,
		PriceYearly:  req.PriceYearly,
		MaxUsers:     req.MaxUsers,
		MaxBranches:  req.MaxBranches,
		Features:     req.Features,
		IsActive:     true, // Default to active
	}
}

// ToSubscriptionModel converts create subscription request DTO to model
func (m *SubscriptionMapper) ToSubscriptionModel(req *dto.CreateSubscriptionRequest) (*models.Subscription, error) {
	if req == nil {
		return nil, nil
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, err
	}

	return &models.Subscription{
		CompanyID:    req.CompanyID,
		PlanID:       req.PlanID,
		Status:       "active", // Default status
		BillingCycle: req.BillingCycle,
		StartDate:    startDate,
		EndDate:      endDate,
		Price:        req.Price,
		Currency:     req.Currency,
		AutoRenew:    req.AutoRenew,
	}, nil
}

// UpdatePlanModel updates plan model with update request DTO
func (m *SubscriptionMapper) UpdatePlanModel(plan *models.SubscriptionPlan, req *dto.UpdateSubscriptionPlanRequest) {
	if plan == nil || req == nil {
		return
	}

	if req.Name != "" {
		plan.Name = req.Name
	}
	if req.DisplayName != "" {
		plan.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		plan.Description = req.Description
	}
	if req.PriceMonthly != nil {
		plan.PriceMonthly = *req.PriceMonthly
	}
	if req.PriceYearly != nil {
		plan.PriceYearly = *req.PriceYearly
	}
	if req.MaxUsers != nil {
		plan.MaxUsers = req.MaxUsers
	}
	if req.MaxBranches != nil {
		plan.MaxBranches = req.MaxBranches
	}
	if req.Features != nil {
		plan.Features = req.Features
	}
	if req.IsActive != nil {
		plan.IsActive = *req.IsActive
	}
}

// UpdateSubscriptionModel updates subscription model with update request DTO
func (m *SubscriptionMapper) UpdateSubscriptionModel(subscription *models.Subscription, req *dto.UpdateSubscriptionRequest) error {
	if subscription == nil || req == nil {
		return nil
	}

	if req.Status != "" {
		subscription.Status = req.Status
	}
	if req.BillingCycle != "" {
		subscription.BillingCycle = req.BillingCycle
	}
	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return err
		}
		subscription.EndDate = endDate
	}
	if req.Price != nil {
		subscription.Price = *req.Price
	}
	if req.PaymentStatus != "" {
		subscription.PaymentStatus = req.PaymentStatus
	}
	if req.NextPaymentDate != "" {
		nextPaymentDate, err := time.Parse("2006-01-02", req.NextPaymentDate)
		if err != nil {
			return err
		}
		subscription.NextPaymentDate = &nextPaymentDate
	}
	if req.AutoRenew != nil {
		subscription.AutoRenew = *req.AutoRenew
	}

	return nil
}

// ToPlanListResponse creates paginated plan list response
func (m *SubscriptionMapper) ToPlanListResponse(plans []*models.SubscriptionPlan, total int64, limit, offset int) *dto.SubscriptionPlanListResponse {
	return &dto.SubscriptionPlanListResponse{
		Data:    m.ToPlanResponseList(plans),
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(plans)) < total,
	}
}

// ToSubscriptionListResponse creates paginated subscription list response
func (m *SubscriptionMapper) ToSubscriptionListResponse(subscriptions []*models.Subscription, total int64, limit, offset int) *dto.SubscriptionListResponse {
	return &dto.SubscriptionListResponse{
		Data:    m.ToSubscriptionResponseList(subscriptions),
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(subscriptions)) < total,
	}
}
