package audit

import (
	"github.com/go-playground/validator/v10"
)

type AuditListRequest struct {
	Limit      int    `form:"limit"`
	Offset     int    `form:"offset"`
	Search     string `form:"search"`
	UserID     *int64 `form:"user_id"`
	Action     string `form:"action"`
	Resource   string `form:"resource"`
	Status     string `form:"status"`
	Method     string `form:"method"`
	DateFrom   string `form:"date_from"`
	DateTo     string `form:"date_to"`
	StatusCode *int   `form:"status_code"`
}

type CreateAuditLogRequest struct {
	UserID       *int64                 `json:"user_id"`
	UserIdentity *string                `json:"user_identity"`
	Action       string                 `json:"action" validate:"required"`
	Resource     string                 `json:"resource" validate:"required"`
	ResourceID   *string                `json:"resource_id"`
	Method       string                 `json:"method" validate:"required"`
	URL          string                 `json:"url" validate:"required"`
	UserAgent    *string                `json:"user_agent"`
	IP           *string                `json:"ip"`
	Status       string                 `json:"status" validate:"required"`
	StatusCode   int                    `json:"status_code" validate:"required"`
	Message      string                 `json:"message"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type AuditLogResponse struct {
	ID           int64                  `json:"id"`
	UserID       *int64                 `json:"user_id"`
	UserIdentity *string                `json:"user_identity"`
	Action       string                 `json:"action"`
	Resource     string                 `json:"resource"`
	ResourceID   *string                `json:"resource_id"`
	Method       string                 `json:"method"`
	URL          string                 `json:"url"`
	UserAgent    *string                `json:"user_agent"`
	IP           *string                `json:"ip"`
	Status       string                 `json:"status"`
	StatusCode   int                    `json:"status_code"`
	Message      string                 `json:"message"`
	Metadata     map[string]interface{} `json:"metadata"`
	CreatedAt    string                 `json:"created_at"`
	UserName     *string                `json:"user_name,omitempty"`
	UserEmail    *string                `json:"user_email,omitempty"`
}

type AuditListResponse struct {
	Data    []*AuditLogResponse `json:"data"`
	Total   int64               `json:"total"`
	Limit   int                 `json:"limit"`
	Offset  int                 `json:"offset"`
	HasMore bool                `json:"has_more"`
}

type AuditStatsResponse struct {
	TotalLogs       int64                       `json:"total_logs"`
	TodayLogs       int64                       `json:"today_logs"`
	SuccessLogs     int64                       `json:"success_logs"`
	ErrorLogs       int64                       `json:"error_logs"`
	TopActions      []ActionCountResponse       `json:"top_actions"`
	TopUsers        []UserActivityCountResponse `json:"top_users"`
	ActivityByHour  []HourlyActivityResponse    `json:"activity_by_hour"`
	StatusBreakdown []StatusCountResponse       `json:"status_breakdown"`
}

type ActionCountResponse struct {
	Action string `json:"action"`
	Count  int64  `json:"count"`
}

type UserActivityCountResponse struct {
	UserID       *int64  `json:"user_id"`
	UserIdentity *string `json:"user_identity"`
	UserName     *string `json:"user_name"`
	Count        int64   `json:"count"`
}

type HourlyActivityResponse struct {
	Hour  int   `json:"hour"`
	Count int64 `json:"count"`
}

type StatusCountResponse struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// Validation functions
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateAuditListRequest validates audit list request
func ValidateAuditListRequest(req *AuditListRequest) error {
	return validate.Struct(req)
}

// ValidateCreateAuditLogRequest validates create audit log request
func ValidateCreateAuditLogRequest(req *CreateAuditLogRequest) error {
	return validate.Struct(req)
}
