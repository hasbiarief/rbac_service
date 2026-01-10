package service

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/interfaces"
	"gin-scalable-api/internal/mapper"
	"gin-scalable-api/internal/models"
	"time"
)

type SubscriptionService struct {
	subscriptionRepo   interfaces.SubscriptionRepositoryInterface
	subscriptionMapper *mapper.SubscriptionMapper
}

func NewSubscriptionService(subscriptionRepo interfaces.SubscriptionRepositoryInterface) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo:   subscriptionRepo,
		subscriptionMapper: mapper.NewSubscriptionMapper(),
	}
}

// Paket Langganan
func (s *SubscriptionService) GetSubscriptionPlans() ([]*dto.SubscriptionPlanResponse, error) {
	plans, err := s.subscriptionRepo.GetAllPlans()
	if err != nil {
		return nil, err
	}

	var response []*dto.SubscriptionPlanResponse
	for _, plan := range plans {
		response = append(response, s.subscriptionMapper.ToPlanResponse(plan))
	}

	return response, nil
}

func (s *SubscriptionService) GetSubscriptionPlanByID(id int64) (*dto.SubscriptionPlanResponse, error) {
	plan, err := s.subscriptionRepo.GetPlanByID(id)
	if err != nil {
		return nil, err
	}

	return s.subscriptionMapper.ToPlanResponse(plan), nil
}

func (s *SubscriptionService) CreateSubscriptionPlan(req *dto.CreateSubscriptionPlanRequest) (*dto.SubscriptionPlanResponse, error) {
	// Konversi DTO ke model menggunakan mapper
	plan := s.subscriptionMapper.ToPlanModel(req)

	if err := s.subscriptionRepo.CreatePlan(plan); err != nil {
		return nil, err
	}

	return s.subscriptionMapper.ToPlanResponse(plan), nil
}

func (s *SubscriptionService) UpdateSubscriptionPlan(id int64, req *dto.UpdateSubscriptionPlanRequest) (*dto.SubscriptionPlanResponse, error) {
	plan, err := s.subscriptionRepo.GetPlanByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields menggunakan mapper
	s.subscriptionMapper.UpdatePlanModel(plan, req)

	if err := s.subscriptionRepo.UpdatePlan(plan); err != nil {
		return nil, err
	}

	return s.subscriptionMapper.ToPlanResponse(plan), nil
}

func (s *SubscriptionService) DeleteSubscriptionPlan(id int64) error {
	return s.subscriptionRepo.DeletePlan(id)
}

// Langganan
func (s *SubscriptionService) GetSubscriptions(req *dto.SubscriptionListRequest) (*dto.SubscriptionListResponse, error) {
	// Set default values jika tidak disediakan
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	// Build filters map dari request
	filters := make(map[string]interface{})
	if req.CompanyID != nil {
		filters["company_id"] = *req.CompanyID
	}
	if req.Status != "" {
		filters["status"] = req.Status
	}

	subscriptions, err := s.subscriptionRepo.GetAll(limit, offset, filters)
	if err != nil {
		return nil, err
	}

	// Dapatkan total count untuk pagination
	total, err := s.subscriptionRepo.Count(filters)
	if err != nil {
		return nil, err
	}

	// Konversi ke DTO menggunakan mapper
	var subscriptionResponses []*dto.SubscriptionResponse
	for _, sub := range subscriptions {
		subscriptionResponses = append(subscriptionResponses, s.subscriptionMapper.ToSubscriptionResponse(sub))
	}

	return &dto.SubscriptionListResponse{
		Data:    subscriptionResponses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(subscriptionResponses)) < total,
	}, nil
}

func (s *SubscriptionService) GetSubscriptionByID(id int64) (*dto.SubscriptionResponse, error) {
	sub, err := s.subscriptionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.subscriptionMapper.ToSubscriptionResponse(sub), nil
}

func (s *SubscriptionService) GetCompanySubscription(companyID int64) (*dto.SubscriptionResponse, error) {
	sub, err := s.subscriptionRepo.GetByCompanyID(companyID)
	if err != nil {
		return nil, err
	}

	return s.subscriptionMapper.ToSubscriptionResponse(sub), nil
}

