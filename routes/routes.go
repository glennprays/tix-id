package routes

import (
	"net/http"
	"tix-id/controller"

	"github.com/gin-gonic/gin" // swagger embed files
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to TIX-ID",
		})
	})
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			customer := v1.Group("/customer")
			{
				customer.POST("/registration", controller.AddCustomer)
				customer.POST("/auth/login", controller.LoginCustomer)
			}

			admin := v1.Group("/admin")
			{
				admin.POST("/auth/login", controller.LoginAdmin)
			}
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
