package model

import "time"

type EmailTracking struct {
    ID            uint       `gorm:"primaryKey"`
    RecordID      uint       `gorm:"index;not null"`
    OpenTime      *time.Time
    OpenCount     int        `gorm:"default:0"`
    ReadDuration  int        `gorm:"default:0"`
    LastOpenTime  *time.Time
    CreatedAt     time.Time
    UpdatedAt     time.Time
}