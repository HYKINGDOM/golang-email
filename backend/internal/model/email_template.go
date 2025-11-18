package model

import "time"

type EmailTemplate struct {
    ID             uint      `gorm:"primaryKey"`
    Name           string    `gorm:"size:100;not null"`
    Subject        string    `gorm:"size:255;not null"`
    Content        string    `gorm:"type:text;not null"`
    IsRichText     bool      `gorm:"default:true"`
    TrackingEnabled bool     `gorm:"default:false"`
    CreatedAt      time.Time
    UpdatedAt      time.Time
}