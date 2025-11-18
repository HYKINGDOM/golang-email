package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"

    "github.com/example/golang-email/internal/model"
    "github.com/example/golang-email/internal/repository"
)

type TemplateRequest struct {
    Name            string `json:"name" binding:"required"`
    Subject         string `json:"subject" binding:"required"`
    Content         string `json:"content" binding:"required"`
    IsRichText      bool   `json:"is_rich_text"`
    TrackingEnabled bool   `json:"tracking_enabled"`
}

func RegisterTemplate(rg *gin.RouterGroup, repo repository.EmailTemplateRepository) {
    rg.GET("/templates", func(c *gin.Context) {
        items, err := repo.List()
        if err != nil { logrus.WithError(err).Error("模板列表失败"); c.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败"}); return }
        c.JSON(http.StatusOK, items)
    })
    rg.POST("/templates", func(c *gin.Context) {
        var req TemplateRequest
        if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"}); return }
        t := &model.EmailTemplate{Name: req.Name, Subject: req.Subject, Content: req.Content, IsRichText: req.IsRichText, TrackingEnabled: req.TrackingEnabled}
        if err := repo.Create(t); err != nil { logrus.WithError(err).Error("模板创建失败"); c.JSON(http.StatusInternalServerError, gin.H{"message": "创建失败"}); return }
        c.JSON(http.StatusOK, t)
    })
    rg.PUT("/templates/:id", func(c *gin.Context) {
        var req TemplateRequest
        if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"}); return }
        id, _ := strconv.Atoi(c.Param("id"))
        t, err := repo.GetByID(uint(id))
        if err != nil { c.JSON(http.StatusNotFound, gin.H{"message": "未找到"}); return }
        t.Name = req.Name; t.Subject = req.Subject; t.Content = req.Content; t.IsRichText = req.IsRichText; t.TrackingEnabled = req.TrackingEnabled
        if err := repo.Update(t); err != nil { logrus.WithError(err).Error("模板更新失败"); c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"}); return }
        c.JSON(http.StatusOK, t)
    })
    rg.DELETE("/templates/:id", func(c *gin.Context) {
        if err := repo.Delete(uint(parseID(c.Param("id")))); err != nil { logrus.WithError(err).Error("模板删除失败"); c.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败"}); return }
        c.Status(http.StatusNoContent)
    })
}

func parseID(s string) int { i, _ := strconv.Atoi(s); return i }