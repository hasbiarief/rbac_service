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

// GetAll retrieves audit logs with pagination and filtering
func (r *AuditRepository) GetAll(limit, offset int, userID *int64, action, resource, status string) ([]*models.AuditLog, error) {
	query := `
		SELECT id, user_id, action, resource, resource_id, 
		       details, ip_address, user_agent, success, error_message, created_at
		FROM audit_logs
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if userID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *userID)
		argIndex++
	}

	if action != "" {
		query += fmt.Sprintf(" AND action ILIKE $%d", argIndex)
		args = append(args, "%"+action+"%")
		argIndex++
	}

	if resource != "" {
		query += fmt.Sprintf(" AND resource ILIKE $%d", argIndex)
		args = append(args, "%"+resource+"%")
		argIndex++
	}

	if status != "" {
		successValue := status == "success"
		query += fmt.Sprintf(" AND success = $%d", argIndex)
		args = append(args, successValue)
		argIndex++
	}

	query += " ORDER BY created_at DESC"

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

	var logs []*models.AuditLog
	for rows.Next() {
		log := &models.AuditLog{}
		var detailsJSON []byte
		var resourceID *int64
		var ipAddress *string
		var userAgent *string
		var success bool
		var errorMessage *string

		err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.Resource, &resourceID,
			&detailsJSON, &ipAddress, &userAgent, &success, &errorMessage,
			&log.CreatedAt,
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

// GetUserLogs retrieves audit logs for a specific user
func (r *AuditRepository) GetUserLogs(userID int64, limit int) ([]*models.AuditLog, error) {
	query := `
		SELECT id, user_id, action, resource, resource_id, 
		       details, ip_address, user_agent, success, error_message, created_at
		FROM audit_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	args := []interface{}{userID}

	if limit > 0 {
		query += " LIMIT $2"
		args = append(args, limit)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user audit logs: %w", err)
	}
	defer rows.Close()

	var logs []*models.AuditLog
	for rows.Next() {
		log := &models.AuditLog{}
		var detailsJSON []byte
		var resourceID *int64
		var ipAddress *string
		var userAgent *string
		var success bool
		var errorMessage *string

		err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.Resource, &resourceID,
			&detailsJSON, &ipAddress, &userAgent, &success, &errorMessage,
			&log.CreatedAt,
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

// GetUserLogsByIdentity retrieves audit logs for a user by identity
func (r *AuditRepository) GetUserLogsByIdentity(userIdentity string, limit int) ([]*models.AuditLog, error) {
	query := `
		SELECT al.id, al.user_id, al.action, al.resource, al.resource_id, 
		       al.details, al.ip_address, al.user_agent, al.success, al.error_message, al.created_at
		FROM audit_logs al
		JOIN users u ON al.user_id = u.id
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

	var logs []*models.AuditLog
	for rows.Next() {
		log := &models.AuditLog{}
		var detailsJSON []byte
		var resourceID *int64
		var ipAddress *string
		var userAgent *string
		var success bool
		var errorMessage *string

		err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.Resource, &resourceID,
			&detailsJSON, &ipAddress, &userAgent, &success, &errorMessage,
			&log.CreatedAt,
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

		// Set user identity from the join
		log.UserIdentity = &userIdentity

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
