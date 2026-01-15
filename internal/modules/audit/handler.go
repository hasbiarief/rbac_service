package audit

import (
	"net/http"
	"strconv"

	"gin-scalable-api/internal/constants"
	"gin-scalable-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

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
