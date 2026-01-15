package subscription

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Repository interface {
	// Plan methods
	GetAllPlans() ([]*SubscriptionPlan, error)
	GetPlanByID(id int64) (*SubscriptionPlan, error)
	CreatePlan(plan *SubscriptionPlan) error
	UpdatePlan(plan *SubscriptionPlan) error
	DeletePlan(id int64) error

	// Subscription methods
	GetAll(limit, offset int, filters map[string]interface{}) ([]*Subscription, error)
	Count(filters map[string]interface{}) (int64, error)
	GetByID(id int64) (*Subscription, error)
	GetByCompanyID(companyID int64) (*Subscription, error)
	Create(sub *Subscription) error
	Update(sub *Subscription) error
	CheckModuleAccess(companyID int64, moduleID int64) (bool, error)
	GetExpiring(days int) ([]*Subscription, error)
	UpdateExpired() error
	GetStats() (map[string]interface{}, error)
	MarkPaymentPaid(id int64) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAllPlans() ([]*SubscriptionPlan, error) {
	query := `SELECT id, name, display_name, description, price_monthly, price_yearly, 
		max_users, max_branches, features, is_active, created_at, updated_at 
		FROM subscription_plans WHERE is_active = true ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*SubscriptionPlan
	for rows.Next() {
		plan := &SubscriptionPlan{}
		err := rows.Scan(&plan.ID, &plan.Name, &plan.DisplayName, &plan.Description,
			&plan.PriceMonthly, &plan.PriceYearly, &plan.MaxUsers, &plan.MaxBranches,
			&plan.Features, &plan.IsActive, &plan.CreatedAt, &plan.UpdatedAt)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}

	return plans, nil
}

func (r *repository) GetPlanByID(id int64) (*SubscriptionPlan, error) {
	query := `SELECT id, name, display_name, description, price_monthly, price_yearly, 
		max_users, max_branches, features, is_active, created_at, updated_at 
		FROM subscription_plans WHERE id = $1`

	plan := &SubscriptionPlan{}
	err := r.db.QueryRow(query, id).Scan(&plan.ID, &plan.Name, &plan.DisplayName,
		&plan.Description, &plan.PriceMonthly, &plan.PriceYearly, &plan.MaxUsers,
		&plan.MaxBranches, &plan.Features, &plan.IsActive, &plan.CreatedAt, &plan.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return plan, err
}

func (r *repository) CreatePlan(plan *SubscriptionPlan) error {
	query := `INSERT INTO subscription_plans (name, display_name, description, price_monthly, 
		price_yearly, max_users, max_branches, features, is_active) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, plan.Name, plan.DisplayName, plan.Description,
		plan.PriceMonthly, plan.PriceYearly, plan.MaxUsers, plan.MaxBranches,
		plan.Features, plan.IsActive).Scan(&plan.ID, &plan.CreatedAt, &plan.UpdatedAt)
}

func (r *repository) UpdatePlan(plan *SubscriptionPlan) error {
	query := `UPDATE subscription_plans SET name = $2, display_name = $3, description = $4, 
		price_monthly = $5, price_yearly = $6, max_users = $7, max_branches = $8, 
		features = $9, is_active = $10, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1 RETURNING updated_at`

	return r.db.QueryRow(query, plan.ID, plan.Name, plan.DisplayName, plan.Description,
		plan.PriceMonthly, plan.PriceYearly, plan.MaxUsers, plan.MaxBranches,
		plan.Features, plan.IsActive).Scan(&plan.UpdatedAt)
}

