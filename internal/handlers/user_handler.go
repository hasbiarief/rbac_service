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

	// Convert to service request
	checkReq := &service.CheckAccessRequest{
		UserIdentity: req.UserIdentity,
		ModuleURL:    req.ModuleURL,
	}

	hasAccess, err := h.moduleService.CheckUserAccess(checkReq)
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

func (h *UserHandler) ChangeMyPassword(c *gin.Context) {
	// For now, we'll use a placeholder user ID
	// In a real implementation, you'd get this from the JWT token
	userID := int64(1)

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

	if err := h.userService.ChangePassword(userID, changeReq); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Password changed successfully", nil)
}
