package user

import (
	"gin-scalable-api/internal/constants"
	"gin-scalable-api/middleware"
	"gin-scalable-api/pkg/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler struct
type Handler struct {
	service  *Service
	userRepo *UserRepository
}

func NewHandler(service *Service, userRepo *UserRepository) *Handler {
	return &Handler{
		service:  service,
		userRepo: userRepo,
	}
}

// Handler methods
func (h *Handler) GetUsers(c *gin.Context) {
	var req UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := h.service.GetUsersWithRoles(&req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get users", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgUsersRetrieved, result)
}

func (h *Handler) GetUserByID(c *gin.Context) {
	if strings.Contains(c.Request.URL.Path, "/modules") {
		c.Next()
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	result, err := h.service.GetUserByIDWithRoles(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgUserNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgUserRetrieved, result)
}

func (h *Handler) CreateUser(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*CreateUserRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.CreateUser(req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to create user", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.MsgUserCreated, result)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*UpdateUserRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	result, err := h.service.UpdateUser(id, req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgUserUpdated, result)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgUserDeleted, nil)
}

func (h *Handler) GetUserModules(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

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

	result, err := h.userRepo.GetUserModulesWithSubscription(id)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgUserModulesRetrieved, result)
}

func (h *Handler) GetUserModulesByIdentity(c *gin.Context) {
	identity := c.Param("identity")

	user, err := h.userRepo.GetByUserIdentity(identity)
	if err != nil {
		response.Error(c, http.StatusNotFound, constants.MsgUserNotFound, "User not found")
		return
	}

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

	result, err := h.userRepo.GetUserModulesWithSubscription(user.ID)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgUserModulesRetrieved, result)
}

func (h *Handler) CheckAccess(c *gin.Context) {
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

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*AccessCheckRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

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

func (h *Handler) checkUserModuleAccess(userID int64, moduleURL string) (bool, error) {
	modules, err := h.userRepo.GetUserModulesWithSubscription(userID)
	if err != nil {
		return false, err
	}

	for _, module := range modules {
		if module == moduleURL {
			return true, nil
		}
	}

	return false, nil
}

func (h *Handler) ChangeUserPassword(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "Invalid user ID")
		return
	}

	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
		return
	}

	req, ok := validatedBody.(*ChangePasswordRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
		return
	}

	if err := h.service.ChangeUserPassword(id, req); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgPasswordChanged, nil)
}

// Route registration
func RegisterRoutes(api *gin.RouterGroup, handler *Handler) {
	users := api.Group("/users")
	{
		// GET /api/v1/users - Get all users with optional filters and pagination
		users.GET("", handler.GetUsers)

		// GET /api/v1/users/:id - Get user by ID with roles and permissions
		users.GET("/:id", handler.GetUserByID)

		// POST /api/v1/users - Create new user with role assignments
		users.POST("",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &CreateUserRequest{},
			}),
			handler.CreateUser,
		)

		// PUT /api/v1/users/:id - Update user information and role assignments
		users.PUT("/:id",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &UpdateUserRequest{},
			}),
			handler.UpdateUser,
		)

		// DELETE /api/v1/users/:id - Delete user and remove all associations
		users.DELETE("/:id", handler.DeleteUser)

		// GET /api/v1/users/:id/modules - Get user accessible modules with subscription check
		users.GET("/:id/modules", handler.GetUserModules)

		// GET /api/v1/users/identity/:identity/modules - Get user modules by identity with subscription check
		users.GET("/identity/:identity/modules", handler.GetUserModulesByIdentity)

		// POST /api/v1/users/check-access - Check user access to specific module URL
		users.POST("/check-access",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &AccessCheckRequest{},
			}),
			handler.CheckAccess,
		)

		// PUT /api/v1/users/:id/password - Change user password with validation
		users.PUT("/:id/password",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &ChangePasswordRequest{},
			}),
			handler.ChangeUserPassword,
		)
	}
}
