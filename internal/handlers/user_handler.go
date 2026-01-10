package handlers

import (
	"gin-scalable-api/internal/constants"
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/interfaces"
	"strconv"
	"strings"

	"gin-scalable-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService   interfaces.UserServiceInterface
	moduleService interfaces.ModuleServiceInterface
	userRepo      interfaces.UserRepositoryInterface
}

func NewUserHandler(userService interfaces.UserServiceInterface, moduleService interfaces.ModuleServiceInterface, userRepo interfaces.UserRepositoryInterface) *UserHandler {
	return &UserHandler{
		userService:   userService,
		moduleService: moduleService,
		userRepo:      userRepo,
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	var req dto.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Internal error", "Invalid user ID type")
		return
	}

	// Use filtered method to get users based on requesting user's permissions
	result, err := h.userService.GetUsersFiltered(userIDInt64, &req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get users", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgUsersRetrieved, result)
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
		response.Error(c, http.StatusNotFound, constants.MsgUserNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgUserRetrieved, result)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert to the expected DTO struct
	req, ok := validatedBody.(*dto.CreateUserRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.userService.CreateUser(req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create user", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgUserCreated, result)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert to the expected DTO struct
	req, ok := validatedBody.(*dto.UpdateUserRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.userService.UpdateUser(id, req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgUserUpdated, result)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	if err := h.userService.DeleteUser(id); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgUserDeleted, nil)
}

// User Module Access Methods
func (h *UserHandler) GetUserModules(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	// Check if grouped format is requested
	grouped := c.Query("grouped")
	if grouped == "true" {
		result, err := h.userRepo.GetUserModulesGroupedWithSubscription(id)
		if err != nil {
			response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
			return
		}
		response.Success(c, http.StatusOK, constants.MsgUserModulesRetrieved, result)
		return
	}

	// Default: return old format for backward compatibility
	category := c.Query("category")
	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	result, err := h.moduleService.GetUserModules(id, category, limit)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgUserModulesRetrieved, result)
}

func (h *UserHandler) GetUserModulesByIdentity(c *gin.Context) {
	identity := c.Param("identity")

	// Get user by identity first
	user, err := h.userRepo.GetByUserIdentity(identity)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgUserNotFound, "User not found")
		return
	}

	// Check if grouped format is requested
	grouped := c.Query("grouped")
	if grouped == "true" {
		result, err := h.userRepo.GetUserModulesGroupedWithSubscription(user.ID)
		if err != nil {
			response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
			return
		}
		response.Success(c, http.StatusOK, constants.MsgUserModulesRetrieved, result)
		return
	}

	// Default: return old format for backward compatibility
	category := c.Query("category")
	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	result, err := h.moduleService.GetUserModules(user.ID, category, limit)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgUserModulesRetrieved, result)
}

func (h *UserHandler) CheckAccess(c *gin.Context) {
	// Get user ID from context (set by auth middleware) - more efficient than using user_identity
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Internal error", "Invalid user ID type")
		return
	}

	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert to the expected DTO struct
	req, ok := validatedBody.(*dto.AccessCheckRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	// Use more efficient check with user ID and RBAC service
	hasAccess, err := h.checkUserModuleAccess(userIDInt64, req.ModuleURL)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	c.JSON(200, gin.H{
		"success":    true,
		"has_access": hasAccess,
	})
}

// checkUserModuleAccess checks if user has access to module using RBAC service (more efficient)
func (h *UserHandler) checkUserModuleAccess(userID int64, moduleURL string) (bool, error) {
	// Get user modules with subscription filtering
	modules, err := h.userRepo.GetUserModulesWithSubscription(userID)
	if err != nil {
		return false, err
	}

	// Check if moduleURL exists in user's accessible modules
	for _, module := range modules {
		if module == moduleURL {
			return true, nil
		}
	}

	return false, nil
}

// Password Management Methods
func (h *UserHandler) ChangeUserPassword(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert to the expected DTO struct
	req, ok := validatedBody.(*dto.ChangePasswordRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	if err := h.userService.ChangeUserPassword(id, req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgPasswordChanged, nil)
}
