package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"recruit/database"
	"recruit/models"
)

type CreateTaskRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Budget      float64 `json:"budget" binding:"required,min=0"`
}

type UpdateTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Budget      float64 `json:"budget" binding:"min=0"`
}

func CreateTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		Budget:      req.Budget,
		EmployerID:  userID.(uint),
		Status:      models.TaskStatusOpen,
	}

	if err := database.GetDB().Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func GetTaskByID(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.Task
	if err := database.GetDB().Preload("Employer").Preload("Bids.Employee").First(&task, taskID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func GetAllTasks(c *gin.Context) {
	var tasks []models.Task

	query := database.GetDB().Preload("Employer")

	// 可选过滤参数
	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	employerIDStr := c.Query("employer_id")
	if employerIDStr != "" {
		employerID, err := strconv.ParseUint(employerIDStr, 10, 64)
		if err == nil {
			query = query.Where("employer_id = ?", employerID)
		}
	}

	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func GetMyTasks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var tasks []models.Task
	if err := database.GetDB().Preload("Employer").Preload("Bids.Employee").Where("employer_id = ?", userID).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func UpdateTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.Task
	if err := database.GetDB().First(&task, taskID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 检查是否是任务发布者
	if task.EmployerID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this task"})
		return
	}

	// 只能更新未开始的任务
	if task.Status != models.TaskStatusOpen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update task that is not in open status"})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Budget > 0 {
		updates["budget"] = req.Budget
	}

	if err := database.GetDB().Model(&task).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.Task
	if err := database.GetDB().First(&task, taskID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 检查是否是任务发布者
	if task.EmployerID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this task"})
		return
	}

	// 只能删除未开始的任务
	if task.Status != models.TaskStatusOpen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete task that is not in open status"})
		return
	}

	if err := database.GetDB().Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
