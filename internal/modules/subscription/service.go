package subscription

import (
	"errors"
	"fmt"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetSubscriptionPlans() ([]*SubscriptionPlanResponse, error) {
	plans, err := s.repo.GetAllPlans()
	if err != nil {
		return nil, err
	}

	var responses []*SubscriptionPlanResponse
	for _, plan := range plans {
		responses = append(responses, toPlanResponse(plan))
	}

	return responses, nil
}

func (s *Service) GetSubscriptionPlanByID(id int64) (*SubscriptionPlanResponse, error) {
	plan, err := s.repo.GetPlanByID(id)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, errors.New("subscription plan not found")
	}

	return toPlanResponse(plan), nil
}

func (s *Service) CreateSubscriptionPlan(req *CreateSubscriptionPlanRequest) (*SubscriptionPlanResponse, error) {
	plan := &SubscriptionPlan{
		Name:         req.Name,
		DisplayName:  req.DisplayName,
		Description:  req.Description,
		PriceMonthly: req.PriceMonthly,
		PriceYearly:  req.PriceYearly,
		MaxUsers:     req.MaxUsers,
		MaxBranches:  req.MaxBranches,
		Features:     req.Features,
		IsActive:     true,
	}

	if err := s.repo.CreatePlan(plan); err != nil {
		return nil, err
	}

	return toPlanResponse(plan), nil
}

func (s *Service) UpdateSubscriptionPlan(id int64, req *UpdateSubscriptionPlanRequest) (*SubscriptionPlanResponse, error) {
	plan, err := s.repo.GetPlanByID(id)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, errors.New("subscription plan not found")
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

	if err := s.repo.UpdatePlan(plan); err != nil {
		return nil, err
	}

	return toPlanResponse(plan), nil
}

func (s *Service) DeleteSubscriptionPlan(id int64) error {
	return s.repo.DeletePlan(id)
}

func (s *Service) GetSubscriptions(req *SubscriptionListRequest) (*SubscriptionListResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	filters := make(map[string]interface{})
	if req.CompanyID != nil {
		filters["company_id"] = *req.CompanyID
	}
	if req.PlanID != nil {
		filters["plan_id"] = *req.PlanID
	}
	if req.Status != "" {
		filters["status"] = req.Status
	}
	if req.Search != "" {
		filters["search"] = req.Search
	}

	subscriptions, err := s.repo.GetAll(limit, offset, filters)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(filters)
	if err != nil {
		return nil, err
	}

	var responses []*SubscriptionResponse
	for _, sub := range subscriptions {
		responses = append(responses, toSubscriptionResponse(sub))
	}

	return &SubscriptionListResponse{
		Data:    responses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(responses)) < total,
	}, nil
}

func (s *Service) GetSubscriptionByID(id int64) (*SubscriptionResponse, error) {
	sub, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, errors.New("subscription not found")
	}

	return toSubscriptionResponse(sub), nil
}

func (s *Service) CreateSubscription(req *CreateSubscriptionRequest) (*SubscriptionResponse, error) {
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	sub := &Subscription{
		CompanyID:     req.CompanyID,
		PlanID:        req.PlanID,
		Status:        "active",
		BillingCycle:  req.BillingCycle,
		StartDate:     startDate,
		EndDate:       endDate,
		Price:         req.Price,
		Currency:      req.Currency,
		PaymentStatus: "pending",
		AutoRenew:     req.AutoRenew,
	}

	if err := s.repo.Create(sub); err != nil {
		return nil, err
	}

	return toSubscriptionResponse(sub), nil
}

func (s *Service) UpdateSubscription(id int64, req *UpdateSubscriptionRequest) (*SubscriptionResponse, error) {
	sub, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, errors.New("subscription not found")
	}

	if req.PlanID != nil {
		sub.PlanID = *req.PlanID
	}
	if req.Status != "" {
		sub.Status = req.Status
	}
	if req.BillingCycle != "" {
		sub.BillingCycle = req.BillingCycle
	}
	if req.EndDate != "" {
		endDate, _ := time.Parse("2006-01-02", req.EndDate)
		sub.EndDate = endDate
	}
	if req.Price != nil {
		sub.Price = *req.Price
	}
	if req.PaymentStatus != "" {
		sub.PaymentStatus = req.PaymentStatus
	}
	if req.AutoRenew != nil {
		sub.AutoRenew = *req.AutoRenew
	}

	if err := s.repo.Update(sub); err != nil {
		return nil, err
	}

	return toSubscriptionResponse(sub), nil
}

func (s *Service) RenewSubscription(id int64, planID *int64, billingCycle string) (*SubscriptionResponse, error) {
	sub, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, errors.New("subscription not found")
	}

	if planID != nil {
		sub.PlanID = *planID
	}
	sub.BillingCycle = billingCycle
	sub.Status = "active"
	sub.StartDate = time.Now()

	if billingCycle == "monthly" {
		sub.EndDate = time.Now().AddDate(0, 1, 0)
	} else {
		sub.EndDate = time.Now().AddDate(1, 0, 0)
	}

	if err := s.repo.Update(sub); err != nil {
		return nil, err
	}

	return toSubscriptionResponse(sub), nil
}

