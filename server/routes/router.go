package routes

import (
	"github.com/Sugarwithfish/Software-System-Design/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/", func(c *gin.Context) {
			c.JSON(200, "hello world")
		})

		test := apiV1.Group("/test")
		{
			test.POST("/", controllers.TestControllers)
		}

		questions := apiV1.Group("/questions")
		{
			questions.POST("/ai-generate", controllers.AIGenQuestionHandler()) // AI生成题目
			questions.GET("", controllers.ListQuestionHandler())                        // 获取题目列表
			questions.DELETE("", controllers.DeleteQuestionHandler())                   // 删除题目 (批量)
			questions.GET("/:id", controllers.GetQuestionHandler())                     // 获取题目详情
			questions.POST("", controllers.SaveQuestionHandler())                       // 创建题目
			questions.PUT("/:id", controllers.EditQuestionHandler())                    // 编辑题目
			questions.POST("/batch", controllers.BatchSaveQuestionHandler())            // 批量创建题目
		}

	}

	return r
}
