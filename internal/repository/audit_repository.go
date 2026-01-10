package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gin-scalable-api/internal/models"
	"gin-scalable-api/pkg/model"
	"strconv"
)

type AuditRepository struct {
	*model.Repository
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{
		Repository: model.NewRepository(db),
		db:         db,
	}
}

// Create creates a new audit log entry
func (r *AuditRepository) Create(log *models.AuditLog) error {
	detailsJSON, err := json.Marshal(log.Metadata)
	if err != nil {
		detailsJSON = []byte("{}")
	}

	// Map model fields to database fields
	var resourceID *int64
	if log.ResourceID != nil {
		if id, err := strconv.ParseInt(*log.ResourceID, 10, 64); err == nil {
			resourceID = &id
		}
	}

	// Use IP field for IP address, not URL
	var ipAddress *string
	if log.IP != nil && *log.IP != "" {
		ipAddress = log.IP
	}

	// Use UserAgent field for user agent, not Method
	var userAgent *string
	if log.UserAgent != nil && *log.UserAgent != "" {
		userAgent = log.UserAgent
	}

	success := log.Status == "success"
	var errorMessage *string
	if log.Message != "" && !success {
		errorMessage = &log.Message
	}

	query := `
		INSERT INTO audit_logs (user_id, action, resource, resource_id, details, ip_address, user_agent, success, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	err = r.db.QueryRow(query,
		log.UserID, log.Action, log.Resource, resourceID, detailsJSON,
		ipAddress, userAgent, success, errorMessage,
	).Scan(&log.ID, &log.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// GetStats retrieves audit statistics
func (r *AuditRepository) GetStats() (*models.AuditStats, error) {
	stats := &models.AuditStats{}

	// Get total logs
	err := r.db.QueryRow("SELECT COUNT(*) FROM audit_logs").Scan(&stats.TotalLogs)
	if err != nil {
		return nil, fmt.Errorf("failed to get total logs: %w", err)
	}

	// Get successful logs (success = true) - map to SuccessLogs
	err = r.db.QueryRow("SELECT COUNT(*) FROM audit_logs WHERE success = true").Scan(&stats.SuccessLogs)
	if err != nil {
		return nil, fmt.Errorf("failed to get successful logs: %w", err)
	}

	// Get failed logs (success = false) - map to ErrorLogs
	err = r.db.QueryRow("SELECT COUNT(*) FROM audit_logs WHERE success = false").Scan(&stats.ErrorLogs)
	if err != nil {
		return nil, fmt.Errorf("failed to get failed logs: %w", err)
	}

	// Get today's logs
	err = r.db.QueryRow("SELECT COUNT(*) FROM audit_logs WHERE DATE(created_at) = CURRENT_DATE").Scan(&stats.TodayLogs)
	if err != nil {
		return nil, fmt.Errorf("failed to get today logs: %w", err)
	}

	return stats, nil
}

// GetAll retrieves audit logs with pagination and filtering - interface compatible
func (r *AuditRepository) GetAll(limit, offset int, filters map[string]interface{}) ([]*models.AuditLogWithUser, error) {
	query := `
		SELECT al.id, al.user_id, al.action, al.resource, al.resource_id, 
		       al.details, al.ip_address, al.user_agent, al.success, al.error_message, al.created_at,
		       u.name as user_name, u.email as user_email, u.user_identity
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	// Apply filters
	if userID, ok := filters["user_id"]; ok {
		query += fmt.Sprintf(" AND al.user_id = $%d", argIndex)
		args = append(args, userID)
		argIndex++
	}

	if action, ok := filters["action"]; ok {
		query += fmt.Sprintf(" AND al.action ILIKE $%d", argIndex)
		args = append(args, "%"+action.(string)+"%")
		argIndex++
	}

	if resource, ok := filters["resource"]; ok {
		query += fmt.Sprintf(" AND al.resource ILIKE $%d", argIndex)
		args = append(args, "%"+resource.(string)+"%")
		argIndex++
	}

	if status, ok := filters["status"]; ok {
		successValue := status.(string) == "success"
		query += fmt.Sprintf(" AND al.success = $%d", argIndex)
		args = append(args, successValue)
		argIndex++
	}

	query += " ORDER BY al.created_at DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
		argIndex++
	}

	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}
	defer rows.Close()

	var logs []*models.AuditLogWithUser
	for rows.Next() {
		log := &models.AuditLogWithUser{}
		var detailsJSON []byte
		var resourceID *int64
		var ipAddress *string
		var userAgent *string
		var success bool
		var errorMessage *string
		var userName, userEmail *string
		var userIdentity *string

		err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.Resource, &resourceID,
			&detailsJSON, &ipAddress, &userAgent, &success, &errorMessage,
			&log.CreatedAt, &userName, &userEmail, &userIdentity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}

		// Map database fields to model fields
		if resourceID != nil {
			resourceIDStr := fmt.Sprintf("%d", *resourceID)
			log.ResourceID = &resourceIDStr
		}

		log.IP = ipAddress
		log.UserAgent = userAgent

		if success {
			log.Status = "success"
			log.StatusCode = 200
		} else {
			log.Status = "error"
			log.StatusCode = 500
		}

		if errorMessage != nil {
			log.Message = *errorMessage
		}

		// Set user information
		log.UserName = userName
		log.UserEmail = userEmail
		log.UserIdentity = userIdentity

		if len(detailsJSON) > 0 {
			if err := json.Unmarshal(detailsJSON, &log.Metadata); err != nil {
				log.Metadata = make(map[string]interface{})
			}
		} else {
			log.Metadata = make(map[string]interface{})
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (r *AuditRepository) GetByID(id int64) (*models.AuditLogWithUser, error) {
	query := `
		SELECT al.id, al.user_id, al.action, al.resource, al.resource_id, 
		       al.details, al.ip_address, al.user_agent, al.success, al.error_message, al.created_at,
		       u.name as user_name, u.email as user_email, u.user_identity
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		WHERE al.id = $1
	`

	log := &models.AuditLogWithUser{}
	var detailsJSON []byte
	var resourceID *int64
	var ipAddress *string
	var userAgent *string
	var success bool
	var errorMessage *string
	var userName, userEmail *string
	var userIdentity *string

	err := r.db.QueryRow(query, id).Scan(
		&log.ID, &log.UserID, &log.Action, &log.Resource, &resourceID,
		&detailsJSON, &ipAddress, &userAgent, &success, &errorMessage,
		&log.CreatedAt, &userName, &userEmail, &userIdentity,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("audit log not found")
		}
		return nil, fmt.Errorf("failed to get audit log: %w", err)
	}

	// Map database fields to model fields
	if resourceID != nil {
		resourceIDStr := fmt.Sprintf("%d", *resourceID)
		log.ResourceID = &resourceIDStr
	}

	log.IP = ipAddress
	log.UserAgent = userAgent

	if success {
		log.Status = "success"
		log.StatusCode = 200
	} else {
		log.Status = "error"
		log.StatusCode = 500
	}

	if errorMessage != nil {
		log.Message = *errorMessage
	}

	// Set user information
	log.UserName = userName
	log.UserEmail = userEmail
	log.UserIdentity = userIdentity

	if len(detailsJSON) > 0 {
		if err := json.Unmarshal(detailsJSON, &log.Metadata); err != nil {
			log.Metadata = make(map[string]interface{})
		}
	} else {
		log.Metadata = make(map[string]interface{})
	}

	return log, nil
}

