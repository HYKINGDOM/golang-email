package model

import "time"

type SendRecordStatus string

const (
    RecordStatusPending SendRecordStatus = "pending"
    RecordStatusSent    SendRecordStatus = "sent"
    RecordStatusFailed  SendRecordStatus = "failed"
)

type SendRecord struct {
    ID             uint             `gorm:"primaryKey"`
    TaskID         uint             `gorm:"index;not null"`
    SenderEmail    string           `gorm:"size:255;not null"`
    RecipientEmail string           `gorm:"size:255;not null"`
    SendTime       *time.Time       `gorm:"index"`
    Status         SendRecordStatus `gorm:"size:20;default:'pending'"`
    ErrorMessage   string           `gorm:"type:text"`
    RetryCount     int              `gorm:"default:0"`
    CreatedAt      time.Time
    UpdatedAt      time.Time
}