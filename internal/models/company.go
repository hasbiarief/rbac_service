package models

import "time"

type Company struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Code      string    `json:"code" db:"code"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (Company) TableName() string {
	return "companies"
}

type CompanyWithStats struct {
	Company
	TotalUsers    int `json:"total_users" db:"total_users"`
	TotalBranches int `json:"total_branches" db:"total_branches"`
}
