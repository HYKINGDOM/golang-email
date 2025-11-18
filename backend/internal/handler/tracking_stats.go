package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "github.com/example/golang-email/internal/model"
)

func RegisterTrackingStats(rg *gin.RouterGroup, db *gorm.DB) {
    rg.GET("/tracking/stats", func(c *gin.Context) {
        tid := c.Query("task_id")
        id, _ := strconv.Atoi(tid)
        var total int64
        var opens int64
        db.Model(&model.SendRecord{}).Where("task_id = ?", id).Count(&total)
        db.Table("email_tracking").Joins("JOIN send_records sr ON sr.id = email_tracking.record_id").Where("sr.task_id = ? AND email_tracking.open_count > 0", id).Count(&opens)
        c.JSON(http.StatusOK, gin.H{"total": total, "opened": opens})
    })
}