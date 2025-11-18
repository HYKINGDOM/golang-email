package model

import "time"

type SendTaskStatus string

const (
    TaskStatusPending  SendTaskStatus = "pending"
    TaskStatusRunning  SendTaskStatus = "running"
    TaskStatusPaused   SendTaskStatus = "paused"
    TaskStatusFailed   SendTaskStatus = "failed"
    TaskStatusFinished SendTaskStatus = "finished"
)

type SendTask struct {
    ID             uint           `gorm:"primaryKey"`
    Name           string         `gorm:"size:100;not null"`
    SenderConfigs  string         `gorm:"type:text;not null"`
    RecipientList  string         `gorm:"type:text;not null"`
    TemplateID     uint           `gorm:"not null"`
    ScheduledTime  *time.Time     `gorm:"index"`
    Status         SendTaskStatus `gorm:"size:20;default:'pending'"`
    CreatedAt      time.Time
    UpdatedAt      time.Time
}