func (r *AuditRepository) Count(filters map[string]interface{}) (int64, error) {
	query := "SELECT COUNT(*) FROM audit_logs al WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	// Apply filters
	if userID, ok := filters["user_id"]; ok {
		query += fmt.Sprintf(" AND al.user_id = $%d", argIndex)
		args = append(args, userID)
		argIndex++
	}

	if action, ok := filters["action"]; ok {
		query += fmt.Sprintf(" AND al.action ILIKE $%d", argIndex)
		args = append(args, "%"+action.(string)+"%")
		argIndex++
	}

	if resource, ok := filters["resource"]; ok {
		query += fmt.Sprintf(" AND al.resource ILIKE $%d", argIndex)
		args = append(args, "%"+resource.(string)+"%")
		argIndex++
	}

	if status, ok := filters["status"]; ok {
		successValue := status.(string) == "success"
		query += fmt.Sprintf(" AND al.success = $%d", argIndex)
		args = append(args, successValue)
	}

	var count int64
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get audit log count: %w", err)
	}

	return count, nil
}

func (r *AuditRepository) GetUserLogs(userID int64, limit int) ([]*models.AuditLogWithUser, error) {
	filters := map[string]interface{}{
		"user_id": userID,
	}
	return r.GetAll(limit, 0, filters)
}

func (r *AuditRepository) GetUserLogsByIdentity(userIdentity string, limit int) ([]*models.AuditLogWithUser, error) {
	query := `
		SELECT al.id, al.user_id, al.action, al.resource, al.resource_id, 
		       al.details, al.ip_address, al.user_agent, al.success, al.error_message, al.created_at,
		       u.name as user_name, u.email as user_email, u.user_identity
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		WHERE u.user_identity = $1
		ORDER BY al.created_at DESC
	`
	args := []interface{}{userIdentity}

	if limit > 0 {
		query += " LIMIT $2"
		args = append(args, limit)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user audit logs by identity: %w", err)
	}
	defer rows.Close()

	var logs []*models.AuditLogWithUser
	for rows.Next() {
		log := &models.AuditLogWithUser{}
		var detailsJSON []byte
		var resourceID *int64
		var ipAddress *string
		var userAgent *string
		var success bool
		var errorMessage *string
		var userName, userEmail *string
		var userIdentityResult *string

		err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.Resource, &resourceID,
			&detailsJSON, &ipAddress, &userAgent, &success, &errorMessage,
			&log.CreatedAt, &userName, &userEmail, &userIdentityResult,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}

		// Map database fields to model fields
		if resourceID != nil {
			resourceIDStr := fmt.Sprintf("%d", *resourceID)
			log.ResourceID = &resourceIDStr
		}

		log.IP = ipAddress
		log.UserAgent = userAgent

		if success {
			log.Status = "success"
			log.StatusCode = 200
		} else {
			log.Status = "error"
			log.StatusCode = 500
		}

		if errorMessage != nil {
			log.Message = *errorMessage
		}

		// Set user information
		log.UserName = userName
		log.UserEmail = userEmail
		log.UserIdentity = userIdentityResult

		if len(detailsJSON) > 0 {
			if err := json.Unmarshal(detailsJSON, &log.Metadata); err != nil {
				log.Metadata = make(map[string]interface{})
			}
		} else {
			log.Metadata = make(map[string]interface{})
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (r *AuditRepository) CleanupOldLogs(daysToKeep int) error {
	query := `DELETE FROM audit_logs WHERE created_at < CURRENT_DATE - INTERVAL '%d days'`
	_, err := r.db.Exec(fmt.Sprintf(query, daysToKeep))
	if err != nil {
		return fmt.Errorf("failed to cleanup old audit logs: %w", err)
	}
	return nil
}
