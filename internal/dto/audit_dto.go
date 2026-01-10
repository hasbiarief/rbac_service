package dto

// Audit Request DTO
type AuditListRequest struct {
	Limit      int    `form:"limit"`
	Offset     int    `form:"offset"`
	Search     string `form:"search"`
	UserID     *int64 `form:"user_id"`
	Action     string `form:"action"`
	Resource   string `form:"resource"`
	Status     string `form:"status"`
	Method     string `form:"method"`
	DateFrom   string `form:"date_from"` // Format: 2006-01-02
	DateTo     string `form:"date_to"`   // Format: 2006-01-02
	StatusCode *int   `form:"status_code"`
}

// Create Audit Log Request DTO
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

// Audit Response DTO
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
}

// Audit Log With User Response DTO
type AuditLogWithUserResponse struct {
	AuditLogResponse
	UserName  *string `json:"user_name"`
	UserEmail *string `json:"user_email"`
}

// Audit Stats Response DTO
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

// Action Count Response DTO
type ActionCountResponse struct {
	Action string `json:"action"`
	Count  int64  `json:"count"`
}

// User Activity Count Response DTO
type UserActivityCountResponse struct {
	UserID       *int64  `json:"user_id"`
	UserIdentity *string `json:"user_identity"`
	UserName     *string `json:"user_name"`
	Count        int64   `json:"count"`
}

// Hourly Activity Response DTO
type HourlyActivityResponse struct {
	Hour  int   `json:"hour"`
	Count int64 `json:"count"`
}

// Status Count Response DTO
type StatusCountResponse struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// Audit List Response DTO
type AuditListResponse struct {
	Data    []*AuditLogWithUserResponse `json:"data"`
	Total   int64                       `json:"total"`
	Limit   int                         `json:"limit"`
	Offset  int                         `json:"offset"`
	HasMore bool                        `json:"has_more"`
}
