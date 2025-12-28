package validation

import "gin-scalable-api/middleware"

// Audit validation rules
var CreateAuditLogValidation = middleware.ValidationRules{
	Body: &struct {
		UserID       *int64                 `json:"user_id"`
		UserIdentity *string                `json:"user_identity"`
		Action       string                 `json:"action" validate:"required"`
		Resource     string                 `json:"resource" validate:"required"`
		ResourceID   *string                `json:"resource_id"`
		Method       string                 `json:"method" validate:"required"`
		URL          string                 `json:"url" validate:"required"`
		Status       string                 `json:"status" validate:"required"`
		StatusCode   int                    `json:"status_code" validate:"required"`
		Message      string                 `json:"message"`
		Metadata     map[string]interface{} `json:"metadata"`
	}{},
}
