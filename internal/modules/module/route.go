package module

import (
	"gin-scalable-api/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, handler *Handler) {
	modules := api.Group("/modules")
	{
		modules.GET("", handler.GetModules)
		modules.GET("/:id", handler.GetModuleByID)
		modules.POST("",
			middleware.ValidateJSON(&CreateModuleRequest{}),
			handler.CreateModule,
		)
		modules.PUT("/:id",
			middleware.ValidateJSON(&UpdateModuleRequest{}),
			handler.UpdateModule,
		)
		modules.DELETE("/:id", handler.DeleteModule)
		modules.GET("/tree", handler.GetModuleTree)
		modules.GET("/:id/children", handler.GetModuleChildren)
		modules.GET("/:id/ancestors", handler.GetModuleAncestors)
	}
}
