package routes

import (
	"minitube-api/controllers"
	"minitube-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	api := router.Group("/api")

	// =========================
	// AUTH
	// =========================

	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)

	auth := api.Group("")
	auth.Use(middleware.AuthMiddleware())

	auth.GET("/me", controllers.Me)

	// =========================
	// CATEGORIES
	// =========================

	api.GET("/categories", controllers.GetCategories)
	auth.POST("/categories", controllers.CreateCategory)

	// =========================
	// VIDEOS
	// =========================

	api.GET("/videos", controllers.GetVideos)
	api.GET("/videos/:id", controllers.GetVideo)
	api.GET("/videos/:id/stream", controllers.StreamVideo)
	api.GET("/videos/stream/:id", controllers.StreamVideo)

	auth.POST("/videos", controllers.UploadVideo)
	auth.DELETE("/videos/:id", controllers.DeleteVideo)

	// =========================
	// COMMENTS
	// =========================

	api.GET("/videos/:id/comments", controllers.GetVideoComments)
	auth.POST("/videos/:id/comments", controllers.CreateComment)
	api.POST("/comments/:id/like", controllers.LikeComment)
	auth.DELETE("/comments/:id", controllers.DeleteComment)

	// =========================
	// SUBSCRIPTIONS
	// =========================

	api.GET("/channels/:id/subscribe", controllers.GetSubscriptionStatus)
	auth.POST("/channels/:id/subscribe", controllers.ToggleSubscribe)
}