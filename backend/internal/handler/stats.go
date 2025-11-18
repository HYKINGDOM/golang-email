package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "github.com/example/golang-email/internal/model"
)

func RegisterStats(rg *gin.RouterGroup, db *gorm.DB) {
    rg.GET("/stats/send", func(c *gin.Context) {
        var total, sent, failed int64
        db.Model(&model.SendRecord{}).Count(&total)
        db.Model(&model.SendRecord{}).Where("status = ?", model.RecordStatusSent).Count(&sent)
        db.Model(&model.SendRecord{}).Where("status = ?", model.RecordStatusFailed).Count(&failed)
        c.JSON(http.StatusOK, gin.H{"total": total, "sent": sent, "failed": failed})
    })
}