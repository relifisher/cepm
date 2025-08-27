package models

import (
	"time"

	"gorm.io/gorm"
)

// Department 部门表
type Department struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	ParentID  *uint
	Parent    *Department `gorm:"foreignKey:ParentID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Role 角色表
type Role struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null;unique"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// User 员工表
type User struct {
	ID            uint   `gorm:"primaryKey"`
	WechatUserid  string `gorm:"unique"`
	Name          string `gorm:"not null"`
	EnglishName   string // New field: English Name
	Email         string `gorm:"unique"`
	Avatar        string
	DepartmentID  *uint
	Department    Department `gorm:"foreignKey:DepartmentID"`
	RoleID        *uint
	Role          Role `gorm:"foreignKey:RoleID"`
	ManagerID     *uint
	Manager       *User `gorm:"foreignKey:ManagerID"`
	IsActive      bool  `gorm:"default:true"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// PerformanceReview 月度绩效评估主表
type PerformanceReview struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"not null;uniqueIndex:idx_user_period,priority:1"`
	User         User      `gorm:"foreignKey:UserID"`
	Period       string    `gorm:"not null;uniqueIndex:idx_user_period,priority:2"`
	Status       string    `gorm:"not null;default:'Draft'"` // Status: Draft, PendingApproval, Approved, PendingScore, Completed, PendingHRConfirmation, Archived, Rejected
	TotalScore   *float64  `gorm:"type:numeric(5,2)"`
	GradePoint   *float64  `gorm:"type:numeric(5,2)"` // New field: Performance Grade Point
	FinalComment string
	Items        []PerformanceItem `gorm:"foreignKey:ReviewID"`
	Approvals    []ApprovalHistory `gorm:"foreignKey:ReviewID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// PerformanceItem 绩效评估项表
type PerformanceItem struct {
	ID                 uint     `gorm:"primaryKey"`
	ReviewID           uint     `gorm:"not null"`
	Category           string   `gorm:"not null;default:'工作业绩'"` // 工作业绩, 大模型, 价值观
	Title              string   `gorm:"not null"`
	Description        string
	Weight             float64  `gorm:"not null;type:numeric(5,2)"`
	Target             string
	CompletionDetails  string
	Score              *float64 `gorm:"type:numeric(5,2)"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// ApprovalHistory 审批流转历史表
type ApprovalHistory struct {
	ID         uint   `gorm:"primaryKey"`
	ReviewID   uint   `gorm:"not null"`
	ApproverID uint   `gorm:"not null"`
	Approver   User   `gorm:"foreignKey:ApproverID"`
	Status     string `gorm:"not null"`
	Comment    string
	CreatedAt  time.Time
}

// SystemSetting 系统设置表
type SystemSetting struct {
	ID        uint   `gorm:"primaryKey"`
	Key       string `gorm:"not null;unique"`
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AutoMigrate will automatically migrate the schema, creating tables and columns
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&Department{}, &Role{}, &User{}, &PerformanceReview{}, &PerformanceItem{}, &ApprovalHistory{}, &SystemSetting{})
}
