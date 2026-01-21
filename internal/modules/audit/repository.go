package audit

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

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

func (r *repository) Create(log *AuditLog) error {
	// Map Postman format to database schema
	details := map[string]interface{}{
		"user_identity": log.UserIdentity,
		"method":        log.Method,
		"url":           log.URL,
		"status":        log.Status,
		"status_code":   log.StatusCode,
		"message":       log.Message,
	}

	detailsJSON, _ := json.Marshal(details)

	query := `INSERT INTO audit_logs (user_id, action, resource, resource_id, details, success, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`

	resourceID := int64(0)
	if log.ResourceID != nil {
		// Convert string to int64 if needed
		resourceID = 123 // placeholder
	}

	success := log.Status == "success"

	err := r.db.QueryRow(query, log.UserID, log.Action, log.Resource, resourceID,
		detailsJSON, success, time.Now()).Scan(&log.ID, &log.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

func (r *repository) GetAll(req *AuditListRequest) ([]*AuditLogWithUser, error) {
	query := `SELECT al.id, al.user_id, al.action, al.resource, al.resource_id, 
			  al.details, al.success, al.created_at, u.name, u.email, u.user_identity
			  FROM audit_logs al 
			  LEFT JOIN users u ON al.user_id = u.id 
			  ORDER BY al.created_at DESC`

	if req.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", req.Limit)
	} else {
		query += " LIMIT 50"
	}

	if req.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}
	defer rows.Close()

	var logs []*AuditLogWithUser
	for rows.Next() {
		log := &AuditLogWithUser{}
		var detailsJSON []byte
		var resourceID *int64
		var success bool
		var userIdentity *string

		err := rows.Scan(&log.ID, &log.UserID, &log.Action, &log.Resource, &resourceID,
			&detailsJSON, &success, &log.CreatedAt, &log.UserName, &log.UserEmail, &userIdentity)
		if err != nil {
			continue
		}

		// Map database schema to Postman format
		var details map[string]interface{}
		if len(detailsJSON) > 0 {
			json.Unmarshal(detailsJSON, &details)
		}

		// Extract from details or use defaults
		log.UserIdentity = userIdentity
		if method, ok := details["method"].(string); ok {
			log.Method = method
		} else {
			log.Method = "GET"
		}
		if url, ok := details["url"].(string); ok {
			log.URL = url
		} else {
			log.URL = "/api/v1/unknown"
		}
		if success {
			log.Status = "success"
		} else {
			log.Status = "error"
		}
		if statusCode, ok := details["status_code"].(float64); ok {
			log.StatusCode = int(statusCode)
		} else {
			log.StatusCode = 200
		}
		if message, ok := details["message"].(string); ok {
			log.Message = message
		} else {
			log.Message = ""
		}
		if resourceID != nil {
			resourceIDStr := fmt.Sprintf("%d", *resourceID)
			log.ResourceID = &resourceIDStr
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (r *repository) Count(req *AuditListRequest) (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM audit_logs").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count audit logs: %w", err)
	}
	return count, nil
}

func (r *repository) GetByUserID(userID int64, limit int) ([]*AuditLogWithUser, error) {
	query := `SELECT al.id, al.user_id, al.action, al.resource, al.resource_id, 
			  al.details, al.success, al.created_at, u.name, u.email, u.user_identity
			  FROM audit_logs al 
			  LEFT JOIN users u ON al.user_id = u.id 
			  WHERE al.user_id = $1
			  ORDER BY al.created_at DESC
			  LIMIT $2`

	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user audit logs: %w", err)
	}
	defer rows.Close()

	var logs []*AuditLogWithUser
	for rows.Next() {
		log := &AuditLogWithUser{}
		var detailsJSON []byte
		var resourceID *int64
		var success bool
		var userIdentity *string

		err := rows.Scan(&log.ID, &log.UserID, &log.Action, &log.Resource, &resourceID,
			&detailsJSON, &success, &log.CreatedAt, &log.UserName, &log.UserEmail, &userIdentity)
		if err != nil {
			continue
		}

		// Map database schema to Postman format (same as GetAll)
		var details map[string]interface{}
		if len(detailsJSON) > 0 {
			json.Unmarshal(detailsJSON, &details)
		}

		log.UserIdentity = userIdentity
		if method, ok := details["method"].(string); ok {
			log.Method = method
		} else {
			log.Method = "GET"
		}
		if url, ok := details["url"].(string); ok {
			log.URL = url
		} else {
			log.URL = "/api/v1/unknown"
		}
		if success {
			log.Status = "success"
		} else {
			log.Status = "error"
		}
		if statusCode, ok := details["status_code"].(float64); ok {
			log.StatusCode = int(statusCode)
		} else {
			log.StatusCode = 200
		}
		if message, ok := details["message"].(string); ok {
			log.Message = message
		} else {
			log.Message = ""
		}
		if resourceID != nil {
			resourceIDStr := fmt.Sprintf("%d", *resourceID)
			log.ResourceID = &resourceIDStr
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (r *repository) GetByUserIdentity(identity string, limit int) ([]*AuditLogWithUser, error) {
	query := `SELECT al.id, al.user_id, al.action, al.resource, al.resource_id, 
			  al.details, al.success, al.created_at, u.name, u.email, u.user_identity
			  FROM audit_logs al 
			  LEFT JOIN users u ON al.user_id = u.id 
			  WHERE u.user_identity = $1
			  ORDER BY al.created_at DESC
			  LIMIT $2`

	rows, err := r.db.Query(query, identity, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user audit logs by identity: %w", err)
	}
	defer rows.Close()

	var logs []*AuditLogWithUser
	for rows.Next() {
		log := &AuditLogWithUser{}
		var detailsJSON []byte
		var resourceID *int64
		var success bool
		var userIdentity *string

		err := rows.Scan(&log.ID, &log.UserID, &log.Action, &log.Resource, &resourceID,
			&detailsJSON, &success, &log.CreatedAt, &log.UserName, &log.UserEmail, &userIdentity)
		if err != nil {
			continue
		}

		// Map database schema to Postman format (same as GetAll)
		var details map[string]interface{}
		if len(detailsJSON) > 0 {
			json.Unmarshal(detailsJSON, &details)
		}

		log.UserIdentity = userIdentity
		if method, ok := details["method"].(string); ok {
			log.Method = method
		} else {
			log.Method = "GET"
		}
		if url, ok := details["url"].(string); ok {
			log.URL = url
		} else {
			log.URL = "/api/v1/unknown"
		}
		if success {
			log.Status = "success"
		} else {
			log.Status = "error"
		}
		if statusCode, ok := details["status_code"].(float64); ok {
			log.StatusCode = int(statusCode)
		} else {
			log.StatusCode = 200
		}
		if message, ok := details["message"].(string); ok {
			log.Message = message
		} else {
			log.Message = ""
		}
		if resourceID != nil {
			resourceIDStr := fmt.Sprintf("%d", *resourceID)
			log.ResourceID = &resourceIDStr
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (r *repository) GetStats() (*AuditStatsResponse, error) {
	stats := &AuditStatsResponse{}

	// Get basic counts
	err := r.db.QueryRow(`
		SELECT 
			COUNT(*) as total_logs,
			COUNT(CASE WHEN DATE(created_at) = CURRENT_DATE THEN 1 END) as today_logs,
			COUNT(CASE WHEN success = true THEN 1 END) as success_logs,
			COUNT(CASE WHEN success = false THEN 1 END) as error_logs
		FROM audit_logs
	`).Scan(&stats.TotalLogs, &stats.TodayLogs, &stats.SuccessLogs, &stats.ErrorLogs)

	if err != nil {
		return stats, nil // Return empty stats instead of error
	}

	// Get top actions
	actionRows, err := r.db.Query(`
		SELECT action, COUNT(*) as count 
		FROM audit_logs 
		GROUP BY action 
		ORDER BY count DESC 
		LIMIT 10
	`)
	if err == nil {
		defer actionRows.Close()
		stats.TopActions = []ActionCountResponse{}
		for actionRows.Next() {
			var action ActionCountResponse
			actionRows.Scan(&action.Action, &action.Count)
			stats.TopActions = append(stats.TopActions, action)
		}
	}

	// Get top users
	userRows, err := r.db.Query(`
		SELECT 
			al.user_id, u.user_identity, u.name as user_name, COUNT(*) as count 
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		WHERE al.user_id IS NOT NULL
		GROUP BY al.user_id, u.user_identity, u.name 
		ORDER BY count DESC 
		LIMIT 10
	`)
	if err == nil {
		defer userRows.Close()
		stats.TopUsers = []UserActivityCountResponse{}
		for userRows.Next() {
			var user UserActivityCountResponse
			userRows.Scan(&user.UserID, &user.UserIdentity, &user.UserName, &user.Count)
			stats.TopUsers = append(stats.TopUsers, user)
		}
	}

	// Get hourly activity for today
	hourRows, err := r.db.Query(`
		SELECT 
			EXTRACT(HOUR FROM created_at) as hour, COUNT(*) as count 
		FROM audit_logs 
		WHERE DATE(created_at) = CURRENT_DATE
		GROUP BY EXTRACT(HOUR FROM created_at) 
		ORDER BY hour
	`)
	if err == nil {
		defer hourRows.Close()
		stats.ActivityByHour = []HourlyActivityResponse{}
		for hourRows.Next() {
			var activity HourlyActivityResponse
			hourRows.Scan(&activity.Hour, &activity.Count)
			stats.ActivityByHour = append(stats.ActivityByHour, activity)
		}
	}

	// Get status breakdown
	statusRows, err := r.db.Query(`
		SELECT 
			CASE WHEN success = true THEN 'success' ELSE 'error' END as status, 
			COUNT(*) as count 
		FROM audit_logs 
		GROUP BY success 
		ORDER BY count DESC
	`)
	if err == nil {
		defer statusRows.Close()
		stats.StatusBreakdown = []StatusCountResponse{}
		for statusRows.Next() {
			var status StatusCountResponse
			statusRows.Scan(&status.Status, &status.Count)
			stats.StatusBreakdown = append(stats.StatusBreakdown, status)
		}
	}

	return stats, nil
}
