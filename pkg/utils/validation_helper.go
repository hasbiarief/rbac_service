package utils

import (
	"fmt"
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/pkg/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ValidationHelper provides utility methods for validation
type ValidationHelper struct{}

// NewValidationHelper creates a new validation helper
func NewValidationHelper() *ValidationHelper {
	return &ValidationHelper{}
}

// GetValidatedBody extracts and type-asserts validated body from context
func (h *ValidationHelper) GetValidatedBody(c *gin.Context, target interface{}) error {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		return errors.NewValidationError("validation failed")
	}

	// Use type assertion based on target type
	switch v := target.(type) {
	case **dto.CreateUserRequest:
		if body, ok := validatedBody.(*dto.CreateUserRequest); ok {
			*v = body
			return nil
		}
	case **dto.UpdateUserRequest:
		if body, ok := validatedBody.(*dto.UpdateUserRequest); ok {
			*v = body
			return nil
		}
	case **dto.CreateCompanyRequest:
		if body, ok := validatedBody.(*dto.CreateCompanyRequest); ok {
			*v = body
			return nil
		}
	case **dto.UpdateCompanyRequest:
		if body, ok := validatedBody.(*dto.UpdateCompanyRequest); ok {
			*v = body
			return nil
		}
	case **dto.CreateBranchRequest:
		if body, ok := validatedBody.(*dto.CreateBranchRequest); ok {
			*v = body
			return nil
		}
	case **dto.UpdateBranchRequest:
		if body, ok := validatedBody.(*dto.UpdateBranchRequest); ok {
			*v = body
			return nil
		}
	case **dto.CreateRoleRequest:
		if body, ok := validatedBody.(*dto.CreateRoleRequest); ok {
			*v = body
			return nil
		}
	case **dto.UpdateRoleRequest:
		if body, ok := validatedBody.(*dto.UpdateRoleRequest); ok {
			*v = body
			return nil
		}
	case **dto.LoginRequest:
		if body, ok := validatedBody.(*dto.LoginRequest); ok {
			*v = body
			return nil
		}
	case **dto.RegisterRequest:
		if body, ok := validatedBody.(*dto.RegisterRequest); ok {
			*v = body
			return nil
		}
	case **dto.ChangePasswordRequest:
		if body, ok := validatedBody.(*dto.ChangePasswordRequest); ok {
			*v = body
			return nil
		}
	default:
		return errors.NewValidationError("unsupported body type")
	}

	return errors.NewValidationError("invalid body structure")
}

// GetIDParam extracts and validates ID parameter from URL
func (h *ValidationHelper) GetIDParam(c *gin.Context, paramName string) (int64, error) {
	idStr := c.Param(paramName)
	if idStr == "" {
		return 0, errors.NewBadRequestError("Missing parameter", fmt.Sprintf("%s parameter is required", paramName))
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.NewBadRequestError("Invalid parameter", fmt.Sprintf("%s must be a valid number", paramName))
	}

	if id <= 0 {
		return 0, errors.NewBadRequestError("Invalid parameter", fmt.Sprintf("%s must be greater than 0", paramName))
	}

	return id, nil
}

// GetQueryParam extracts query parameter with default value
func (h *ValidationHelper) GetQueryParam(c *gin.Context, key, defaultValue string) string {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetQueryParamInt extracts integer query parameter with default value
func (h *ValidationHelper) GetQueryParamInt(c *gin.Context, key string, defaultValue int) (int, error) {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, errors.NewBadRequestError("Invalid query parameter", fmt.Sprintf("%s must be a valid integer", key))
	}

	return value, nil
}

// GetQueryParamBool extracts boolean query parameter
func (h *ValidationHelper) GetQueryParamBool(c *gin.Context, key string) (*bool, error) {
	valueStr := c.Query(key)
	if valueStr == "" {
		return nil, nil
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return nil, errors.NewBadRequestError("Invalid query parameter", fmt.Sprintf("%s must be a valid boolean", key))
	}

	return &value, nil
}

// GetQueryParamInt64 extracts int64 query parameter
func (h *ValidationHelper) GetQueryParamInt64(c *gin.Context, key string) (*int64, error) {
	valueStr := c.Query(key)
	if valueStr == "" {
		return nil, nil
	}

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return nil, errors.NewBadRequestError("Invalid query parameter", fmt.Sprintf("%s must be a valid integer", key))
	}

	return &value, nil
}
