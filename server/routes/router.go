package routes

import (
	"github.com/Sugarwithfish/Software-System-Design/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/", func(c *gin.Context) {
			c.JSON(200, "hello world")
		})

		test := apiV1.Group("/test")
		{
			test.POST("/", controllers.TestControllers)
		}
	}

	return r
}
