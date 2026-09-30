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

	auth.POST(
		"/categories",
		controllers.CreateCategory,
	)

	// =========================
	// VIDEOS
	// =========================

	api.GET("/videos", controllers.GetVideos)

	api.GET(
		"/videos/:id",
		controllers.GetVideo,
	)

	api.GET(
		"/videos/:id/stream",
		controllers.StreamVideo,
	)

	auth.POST(
		"/videos",
		controllers.UploadVideo,
	)

	auth.DELETE(
		"/videos/:id",
		controllers.DeleteVideo,
	)
}