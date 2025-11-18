package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterHealth(rg *gin.RouterGroup, db *gorm.DB) {
    rg.GET("/health", func(c *gin.Context) {
        sqlDB, err := db.DB()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "数据库不可用"})
            return
        }
        if err := sqlDB.Ping(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "数据库未连接"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })
}