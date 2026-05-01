package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"recruit/database"
	"recruit/models"
)

type CreateBidRequest struct {
	ProposedPrice float64 `json:"proposed_price" binding:"required,min=0"`
	Proposal      string  `json:"proposal" binding:"required"`
}

func CreateBid(c *gin.Context) {
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

	// 检查任务是否存在且状态为开放
	var task models.Task
	if err := database.GetDB().First(&task, taskID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if task.Status != models.TaskStatusOpen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task is not open for bidding"})
		return
	}

	// 检查用户是否是任务发布者（不能投标自己的任务）
	if task.EmployerID == userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot bid on your own task"})
		return
	}

	// 检查是否已经投标过
	var existingBid models.Bid
	if err := database.GetDB().Where("task_id = ? AND employee_id = ?", taskID, userID).First(&existingBid).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "You have already bid on this task"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var req CreateBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bid := models.Bid{
		TaskID:        uint(taskID),
		EmployeeID:    userID.(uint),
		ProposedPrice: req.ProposedPrice,
		Proposal:      req.Proposal,
		Status:        models.BidStatusPending,
	}

	if err := database.GetDB().Create(&bid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bid"})
		return
	}

	c.JSON(http.StatusCreated, bid)
}

func GetBidsByTask(c *gin.Context) {
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var bids []models.Bid
	if err := database.GetDB().Preload("Employee").Where("task_id = ?", taskID).Find(&bids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bids"})
		return
	}

	c.JSON(http.StatusOK, bids)
}

func GetMyBids(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var bids []models.Bid
	if err := database.GetDB().Preload("Task").Preload("Task.Employer").Where("employee_id = ?", userID).Find(&bids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bids"})
		return
	}

	c.JSON(http.StatusOK, bids)
}

func AcceptBid(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	bidIDStr := c.Param("bid_id")
	bidID, err := strconv.ParseUint(bidIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
		return
	}

	var bid models.Bid
	if err := database.GetDB().Preload("Task").First(&bid, bidID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bid not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 检查是否是任务发布者
	if bid.Task.EmployerID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to accept this bid"})
		return
	}

	// 检查任务状态
	if bid.Task.Status != models.TaskStatusOpen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task is not open for bidding"})
		return
	}

	// 开始事务
	tx := database.GetDB().Begin()

	// 更新投标状态为已接受
	if err := tx.Model(&bid).Update("status", models.BidStatusAccepted).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept bid"})
		return
	}

	// 拒绝其他所有投标
	if err := tx.Model(&models.Bid{}).Where("task_id = ? AND id != ?", bid.TaskID, bidID).Update("status", models.BidStatusRejected).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject other bids"})
		return
	}

	// 更新任务状态和中标者
	if err := tx.Model(&models.Task{}).Where("id = ?", bid.TaskID).Updates(map[string]interface{}{
		"status":    models.TaskStatusInProgress,
		"winner_id": bid.EmployeeID,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Bid accepted successfully", "bid": bid})
}
