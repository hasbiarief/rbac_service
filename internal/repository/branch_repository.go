package repository

import (
	"database/sql"
	"fmt"
	"gin-scalable-api/internal/models"
	"gin-scalable-api/pkg/model"
	"gin-scalable-api/pkg/query"
)

type BranchRepository struct {
	*model.Repository
	db *sql.DB
}

func NewBranchRepository(db *sql.DB) *BranchRepository {
	return &BranchRepository{
		Repository: model.NewRepository(db),
		db:         db,
	}
}

// GetAll retrieves all branches with pagination and filtering
func (r *BranchRepository) GetAll(limit, offset int, search string, companyID *int64, isActive *bool) ([]*models.Branch, error) {
	baseQuery := `
		SELECT id, company_id, name, code, parent_id, level, path, is_active, created_at, updated_at
		FROM branches
	`

	qb := query.NewQueryBuilder(baseQuery)

	if companyID != nil {
		qb.AddCondition("company_id = $%d", *companyID)
	}

	if search != "" {
		qb.AddLikeCondition([]string{"name", "code"}, search)
	}

	if isActive != nil {
		qb.AddCondition("is_active = $%d", *isActive)
	}

	qb.AddOrderBy("company_id, level, name").
		AddLimit(limit).
		AddOffset(offset)

	query, args := qb.Build()

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}
	defer rows.Close()

	var branches []*models.Branch
	for rows.Next() {
		branch := &models.Branch{}
		err := rows.Scan(
			&branch.ID, &branch.CompanyID, &branch.Name, &branch.Code,
			&branch.ParentID, &branch.Level, &branch.Path, &branch.IsActive,
			&branch.CreatedAt, &branch.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan branch: %w", err)
		}
		branches = append(branches, branch)
	}

	return branches, nil
}

// GetByID retrieves a branch by ID
func (r *BranchRepository) GetByID(id int64) (*models.Branch, error) {
	query := `
		SELECT id, company_id, name, code, parent_id, level, path, is_active, created_at, updated_at
		FROM branches
		WHERE id = $1
	`

	branch := &models.Branch{}
	err := r.db.QueryRow(query, id).Scan(
		&branch.ID, &branch.CompanyID, &branch.Name, &branch.Code,
		&branch.ParentID, &branch.Level, &branch.Path, &branch.IsActive,
		&branch.CreatedAt, &branch.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("branch not found")
		}
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}

	return branch, nil
}

// Create creates a new branch
func (r *BranchRepository) Create(branch *models.Branch) error {
	// Calculate level and path based on parent
	if branch.ParentID != nil {
		parent, err := r.GetByID(*branch.ParentID)
		if err != nil {
			return fmt.Errorf("failed to get parent branch: %w", err)
		}
		branch.Level = parent.Level + 1
		branch.Path = fmt.Sprintf("%s/%d", parent.Path, *branch.ParentID)
	} else {
		branch.Level = 1
		branch.Path = "/"
	}

	query := `
		INSERT INTO branches (company_id, name, code, parent_id, level, path, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query,
		branch.CompanyID, branch.Name, branch.Code, branch.ParentID,
		branch.Level, branch.Path, branch.IsActive,
	).Scan(&branch.ID, &branch.CreatedAt, &branch.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	return nil
}

// Update updates a branch
func (r *BranchRepository) Update(branch *models.Branch) error {
	query, values := r.BuildUpdateQuery(branch, branch.ID)

	err := r.db.QueryRow(query, values...).Scan(&branch.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update branch: %w", err)
	}

	return nil
}

// Delete deletes a branch
func (r *BranchRepository) Delete(id int64) error {
	// Check if branch has children
	var childCount int
	err := r.db.QueryRow("SELECT COUNT(*) FROM branches WHERE parent_id = $1", id).Scan(&childCount)
	if err != nil {
		return fmt.Errorf("failed to check branch children: %w", err)
	}

	if childCount > 0 {
		return fmt.Errorf("cannot delete branch with children")
	}

	query := `DELETE FROM branches WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("branch not found")
	}

	return nil
}

// GetByCompany retrieves branches for a specific company
func (r *BranchRepository) GetByCompany(companyID int64, includeHierarchy bool) ([]*models.Branch, error) {
	query := `
		SELECT id, company_id, name, code, parent_id, level, path, is_active, created_at, updated_at
		FROM branches
		WHERE company_id = $1 AND is_active = true
	`

	if includeHierarchy {
		query += " ORDER BY level, name"
	} else {
		query += " ORDER BY name"
	}

	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company branches: %w", err)
	}
	defer rows.Close()

	var branches []*models.Branch
	for rows.Next() {
		branch := &models.Branch{}
		err := rows.Scan(
			&branch.ID, &branch.CompanyID, &branch.Name, &branch.Code,
			&branch.ParentID, &branch.Level, &branch.Path, &branch.IsActive,
			&branch.CreatedAt, &branch.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan branch: %w", err)
		}
		branches = append(branches, branch)
	}

	return branches, nil
}

// GetChildren retrieves child branches of a parent branch
func (r *BranchRepository) GetChildren(parentID int64) ([]*models.Branch, error) {
	query := `
		SELECT id, company_id, name, code, parent_id, level, path, is_active, created_at, updated_at
		FROM branches
		WHERE parent_id = $1 AND is_active = true
		ORDER BY name
	`

	rows, err := r.db.Query(query, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch children: %w", err)
	}
	defer rows.Close()

	var branches []*models.Branch
	for rows.Next() {
		branch := &models.Branch{}
		err := rows.Scan(
			&branch.ID, &branch.CompanyID, &branch.Name, &branch.Code,
			&branch.ParentID, &branch.Level, &branch.Path, &branch.IsActive,
			&branch.CreatedAt, &branch.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan branch: %w", err)
		}
		branches = append(branches, branch)
	}

	return branches, nil
}
