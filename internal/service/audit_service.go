package service

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/interfaces"
	"gin-scalable-api/internal/mapper"
)

type AuditService struct {
	auditRepo   interfaces.AuditRepositoryInterface
	auditMapper *mapper.AuditMapper
}

func NewAuditService(auditRepo interfaces.AuditRepositoryInterface) *AuditService {
	return &AuditService{
		auditRepo:   auditRepo,
		auditMapper: mapper.NewAuditMapper(),
	}
}

func (s *AuditService) GetAuditLogs(req *dto.AuditListRequest) (*dto.AuditListResponse, error) {
	// Build filters map from request
	filters := make(map[string]interface{})
	if req.UserID != nil {
		filters["user_id"] = *req.UserID
	}
	if req.Action != "" {
		filters["action"] = req.Action
	}
	if req.Resource != "" {
		filters["resource"] = req.Resource
	}
	if req.Status != "" {
		filters["status"] = req.Status
	}
	if req.Method != "" {
		filters["method"] = req.Method
	}
	if req.DateFrom != "" {
		filters["date_from"] = req.DateFrom
	}
	if req.DateTo != "" {
		filters["date_to"] = req.DateTo
	}
	if req.StatusCode != nil {
		filters["status_code"] = *req.StatusCode
	}
	if req.Search != "" {
		filters["search"] = req.Search
	}

	logs, err := s.auditRepo.GetAll(req.Limit, req.Offset, filters)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.auditRepo.Count(filters)
	if err != nil {
		return nil, err
	}

	return s.auditMapper.ToListResponse(logs, total, req.Limit, req.Offset), nil
}

func (s *AuditService) GetUserAuditLogs(userID int64, limit int) (*dto.AuditListResponse, error) {
	logs, err := s.auditRepo.GetUserLogs(userID, limit)
	if err != nil {
		return nil, err
	}

	return s.auditMapper.ToListResponse(logs, int64(len(logs)), limit, 0), nil
}

func (s *AuditService) GetUserAuditLogsByIdentity(userIdentity string, limit int) (*dto.AuditListResponse, error) {
	logs, err := s.auditRepo.GetUserLogsByIdentity(userIdentity, limit)
	if err != nil {
		return nil, err
	}

	return s.auditMapper.ToListResponse(logs, int64(len(logs)), limit, 0), nil
}

func (s *AuditService) CreateAuditLog(req *dto.CreateAuditLogRequest) (*dto.AuditLogResponse, error) {
	log := s.auditMapper.ToModel(req)

	if log.Metadata == nil {
		log.Metadata = make(map[string]interface{})
	}

	err := s.auditRepo.Create(log)
	if err != nil {
		return nil, err
	}

	return s.auditMapper.ToResponse(log), nil
}

func (s *AuditService) GetAuditLogByID(id int64) (*dto.AuditLogWithUserResponse, error) {
	log, err := s.auditRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.auditMapper.ToWithUserResponse(log), nil
}

func (s *AuditService) GetAuditStats() (*dto.AuditStatsResponse, error) {
	stats, err := s.auditRepo.GetStats()
	if err != nil {
		return nil, err
	}

	return s.auditMapper.ToStatsResponse(stats), nil
}

func (s *AuditService) CleanupOldLogs(daysToKeep int) error {
	return s.auditRepo.CleanupOldLogs(daysToKeep)
}
