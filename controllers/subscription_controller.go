package controllers

import (
	"net/http"
	"strconv"

	"minitube-api/config"
	"minitube-api/models"

	"github.com/gin-gonic/gin"
)

func ToggleSubscribe(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	channelIDStr := c.Param("id")
	channelID, err := strconv.ParseUint(channelIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid channel ID",
		})
		return
	}

	if userID == uint(channelID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "You cannot subscribe to your own channel",
		})
		return
	}

	var sub models.Subscription
	err = config.DB.
		Where("subscriber_id = ? AND channel_id = ?", userID, channelID).
		First(&sub).Error

	isSubscribed := false
	if err == nil {
		// Already subscribed, unsubscribe
		config.DB.Delete(&sub)
		isSubscribed = false
	} else {
		// Not subscribed, subscribe
		newSub := models.Subscription{
			SubscriberID: userID,
			ChannelID:    uint(channelID),
		}
		if err := config.DB.Create(&newSub).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Could not update subscription",
			})
			return
		}
		isSubscribed = true
	}

	var count int64
	config.DB.Model(&models.Subscription{}).
		Where("channel_id = ?", channelID).
		Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"is_subscribed": isSubscribed,
		"subscribers":   count,
	})
}

func GetSubscriptionStatus(c *gin.Context) {
	channelIDStr := c.Param("id")
	channelID, err := strconv.ParseUint(channelIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid channel ID",
		})
		return
	}

	var count int64
	config.DB.Model(&models.Subscription{}).
		Where("channel_id = ?", channelID).
		Count(&count)

	isSubscribed := false
	if val, exists := c.Get("user_id"); exists {
		if uid, ok := val.(uint); ok {
			var sub models.Subscription
			if err := config.DB.
				Where("subscriber_id = ? AND channel_id = ?", uid, channelID).
				First(&sub).Error; err == nil {
				isSubscribed = true
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"is_subscribed": isSubscribed,
		"subscribers":   count,
	})
}