func (r *repository) DeletePlan(id int64) error {
	query := `UPDATE subscription_plans SET is_active = false WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *repository) GetAll(limit, offset int, filters map[string]interface{}) ([]*Subscription, error) {
	query := `SELECT s.id, s.company_id, s.plan_id, s.status, s.billing_cycle, s.start_date, 
		s.end_date, s.price, s.currency, s.payment_status, s.last_payment_date, 
		s.next_payment_date, s.auto_renew, s.created_at, s.updated_at,
		c.name as company_name, sp.display_name as plan_display_name
		FROM subscriptions s
		JOIN companies c ON s.company_id = c.id
		JOIN subscription_plans sp ON s.plan_id = sp.id
		WHERE 1=1`

	var args []interface{}
	argCount := 1

	if companyID, ok := filters["company_id"]; ok {
		query += fmt.Sprintf(` AND s.company_id = $%d`, argCount)
		args = append(args, companyID)
		argCount++
	}

	if planID, ok := filters["plan_id"]; ok {
		query += fmt.Sprintf(` AND s.plan_id = $%d`, argCount)
		args = append(args, planID)
		argCount++
	}

	if status, ok := filters["status"]; ok {
		query += fmt.Sprintf(` AND s.status = $%d`, argCount)
		args = append(args, status)
		argCount++
	}

	if search, ok := filters["search"]; ok {
		searchPattern := "%" + strings.ToLower(search.(string)) + "%"
		query += fmt.Sprintf(` AND LOWER(c.name) LIKE $%d`, argCount)
		args = append(args, searchPattern)
		argCount++
	}

	query += ` ORDER BY s.created_at DESC`

	if limit > 0 {
		query += fmt.Sprintf(` LIMIT $%d`, argCount)
		args = append(args, limit)
		argCount++
	}

	if offset > 0 {
		query += fmt.Sprintf(` OFFSET $%d`, argCount)
		args = append(args, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []*Subscription
	for rows.Next() {
		sub := &Subscription{}
		err := rows.Scan(&sub.ID, &sub.CompanyID, &sub.PlanID, &sub.Status, &sub.BillingCycle,
			&sub.StartDate, &sub.EndDate, &sub.Price, &sub.Currency, &sub.PaymentStatus,
			&sub.LastPaymentDate, &sub.NextPaymentDate, &sub.AutoRenew, &sub.CreatedAt,
			&sub.UpdatedAt, &sub.CompanyName, &sub.PlanDisplayName)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	return subs, nil
}

func (r *repository) Count(filters map[string]interface{}) (int64, error) {
	query := `SELECT COUNT(*) FROM subscriptions s WHERE 1=1`
	var args []interface{}
	argCount := 1

	if companyID, ok := filters["company_id"]; ok {
		query += fmt.Sprintf(` AND s.company_id = $%d`, argCount)
		args = append(args, companyID)
		argCount++
	}

	if status, ok := filters["status"]; ok {
		query += fmt.Sprintf(` AND s.status = $%d`, argCount)
		args = append(args, status)
	}

	var count int64
	err := r.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

func (r *repository) GetByID(id int64) (*Subscription, error) {
	query := `SELECT s.id, s.company_id, s.plan_id, s.status, s.billing_cycle, s.start_date, 
		s.end_date, s.price, s.currency, s.payment_status, s.last_payment_date, 
		s.next_payment_date, s.auto_renew, s.created_at, s.updated_at,
		c.name as company_name, sp.display_name as plan_display_name
		FROM subscriptions s
		JOIN companies c ON s.company_id = c.id
		JOIN subscription_plans sp ON s.plan_id = sp.id
		WHERE s.id = $1`

	sub := &Subscription{}
	err := r.db.QueryRow(query, id).Scan(&sub.ID, &sub.CompanyID, &sub.PlanID, &sub.Status,
		&sub.BillingCycle, &sub.StartDate, &sub.EndDate, &sub.Price, &sub.Currency,
		&sub.PaymentStatus, &sub.LastPaymentDate, &sub.NextPaymentDate, &sub.AutoRenew,
		&sub.CreatedAt, &sub.UpdatedAt, &sub.CompanyName, &sub.PlanDisplayName)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return sub, err
}

func (r *repository) GetByCompanyID(companyID int64) (*Subscription, error) {
	query := `SELECT s.id, s.company_id, s.plan_id, s.status, s.billing_cycle, s.start_date, 
		s.end_date, s.price, s.currency, s.payment_status, s.last_payment_date, 
		s.next_payment_date, s.auto_renew, s.created_at, s.updated_at,
		c.name as company_name, sp.display_name as plan_display_name
		FROM subscriptions s
		JOIN companies c ON s.company_id = c.id
		JOIN subscription_plans sp ON s.plan_id = sp.id
		WHERE s.company_id = $1 AND s.status = 'active'
		ORDER BY s.created_at DESC LIMIT 1`

	sub := &Subscription{}
	err := r.db.QueryRow(query, companyID).Scan(&sub.ID, &sub.CompanyID, &sub.PlanID,
		&sub.Status, &sub.BillingCycle, &sub.StartDate, &sub.EndDate, &sub.Price,
		&sub.Currency, &sub.PaymentStatus, &sub.LastPaymentDate, &sub.NextPaymentDate,
		&sub.AutoRenew, &sub.CreatedAt, &sub.UpdatedAt, &sub.CompanyName, &sub.PlanDisplayName)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return sub, err
}

func (r *repository) Create(sub *Subscription) error {
	query := `INSERT INTO subscriptions (company_id, plan_id, status, billing_cycle, start_date, 
		end_date, price, currency, payment_status, auto_renew) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, sub.CompanyID, sub.PlanID, sub.Status, sub.BillingCycle,
		sub.StartDate, sub.EndDate, sub.Price, sub.Currency, sub.PaymentStatus,
		sub.AutoRenew).Scan(&sub.ID, &sub.CreatedAt, &sub.UpdatedAt)
}

