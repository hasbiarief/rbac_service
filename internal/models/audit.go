package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID           int64                  `json:"id" db:"id"`
	UserID       *int64                 `json:"user_id" db:"user_id"`
	UserIdentity *string                `json:"user_identity" db:"user_identity"`
	Action       string                 `json:"action" db:"action"`
	Resource     string                 `json:"resource" db:"resource"`
	ResourceID   *string                `json:"resource_id" db:"resource_id"`
	Method       string                 `json:"method" db:"method"`
	URL          string                 `json:"url" db:"url"`
	UserAgent    *string                `json:"user_agent" db:"user_agent"`
	IP           *string                `json:"ip" db:"ip"`
	Status       string                 `json:"status" db:"status"`
	StatusCode   int                    `json:"status_code" db:"status_code"`
	Message      string                 `json:"message" db:"message"`
	Metadata     map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

type AuditLogWithUser struct {
	AuditLog
	UserName  *string `json:"user_name" db:"user_name"`
	UserEmail *string `json:"user_email" db:"user_email"`
}

type AuditStats struct {
	TotalLogs       int64               `json:"total_logs"`
	TodayLogs       int64               `json:"today_logs"`
	SuccessLogs     int64               `json:"success_logs"`
	ErrorLogs       int64               `json:"error_logs"`
	TopActions      []ActionCount       `json:"top_actions"`
	TopUsers        []UserActivityCount `json:"top_users"`
	ActivityByHour  []HourlyActivity    `json:"activity_by_hour"`
	StatusBreakdown []StatusCount       `json:"status_breakdown"`
}

type ActionCount struct {
	Action string `json:"action" db:"action"`
	Count  int64  `json:"count" db:"count"`
}

type UserActivityCount struct {
	UserID       *int64  `json:"user_id" db:"user_id"`
	UserIdentity *string `json:"user_identity" db:"user_identity"`
	UserName     *string `json:"user_name" db:"user_name"`
	Count        int64   `json:"count" db:"count"`
}

type HourlyActivity struct {
	Hour  int   `json:"hour" db:"hour"`
	Count int64 `json:"count" db:"count"`
}

type StatusCount struct {
	Status string `json:"status" db:"status"`
	Count  int64  `json:"count" db:"count"`
}

// JSONB type for PostgreSQL JSON fields
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}
