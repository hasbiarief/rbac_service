package branch

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, handler *Handler) {
	branches := api.Group("/branches")
	{
		branches.GET("", handler.GetBranches)
		branches.GET("/:id", handler.GetBranchByID)
		branches.GET("/:id/hierarchy", handler.GetBranchHierarchy)
		branches.GET("/:id/children", handler.GetBranchChildren)
		branches.POST("", handler.CreateBranch)
		branches.PUT("/:id", handler.UpdateBranch)
		branches.DELETE("/:id", handler.DeleteBranch)
		branches.GET("/company/:companyId", handler.GetCompanyBranches)
	}
}
