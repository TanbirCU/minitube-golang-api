package main

import (
	"os"

	"minitube-api/config"
	"minitube-api/models"
	"minitube-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// Database
	config.ConnectDatabase()

	// Migration
	config.DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Video{},
	)

	// Create upload directory
	os.MkdirAll(
		"uploads/videos",
		0755,
	)

	// Gin
	router := gin.Default()

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
	}))

	// Routes
	routes.SetupRoutes(router)

	// Server
	router.Run(":8080")
}