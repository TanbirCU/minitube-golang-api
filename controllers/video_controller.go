package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"minitube-api/config"
	"minitube-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadVideo(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)

	title := c.PostForm("title")
	description := c.PostForm("description")
	categoryID := c.PostForm("category_id")

	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Title is required",
		})
		return
	}

	file, err := c.FormFile("video")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Video file is required",
		})
		return
	}

	extension := strings.ToLower(
		filepath.Ext(file.Filename),
	)

	if extension != ".mp4" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Only MP4 videos are supported",
		})
		return
	}

	if file.Size > 500*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Maximum video size is 500MB",
		})
		return
	}

	uploadDir := "uploads/videos"

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not create upload directory",
		})
		return
	}

	filename := uuid.New().String() + extension

	filePath := filepath.Join(
		uploadDir,
		filename,
	)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Video upload failed",
		})
		return
	}

	catID, _ := strconv.ParseUint(categoryID, 10, 32)

	video := models.Video{
		UserID:      userID,
		CategoryID:  uint(catID),
		Title:       title,
		Description: description,
		Filename:    filename,
		VideoPath:   filePath,
	}

	if err := config.DB.Create(&video).Error; err != nil {

		os.Remove(filePath)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not save video information",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Video uploaded successfully",
		"data":    video,
		"stream_url": fmt.Sprintf(
			"/api/videos/%d/stream",
			video.ID,
		),
	})
}

func GetVideos(c *gin.Context) {

	var videos []models.Video

	query := config.DB.
		Preload("User").
		Preload("Category")

	search := c.Query("search")
	category := c.Query("category")

	if search != "" {

		query = query.Where(
			"title LIKE ?",
			"%"+search+"%",
		)
	}

	if category != "" {

		query = query.Where(
			"category_id = ?",
			category,
		)
	}

	query.
		Order("created_at DESC").
		Find(&videos)

	c.JSON(http.StatusOK, gin.H{
		"data": videos,
	})
}

func GetVideo(c *gin.Context) {

	id := c.Param("id")

	var video models.Video

	if err := config.DB.
		Preload("User").
		Preload("Category").
		First(&video, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Video not found",
		})
		return
	}

	config.DB.
		Model(&video).
		UpdateColumn("views", video.Views+1)

	c.JSON(http.StatusOK, gin.H{
		"data": video,
	})
}

func StreamVideo(c *gin.Context) {

	id := c.Param("id")

	var video models.Video

	if err := config.DB.First(&video, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Video not found",
		})
		return
	}

	file, err := os.Open(video.VideoPath)

	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Video file not found",
		})
		return
	}

	defer file.Close()

	stat, err := file.Stat()

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not read video",
		})
		return
	}

	http.ServeContent(
		c.Writer,
		c.Request,
		stat.Name(),
		stat.ModTime(),
		file,
	)
}


func DeleteVideo(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)

	id := c.Param("id")

	var video models.Video

	if err := config.DB.First(&video, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Video not found",
		})
		return
	}

	if video.UserID != userID {

		c.JSON(http.StatusForbidden, gin.H{
			"message": "You cannot delete this video",
		})
		return
	}

	os.Remove(video.VideoPath)

	config.DB.Delete(&video)

	c.JSON(http.StatusOK, gin.H{
		"message": "Video deleted successfully",
	})
}