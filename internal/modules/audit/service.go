package audit

import "time"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAuditLogs(req *AuditListRequest) (*AuditListResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	logs, err := s.repo.GetAll(req)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(req)
	if err != nil {
		return nil, err
	}

	var responses []*AuditLogResponse
	for _, log := range logs {
		responses = append(responses, toAuditLogResponse(log))
	}

	return &AuditListResponse{
		Data:    responses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(responses)) < total,
	}, nil
}

func (s *Service) CreateAuditLog(req *CreateAuditLogRequest) (*AuditLogResponse, error) {
	log := &AuditLog{
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

	if err := s.repo.Create(log); err != nil {
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
		UserAgent:    log.UserAgent,
		IP:           log.IP,
		Status:       log.Status,
		StatusCode:   log.StatusCode,
		Message:      log.Message,
		Metadata:     log.Metadata,
		CreatedAt:    log.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Service) GetUserAuditLogs(userID int64, limit int) ([]*AuditLogResponse, error) {
	if limit <= 0 {
		limit = 50
	}

	logs, err := s.repo.GetByUserID(userID, limit)
	if err != nil {
		return nil, err
	}

	var responses []*AuditLogResponse
	for _, log := range logs {
		responses = append(responses, toAuditLogResponse(log))
	}

	return responses, nil
}

func (s *Service) GetUserAuditLogsByIdentity(identity string, limit int) ([]*AuditLogResponse, error) {
	if limit <= 0 {
		limit = 50
	}

	logs, err := s.repo.GetByUserIdentity(identity, limit)
	if err != nil {
		return nil, err
	}

	var responses []*AuditLogResponse
	for _, log := range logs {
		responses = append(responses, toAuditLogResponse(log))
	}

	return responses, nil
}

func (s *Service) GetAuditStats() (*AuditStatsResponse, error) {
	return s.repo.GetStats()
}

func toAuditLogResponse(log *AuditLogWithUser) *AuditLogResponse {
	if log == nil {
		return nil
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
		UserAgent:    log.UserAgent,
		IP:           log.IP,
		Status:       log.Status,
		StatusCode:   log.StatusCode,
		Message:      log.Message,
		Metadata:     log.Metadata,
		CreatedAt:    log.CreatedAt.Format(time.RFC3339),
		UserName:     log.UserName,
		UserEmail:    log.UserEmail,
	}
}
