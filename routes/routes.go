package routes

import (
	"net/http"
	"tix-id/controller"
	"tix-id/middleware"

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
			v1.DELETE("/auth/logout", middleware.AuthMiddleware("admin", "customer"), controller.LogoutAccount)
			customer := v1.Group("/customer")

			{
				customer.POST("/registration", controller.AddCustomer)
				customer.POST("/auth/login", controller.LoginCustomer)
				customerId := customer.Group("/:customerId")
				{
					customerId.POST("/tickets", controller.CreateTicket)
					customerId.GET("/tickets", controller.GetTickets)
					customerId.GET("/tickets/:ticketId", controller.GetTicket)
					customerId.POST("/tickets/:ticketId/payment", controller.ConfirmPayment)
				}
			}

			admin := v1.Group("/admin")
			{
				admin.POST("/auth/login", controller.LoginAdmin)
			}

			movie := v1.Group("/movies")
			{
				movieId := movie.Group("/:movieId")
				{
					movieId.GET("/schedules", controller.GetSchedules)
				}
			}
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
