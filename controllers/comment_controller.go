package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"minitube-api/config"
	"minitube-api/models"

	"github.com/gin-gonic/gin"
)

type CreateCommentInput struct {
	Text string `json:"comment" binding:"required"`
}

func GetVideoComments(c *gin.Context) {
	videoID := c.Param("id")

	var comments []models.Comment
	if err := config.DB.
		Preload("User").
		Where("video_id = ?", videoID).
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch comments",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": comments,
	})
}

func CreateComment(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	videoIDStr := c.Param("id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid video ID",
		})
		return
	}

	var input CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Comment text is required",
		})
		return
	}

	text := strings.TrimSpace(input.Text)
	if text == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Comment cannot be empty",
		})
		return
	}

	comment := models.Comment{
		VideoID: uint(videoID),
		UserID:  userID,
		Text:    text,
		Likes:   0,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not save comment",
		})
		return
	}

	config.DB.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment added successfully",
		"data":    comment,
	})
}

func LikeComment(c *gin.Context) {
	commentID := c.Param("id")

	var comment models.Comment
	if err := config.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Comment not found",
		})
		return
	}

	config.DB.Model(&comment).UpdateColumn("likes", comment.Likes+1)

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment liked",
		"likes":   comment.Likes + 1,
	})
}

func DeleteComment(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	commentID := c.Param("id")

	var comment models.Comment
	if err := config.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Comment not found",
		})
		return
	}

	if comment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "You can only delete your own comments",
		})
		return
	}

	if err := config.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not delete comment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})
}
