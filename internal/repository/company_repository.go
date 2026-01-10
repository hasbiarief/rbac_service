package repository

import (
	"database/sql"
	"fmt"
	"gin-scalable-api/internal/models"
	"gin-scalable-api/pkg/model"
)

type CompanyRepository struct {
	*model.Repository
	db *sql.DB
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{
		Repository: model.NewRepository(db),
		db:         db,
	}
}

// GetAll retrieves all companies with pagination and filtering
func (r *CompanyRepository) GetAll(limit, offset int, search string, isActive *bool) ([]*models.Company, error) {
	query := `
		SELECT id, name, code, is_active, created_at, updated_at
		FROM companies
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR code ILIKE $%d)", argIndex, argIndex+1)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	if isActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *isActive)
		argIndex++
	}

	query += " ORDER BY name"

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
		return nil, fmt.Errorf("failed to get companies: %w", err)
	}
	defer rows.Close()

	var companies []*models.Company
	for rows.Next() {
		company := &models.Company{}
		err := rows.Scan(
			&company.ID, &company.Name, &company.Code,
			&company.IsActive, &company.CreatedAt, &company.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan company: %w", err)
		}
		companies = append(companies, company)
	}

	return companies, nil
}

// GetByID retrieves a company by ID
func (r *CompanyRepository) GetByID(id int64) (*models.Company, error) {
	query := `
		SELECT id, name, code, is_active, created_at, updated_at
		FROM companies
		WHERE id = $1
	`

	company := &models.Company{}
	err := r.db.QueryRow(query, id).Scan(
		&company.ID, &company.Name, &company.Code,
		&company.IsActive, &company.CreatedAt, &company.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("company not found")
		}
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	return company, nil
}

// Create creates a new company
func (r *CompanyRepository) Create(company *models.Company) error {
	query, values := r.BuildInsertQuery(company)

	err := r.db.QueryRow(query, values...).Scan(
		&company.ID, &company.CreatedAt, &company.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}

	return nil
}

// Update updates a company
func (r *CompanyRepository) Update(company *models.Company) error {
	query, values := r.BuildUpdateQuery(company, company.ID)

	err := r.db.QueryRow(query, values...).Scan(&company.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}

	return nil
}

// Delete soft deletes a company
func (r *CompanyRepository) Delete(id int64) error {
	query := `UPDATE companies SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("company not found")
	}

	return nil
}

// GetByCode retrieves a company by code
func (r *CompanyRepository) GetByCode(code string) (*models.Company, error) {
	query := `
		SELECT id, name, code, is_active, created_at, updated_at
		FROM companies
		WHERE code = $1
	`

	company := &models.Company{}
	err := r.db.QueryRow(query, code).Scan(
		&company.ID, &company.Name, &company.Code,
		&company.IsActive, &company.CreatedAt, &company.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("company not found")
		}
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	return company, nil
}

// GetWithStats retrieves a company with statistics
func (r *CompanyRepository) GetWithStats(id int64) (*models.CompanyWithStats, error) {
	// First get the company
	company, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Get user count
	var userCount int
	userQuery := `SELECT COUNT(*) FROM users u JOIN user_roles ur ON u.id = ur.user_id WHERE ur.company_id = $1 AND u.deleted_at IS NULL`
	err = r.db.QueryRow(userQuery, id).Scan(&userCount)
	if err != nil {
		userCount = 0 // Default to 0 if error
	}

	// Get branch count
	var branchCount int
	branchQuery := `SELECT COUNT(*) FROM branches WHERE company_id = $1 AND is_active = true`
	err = r.db.QueryRow(branchQuery, id).Scan(&branchCount)
	if err != nil {
		branchCount = 0 // Default to 0 if error
	}

	return &models.CompanyWithStats{
		Company:       *company,
		TotalUsers:    userCount,
		TotalBranches: branchCount,
	}, nil
}

// Count returns total count of companies with filtering
func (r *CompanyRepository) Count(search string, isActive *bool) (int64, error) {
	query := "SELECT COUNT(*) FROM companies WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR code ILIKE $%d)", argIndex, argIndex+1)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	if isActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *isActive)
	}

	var count int64
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get company count: %w", err)
	}

	return count, nil
}
