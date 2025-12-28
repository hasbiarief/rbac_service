package handlers

import (
	"gin-scalable-api/internal/repository"
	"gin-scalable-api/internal/service"
	"strconv"
	"strings"

	"gin-scalable-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService   *service.UserService
	moduleService *service.ModuleService
	userRepo      *repository.UserRepository
}

func NewUserHandler(userService *service.UserService, moduleService *service.ModuleService, userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		userService:   userService,
		moduleService: moduleService,
		userRepo:      userRepo,
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	var req service.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.userService.GetUsers(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Check if this is a module-related route
	if strings.Contains(c.Request.URL.Path, "/modules") {
		c.Next()
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	result, err := h.userService.GetUserByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.userService.CreateUser(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Created successfully", result)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.userService.UpdateUser(id, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	if err := h.userService.DeleteUser(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}

// User Module Access Methods
func (h *UserHandler) GetUserModules(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	category := c.Query("category")
	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	result, err := h.moduleService.GetUserModules(id, category, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *UserHandler) GetUserModulesByIdentity(c *gin.Context) {
	identity := c.Param("identity")

	// Get user by identity first
	user, err := h.userRepo.GetByUserIdentity(identity)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Not found", "User not found")
		return
	}

	category := c.Query("category")
	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	result, err := h.moduleService.GetUserModules(user.ID, category, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Success", result)
}

func (h *UserHandler) CheckAccess(c *gin.Context) {
	var req service.CheckAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	hasAccess, err := h.moduleService.CheckUserAccess(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	c.JSON(200, gin.H{
		"success":    true,
		"has_access": hasAccess,
	})
}

// Password Management Methods
func (h *UserHandler) ChangeUserPassword(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	if err := h.userService.ChangeUserPassword(id, &req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Password changed successfully", nil)
}

func (h *UserHandler) ChangeMyPassword(c *gin.Context) {
	// For now, we'll use a placeholder user ID
	// In a real implementation, you'd get this from the JWT token
	userID := int64(1)

	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	if err := h.userService.ChangePassword(userID, &req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Password changed successfully", nil)
}
