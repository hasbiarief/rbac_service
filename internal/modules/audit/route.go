package audit

import (
	"gin-scalable-api/internal/constants"
	"gin-scalable-api/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler struct
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Handler methods
func (h *Handler) GetAuditLogs(c *gin.Context) {
	var req AuditListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request parameters", err.Error())
		return
	}

	auditResponse, err := h.service.GetAuditLogs(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get audit logs", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgAuditLogsRetrieved, auditResponse)
}

func (h *Handler) CreateAuditLog(c *gin.Context) {
	var req CreateAuditLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	auditResponse, err := h.service.CreateAuditLog(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create audit log", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgAuditLogCreated, auditResponse)
}

func (h *Handler) GetUserAuditLogs(c *gin.Context) {
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

	auditResponse, err := h.service.GetUserAuditLogs(userID, limit)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get user audit logs", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgAuditLogsRetrieved, auditResponse)
}

func (h *Handler) GetUserAuditLogsByIdentity(c *gin.Context) {
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

	auditResponse, err := h.service.GetUserAuditLogsByIdentity(identity, limit)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get user audit logs", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgAuditLogsRetrieved, auditResponse)
}

func (h *Handler) GetAuditStats(c *gin.Context) {
	statsResponse, err := h.service.GetAuditStats()
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get audit statistics", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgAuditStatsRetrieved, statsResponse)
}

// Route registration
func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	audit := router.Group("/audit")
	{
		// GET /api/v1/audit/logs - Get audit logs
		audit.GET("/logs", handler.GetAuditLogs)

		// POST /api/v1/audit/logs - Create audit log
		audit.POST("/logs", handler.CreateAuditLog)

		// GET /api/v1/audit/stats - Get audit statistics
		audit.GET("/stats", handler.GetAuditStats)

		// GET /api/v1/audit/users/:userId/logs - Get user audit logs by ID
		audit.GET("/users/:userId/logs", handler.GetUserAuditLogs)

		// GET /api/v1/audit/users/identity/:identity/logs - Get user audit logs by identity
		audit.GET("/users/identity/:identity/logs", handler.GetUserAuditLogsByIdentity)
	}
}
