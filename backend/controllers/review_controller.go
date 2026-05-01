package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"recruit/database"
	"recruit/models"
)

type CreateReviewRequest struct {
	Rating  uint8  `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"required"`
}

func CreateReview(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 检查任务是否存在且已完成
	var task models.Task
	if err := database.GetDB().Preload("Winner").First(&task, taskID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if task.Status != models.TaskStatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task is not completed yet"})
		return
	}

	// 检查是否是任务发布者
	if task.EmployerID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to review this task"})
		return
	}

	// 检查是否已经评价过
	var existingReview models.Review
	if err := database.GetDB().Where("task_id = ?", taskID).First(&existingReview).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "You have already reviewed this task"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 检查是否有中标者
	if task.WinnerID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No winner for this task"})
		return
	}

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review := models.Review{
		TaskID:     uint(taskID),
		EmployerID: userID.(uint),
		EmployeeID: *task.WinnerID,
		Rating:     req.Rating,
		Comment:    req.Comment,
	}

	if err := database.GetDB().Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func GetReviewsByEmployee(c *gin.Context) {
	employeeIDStr := c.Param("employee_id")
	employeeID, err := strconv.ParseUint(employeeIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var reviews []models.Review
	if err := database.GetDB().Preload("Task").Preload("Employer").Where("employee_id = ?", employeeID).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func GetMyReviews(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var reviews []models.Review
	if err := database.GetDB().Preload("Task").Preload("Employer").Where("employee_id = ?", userID).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}
