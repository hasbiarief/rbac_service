package models

import (
	"gin-scalable-api/pkg/model"
)

type User struct {
	model.BaseModel
	Name         string  `json:"name" db:"name"`
	Email        string  `json:"email" db:"email"`
	UserIdentity *string `json:"user_identity" db:"user_identity"`
	PasswordHash string  `json:"-" db:"password_hash"`
	IsActive     bool    `json:"is_active" db:"is_active"`
}
