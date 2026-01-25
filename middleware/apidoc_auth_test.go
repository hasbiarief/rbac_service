package middleware

import (
	"gin-scalable-api/internal/constants"
	"gin-scalable-api/pkg/rbac"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRBACService for testing middleware
type MockRBACService struct {
	mock.Mock
}

func (m *MockRBACService) HasPermission(userID int64, moduleID int64, permission string) (bool, error) {
	args := m.Called(userID, moduleID, permission)
	return args.Bool(0), args.Error(1)
}

func (m *MockRBACService) GetUserPermissions(userID int64) (*rbac.UserPermissions, error) {
	args := m.Called(userID)
	return args.Get(0).(*rbac.UserPermissions), args.Error(1)
}

func (m *MockRBACService) HasRole(userID int64, roleName string) (bool, error) {
	args := m.Called(userID, roleName)
	return args.Bool(0), args.Error(1)
}

func (m *MockRBACService) GetAccessibleModules(userID int64, permission string) ([]int64, error) {
	args := m.Called(userID, permission)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *MockRBACService) IsSuperAdmin(userID int64) (bool, error) {
	args := m.Called(userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRBACService) GetFilteredModules(userID int64, permission string, limit, offset int, category, search string, isActive *bool) ([]*rbac.ModuleInfo, error) {
	args := m.Called(userID, permission, limit, offset, category, search, isActive)
	return args.Get(0).([]*rbac.ModuleInfo), args.Error(1)
}

func TestAPIDocAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Sets RBAC service and user ID in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set user ID in context (normally set by AuthMiddleware)
		c.Set("user_id", int64(1))

		// Create middleware with nil DB (we're just testing context setting)
		middleware := APIDocAuthMiddleware(nil)

		// Execute middleware
		middleware(c)

		// Check that context was set correctly
		rbacService, exists := c.Get("rbac_service")
		assert.True(t, exists)
		assert.NotNil(t, rbacService)

		userID, exists := c.Get("user_id_int64")
		assert.True(t, exists)
		assert.Equal(t, int64(1), userID)
	})

	t.Run("Error - User not authenticated", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Don't set user ID in context
		middleware := APIDocAuthMiddleware(nil)

		// Execute middleware
		middleware(c)

		// Check that request was aborted
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Error - Invalid user ID type", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set invalid user ID type
		c.Set("user_id", "invalid")

		middleware := APIDocAuthMiddleware(nil)

		// Execute middleware
		middleware(c)

		// Check that request was aborted
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestRequireAPIDocPermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - User has permission", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set up context
		c.Set("user_id_int64", int64(1))
		mockRBAC := new(MockRBACService)
		c.Set("rbac_service", mockRBAC)

		// Mock permission check
		mockRBAC.On("HasPermission", int64(1), int64(140), "read").Return(true, nil)

		middleware := RequireAPIDocPermission(constants.ModuleAPIDocCollections, "read")

		// Execute middleware
		middleware(c)

		// Check that request was not aborted
		assert.False(t, c.IsAborted())
		mockRBAC.AssertExpectations(t)
	})

	t.Run("Error - User lacks permission", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set up context
		c.Set("user_id_int64", int64(1))
		mockRBAC := new(MockRBACService)
		c.Set("rbac_service", mockRBAC)

		// Mock permission check - user doesn't have permission
		mockRBAC.On("HasPermission", int64(1), int64(140), "write").Return(false, nil)

		middleware := RequireAPIDocPermission(constants.ModuleAPIDocCollections, "write")

		// Execute middleware
		middleware(c)

		// Check that request was aborted
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusForbidden, w.Code)
		mockRBAC.AssertExpectations(t)
	})

	t.Run("Error - User not authenticated", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Don't set user ID in context
		middleware := RequireAPIDocPermission(constants.ModuleAPIDocCollections, "read")

		// Execute middleware
		middleware(c)

		// Check that request was aborted
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Error - RBAC service not available", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set user ID but not RBAC service
		c.Set("user_id_int64", int64(1))

		middleware := RequireAPIDocPermission(constants.ModuleAPIDocCollections, "read")

		// Execute middleware
		middleware(c)

		// Check that request was aborted
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestValidateCollectionID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Valid collection ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set collection ID parameter
		c.Params = gin.Params{
			{Key: "collection_id", Value: "123"},
		}

		middleware := ValidateCollectionID()

		// Execute middleware
		middleware(c)

		// Check that request was not aborted and collection ID was set
		assert.False(t, c.IsAborted())
		collectionID, exists := c.Get("collection_id")
		assert.True(t, exists)
		assert.Equal(t, int64(123), collectionID)
	})

	t.Run("Success - Valid ID parameter fallback", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set generic ID parameter
		c.Params = gin.Params{
			{Key: "id", Value: "456"},
		}

		middleware := ValidateCollectionID()

		// Execute middleware
		middleware(c)

		// Check that request was not aborted and collection ID was set
		assert.False(t, c.IsAborted())
		collectionID, exists := c.Get("collection_id")
		assert.True(t, exists)
		assert.Equal(t, int64(456), collectionID)
	})

	t.Run("Error - Missing collection ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Don't set any ID parameter
		middleware := ValidateCollectionID()

		// Execute middleware
		middleware(c)

		// Check that request was aborted
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error - Invalid collection ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set invalid collection ID parameter
		c.Params = gin.Params{
			{Key: "collection_id", Value: "invalid"},
		}

		middleware := ValidateCollectionID()

		// Execute middleware
		middleware(c)

		// Check that request was aborted
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error - Zero collection ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set zero collection ID parameter
		c.Params = gin.Params{
			{Key: "collection_id", Value: "0"},
		}

		middleware := ValidateCollectionID()

		// Execute middleware
		middleware(c)

		// Check that request was aborted
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error - Negative collection ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set negative collection ID parameter
		c.Params = gin.Params{
			{Key: "collection_id", Value: "-1"},
		}

		middleware := ValidateCollectionID()

		// Execute middleware
		middleware(c)

		// Check that request was aborted
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAPIDocPermissionHelpers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("APIDocCollectionPermission", func(t *testing.T) {
		middleware := APIDocCollectionPermission("read")
		assert.NotNil(t, middleware)
	})

	t.Run("APIDocEndpointPermission", func(t *testing.T) {
		middleware := APIDocEndpointPermission("write")
		assert.NotNil(t, middleware)
	})

	t.Run("APIDocEnvironmentPermission", func(t *testing.T) {
		middleware := APIDocEnvironmentPermission("delete")
		assert.NotNil(t, middleware)
	})

	t.Run("APIDocExportPermission", func(t *testing.T) {
		middleware := APIDocExportPermission("read")
		assert.NotNil(t, middleware)
	})
}

func TestRequireResourceOwnership(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Sets current user ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set user ID in context
		c.Set("user_id_int64", int64(1))

		middleware := RequireResourceOwnership()

		// Execute middleware
		middleware(c)

		// Check that request was not aborted and current user ID was set
		assert.False(t, c.IsAborted())
		currentUserID, exists := c.Get("current_user_id")
		assert.True(t, exists)
		assert.Equal(t, int64(1), currentUserID)
	})

	t.Run("Error - User not authenticated", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Don't set user ID in context
		middleware := RequireResourceOwnership()

		// Execute middleware
		middleware(c)

		// Check that request was aborted
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
