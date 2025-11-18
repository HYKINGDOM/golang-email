package model

import "time"

type LinkClick struct {
    ID        uint      `gorm:"primaryKey"`
    RecordID  uint      `gorm:"index;not null"`
    URL       string    `gorm:"type:text;not null"`
    ClickTime time.Time `gorm:"index"`
    CreatedAt time.Time
}