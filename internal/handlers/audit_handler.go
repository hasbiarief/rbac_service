package handlers

import (
	"net/http"
	"strconv"

	"gin-scalable-api/internal/service"
	"gin-scalable-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	auditService *service.AuditService
}

func NewAuditHandler(auditService *service.AuditService) *AuditHandler {
	return &AuditHandler{
		auditService: auditService,
	}
}

func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	var req service.AuditLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request parameters", err.Error())
		return
	}

	auditResponse, err := h.auditService.GetAuditLogs(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get audit logs", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Audit logs retrieved successfully", auditResponse)
}

func (h *AuditHandler) CreateAuditLog(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	// Type assert to the expected struct
	req, ok := validatedBody.(*struct {
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
	})
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "invalid body structure")
		return
	}

	// Convert to service request
	createReq := &service.CreateAuditLogRequest{
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

	auditResponse, err := h.auditService.CreateAuditLog(createReq)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create audit log", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Audit log created successfully", auditResponse)
}

func (h *AuditHandler) GetUserAuditLogs(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID", "User ID must be a valid number")
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	auditResponse, err := h.auditService.GetUserAuditLogs(userID, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get user audit logs", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User audit logs retrieved successfully", auditResponse)
}

func (h *AuditHandler) GetUserAuditLogsByIdentity(c *gin.Context) {
	identity := c.Param("identity")
	if identity == "" {
		response.Error(c, http.StatusBadRequest, "Invalid user identity", "User identity is required")
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	auditResponse, err := h.auditService.GetUserAuditLogsByIdentity(identity, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get user audit logs", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User audit logs retrieved successfully", auditResponse)
}

func (h *AuditHandler) GetAuditStats(c *gin.Context) {
	statsResponse, err := h.auditService.GetAuditStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get audit statistics", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Audit statistics retrieved successfully", statsResponse)
}