func (r *repository) Update(sub *Subscription) error {
	query := `UPDATE subscriptions SET plan_id = $2, status = $3, billing_cycle = $4, 
		end_date = $5, price = $6, payment_status = $7, auto_renew = $8, 
		updated_at = CURRENT_TIMESTAMP WHERE id = $1 RETURNING updated_at`

	return r.db.QueryRow(query, sub.ID, sub.PlanID, sub.Status, sub.BillingCycle,
		sub.EndDate, sub.Price, sub.PaymentStatus, sub.AutoRenew).Scan(&sub.UpdatedAt)
}

func (r *repository) CheckModuleAccess(companyID int64, moduleID int64) (bool, error) {
	query := `SELECT EXISTS(
		SELECT 1 FROM subscriptions s
		JOIN plan_modules pm ON s.plan_id = pm.plan_id
		WHERE s.company_id = $1 AND pm.module_id = $2 
		AND s.status = 'active' AND pm.is_included = true
	)`

	var hasAccess bool
	err := r.db.QueryRow(query, companyID, moduleID).Scan(&hasAccess)
	return hasAccess, err
}

func (r *repository) GetExpiring(days int) ([]*Subscription, error) {
	query := `SELECT s.id, s.company_id, s.plan_id, s.status, s.billing_cycle, s.start_date, 
		s.end_date, s.price, s.currency, s.payment_status, s.last_payment_date, 
		s.next_payment_date, s.auto_renew, s.created_at, s.updated_at,
		c.name as company_name, sp.display_name as plan_display_name
		FROM subscriptions s
		JOIN companies c ON s.company_id = c.id
		JOIN subscription_plans sp ON s.plan_id = sp.id
		WHERE s.status = 'active' AND s.end_date <= $1`

	expiryDate := time.Now().AddDate(0, 0, days)
	rows, err := r.db.Query(query, expiryDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []*Subscription
	for rows.Next() {
		sub := &Subscription{}
		err := rows.Scan(&sub.ID, &sub.CompanyID, &sub.PlanID, &sub.Status, &sub.BillingCycle,
			&sub.StartDate, &sub.EndDate, &sub.Price, &sub.Currency, &sub.PaymentStatus,
			&sub.LastPaymentDate, &sub.NextPaymentDate, &sub.AutoRenew, &sub.CreatedAt,
			&sub.UpdatedAt, &sub.CompanyName, &sub.PlanDisplayName)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	return subs, nil
}

func (r *repository) UpdateExpired() error {
	query := `UPDATE subscriptions SET status = 'expired' 
		WHERE status = 'active' AND end_date < CURRENT_DATE`
	_, err := r.db.Exec(query)
	return err
}

func (r *repository) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var total, active, expired int64
	r.db.QueryRow(`SELECT COUNT(*) FROM subscriptions`).Scan(&total)
	r.db.QueryRow(`SELECT COUNT(*) FROM subscriptions WHERE status = 'active'`).Scan(&active)
	r.db.QueryRow(`SELECT COUNT(*) FROM subscriptions WHERE status = 'expired'`).Scan(&expired)

	stats["total"] = total
	stats["active"] = active
	stats["expired"] = expired

	return stats, nil
}

func (r *repository) MarkPaymentPaid(id int64) error {
	query := `UPDATE subscriptions SET payment_status = 'paid', 
		last_payment_date = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
