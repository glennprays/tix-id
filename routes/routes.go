package routes

import (
	"net/http"
	"tix-id/controller"
	"tix-id/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin" // swagger embed files
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func corsPolicy() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{c.Request.Host}
		config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
		config.AllowHeaders = []string{"Origin"}
		config.AllowCredentials = true
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Host)
		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin")
			c.AbortWithStatus(http.StatusOK)
			return
		}
	}
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(corsPolicy())

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
				customerId.Use(middleware.AuthMiddleware("customer"))

				{
					customerId.POST("/tickets", controller.CreateTicket)
					customerId.GET("/tickets", controller.GetTickets)
					customerId.GET("/tickets/:ticketId", controller.GetTicket)
					customerId.POST("/tickets/:ticketId/payment", controller.ConfirmPayment)
					customerId.GET("/profile", controller.GetCustomer)
					customerId.PUT("/profile", controller.UpdateCustomer)
				}
			}

			admin := v1.Group("/admin")
			{
				admin.POST("/auth/login", controller.LoginAdmin)
			}

			movie := v1.Group("/movies")
			{
				movie.GET("/", controller.GetMovies)
				movie.GET("/search", controller.SearchMovies)
				movie.POST("/", controller.CreateMovie)
				movieId := movie.Group("/:movieId")
				{
					movieId.POST("/schedules/", middleware.AuthMiddleware("admin"), controller.CreateMovieSchedule)
					movieId.GET("/schedules", controller.GetSchedules)
					movieId.GET("/schedules/:scheduleId", controller.GetSchedule)
					movieId.PUT("/", middleware.AuthMiddleware("admin"), controller.UpdateMovie)
					movieId.DELETE("/", middleware.AuthMiddleware("admin"), controller.DeleteMovie)
					movieId.PUT("/schedules/:scheduleId", middleware.AuthMiddleware("admin"), controller.UpdateMovieSchedule)
					movieId.DELETE("/schedules/:scheduleId", middleware.AuthMiddleware("admin"), controller.DeleteSchedule)
					movieId.POST("/schedules/:scheduleId/seats", middleware.AuthMiddleware("admin"), controller.AddScheduleSeats)
					movieId.GET("/", controller.GetMovieById)
				}
			}

			branches := v1.Group("/branches")
			branches.Use(middleware.AuthMiddleware("admin"))
			{
				branches.GET("/", controller.GetBranches)
				branches.POST("/", controller.CreateBranch)
				branchId := branches.Group("/:branchId")
				{
					branchId.GET("/branch", controller.GetBranch)
					branchId.PUT("/", controller.UpdateBranch)
					branchId.DELETE("/", controller.DeleteBranch)
					branchId.POST("/theatres", controller.CreateTheatre)
					branchId.PUT("/theatres/:theatreId", controller.UpdateTheatre)
					branchId.DELETE("/theatres/:theatreId", controller.DeleteTheatre)
				}
			}

		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