func (s *SubscriptionService) CreateSubscription(req *dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	// Dapatkan plan untuk menentukan harga
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
		nextPaymentDate = time.Now().AddDate(1, 0, 0) // Pembayaran berikutnya dalam 1 tahun
	} else {
		endDate = time.Now().AddDate(0, 1, 0)
		price = plan.PriceMonthly
		nextPaymentDate = time.Now().AddDate(0, 1, 0) // Pembayaran berikutnya dalam 1 bulan
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

	// Dapatkan subscription yang dibuat dengan nama company dan plan
	createdSub, err := s.subscriptionRepo.GetByID(subscription.ID)
	if err != nil {
		return nil, err
	}

	return s.subscriptionMapper.ToSubscriptionResponse(createdSub), nil
}

func (s *SubscriptionService) UpdateSubscription(id int64, req *dto.UpdateSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	subscription, err := s.subscriptionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Jika plan sedang diubah, update harga sesuai
	if req.Status != "" {
		subscription.Status = req.Status
	}
	if req.BillingCycle != "" {
		subscription.BillingCycle = req.BillingCycle
	}

	if req.AutoRenew != nil {
		subscription.AutoRenew = *req.AutoRenew
	}

	// Pastikan field yang diperlukan diset (jika hilang dari DB)
	if subscription.Currency == "" {
		subscription.Currency = "IDR"
	}
	if subscription.PaymentStatus == "" {
		subscription.PaymentStatus = "pending"
	}

	if err := s.subscriptionRepo.Update(subscription); err != nil {
		return nil, err
	}

	// Dapatkan subscription yang diupdate
	updatedSub, err := s.subscriptionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.subscriptionMapper.ToSubscriptionResponse(updatedSub), nil
}

func (s *SubscriptionService) RenewSubscription(subscriptionID int64, planID *int64, billingCycle string) (*dto.SubscriptionResponse, error) {
	// Gunakan repository method yang sudah ada
	if err := s.subscriptionRepo.RenewSubscription(subscriptionID, planID, billingCycle); err != nil {
		return nil, err
	}

	// Dapatkan subscription yang diperpanjang
	renewedSub, err := s.subscriptionRepo.GetByID(subscriptionID)
	if err != nil {
		return nil, err
	}

	return s.subscriptionMapper.ToSubscriptionResponse(renewedSub), nil
}

func (s *SubscriptionService) CancelSubscription(id int64) error {
	// Repository interface hanya menerima ID, parameter tambahan diabaikan untuk saat ini
	return s.subscriptionRepo.Delete(id) // Menggunakan Delete sebagai equivalent CancelSubscription
}

// Additional methods used by handler
func (s *SubscriptionService) CheckModuleAccess(companyID, moduleID int64) (bool, error) {
	return s.subscriptionRepo.CheckModuleAccess(companyID, moduleID)
}

func (s *SubscriptionService) GetExpiringSubscriptions(days int) (interface{}, error) {
	subscriptions, err := s.subscriptionRepo.GetExpiringSubscriptions(days)
	if err != nil {
		return nil, err
	}

	var response []*dto.SubscriptionResponse
	for _, sub := range subscriptions {
		response = append(response, s.subscriptionMapper.ToSubscriptionResponse(sub))
	}

	return response, nil
}

func (s *SubscriptionService) UpdateExpiredSubscriptions() error {
	return s.subscriptionRepo.UpdateExpiredSubscriptions()
}

func (s *SubscriptionService) GetSubscriptionStats() (interface{}, error) {
	// Dapatkan statistik langganan dasar
	stats := map[string]interface{}{
		"total_subscriptions":     0,
		"active_subscriptions":    0,
		"expired_subscriptions":   0,
		"cancelled_subscriptions": 0,
		"total_revenue":           0.0,
	}

	// Ini adalah implementasi placeholder
	// Dalam implementasi nyata, Anda akan query database untuk statistik aktual
	return stats, nil
}

// MarkPaymentAsPaid menandai pembayaran langganan sebagai dibayar
func (s *SubscriptionService) MarkPaymentAsPaid(subscriptionID int64) error {
	return s.subscriptionRepo.MarkPaymentAsPaid(subscriptionID)
}
