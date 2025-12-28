package validation

import "gin-scalable-api/middleware"

// Helper function to create int pointer
func IntPtr(i int) *int {
	return &i
}

// Common validation rules
var IDValidation = middleware.ValidationRules{
	Params: []middleware.ParamValidation{
		{Name: "id", Type: "int", Required: true, Min: IntPtr(1)},
	},
}

var UserIDValidation = middleware.ValidationRules{
	Params: []middleware.ParamValidation{
		{Name: "userId", Type: "int", Required: true, Min: IntPtr(1)},
	},
}

var RoleIDValidation = middleware.ValidationRules{
	Params: []middleware.ParamValidation{
		{Name: "roleId", Type: "int", Required: true, Min: IntPtr(1)},
	},
}

var CompanyIDValidation = middleware.ValidationRules{
	Params: []middleware.ParamValidation{
		{Name: "companyId", Type: "int", Required: true, Min: IntPtr(1)},
	},
}

var IdentityValidation = middleware.ValidationRules{
	Params: []middleware.ParamValidation{
		{Name: "identity", Type: "string", Required: true, Min: IntPtr(3), Max: IntPtr(50)},
	},
}

var ListValidation = middleware.ValidationRules{
	Query: []middleware.QueryValidation{
		{Name: "page", Type: "int", Default: 1, Min: IntPtr(1)},
		{Name: "limit", Type: "int", Default: 10, Min: IntPtr(1), Max: IntPtr(100)},
		{Name: "search", Type: "string"},
	},
}