func (s *Service) CancelSubscription(id int64) error {
	sub, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if sub == nil {
		return errors.New("subscription not found")
	}

	sub.Status = "cancelled"
	return s.repo.Update(sub)
}

func (s *Service) GetCompanySubscription(companyID int64) (*SubscriptionResponse, error) {
	sub, err := s.repo.GetByCompanyID(companyID)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, errors.New("subscription not found")
	}

	return toSubscriptionResponse(sub), nil
}

func (s *Service) CheckModuleAccess(companyID int64, moduleID int64) (bool, error) {
	return s.repo.CheckModuleAccess(companyID, moduleID)
}

func (s *Service) GetExpiringSubscriptions(days int) ([]*SubscriptionResponse, error) {
	subs, err := s.repo.GetExpiring(days)
	if err != nil {
		return nil, err
	}

	var responses []*SubscriptionResponse
	for _, sub := range subs {
		responses = append(responses, toSubscriptionResponse(sub))
	}

	return responses, nil
}

func (s *Service) UpdateExpiredSubscriptions() error {
	return s.repo.UpdateExpired()
}

func (s *Service) GetSubscriptionStats() (map[string]interface{}, error) {
	return s.repo.GetStats()
}

func (s *Service) MarkPaymentAsPaid(id int64) error {
	return s.repo.MarkPaymentPaid(id)
}

func toPlanResponse(plan *SubscriptionPlan) *SubscriptionPlanResponse {
	if plan == nil {
		return nil
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
		Features:     plan.Features,
		IsActive:     plan.IsActive,
		CreatedAt:    plan.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    plan.UpdatedAt.Format(time.RFC3339),
	}
}

func toSubscriptionResponse(sub *Subscription) *SubscriptionResponse {
	if sub == nil {
		return nil
	}

	var lastPayment, nextPayment *string
	if sub.LastPaymentDate != nil {
		lp := sub.LastPaymentDate.Format(time.RFC3339)
		lastPayment = &lp
	}
	if sub.NextPaymentDate != nil {
		np := sub.NextPaymentDate.Format(time.RFC3339)
		nextPayment = &np
	}

	return &SubscriptionResponse{
		ID:              sub.ID,
		CompanyID:       sub.CompanyID,
		PlanID:          sub.PlanID,
		Status:          sub.Status,
		BillingCycle:    sub.BillingCycle,
		StartDate:       sub.StartDate.Format(time.RFC3339),
		EndDate:         sub.EndDate.Format(time.RFC3339),
		Price:           sub.Price,
		Currency:        sub.Currency,
		PaymentStatus:   sub.PaymentStatus,
		LastPaymentDate: lastPayment,
		NextPaymentDate: nextPayment,
		AutoRenew:       sub.AutoRenew,
		CompanyName:     sub.CompanyName,
		PlanDisplayName: sub.PlanDisplayName,
		CreatedAt:       sub.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       sub.UpdatedAt.Format(time.RFC3339),
	}
}

// Plan Modules Management Methods

func (s *Service) GetPlanModules(planID int64) (*PlanModulesListResponse, error) {
	// Check if plan exists
	exists, err := s.repo.CheckPlanExists(planID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("subscription plan not found")
	}

	// Get plan modules
	modules, err := s.repo.GetPlanModules(planID)
	if err != nil {
		return nil, err
	}

	// Get plan name only (avoid features field issue)
	var planName string
	if len(modules) > 0 {
		planName = "Plan " + fmt.Sprintf("%d", planID) // Simple fallback
	} else {
		planName = "Plan " + fmt.Sprintf("%d", planID)
	}

	return &PlanModulesListResponse{
		Data:     modules,
		Total:    int64(len(modules)),
		PlanID:   planID,
		PlanName: planName,
	}, nil
}

func (s *Service) AddModulesToPlan(planID int64, req *AddModulesToPlanRequest) error {
	// Validate request
	if err := ValidateAddModulesToPlanRequest(req); err != nil {
		return err
	}

	// Check if plan exists
	exists, err := s.repo.CheckPlanExists(planID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("subscription plan not found")
	}

	// Check if all modules exist
	for _, moduleID := range req.ModuleIDs {
		exists, err := s.repo.CheckModuleExists(moduleID)
		if err != nil {
			return err
		}
		if !exists {
			return errors.New("one or more modules not found")
		}
	}

	// Add modules to plan
	return s.repo.AddModulesToPlan(planID, req.ModuleIDs)
}

func (s *Service) RemoveModuleFromPlan(planID int64, moduleID int64) error {
	// Check if plan exists
	exists, err := s.repo.CheckPlanExists(planID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("subscription plan not found")
	}

	// Check if module exists
	exists, err = s.repo.CheckModuleExists(moduleID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("module not found")
	}

	// Remove module from plan
	return s.repo.RemoveModuleFromPlan(planID, moduleID)
}
