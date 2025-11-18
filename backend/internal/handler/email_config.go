package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"

    "github.com/example/golang-email/internal/model"
    "github.com/example/golang-email/internal/service"
)

type EmailConfigRequest struct {
    Provider   string `json:"provider" binding:"required"`
    Host       string `json:"host" binding:"required"`
    Port       int    `json:"port" binding:"required"`
    Username   string `json:"username" binding:"required,email"`
    Password   string `json:"password" binding:"required"`
    DailyLimit int    `json:"daily_limit"`
    IsActive   bool   `json:"is_active"`
}

func RegisterEmailConfig(rg *gin.RouterGroup, svc service.EmailConfigService) {
    rg.GET("/email-configs", func(c *gin.Context) {
        items, err := svc.List()
        if err != nil {
            logrus.WithError(err).Error("邮箱配置列表查询失败")
            c.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败"})
            return
        }
        for i := range items {
            items[i].Password = ""
        }
        c.JSON(http.StatusOK, items)
    })

    rg.GET("/email-configs/:id", func(c *gin.Context) {
        idStr := c.Param("id")
        id, _ := strconv.Atoi(idStr)
        item, err := svc.Get(uint(id))
        if err != nil {
            logrus.WithError(err).Error("邮箱配置查询失败")
            c.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败"})
            return
        }
        if item == nil {
            c.JSON(http.StatusNotFound, gin.H{"message": "未找到"})
            return
        }
        item.Password = ""
        c.JSON(http.StatusOK, item)
    })

    rg.POST("/email-configs", func(c *gin.Context) {
        var req EmailConfigRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
            return
        }
        cfg := &model.EmailConfig{
            Provider:   req.Provider,
            Host:       req.Host,
            Port:       req.Port,
            Username:   req.Username,
            Password:   req.Password,
            DailyLimit: req.DailyLimit,
            IsActive:   req.IsActive,
        }
        res, err := svc.Create(cfg)
        if err != nil {
            logrus.WithError(err).Error("邮箱配置创建失败")
            c.JSON(http.StatusInternalServerError, gin.H{"message": "创建失败"})
            return
        }
        c.JSON(http.StatusOK, res)
    })

    rg.PUT("/email-configs/:id", func(c *gin.Context) {
        idStr := c.Param("id")
        id, _ := strconv.Atoi(idStr)
        var req EmailConfigRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
            return
        }
        upd := &model.EmailConfig{
            Provider:   req.Provider,
            Host:       req.Host,
            Port:       req.Port,
            Username:   req.Username,
            Password:   req.Password,
            DailyLimit: req.DailyLimit,
            IsActive:   req.IsActive,
        }
        res, err := svc.Update(uint(id), upd)
        if err != nil {
            logrus.WithError(err).Error("邮箱配置更新失败")
            c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
            return
        }
        c.JSON(http.StatusOK, res)
    })

    rg.DELETE("/email-configs/:id", func(c *gin.Context) {
        idStr := c.Param("id")
        id, _ := strconv.Atoi(idStr)
        if err := svc.Delete(uint(id)); err != nil {
            logrus.WithError(err).Error("邮箱配置删除失败")
            c.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败"})
            return
        }
        c.Status(http.StatusNoContent)
    })
}