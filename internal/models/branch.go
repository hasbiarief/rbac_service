package models

import "time"

type Branch struct {
	ID        int64     `json:"id" db:"id"`
	CompanyID int64     `json:"company_id" db:"company_id"`
	ParentID  *int64    `json:"parent_id" db:"parent_id"`
	Name      string    `json:"name" db:"name"`
	Code      string    `json:"code" db:"code"`
	Level     int       `json:"level" db:"level"`
	Path      string    `json:"path" db:"path"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type BranchWithCompany struct {
	Branch
	CompanyName string `json:"company_name" db:"company_name"`
	CompanyCode string `json:"company_code" db:"company_code"`
}

type BranchHierarchy struct {
	Branch
	Children []BranchHierarchy `json:"children,omitempty"`
}

type BranchWithStats struct {
	Branch
	TotalUsers       int `json:"total_users" db:"total_users"`
	TotalSubBranches int `json:"total_sub_branches" db:"total_sub_branches"`
}

func (Branch) TableName() string {
	return "branches"
}
