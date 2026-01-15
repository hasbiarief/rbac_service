package audit

import (
	"github.com/go-playground/validator/v10"
)

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
