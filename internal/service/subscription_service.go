package service

import (
	"gin-scalable-api/internal/models"
	"gin-scalable-api/internal/repository"
	"time"
)

type SubscriptionService struct {
	subscriptionRepo *repository.SubscriptionRepository
}

func NewSubscriptionService(subscriptionRepo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
	}
}

type SubscriptionPlanResponse struct {
	ID           int64                  `json:"id"`
	Name         string                 `json:"name"`
	DisplayName  string                 `json:"display_name"`
	Description  string                 `json:"description"`
	PriceMonthly float64                `json:"price_monthly"`
	PriceYearly  float64                `json:"price_yearly"`
	MaxUsers     *int                   `json:"max_users"`
	MaxBranches  *int                   `json:"max_branches"`
	Features     map[string]interface{} `json:"features"`
	IsActive     bool                   `json:"is_active"`
}

type SubscriptionResponse struct {
	ID              int64     `json:"id"`
	CompanyID       int64     `json:"company_id"`
	PlanID          int64     `json:"plan_id"`
	Status          string    `json:"status"`
	BillingCycle    string    `json:"billing_cycle"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	AutoRenew       bool      `json:"auto_renew"`
	CompanyName     string    `json:"company_name,omitempty"`
	PlanDisplayName string    `json:"plan_display_name,omitempty"`
}

type CreateSubscriptionRequest struct {
	CompanyID    int64  `json:"company_id" binding:"required"`
	PlanID       int64  `json:"plan_id" binding:"required"`
	BillingCycle string `json:"billing_cycle" binding:"required"`
	AutoRenew    bool   `json:"auto_renew"`
}

type UpdateSubscriptionRequest struct {
	PlanID    *int64 `json:"plan_id"`
	AutoRenew *bool  `json:"auto_renew"`
}

type RenewSubscriptionRequest struct {
	BillingCycle string `json:"billing_cycle" binding:"required"`
	PlanID       *int64 `json:"plan_id"`
}

type CancelSubscriptionRequest struct {
	Reason            string `json:"reason"`
	CancelImmediately bool   `json:"cancel_immediately"`
}

type SubscriptionListRequest struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

func (s *SubscriptionService) GetAllPlans() ([]*SubscriptionPlanResponse, error) {
	plans, err := s.subscriptionRepo.GetAllPlans()
	if err != nil {
		return nil, err
	}

	var response []*SubscriptionPlanResponse
	for _, plan := range plans {
		features := make(map[string]interface{})
		if plan.Features != nil {
			features = plan.Features
		}

		response = append(response, &SubscriptionPlanResponse{
			ID:           plan.ID,
			Name:         plan.Name,
			DisplayName:  plan.DisplayName,
			Description:  plan.Description,
			PriceMonthly: plan.PriceMonthly,
			PriceYearly:  plan.PriceYearly,
			MaxUsers:     plan.MaxUsers,
			MaxBranches:  plan.MaxBranches,
			Features:     features,
			IsActive:     plan.IsActive,
		})
	}

	return response, nil
}

func (s *SubscriptionService) GetPlanByID(id int64) (*SubscriptionPlanResponse, error) {
	plan, err := s.subscriptionRepo.GetPlanByID(id)
	if err != nil {
		return nil, err
	}

	features := make(map[string]interface{})
	if plan.Features != nil {
		features = plan.Features
	}

	return &SubscriptionPlanResponse{
		ID:           plan.ID,
		Name:         plan.Name,
		DisplayName:  plan.DisplayName,
		Description:  plan.Description,
		PriceMonthly: plan.PriceMonthly,
		PriceYearly:  plan.PriceYearly,
		MaxUsers:     plan.MaxUsers,
		MaxBranches:  plan.MaxBranches,
		Features:     features,
		IsActive:     plan.IsActive,
	}, nil
}

func (s *SubscriptionService) GetAllSubscriptions(req *SubscriptionListRequest) ([]*SubscriptionResponse, error) {
	subscriptions, err := s.subscriptionRepo.GetAllSubscriptions(req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	var response []*SubscriptionResponse
	for _, sub := range subscriptions {
		response = append(response, &SubscriptionResponse{
			ID:              sub.ID,
			CompanyID:       sub.CompanyID,
			PlanID:          sub.PlanID,
			Status:          sub.Status,
			BillingCycle:    sub.BillingCycle,
			StartDate:       sub.StartDate,
			EndDate:         sub.EndDate,
			AutoRenew:       sub.AutoRenew,
			CompanyName:     sub.CompanyName,
			PlanDisplayName: sub.PlanDisplayName,
		})
	}

	return response, nil
}

func (s *SubscriptionService) GetSubscriptionByID(id int64) (*SubscriptionResponse, error) {
	sub, err := s.subscriptionRepo.GetSubscriptionByID(id)
	if err != nil {
		return nil, err
	}

	return &SubscriptionResponse{
		ID:              sub.ID,
		CompanyID:       sub.CompanyID,
		PlanID:          sub.PlanID,
		Status:          sub.Status,
		BillingCycle:    sub.BillingCycle,
		StartDate:       sub.StartDate,
		EndDate:         sub.EndDate,
		AutoRenew:       sub.AutoRenew,
		CompanyName:     sub.CompanyName,
		PlanDisplayName: sub.PlanDisplayName,
	}, nil
}

func (s *SubscriptionService) CreateSubscription(req *CreateSubscriptionRequest) (*SubscriptionResponse, error) {
	// Get the plan to determine the price
	plan, err := s.subscriptionRepo.GetPlanByID(req.PlanID)
	if err != nil {
		return nil, err
	}

	var endDate time.Time
	var price float64
	var nextPaymentDate time.Time

	if req.BillingCycle == "yearly" {
		endDate = time.Now().AddDate(1, 0, 0)
		price = plan.PriceYearly
		nextPaymentDate = time.Now().AddDate(1, 0, 0) // Next payment in 1 year
	} else {
		endDate = time.Now().AddDate(0, 1, 0)
		price = plan.PriceMonthly
		nextPaymentDate = time.Now().AddDate(0, 1, 0) // Next payment in 1 month
	}

	subscription := &models.Subscription{
		CompanyID:       req.CompanyID,
		PlanID:          req.PlanID,
		Status:          "active",
		BillingCycle:    req.BillingCycle,
		StartDate:       time.Now(),
		EndDate:         endDate,
		Price:           price,
		Currency:        "IDR",
		PaymentStatus:   "pending",
		NextPaymentDate: &nextPaymentDate,
		AutoRenew:       req.AutoRenew,
	}

	if err := s.subscriptionRepo.Create(subscription); err != nil {
		return nil, err
	}

	// Get the created subscription with company and plan names
	createdSub, err := s.subscriptionRepo.GetSubscriptionByID(subscription.ID)
	if err != nil {
		return nil, err
	}

	return &SubscriptionResponse{
		ID:              createdSub.ID,
		CompanyID:       createdSub.CompanyID,
		PlanID:          createdSub.PlanID,
		Status:          createdSub.Status,
		BillingCycle:    createdSub.BillingCycle,
		StartDate:       createdSub.StartDate,
		EndDate:         createdSub.EndDate,
		AutoRenew:       createdSub.AutoRenew,
		CompanyName:     createdSub.CompanyName,
		PlanDisplayName: createdSub.PlanDisplayName,
	}, nil
}

func (s *SubscriptionService) UpdateSubscription(id int64, req *UpdateSubscriptionRequest) (*SubscriptionResponse, error) {
	subscription, err := s.subscriptionRepo.GetSubscriptionByID(id)
	if err != nil {
		return nil, err
	}

	// If plan is being changed, update the price accordingly
	if req.PlanID != nil && *req.PlanID != subscription.PlanID {
		plan, err := s.subscriptionRepo.GetPlanByID(*req.PlanID)
		if err != nil {
			return nil, err
		}

		subscription.PlanID = *req.PlanID

		// Update price based on current billing cycle
		if subscription.BillingCycle == "yearly" {
			subscription.Price = plan.PriceYearly
		} else {
			subscription.Price = plan.PriceMonthly
		}
	}

	if req.AutoRenew != nil {
		subscription.AutoRenew = *req.AutoRenew
	}

	// Ensure required fields are set (in case they're missing from DB)
	if subscription.Currency == "" {
		subscription.Currency = "IDR"
	}
	if subscription.PaymentStatus == "" {
		subscription.PaymentStatus = "pending"
	}

	if err := s.subscriptionRepo.Update(subscription); err != nil {
		return nil, err
	}

	// Get updated subscription
	updatedSub, err := s.subscriptionRepo.GetSubscriptionByID(id)
	if err != nil {
		return nil, err
	}

	return &SubscriptionResponse{
		ID:              updatedSub.ID,
		CompanyID:       updatedSub.CompanyID,
		PlanID:          updatedSub.PlanID,
		Status:          updatedSub.Status,
		BillingCycle:    updatedSub.BillingCycle,
		StartDate:       updatedSub.StartDate,
		EndDate:         updatedSub.EndDate,
		AutoRenew:       updatedSub.AutoRenew,
		CompanyName:     updatedSub.CompanyName,
		PlanDisplayName: updatedSub.PlanDisplayName,
	}, nil
}

func (s *SubscriptionService) RenewSubscription(id int64, req *RenewSubscriptionRequest) error {
	return s.subscriptionRepo.RenewSubscription(id, req.PlanID, req.BillingCycle)
}

func (s *SubscriptionService) CancelSubscription(id int64, req *CancelSubscriptionRequest) error {
	return s.subscriptionRepo.CancelSubscription(id, req.Reason, req.CancelImmediately)
}

func (s *SubscriptionService) GetCompanySubscription(companyID int64) (*SubscriptionResponse, error) {
	sub, err := s.subscriptionRepo.GetCompanySubscription(companyID)
	if err != nil {
		return nil, err
	}

	return &SubscriptionResponse{
		ID:              sub.ID,
		CompanyID:       sub.CompanyID,
		PlanID:          sub.PlanID,
		Status:          sub.Status,
		BillingCycle:    sub.BillingCycle,
		StartDate:       sub.StartDate,
		EndDate:         sub.EndDate,
		AutoRenew:       sub.AutoRenew,
		CompanyName:     sub.CompanyName,
		PlanDisplayName: sub.PlanDisplayName,
	}, nil
}

func (s *SubscriptionService) CheckModuleAccess(companyID, moduleID int64) (bool, error) {
	return s.subscriptionRepo.CheckModuleAccess(companyID, moduleID)
}

func (s *SubscriptionService) GetExpiringSubscriptions(days int) ([]*SubscriptionResponse, error) {
	subscriptions, err := s.subscriptionRepo.GetExpiringSubscriptions(days)
	if err != nil {
		return nil, err
	}

	var response []*SubscriptionResponse
	for _, sub := range subscriptions {
		response = append(response, &SubscriptionResponse{
			ID:              sub.ID,
			CompanyID:       sub.CompanyID,
			PlanID:          sub.PlanID,
			Status:          sub.Status,
			BillingCycle:    sub.BillingCycle,
			StartDate:       sub.StartDate,
			EndDate:         sub.EndDate,
			AutoRenew:       sub.AutoRenew,
			CompanyName:     sub.CompanyName,
			PlanDisplayName: sub.PlanDisplayName,
		})
	}

	return response, nil
}

func (s *SubscriptionService) UpdateExpiredSubscriptions() error {
	return s.subscriptionRepo.UpdateExpiredSubscriptions()
}
func (s *SubscriptionService) GetSubscriptionStats() (map[string]interface{}, error) {
	// Get basic subscription statistics
	stats := map[string]interface{}{
		"total_subscriptions":     0,
		"active_subscriptions":    0,
		"expired_subscriptions":   0,
		"cancelled_subscriptions": 0,
		"total_revenue":           0.0,
	}

	// This is a placeholder implementation
	// In a real implementation, you would query the database for actual statistics
	return stats, nil
}

// MarkPaymentAsPaid marks a subscription payment as paid
func (s *SubscriptionService) MarkPaymentAsPaid(subscriptionID int64) error {
	return s.subscriptionRepo.MarkPaymentAsPaid(subscriptionID)
}
