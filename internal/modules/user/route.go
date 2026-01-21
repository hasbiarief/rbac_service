package user

import (
	"gin-scalable-api/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, handler *Handler) {
	users := api.Group("/users")
	{
		users.GET("", handler.GetUsers)
		users.GET("/:id", handler.GetUserByID)
		users.POST("",
			middleware.ValidateJSON(&CreateUserRequest{}),
			handler.CreateUser,
		)
		users.PUT("/:id",
			middleware.ValidateJSON(&UpdateUserRequest{}),
			handler.UpdateUser,
		)
		users.DELETE("/:id", handler.DeleteUser)
		users.GET("/:id/modules", handler.GetUserModules)
		users.GET("/identity/:identity/modules", handler.GetUserModulesByIdentity)
		users.POST("/check-access",
			middleware.ValidateJSON(&AccessCheckRequest{}),
			handler.CheckAccess,
		)
		users.PUT("/:id/password",
			middleware.ValidateJSON(&ChangePasswordRequest{}),
			handler.ChangeUserPassword,
		)
	}
}
