package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"

    "github.com/example/golang-email/internal/service"
)

func RegisterRecord(rg *gin.RouterGroup, engine *service.SendEngine) {
    rg.POST("/records/:id/retry", func(c *gin.Context) {
        id, _ := strconv.Atoi(c.Param("id"))
        if err := engine.RetryRecord(uint(id)); err != nil { logrus.WithError(err).Error("重试失败"); c.JSON(http.StatusInternalServerError, gin.H{"message": "重试失败"}); return }
        c.JSON(http.StatusOK, gin.H{"message": "重试已完成"})
    })
}