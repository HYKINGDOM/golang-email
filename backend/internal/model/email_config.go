package model

import "time"

type EmailConfig struct {
    ID         uint      `gorm:"primaryKey"`
    Provider   string    `gorm:"size:50;not null"`
    Host       string    `gorm:"size:255;not null"`
    Port       int       `gorm:"not null"`
    Username   string    `gorm:"size:255;uniqueIndex;not null"`
    Password   string    `gorm:"size:255;not null"`
    DailyLimit int       `gorm:"default:200"`
    IsActive   bool      `gorm:"default:true"`
    CreatedAt  time.Time
    UpdatedAt  time.Time
}