package controllers

import (
	"net/http"
	"strings"

	"minitube-api/config"
	"minitube-api/models"
	"minitube-api/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func bindRegisterRequest(c *gin.Context) (RegisterRequest, bool) {
	contentType := c.GetHeader("Content-Type")
	var req RegisterRequest

	if strings.Contains(contentType, "application/json") {
		if err := c.ShouldBindJSON(&req); err != nil {
			return req, false
		}
	} else {
		// form-data or x-www-form-urlencoded
		req.Name = c.PostForm("name")
		req.Email = c.PostForm("email")
		req.Password = c.PostForm("password")
	}
	return req, true
}

func Register(c *gin.Context) {
	request, ok := bindRegisterRequest(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// Validation
	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(strings.ToLower(request.Email))
	request.Password = strings.TrimSpace(request.Password)

	if request.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Name is required",
		})
		return
	}

	if request.Email == "" || !strings.Contains(request.Email, "@") {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Valid email is required",
		})
		return
	}

	if len(request.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password must be at least 6 characters",
		})
		return
	}

	var existingUser models.User

	if err := config.DB.
		Where("email = ?", request.Email).
		First(&existingUser).Error; err == nil {

		c.JSON(http.StatusConflict, gin.H{
			"message": "Email already exists",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Password hashing failed",
		})
		return
	}

	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "User registration failed",
		})
		return
	}

	token, err := utils.GenerateToken(user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Token generation failed",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var request LoginRequest
	contentType := c.GetHeader("Content-Type")

	if strings.Contains(contentType, "application/json") {
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request body",
			})
			return
		}
	} else {
		request.Email = c.PostForm("email")
		request.Password = c.PostForm("password")
	}

	request.Email = strings.TrimSpace(strings.ToLower(request.Email))

	if request.Email == "" || request.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email and password are required",
		})
		return
	}

	var user models.User

	if err := config.DB.
		Where("email = ?", request.Email).
		First(&user).Error; err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(request.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	token, err := utils.GenerateToken(user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Token generation failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"avatar": user.Avatar,
		},
	})
}

func Me(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)

	var user models.User

	if err := config.DB.First(&user, userID).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"avatar": user.Avatar,
		},
	})
}