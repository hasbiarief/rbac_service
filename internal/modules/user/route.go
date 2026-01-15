package user

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, handler *Handler) {
	users := api.Group("/users")
	{
		users.GET("", handler.GetUsers)
		users.GET("/:id", handler.GetUserByID)
		users.POST("", handler.CreateUser)
		users.PUT("/:id", handler.UpdateUser)
		users.DELETE("/:id", handler.DeleteUser)
		users.GET("/:id/modules", handler.GetUserModules)
		users.GET("/identity/:identity/modules", handler.GetUserModulesByIdentity)
		users.POST("/check-access", handler.CheckAccess)
		users.PUT("/:id/password", handler.ChangeUserPassword)
	}
}
