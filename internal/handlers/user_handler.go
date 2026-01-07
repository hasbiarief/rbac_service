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
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert to the expected struct
	req, ok := validatedBody.(*struct {
		Name         string  `json:"name" validate:"required,min=2,max=100"`
		Email        string  `json:"email" validate:"required,email,max=255"`
		UserIdentity *string `json:"user_identity" validate:"omitempty,min=3,max=50"`
		Password     string  `json:"password" validate:"omitempty,min=6,max=100"`
	})
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	// Convert to service request
	createReq := &service.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if req.UserIdentity != nil {
		createReq.UserIdentity = *req.UserIdentity
	}

	result, err := h.userService.CreateUser(createReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create user", err.Error())
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

	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	// Type assert to the expected struct
	req, ok := validatedBody.(*struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		IsActive *bool  `json:"is_active"`
	})
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	// Convert to service request
	updateReq := &service.UpdateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		IsActive: req.IsActive,
	}

	result, err := h.userService.UpdateUser(id, updateReq)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
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
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
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

	// Check if grouped format is requested
	grouped := c.Query("grouped")
	if grouped == "true" {
		result, err := h.userRepo.GetUserModulesGroupedWithSubscription(id)
		if err != nil {
			response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
			return
		}
		response.Success(c, http.StatusOK, "Success", result)
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

	// Check if grouped format is requested
	grouped := c.Query("grouped")
	if grouped == "true" {
		result, err := h.userRepo.GetUserModulesGroupedWithSubscription(user.ID)
		if err != nil {
			response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
			return
		}
		response.Success(c, http.StatusOK, "Success", result)
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

	response.Success(c, http.StatusOK, "Success", result)
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

	// Type assert to the expected struct
	req, ok := validatedBody.(*struct {
		UserIdentity string `json:"user_identity" validate:"required"`
		ModuleURL    string `json:"module_url" validate:"required"`
	})
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

	// Type assert to the expected struct
	req, ok := validatedBody.(*struct {
		CurrentPassword string `json:"current_password" validate:"required,min=6"`
		NewPassword     string `json:"new_password" validate:"required,min=6,max=100"`
		ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
	})
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	// Convert to service request
	changeReq := &service.ChangePasswordRequest{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	}

	if err := h.userService.ChangeUserPassword(id, changeReq); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Password changed successfully", nil)
}

// func (h *UserHandler) ChangeMyPassword(c *gin.Context) {
// 	// For now, we'll use a placeholder user ID
// 	// In a real implementation, you'd get this from the JWT token
// 	userID := int64(1)

// 	// Get validated body from context (set by validation middleware)
// 	validatedBody, exists := c.Get("validated_body")
// 	if !exists {
// 		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
// 		return
// 	}

// 	// Type assert to the expected struct
// 	req, ok := validatedBody.(*struct {
// 		CurrentPassword string `json:"current_password" validate:"required,min=6"`
// 		NewPassword     string `json:"new_password" validate:"required,min=6,max=100"`
// 		ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
// 	})
// 	if !ok {
// 		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
// 		return
// 	}

// 	// Convert to service request
// 	changeReq := &service.ChangePasswordRequest{
// 		CurrentPassword: req.CurrentPassword,
// 		NewPassword:     req.NewPassword,
// 		ConfirmPassword: req.ConfirmPassword,
// 	}

// 	if err := h.userService.ChangePassword(userID, changeReq); err != nil {
// 		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
// 		return
// 	}

// 	response.Success(c, http.StatusOK, "Password changed successfully", nil)
// }
