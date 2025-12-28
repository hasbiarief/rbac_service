package service

import (
	"gin-scalable-api/internal/models"
	"gin-scalable-api/internal/repository"
	"time"
)

type AuditService struct {
	auditRepo *repository.AuditRepository
}

func NewAuditService(auditRepo *repository.AuditRepository) *AuditService {
	return &AuditService{
		auditRepo: auditRepo,
	}
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
	Status       string                 `json:"status"`
	StatusCode   int                    `json:"status_code"`
	Message      string                 `json:"message"`
	Metadata     map[string]interface{} `json:"metadata"`
	CreatedAt    string                 `json:"created_at"`
}

type AuditLogListRequest struct {
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
	UserID   *int64 `form:"user_id"`
	Action   string `form:"action"`
	Resource string `form:"resource"`
	Status   string `form:"status"`
}

type CreateAuditLogRequest struct {
	UserID       *int64                 `json:"user_id"`
	UserIdentity *string                `json:"user_identity"`
	Action       string                 `json:"action" binding:"required"`
	Resource     string                 `json:"resource" binding:"required"`
	ResourceID   *string                `json:"resource_id"`
	Method       string                 `json:"method" binding:"required"`
	URL          string                 `json:"url" binding:"required"`
	Status       string                 `json:"status" binding:"required"`
	StatusCode   int                    `json:"status_code" binding:"required"`
	Message      string                 `json:"message"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type AuditStatsResponse struct {
	TotalLogs   int64 `json:"total_logs"`
	TodayLogs   int64 `json:"today_logs"`
	SuccessLogs int64 `json:"success_logs"`
	ErrorLogs   int64 `json:"error_logs"`
}

func (s *AuditService) GetAuditLogs(req *AuditLogListRequest) ([]*AuditLogResponse, error) {
	logs, err := s.auditRepo.GetAll(req.Limit, req.Offset, req.UserID, req.Action, req.Resource, req.Status)
	if err != nil {
		return nil, err
	}

	var response []*AuditLogResponse
	for _, log := range logs {
		response = append(response, &AuditLogResponse{
			ID:           log.ID,
			UserID:       log.UserID,
			UserIdentity: log.UserIdentity,
			Action:       log.Action,
			Resource:     log.Resource,
			ResourceID:   log.ResourceID,
			Method:       log.Method,
			URL:          log.URL,
			Status:       log.Status,
			StatusCode:   log.StatusCode,
			Message:      log.Message,
			Metadata:     log.Metadata,
			CreatedAt:    log.CreatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}

func (s *AuditService) GetUserAuditLogs(userID int64, limit int) ([]*AuditLogResponse, error) {
	logs, err := s.auditRepo.GetUserLogs(userID, limit)
	if err != nil {
		return nil, err
	}

	var response []*AuditLogResponse
	for _, log := range logs {
		response = append(response, &AuditLogResponse{
			ID:           log.ID,
			UserID:       log.UserID,
			UserIdentity: log.UserIdentity,
			Action:       log.Action,
			Resource:     log.Resource,
			ResourceID:   log.ResourceID,
			Method:       log.Method,
			URL:          log.URL,
			Status:       log.Status,
			StatusCode:   log.StatusCode,
			Message:      log.Message,
			Metadata:     log.Metadata,
			CreatedAt:    log.CreatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}

func (s *AuditService) GetUserAuditLogsByIdentity(userIdentity string, limit int) ([]*AuditLogResponse, error) {
	logs, err := s.auditRepo.GetUserLogsByIdentity(userIdentity, limit)
	if err != nil {
		return nil, err
	}

	var response []*AuditLogResponse
	for _, log := range logs {
		response = append(response, &AuditLogResponse{
			ID:           log.ID,
			UserID:       log.UserID,
			UserIdentity: log.UserIdentity,
			Action:       log.Action,
			Resource:     log.Resource,
			ResourceID:   log.ResourceID,
			Method:       log.Method,
			URL:          log.URL,
			Status:       log.Status,
			StatusCode:   log.StatusCode,
			Message:      log.Message,
			Metadata:     log.Metadata,
			CreatedAt:    log.CreatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}

func (s *AuditService) CreateAuditLog(req *CreateAuditLogRequest) (*AuditLogResponse, error) {
	log := &models.AuditLog{
		UserID:       req.UserID,
		UserIdentity: req.UserIdentity,
		Action:       req.Action,
		Resource:     req.Resource,
		ResourceID:   req.ResourceID,
		Method:       req.Method,
		URL:          req.URL,
		Status:       req.Status,
		StatusCode:   req.StatusCode,
		Message:      req.Message,
		Metadata:     req.Metadata,
	}

	if log.Metadata == nil {
		log.Metadata = make(map[string]interface{})
	}

	err := s.auditRepo.Create(log)
	if err != nil {
		return nil, err
	}

	return &AuditLogResponse{
		ID:           log.ID,
		UserID:       log.UserID,
		UserIdentity: log.UserIdentity,
		Action:       log.Action,
		Resource:     log.Resource,
		ResourceID:   log.ResourceID,
		Method:       log.Method,
		URL:          log.URL,
		Status:       log.Status,
		StatusCode:   log.StatusCode,
		Message:      log.Message,
		Metadata:     log.Metadata,
		CreatedAt:    log.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *AuditService) GetAuditStats() (*AuditStatsResponse, error) {
	stats, err := s.auditRepo.GetStats()
	if err != nil {
		return nil, err
	}

	return &AuditStatsResponse{
		TotalLogs:   stats.TotalLogs,
		TodayLogs:   stats.TodayLogs,
		SuccessLogs: stats.SuccessLogs,
		ErrorLogs:   stats.ErrorLogs,
	}, nil
}
