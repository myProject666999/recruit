package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"recruit/database"
	"recruit/models"
)

func AdminGetAllUsers(c *gin.Context) {
	var users []models.User

	query := database.GetDB()

	// 可选过滤参数
	role := c.Query("role")
	if role != "" {
		query = query.Where("role = ?", role)
	}

	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func AdminGetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := database.GetDB().Preload("Tasks").Preload("Bids").Preload("Reviews").First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func AdminUpdateUserStatus(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 管理员不能被修改
	if user.Role == models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot modify admin user"})
		return
	}

	var req struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role != "" {
		if req.Role != string(models.RoleEmployer) && req.Role != string(models.RoleEmployee) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
			return
		}
		user.Role = models.UserRole(req.Role)
	}

	if err := database.GetDB().Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func AdminDeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 管理员不能被删除
	if user.Role == models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete admin user"})
		return
	}

	if err := database.GetDB().Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func AdminGetAllTasks(c *gin.Context) {
	var tasks []models.Task

	query := database.GetDB().Preload("Employer").Preload("Winner")

	// 可选过滤参数
	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func AdminUpdateTaskStatus(c *gin.Context) {
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

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validStatuses := []models.TaskStatus{
		models.TaskStatusOpen,
		models.TaskStatusInProgress,
		models.TaskStatusCompleted,
		models.TaskStatusCancelled,
	}

	var isValidStatus bool
	for _, s := range validStatuses {
		if string(s) == req.Status {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	task.Status = models.TaskStatus(req.Status)

	if err := database.GetDB().Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func AdminDeleteTask(c *gin.Context) {
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

	if err := database.GetDB().Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func CompleteTask(c *gin.Context) {
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

	var task models.Task
	if err := database.GetDB().First(&task, taskID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 检查权限：只有中标者或任务发布者可以完成任务
	userRole, _ := c.Get("user_role")
	isWinner := task.WinnerID != nil && *task.WinnerID == userID.(uint)
	isEmployer := task.EmployerID == userID.(uint)
	isAdmin := userRole == string(models.RoleAdmin)

	if !isWinner && !isEmployer && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to complete this task"})
		return
	}

	// 检查任务状态
	if task.Status != models.TaskStatusInProgress {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task is not in progress"})
		return
	}

	task.Status = models.TaskStatusCompleted

	if err := database.GetDB().Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task completed successfully", "task": task})
}
