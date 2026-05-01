package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"recruit/database"
	"recruit/models"
)

func ToggleFavorite(c *gin.Context) {
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

	// 检查任务是否存在
	var task models.Task
	if err := database.GetDB().First(&task, taskID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 检查是否已经收藏
	var favorite models.Favorite
	if err := database.GetDB().Where("task_id = ? AND employee_id = ?", taskID, userID).First(&favorite).Error; err == nil {
		// 已经收藏，取消收藏
		if err := database.GetDB().Delete(&favorite).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove favorite"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Favorite removed successfully", "is_favorited": false})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 未收藏，添加收藏
	newFavorite := models.Favorite{
		TaskID:     uint(taskID),
		EmployeeID: userID.(uint),
	}

	if err := database.GetDB().Create(&newFavorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add favorite"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Favorite added successfully", "is_favorited": true})
}

func GetMyFavorites(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var favorites []models.Favorite
	if err := database.GetDB().Preload("Task").Preload("Task.Employer").Where("employee_id = ?", userID).Find(&favorites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch favorites"})
		return
	}

	c.JSON(http.StatusOK, favorites)
}

func CheckFavorite(c *gin.Context) {
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

	var favorite models.Favorite
	if err := database.GetDB().Where("task_id = ? AND employee_id = ?", taskID, userID).First(&favorite).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"is_favorited": true})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_favorited": false})
}
