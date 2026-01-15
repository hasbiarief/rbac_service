package audit

import "time"

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
