package audit

import "database/sql"

type Repository interface {
	GetAll(req *AuditListRequest) ([]*AuditLogWithUser, error)
	Count(req *AuditListRequest) (int64, error)
	Create(log *AuditLog) error
	GetByUserID(userID int64, limit int) ([]*AuditLogWithUser, error)
	GetByUserIdentity(identity string, limit int) ([]*AuditLogWithUser, error)
	GetStats() (*AuditStatsResponse, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(req *AuditListRequest) ([]*AuditLogWithUser, error) {
	return nil, nil
}

func (r *repository) Count(req *AuditListRequest) (int64, error) {
	return 0, nil
}

func (r *repository) Create(log *AuditLog) error {
	return nil
}

func (r *repository) GetByUserID(userID int64, limit int) ([]*AuditLogWithUser, error) {
	return nil, nil
}

func (r *repository) GetByUserIdentity(identity string, limit int) ([]*AuditLogWithUser, error) {
	return nil, nil
}

func (r *repository) GetStats() (*AuditStatsResponse, error) {
	return &AuditStatsResponse{}, nil
}
