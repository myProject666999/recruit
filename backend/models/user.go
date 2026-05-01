package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleEmployer  UserRole = "employer"
	RoleEmployee  UserRole = "employee"
	RoleAdmin     UserRole = "admin"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Role      UserRole       `json:"role" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关系
	Tasks       []Task       `json:"tasks,omitempty" gorm:"foreignKey:EmployerID"`
	Bids        []Bid        `json:"bids,omitempty" gorm:"foreignKey:EmployeeID"`
	Reviews     []Review     `json:"reviews,omitempty" gorm:"foreignKey:EmployeeID"`
	Favorites   []Favorite   `json:"favorites,omitempty" gorm:"foreignKey:EmployeeID"`
}

type TaskStatus string

const (
	TaskStatusOpen       TaskStatus = "open"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

type Task struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"not null"`
	Budget      float64        `json:"budget" gorm:"not null"`
	Status      TaskStatus     `json:"status" gorm:"default:open"`
	EmployerID  uint           `json:"employer_id" gorm:"not null"`
	Employer    User           `json:"employer" gorm:"foreignKey:EmployerID"`
	WinnerID    *uint          `json:"winner_id,omitempty"`
	Winner      *User          `json:"winner,omitempty" gorm:"foreignKey:WinnerID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关系
	Bids      []Bid      `json:"bids,omitempty" gorm:"foreignKey:TaskID"`
}

type BidStatus string

const (
	BidStatusPending   BidStatus = "pending"
	BidStatusAccepted  BidStatus = "accepted"
	BidStatusRejected  BidStatus = "rejected"
)

type Bid struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TaskID      uint           `json:"task_id" gorm:"not null"`
	Task        Task           `json:"task" gorm:"foreignKey:TaskID"`
	EmployeeID  uint           `json:"employee_id" gorm:"not null"`
	Employee    User           `json:"employee" gorm:"foreignKey:EmployeeID"`
	ProposedPrice float64      `json:"proposed_price" gorm:"not null"`
	Proposal    string         `json:"proposal" gorm:"not null"`
	Status      BidStatus      `json:"status" gorm:"default:pending"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Review struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TaskID      uint           `json:"task_id" gorm:"not null"`
	Task        Task           `json:"task" gorm:"foreignKey:TaskID"`
	EmployerID  uint           `json:"employer_id" gorm:"not null"`
	Employer    User           `json:"employer" gorm:"foreignKey:EmployerID"`
	EmployeeID  uint           `json:"employee_id" gorm:"not null"`
	Employee    User           `json:"employee" gorm:"foreignKey:EmployeeID"`
	Rating      uint8          `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Comment     string         `json:"comment" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Favorite struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TaskID      uint           `json:"task_id" gorm:"not null"`
	Task        Task           `json:"task" gorm:"foreignKey:TaskID"`
	EmployeeID  uint           `json:"employee_id" gorm:"not null"`
	Employee    User           `json:"employee" gorm:"foreignKey:EmployeeID"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 唯一约束：同一个雇员不能收藏同一个任务多次
	gorm.Model `json:"-" gorm:"uniqueIndex:idx_task_employee,priority:2"`
}
