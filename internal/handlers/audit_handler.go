package handlers

import (
	"net/http"
	"strconv"

	"gin-scalable-api/internal/constants"
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/interfaces"
	"gin-scalable-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	auditService interfaces.AuditServiceInterface
}

func NewAuditHandler(auditService interfaces.AuditServiceInterface) *AuditHandler {
	return &AuditHandler{
		auditService: auditService,
	}
}

func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	var req dto.AuditListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request parameters", err.Error())
		return
	}

	auditResponse, err := h.auditService.GetAuditLogs(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get audit logs", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgAuditLogsRetrieved, auditResponse)
}

func (h *AuditHandler) CreateAuditLog(c *gin.Context) {
	var req dto.CreateAuditLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	auditResponse, err := h.auditService.CreateAuditLog(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create audit log", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgAuditLogCreated, auditResponse)
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
		response.ErrorWithAutoStatus(c, "Failed to get user audit logs", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgAuditLogsRetrieved, auditResponse)
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
		response.ErrorWithAutoStatus(c, "Failed to get user audit logs", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgAuditLogsRetrieved, auditResponse)
}

func (h *AuditHandler) GetAuditStats(c *gin.Context) {
	statsResponse, err := h.auditService.GetAuditStats()
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get audit statistics", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgAuditStatsRetrieved, statsResponse)
}
