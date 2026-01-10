package mapper

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/models"
	"time"
)

// AuditMapper handles conversion between audit models and DTOs
type AuditMapper struct{}

// NewAuditMapper creates a new audit mapper
func NewAuditMapper() *AuditMapper {
	return &AuditMapper{}
}

// ToResponse converts model to response DTO
func (m *AuditMapper) ToResponse(audit *models.AuditLog) *dto.AuditLogResponse {
	if audit == nil {
		return nil
	}

	return &dto.AuditLogResponse{
		ID:           audit.ID,
		UserID:       audit.UserID,
		UserIdentity: audit.UserIdentity,
		Action:       audit.Action,
		Resource:     audit.Resource,
		ResourceID:   audit.ResourceID,
		Method:       audit.Method,
		URL:          audit.URL,
		UserAgent:    audit.UserAgent,
		IP:           audit.IP,
		Status:       audit.Status,
		StatusCode:   audit.StatusCode,
		Message:      audit.Message,
		Metadata:     audit.Metadata,
		CreatedAt:    audit.CreatedAt.Format(time.RFC3339),
	}
}

// ToWithUserResponse converts model with user to response DTO
func (m *AuditMapper) ToWithUserResponse(audit *models.AuditLogWithUser) *dto.AuditLogWithUserResponse {
	if audit == nil {
		return nil
	}

	return &dto.AuditLogWithUserResponse{
		AuditLogResponse: *m.ToResponse(&audit.AuditLog),
		UserName:         audit.UserName,
		UserEmail:        audit.UserEmail,
	}
}

// ToWithUserResponseList converts model slice to response DTO slice
func (m *AuditMapper) ToWithUserResponseList(audits []*models.AuditLogWithUser) []*dto.AuditLogWithUserResponse {
	if audits == nil {
		return nil
	}

	responses := make([]*dto.AuditLogWithUserResponse, len(audits))
	for i, audit := range audits {
		responses[i] = m.ToWithUserResponse(audit)
	}
	return responses
}

// ToStatsResponse converts stats model to response DTO
func (m *AuditMapper) ToStatsResponse(stats *models.AuditStats) *dto.AuditStatsResponse {
	if stats == nil {
		return nil
	}

	// Convert top actions
	topActions := make([]dto.ActionCountResponse, len(stats.TopActions))
	for i, action := range stats.TopActions {
		topActions[i] = dto.ActionCountResponse{
			Action: action.Action,
			Count:  action.Count,
		}
	}

	// Convert top users
	topUsers := make([]dto.UserActivityCountResponse, len(stats.TopUsers))
	for i, user := range stats.TopUsers {
		topUsers[i] = dto.UserActivityCountResponse{
			UserID:       user.UserID,
			UserIdentity: user.UserIdentity,
			UserName:     user.UserName,
			Count:        user.Count,
		}
	}

	// Convert activity by hour
	activityByHour := make([]dto.HourlyActivityResponse, len(stats.ActivityByHour))
	for i, activity := range stats.ActivityByHour {
		activityByHour[i] = dto.HourlyActivityResponse{
			Hour:  activity.Hour,
			Count: activity.Count,
		}
	}

	// Convert status breakdown
	statusBreakdown := make([]dto.StatusCountResponse, len(stats.StatusBreakdown))
	for i, status := range stats.StatusBreakdown {
		statusBreakdown[i] = dto.StatusCountResponse{
			Status: status.Status,
			Count:  status.Count,
		}
	}

	return &dto.AuditStatsResponse{
		TotalLogs:       stats.TotalLogs,
		TodayLogs:       stats.TodayLogs,
		SuccessLogs:     stats.SuccessLogs,
		ErrorLogs:       stats.ErrorLogs,
		TopActions:      topActions,
		TopUsers:        topUsers,
		ActivityByHour:  activityByHour,
		StatusBreakdown: statusBreakdown,
	}
}

// ToModel converts create request DTO to model
func (m *AuditMapper) ToModel(req *dto.CreateAuditLogRequest) *models.AuditLog {
	if req == nil {
		return nil
	}

	return &models.AuditLog{
		UserID:       req.UserID,
		UserIdentity: req.UserIdentity,
		Action:       req.Action,
		Resource:     req.Resource,
		ResourceID:   req.ResourceID,
		Method:       req.Method,
		URL:          req.URL,
		UserAgent:    req.UserAgent,
		IP:           req.IP,
		Status:       req.Status,
		StatusCode:   req.StatusCode,
		Message:      req.Message,
		Metadata:     req.Metadata,
	}
}

// ToListResponse creates paginated list response
func (m *AuditMapper) ToListResponse(audits []*models.AuditLogWithUser, total int64, limit, offset int) *dto.AuditListResponse {
	return &dto.AuditListResponse{
		Data:    m.ToWithUserResponseList(audits),
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(audits)) < total,
	}
}
