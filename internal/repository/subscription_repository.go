package repository

import (
	"database/sql"
	"fmt"
	"gin-scalable-api/internal/models"
	"gin-scalable-api/pkg/model"
	"time"
)

type SubscriptionRepository struct {
	*model.Repository
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		Repository: model.NewRepository(db),
		db:         db,
	}
}

// GetAllPlans retrieves all subscription plans
func (r *SubscriptionRepository) GetAllPlans() ([]*models.SubscriptionPlan, error) {
	query := `
		SELECT id, name, display_name, description, price_monthly, price_yearly, 
			   max_users, max_branches, features, is_active, created_at, updated_at
		FROM subscription_plans
		WHERE is_active = true
		ORDER BY price_monthly
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription plans: %w", err)
	}
	defer rows.Close()

	var plans []*models.SubscriptionPlan
	for rows.Next() {
		plan := &models.SubscriptionPlan{}
		err := rows.Scan(
			&plan.ID, &plan.Name, &plan.DisplayName, &plan.Description,
			&plan.PriceMonthly, &plan.PriceYearly, &plan.MaxUsers, &plan.MaxBranches,
			&plan.Features, &plan.IsActive, &plan.CreatedAt, &plan.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription plan: %w", err)
		}
		plans = append(plans, plan)
	}

	return plans, nil
}

// GetPlanByID retrieves a subscription plan by ID
func (r *SubscriptionRepository) GetPlanByID(id int64) (*models.SubscriptionPlan, error) {
	query := `
		SELECT id, name, display_name, description, price_monthly, price_yearly, 
			   max_users, max_branches, features, is_active, created_at, updated_at
		FROM subscription_plans
		WHERE id = $1
	`

	plan := &models.SubscriptionPlan{}
	err := r.db.QueryRow(query, id).Scan(
		&plan.ID, &plan.Name, &plan.DisplayName, &plan.Description,
		&plan.PriceMonthly, &plan.PriceYearly, &plan.MaxUsers, &plan.MaxBranches,
		&plan.Features, &plan.IsActive, &plan.CreatedAt, &plan.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subscription plan not found")
		}
		return nil, fmt.Errorf("failed to get subscription plan: %w", err)
	}

	return plan, nil
}

// GetAllSubscriptions retrieves all subscriptions with pagination
func (r *SubscriptionRepository) GetAllSubscriptions(limit, offset int) ([]*models.Subscription, error) {
	query := `
		SELECT s.id, s.company_id, s.plan_id, s.status, s.billing_cycle, s.start_date, 
			   s.end_date, s.auto_renew, s.created_at, s.updated_at,
			   c.name as company_name, sp.display_name as plan_display_name
		FROM subscriptions s
		JOIN companies c ON s.company_id = c.id
		JOIN subscription_plans sp ON s.plan_id = sp.id
		ORDER BY s.created_at DESC
	`
	args := []interface{}{}

	if limit > 0 {
		query += " LIMIT $1"
		args = append(args, limit)
		if offset > 0 {
			query += " OFFSET $2"
			args = append(args, offset)
		}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		subscription := &models.Subscription{}
		var companyName, planDisplayName string
		err := rows.Scan(
			&subscription.ID, &subscription.CompanyID, &subscription.PlanID,
			&subscription.Status, &subscription.BillingCycle, &subscription.StartDate,
			&subscription.EndDate, &subscription.AutoRenew, &subscription.CreatedAt,
			&subscription.UpdatedAt, &companyName, &planDisplayName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}

		// Add company and plan names to metadata
		subscription.CompanyName = companyName
		subscription.PlanDisplayName = planDisplayName
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

// GetSubscriptionByID retrieves a subscription by ID
func (r *SubscriptionRepository) GetSubscriptionByID(id int64) (*models.Subscription, error) {
	query := `
		SELECT s.id, s.company_id, s.plan_id, s.status, s.billing_cycle, s.start_date, 
			   s.end_date, s.auto_renew, s.created_at, s.updated_at,
			   c.name as company_name, sp.display_name as plan_display_name
		FROM subscriptions s
		JOIN companies c ON s.company_id = c.id
		JOIN subscription_plans sp ON s.plan_id = sp.id
		WHERE s.id = $1
	`

	subscription := &models.Subscription{}
	var companyName, planDisplayName string
	err := r.db.QueryRow(query, id).Scan(
		&subscription.ID, &subscription.CompanyID, &subscription.PlanID,
		&subscription.Status, &subscription.BillingCycle, &subscription.StartDate,
		&subscription.EndDate, &subscription.AutoRenew, &subscription.CreatedAt,
		&subscription.UpdatedAt, &companyName, &planDisplayName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	subscription.CompanyName = companyName
	subscription.PlanDisplayName = planDisplayName
	return subscription, nil
}

// GetCompanySubscription retrieves active subscription for a company
func (r *SubscriptionRepository) GetCompanySubscription(companyID int64) (*models.Subscription, error) {
	query := `
		SELECT s.id, s.company_id, s.plan_id, s.status, s.billing_cycle, s.start_date, 
			   s.end_date, s.auto_renew, s.created_at, s.updated_at,
			   c.name as company_name, sp.display_name as plan_display_name
		FROM subscriptions s
		JOIN companies c ON s.company_id = c.id
		JOIN subscription_plans sp ON s.plan_id = sp.id
		WHERE s.company_id = $1 AND s.status = 'active' AND s.end_date > CURRENT_DATE
		ORDER BY s.end_date DESC
		LIMIT 1
	`

	subscription := &models.Subscription{}
	var companyName, planDisplayName string
	err := r.db.QueryRow(query, companyID).Scan(
		&subscription.ID, &subscription.CompanyID, &subscription.PlanID,
		&subscription.Status, &subscription.BillingCycle, &subscription.StartDate,
		&subscription.EndDate, &subscription.AutoRenew, &subscription.CreatedAt,
		&subscription.UpdatedAt, &companyName, &planDisplayName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no active subscription found for company")
		}
		return nil, fmt.Errorf("failed to get company subscription: %w", err)
	}

	subscription.CompanyName = companyName
	subscription.PlanDisplayName = planDisplayName
	return subscription, nil
}

// Create creates a new subscription
func (r *SubscriptionRepository) Create(subscription *models.Subscription) error {
	query, values := r.BuildInsertQuery(subscription)

	err := r.db.QueryRow(query, values...).Scan(
		&subscription.ID, &subscription.CreatedAt, &subscription.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

// Update updates a subscription
func (r *SubscriptionRepository) Update(subscription *models.Subscription) error {
	query, values := r.BuildUpdateQuery(subscription, subscription.ID)

	err := r.db.QueryRow(query, values...).Scan(&subscription.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

// RenewSubscription renews a subscription
func (r *SubscriptionRepository) RenewSubscription(subscriptionID int64, planID *int64, billingCycle string) error {
	var endDate time.Time
	var nextPaymentDate time.Time

	if billingCycle == "yearly" {
		endDate = time.Now().AddDate(1, 0, 0)
		nextPaymentDate = time.Now().AddDate(1, 0, 0)
	} else {
		endDate = time.Now().AddDate(0, 1, 0)
		nextPaymentDate = time.Now().AddDate(0, 1, 0)
	}

	// Get the current subscription to determine the plan
	currentSub, err := r.GetSubscriptionByID(subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get current subscription: %w", err)
	}

	// Determine which plan to use for price calculation
	targetPlanID := currentSub.PlanID
	if planID != nil {
		targetPlanID = *planID
	}

	// Get the plan to calculate the correct price
	plan, err := r.GetPlanByID(targetPlanID)
	if err != nil {
		return fmt.Errorf("failed to get plan for price calculation: %w", err)
	}

	// Calculate price based on billing cycle
	var price float64
	if billingCycle == "yearly" {
		price = plan.PriceYearly
	} else {
		price = plan.PriceMonthly
	}

	query := `
		UPDATE subscriptions 
		SET plan_id = COALESCE($2, plan_id),
			billing_cycle = $3,
			end_date = $4,
			price = $5,
			payment_status = 'pending',
			last_payment_date = CURRENT_DATE,
			next_payment_date = $6,
			status = 'active',
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.Exec(query, subscriptionID, planID, billingCycle, endDate, price, nextPaymentDate)
	if err != nil {
		return fmt.Errorf("failed to renew subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

// CancelSubscription cancels a subscription
func (r *SubscriptionRepository) CancelSubscription(subscriptionID int64, reason string, cancelImmediately bool) error {
	var query string
	if cancelImmediately {
		query = `
			UPDATE subscriptions 
			SET status = 'cancelled',
				end_date = CURRENT_DATE,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`
	} else {
		query = `
			UPDATE subscriptions 
			SET status = 'cancelled',
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`
	}

	result, err := r.db.Exec(query, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

// CheckModuleAccess checks if company has access to a specific module
func (r *SubscriptionRepository) CheckModuleAccess(companyID, moduleID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM subscriptions s
			JOIN plan_modules pm ON s.plan_id = pm.plan_id
			WHERE s.company_id = $1 
				AND pm.module_id = $2
				AND pm.is_included = true
				AND s.status = 'active'
				AND s.end_date > CURRENT_DATE
		)
	`

	var hasAccess bool
	err := r.db.QueryRow(query, companyID, moduleID).Scan(&hasAccess)
	if err != nil {
		return false, fmt.Errorf("failed to check module access: %w", err)
	}

	return hasAccess, nil
}

// GetExpiringSubscriptions retrieves subscriptions expiring within specified days
func (r *SubscriptionRepository) GetExpiringSubscriptions(days int) ([]*models.Subscription, error) {
	query := `
		SELECT s.id, s.company_id, s.plan_id, s.status, s.billing_cycle, s.start_date, 
			   s.end_date, s.auto_renew, s.created_at, s.updated_at,
			   c.name as company_name, sp.display_name as plan_display_name
		FROM subscriptions s
		JOIN companies c ON s.company_id = c.id
		JOIN subscription_plans sp ON s.plan_id = sp.id
		WHERE s.status = 'active' 
			AND s.end_date <= CURRENT_DATE + INTERVAL '%d days'
			AND s.end_date > CURRENT_DATE
		ORDER BY s.end_date
	`

	rows, err := r.db.Query(fmt.Sprintf(query, days))
	if err != nil {
		return nil, fmt.Errorf("failed to get expiring subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		subscription := &models.Subscription{}
		var companyName, planDisplayName string
		err := rows.Scan(
			&subscription.ID, &subscription.CompanyID, &subscription.PlanID,
			&subscription.Status, &subscription.BillingCycle, &subscription.StartDate,
			&subscription.EndDate, &subscription.AutoRenew, &subscription.CreatedAt,
			&subscription.UpdatedAt, &companyName, &planDisplayName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}

		subscription.CompanyName = companyName
		subscription.PlanDisplayName = planDisplayName
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

// UpdateExpiredSubscriptions updates expired subscriptions to inactive status
func (r *SubscriptionRepository) UpdateExpiredSubscriptions() error {
	query := `
		UPDATE subscriptions 
		SET status = 'expired',
			updated_at = CURRENT_TIMESTAMP
		WHERE status = 'active' AND end_date <= CURRENT_DATE
	`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to update expired subscriptions: %w", err)
	}

	return nil
}

// MarkPaymentAsPaid marks a subscription payment as paid and updates payment dates
func (r *SubscriptionRepository) MarkPaymentAsPaid(subscriptionID int64) error {
	// Get current subscription to determine next payment date
	subscription, err := r.GetByID(subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	var nextPaymentDate time.Time
	if subscription.BillingCycle == "yearly" {
		nextPaymentDate = time.Now().AddDate(1, 0, 0)
	} else {
		nextPaymentDate = time.Now().AddDate(0, 1, 0)
	}

	query := `
		UPDATE subscriptions 
		SET payment_status = 'paid',
			last_payment_date = CURRENT_DATE,
			next_payment_date = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.Exec(query, subscriptionID, nextPaymentDate)
	if err != nil {
		return fmt.Errorf("failed to mark payment as paid: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

// Interface-compatible methods
func (r *SubscriptionRepository) GetAll(limit, offset int, filters map[string]interface{}) ([]*models.Subscription, error) {
	query := `
		SELECT s.id, s.company_id, s.plan_id, s.status, s.billing_cycle, s.start_date, 
			   s.end_date, s.price, s.currency, s.payment_status, s.last_payment_date,
			   s.next_payment_date, s.auto_renew, s.created_at, s.updated_at,
			   c.name as company_name, sp.display_name as plan_display_name
		FROM subscriptions s
		JOIN companies c ON s.company_id = c.id
		JOIN subscription_plans sp ON s.plan_id = sp.id
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	// Apply filters
	if companyID, ok := filters["company_id"]; ok {
		query += fmt.Sprintf(" AND s.company_id = $%d", argIndex)
		args = append(args, companyID)
		argIndex++
	}

	if status, ok := filters["status"]; ok {
		query += fmt.Sprintf(" AND s.status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	query += " ORDER BY s.created_at DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
		argIndex++
	}

	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		subscription := &models.Subscription{}
		var companyName, planDisplayName string
		err := rows.Scan(
			&subscription.ID, &subscription.CompanyID, &subscription.PlanID,
			&subscription.Status, &subscription.BillingCycle, &subscription.StartDate,
			&subscription.EndDate, &subscription.Price, &subscription.Currency,
			&subscription.PaymentStatus, &subscription.LastPaymentDate,
			&subscription.NextPaymentDate, &subscription.AutoRenew,
			&subscription.CreatedAt, &subscription.UpdatedAt,
			&companyName, &planDisplayName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}

		subscription.CompanyName = companyName
		subscription.PlanDisplayName = planDisplayName
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) GetByID(id int64) (*models.Subscription, error) {
	return r.GetSubscriptionByID(id)
}

func (r *SubscriptionRepository) GetByCompanyID(companyID int64) (*models.Subscription, error) {
	return r.GetCompanySubscription(companyID)
}

func (r *SubscriptionRepository) Delete(id int64) error {
	return r.CancelSubscription(id, "Deleted via API", true)
}

func (r *SubscriptionRepository) Count(filters map[string]interface{}) (int64, error) {
	query := "SELECT COUNT(*) FROM subscriptions s WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	// Apply filters
	if companyID, ok := filters["company_id"]; ok {
		query += fmt.Sprintf(" AND s.company_id = $%d", argIndex)
		args = append(args, companyID)
		argIndex++
	}

	if status, ok := filters["status"]; ok {
		query += fmt.Sprintf(" AND s.status = $%d", argIndex)
		args = append(args, status)
	}

	var count int64
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get subscription count: %w", err)
	}

	return count, nil
}

// Plan methods
func (r *SubscriptionRepository) GetPlanByName(name string) (*models.SubscriptionPlan, error) {
	query := `
		SELECT id, name, display_name, description, price_monthly, price_yearly, 
			   max_users, max_branches, features, is_active, created_at, updated_at
		FROM subscription_plans
		WHERE name = $1
	`

	plan := &models.SubscriptionPlan{}
	err := r.db.QueryRow(query, name).Scan(
		&plan.ID, &plan.Name, &plan.DisplayName, &plan.Description,
		&plan.PriceMonthly, &plan.PriceYearly, &plan.MaxUsers, &plan.MaxBranches,
		&plan.Features, &plan.IsActive, &plan.CreatedAt, &plan.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subscription plan not found")
		}
		return nil, fmt.Errorf("failed to get subscription plan: %w", err)
	}

	return plan, nil
}

func (r *SubscriptionRepository) CreatePlan(plan *models.SubscriptionPlan) error {
	query, values := r.BuildInsertQuery(plan)

	err := r.db.QueryRow(query, values...).Scan(
		&plan.ID, &plan.CreatedAt, &plan.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create subscription plan: %w", err)
	}

	return nil
}

func (r *SubscriptionRepository) UpdatePlan(plan *models.SubscriptionPlan) error {
	query, values := r.BuildUpdateQuery(plan, plan.ID)

	err := r.db.QueryRow(query, values...).Scan(&plan.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update subscription plan: %w", err)
	}

	return nil
}

func (r *SubscriptionRepository) DeletePlan(id int64) error {
	query := `UPDATE subscription_plans SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription plan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription plan not found")
	}

	return nil
}
