package middleware

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"gin-scalable-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom tag name function to use json tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ParamValidation defines validation rules for route parameters
type ParamValidation struct {
	Name     string `json:"name"`
	Type     string `json:"type"` // "int", "string", "uuid"
	Required bool   `json:"required"`
	Min      *int   `json:"min,omitempty"`
	Max      *int   `json:"max,omitempty"`
	Pattern  string `json:"pattern,omitempty"`
}

// QueryValidation defines validation rules for query parameters
type QueryValidation struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"` // "int", "string", "bool"
	Required bool        `json:"required"`
	Default  interface{} `json:"default,omitempty"`
	Min      *int        `json:"min,omitempty"`
	Max      *int        `json:"max,omitempty"`
	Options  []string    `json:"options,omitempty"` // for enum validation
}

// ValidationRules defines all validation rules for a route
type ValidationRules struct {
	Params []ParamValidation `json:"params"`
	Query  []QueryValidation `json:"query"`
	Body   interface{}       `json:"body,omitempty"` // struct for body validation
}

// ValidateRequest creates a middleware that validates request parameters, query, and body
func ValidateRequest(rules ValidationRules) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate path parameters
		if err := validateParams(c, rules.Params); err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid path parameter", err.Error())
			c.Abort()
			return
		}

		// Validate query parameters
		if err := validateQuery(c, rules.Query); err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid query parameter", err.Error())
			c.Abort()
			return
		}

		// Validate request body if rules provided
		if rules.Body != nil {
			if err := validateBody(c, rules.Body); err != nil {
				response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func validateParams(c *gin.Context, paramRules []ParamValidation) error {
	for _, rule := range paramRules {
		value := c.Param(rule.Name)

		// Check if required parameter is missing
		if rule.Required && value == "" {
			return fmt.Errorf("parameter '%s' is required", rule.Name)
		}

		if value == "" {
			continue // Skip validation for optional empty parameters
		}

		// Validate based on type
		switch rule.Type {
		case "int":
			intVal, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("parameter '%s' must be a valid integer", rule.Name)
			}

			if rule.Min != nil && intVal < *rule.Min {
				return fmt.Errorf("parameter '%s' must be at least %d", rule.Name, *rule.Min)
			}

			if rule.Max != nil && intVal > *rule.Max {
				return fmt.Errorf("parameter '%s' must be at most %d", rule.Name, *rule.Max)
			}

		case "uuid":
			if !isValidUUID(value) {
				return fmt.Errorf("parameter '%s' must be a valid UUID", rule.Name)
			}

		case "string":
			if rule.Min != nil && len(value) < *rule.Min {
				return fmt.Errorf("parameter '%s' must be at least %d characters", rule.Name, *rule.Min)
			}

			if rule.Max != nil && len(value) > *rule.Max {
				return fmt.Errorf("parameter '%s' must be at most %d characters", rule.Name, *rule.Max)
			}
		}

		// Store validated parameter in context for easy access
		c.Set(fmt.Sprintf("param_%s", rule.Name), value)
	}

	return nil
}

func validateQuery(c *gin.Context, queryRules []QueryValidation) error {
	for _, rule := range queryRules {
		value := c.Query(rule.Name)

		// Set default value if parameter is missing
		if value == "" && rule.Default != nil {
			value = fmt.Sprintf("%v", rule.Default)
			c.Set(fmt.Sprintf("query_%s", rule.Name), rule.Default)
		}

		// Check if required parameter is missing
		if rule.Required && value == "" {
			return fmt.Errorf("query parameter '%s' is required", rule.Name)
		}

		if value == "" {
			continue // Skip validation for optional empty parameters
		}

		// Validate based on type
		switch rule.Type {
		case "int":
			intVal, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("query parameter '%s' must be a valid integer", rule.Name)
			}

			if rule.Min != nil && intVal < *rule.Min {
				return fmt.Errorf("query parameter '%s' must be at least %d", rule.Name, *rule.Min)
			}

			if rule.Max != nil && intVal > *rule.Max {
				return fmt.Errorf("query parameter '%s' must be at most %d", rule.Name, *rule.Max)
			}

			c.Set(fmt.Sprintf("query_%s", rule.Name), intVal)

		case "bool":
			boolVal, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("query parameter '%s' must be a valid boolean", rule.Name)
			}
			c.Set(fmt.Sprintf("query_%s", rule.Name), boolVal)

		case "string":
			if rule.Min != nil && len(value) < *rule.Min {
				return fmt.Errorf("query parameter '%s' must be at least %d characters", rule.Name, *rule.Min)
			}

			if rule.Max != nil && len(value) > *rule.Max {
				return fmt.Errorf("query parameter '%s' must be at most %d characters", rule.Name, *rule.Max)
			}

			// Check if value is in allowed options
			if len(rule.Options) > 0 {
				found := false
				for _, option := range rule.Options {
					if value == option {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("query parameter '%s' must be one of: %s", rule.Name, strings.Join(rule.Options, ", "))
				}
			}

			c.Set(fmt.Sprintf("query_%s", rule.Name), value)
		}
	}

	return nil
}

func validateBody(c *gin.Context, bodyStruct interface{}) error {
	// Create a new instance of the struct type
	bodyType := reflect.TypeOf(bodyStruct)
	if bodyType.Kind() == reflect.Ptr {
		bodyType = bodyType.Elem()
	}

	bodyValue := reflect.New(bodyType).Interface()

	// Bind JSON to struct
	if err := c.ShouldBindJSON(bodyValue); err != nil {
		return fmt.Errorf("invalid JSON format: %v", err)
	}

	// Validate struct using validator
	if err := validate.Struct(bodyValue); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, formatValidationError(err))
		}
		return fmt.Errorf("%s", strings.Join(validationErrors, "; "))
	}

	// Store validated body in context
	c.Set("validated_body", bodyValue)
	return nil
}

func formatValidationError(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()

	switch tag {
	case "required":
		return fmt.Sprintf("field '%s' is required", field)
	case "email":
		return fmt.Sprintf("field '%s' must be a valid email", field)
	case "min":
		return fmt.Sprintf("field '%s' must be at least %s characters/value", field, err.Param())
	case "max":
		return fmt.Sprintf("field '%s' must be at most %s characters/value", field, err.Param())
	case "len":
		return fmt.Sprintf("field '%s' must be exactly %s characters", field, err.Param())
	case "oneof":
		return fmt.Sprintf("field '%s' must be one of: %s", field, err.Param())
	default:
		return fmt.Sprintf("field '%s' failed validation for '%s'", field, tag)
	}
}

func isValidUUID(uuid string) bool {
	// Simple UUID validation (you can use a more robust library)
	if len(uuid) != 36 {
		return false
	}

	for i, char := range uuid {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if char != '-' {
				return false
			}
		} else {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
				return false
			}
		}
	}

	return true
}

// Helper functions to get validated values from context
func GetValidatedParam(c *gin.Context, name string) string {
	if value, exists := c.Get(fmt.Sprintf("param_%s", name)); exists {
		return value.(string)
	}
	return c.Param(name)
}

func GetValidatedQuery(c *gin.Context, name string) interface{} {
	if value, exists := c.Get(fmt.Sprintf("query_%s", name)); exists {
		return value
	}
	return c.Query(name)
}

func GetValidatedBody(c *gin.Context) interface{} {
	if value, exists := c.Get("validated_body"); exists {
		return value
	}
	return nil
}
