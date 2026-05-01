package routes

import (
	"github.com/gin-gonic/gin"
	"recruit/controllers"
	"recruit/middleware"
	"recruit/models"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 公开路由
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// 任务公开路由（查看列表）
	tasksPublic := r.Group("/api/tasks")
	{
		tasksPublic.GET("", controllers.GetAllTasks)
		tasksPublic.GET("/:id", controllers.GetTaskByID)
	}

	// 需要认证的路由
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// 用户信息
		api.GET("/me", controllers.GetCurrentUser)

		// 任务相关
		tasks := api.Group("/tasks")
		{
			tasks.POST("", middleware.RoleMiddleware(string(models.RoleEmployer)), controllers.CreateTask)
			tasks.GET("/my", controllers.GetMyTasks)
			tasks.PUT("/:id", middleware.RoleMiddleware(string(models.RoleEmployer)), controllers.UpdateTask)
			tasks.DELETE("/:id", middleware.RoleMiddleware(string(models.RoleEmployer)), controllers.DeleteTask)
		}

		// 投标相关
		bids := api.Group("/bids")
		{
			bids.POST("/task/:task_id", middleware.RoleMiddleware(string(models.RoleEmployee)), controllers.CreateBid)
			bids.GET("/task/:task_id", controllers.GetBidsByTask)
			bids.GET("/my", controllers.GetMyBids)
			bids.POST("/:bid_id/accept", middleware.RoleMiddleware(string(models.RoleEmployer)), controllers.AcceptBid)
		}

		// 评价相关
		reviews := api.Group("/reviews")
		{
			reviews.POST("/task/:task_id", middleware.RoleMiddleware(string(models.RoleEmployer)), controllers.CreateReview)
			reviews.GET("/employee/:employee_id", controllers.GetReviewsByEmployee)
			reviews.GET("/my", controllers.GetMyReviews)
		}

		// 收藏相关
		favorites := api.Group("/favorites")
		{
			favorites.POST("/task/:task_id", middleware.RoleMiddleware(string(models.RoleEmployee)), controllers.ToggleFavorite)
			favorites.GET("/my", middleware.RoleMiddleware(string(models.RoleEmployee)), controllers.GetMyFavorites)
			favorites.GET("/task/:task_id/check", middleware.RoleMiddleware(string(models.RoleEmployee)), controllers.CheckFavorite)
		}

		// 完成任务
		api.POST("/tasks/:task_id/complete", controllers.CompleteTask)
	}

	// 管理员路由
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.RoleMiddleware(string(models.RoleAdmin)))
	{
		// 用户管理
		users := admin.Group("/users")
		{
			users.GET("", controllers.AdminGetAllUsers)
			users.GET("/:id", controllers.AdminGetUserByID)
			users.PUT("/:id", controllers.AdminUpdateUserStatus)
			users.DELETE("/:id", controllers.AdminDeleteUser)
		}

		// 任务管理
		tasks := admin.Group("/tasks")
		{
			tasks.GET("", controllers.AdminGetAllTasks)
			tasks.PUT("/:id/status", controllers.AdminUpdateTaskStatus)
			tasks.DELETE("/:id", controllers.AdminDeleteTask)
		}
	}

	return r
}
