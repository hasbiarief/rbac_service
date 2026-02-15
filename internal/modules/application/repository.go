package application

import (
	"database/sql"
	"fmt"
	"gin-scalable-api/pkg/model"
)

type Repository struct {
	*model.Repository
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Repository: model.NewRepository(db),
		db:         db,
	}
}

// GetAll retrieves all applications with pagination and filtering
func (r *Repository) GetAll(limit, offset int, search string, isActive *bool) ([]*Application, error) {
	query := `
		SELECT id, name, code, description, icon, url, is_active, sort_order, created_at, updated_at
		FROM applications
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR code ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex+1, argIndex+2)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
		argIndex += 3
	}

	if isActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *isActive)
		argIndex++
	}

	query += " ORDER BY sort_order, name"

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
		return nil, fmt.Errorf("failed to get applications: %w", err)
	}
	defer rows.Close()

	var applications []*Application
	for rows.Next() {
		app := &Application{}
		err := rows.Scan(
			&app.ID, &app.Name, &app.Code, &app.Description, &app.Icon, &app.URL,
			&app.IsActive, &app.SortOrder, &app.CreatedAt, &app.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan application: %w", err)
		}
		applications = append(applications, app)
	}

	return applications, nil
}

// GetByID retrieves an application by ID
func (r *Repository) GetByID(id int64) (*Application, error) {
	query := `
		SELECT id, name, code, description, icon, url, is_active, sort_order, created_at, updated_at
		FROM applications
		WHERE id = $1
	`

	app := &Application{}
	err := r.db.QueryRow(query, id).Scan(
		&app.ID, &app.Name, &app.Code, &app.Description, &app.Icon, &app.URL,
		&app.IsActive, &app.SortOrder, &app.CreatedAt, &app.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("application not found")
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	return app, nil
}

// GetByCode retrieves an application by code
func (r *Repository) GetByCode(code string) (*Application, error) {
	query := `
		SELECT id, name, code, description, icon, url, is_active, sort_order, created_at, updated_at
		FROM applications
		WHERE code = $1
	`

	app := &Application{}
	err := r.db.QueryRow(query, code).Scan(
		&app.ID, &app.Name, &app.Code, &app.Description, &app.Icon, &app.URL,
		&app.IsActive, &app.SortOrder, &app.CreatedAt, &app.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("application not found")
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	return app, nil
}

// Create creates a new application
func (r *Repository) Create(app *Application) error {
	query, values := r.BuildInsertQuery(app)

	err := r.db.QueryRow(query, values...).Scan(
		&app.ID, &app.CreatedAt, &app.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}

	return nil
}

// Update updates an application
func (r *Repository) Update(app *Application) error {
	query, values := r.BuildUpdateQuery(app, app.ID)

	err := r.db.QueryRow(query, values...).Scan(&app.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}

	return nil
}

// Delete deletes an application
func (r *Repository) Delete(id int64) error {
	query := `DELETE FROM applications WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete application: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("application not found")
	}

	return nil
}

// Count returns total count of applications with filtering
func (r *Repository) Count(search string, isActive *bool) (int64, error) {
	query := "SELECT COUNT(*) FROM applications WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR code ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex+1, argIndex+2)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
		argIndex += 3
	}

	if isActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *isActive)
	}

	var count int64
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get application count: %w", err)
	}

	return count, nil
}

// GetPlanApplications retrieves applications for a specific plan
func (r *Repository) GetPlanApplications(planID int64) ([]*PlanApplicationResponse, error) {
	query := `
		SELECT 
			pa.id, pa.plan_id, sp.name as plan_name,
			pa.application_id, a.name as application_name, a.code as application_code,
			pa.is_included, pa.created_at
		FROM plan_applications pa
		JOIN subscription_plans sp ON pa.plan_id = sp.id
		JOIN applications a ON pa.application_id = a.id
		WHERE pa.plan_id = $1
		ORDER BY a.sort_order, a.name
	`

	rows, err := r.db.Query(query, planID)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan applications: %w", err)
	}
	defer rows.Close()

	var planApps []*PlanApplicationResponse
	for rows.Next() {
		planApp := &PlanApplicationResponse{}
		err := rows.Scan(
			&planApp.ID, &planApp.PlanID, &planApp.PlanName,
			&planApp.ApplicationID, &planApp.ApplicationName, &planApp.ApplicationCode,
			&planApp.IsIncluded, &planApp.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan plan application: %w", err)
		}
		planApps = append(planApps, planApp)
	}

	return planApps, nil
}

// AddApplicationsToPlan adds applications to a subscription plan
func (r *Repository) AddApplicationsToPlan(planID int64, applicationIDs []int64, isIncluded bool) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, appID := range applicationIDs {
		// Check if already exists
		var exists bool
		checkQuery := "SELECT EXISTS(SELECT 1 FROM plan_applications WHERE plan_id = $1 AND application_id = $2)"
		err = tx.QueryRow(checkQuery, planID, appID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check existing plan application: %w", err)
		}

		if exists {
			// Update existing
			_, err = tx.Exec(`
				UPDATE plan_applications 
				SET is_included = $3
				WHERE plan_id = $1 AND application_id = $2
			`, planID, appID, isIncluded)
			if err != nil {
				return fmt.Errorf("failed to update plan application: %w", err)
			}
		} else {
			// Insert new
			_, err = tx.Exec(`
				INSERT INTO plan_applications (plan_id, application_id, is_included)
				VALUES ($1, $2, $3)
			`, planID, appID, isIncluded)
			if err != nil {
				return fmt.Errorf("failed to insert plan application: %w", err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RemoveApplicationFromPlan removes an application from a subscription plan
func (r *Repository) RemoveApplicationFromPlan(planID, applicationID int64) error {
	query := `DELETE FROM plan_applications WHERE plan_id = $1 AND application_id = $2`

	result, err := r.db.Exec(query, planID, applicationID)
	if err != nil {
		return fmt.Errorf("failed to remove application from plan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("plan application not found")
	}

	return nil
}
