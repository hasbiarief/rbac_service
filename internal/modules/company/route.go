package company

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, handler *Handler) {
	companies := api.Group("/companies")
	{
		companies.GET("", handler.GetCompanies)
		companies.GET("/:id", handler.GetCompanyByID)
		companies.POST("", handler.CreateCompany)
		companies.PUT("/:id", handler.UpdateCompany)
		companies.DELETE("/:id", handler.DeleteCompany)
	}
}
