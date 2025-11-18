package handler

import (
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "gorm.io/gorm"

    "github.com/example/golang-email/internal/model"
    "github.com/example/golang-email/internal/service"
)

type TaskRequest struct {
    Name          string    `json:"name" binding:"required"`
    SenderConfigs []uint    `json:"sender_configs" binding:"required"`
    RecipientList []string  `json:"recipient_list" binding:"required"`
    TemplateID    uint      `json:"template_id" binding:"required"`
    ScheduledTime *time.Time `json:"scheduled_time"`
}

func RegisterTask(rg *gin.RouterGroup, db *gorm.DB, engine *service.SendEngine) {
    rg.POST("/tasks", func(c *gin.Context) {
        var req TaskRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
            return
        }
        recips := ""
        for i, r := range req.RecipientList {
            if i > 0 {
                recips += ","
            }
            recips += r
        }
        senders := ""
        for i, id := range req.SenderConfigs {
            if i > 0 {
                senders += ","
            }
            senders += strconv.Itoa(int(id))
        }
        task := &model.SendTask{Name: req.Name, SenderConfigs: senders, RecipientList: recips, TemplateID: req.TemplateID, Status: model.TaskStatusPending, ScheduledTime: req.ScheduledTime}
        if err := db.Create(task).Error; err != nil {
            logrus.WithError(err).Error("任务创建失败")
            c.JSON(http.StatusInternalServerError, gin.H{"message": "创建失败"})
            return
        }
        if task.ScheduledTime == nil || time.Now().After(*task.ScheduledTime) {
            db.Model(task).Update("status", model.TaskStatusRunning)
            go engine.StartTask(task.ID)
        }
        c.JSON(http.StatusOK, task)
    })

    rg.GET("/tasks/:id/status", func(c *gin.Context) {
        idStr := c.Param("id")
        id, _ := strconv.Atoi(idStr)
        var task model.SendTask
        if err := db.First(&task, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"message": "未找到"})
            return
        }
        var total int64
        var sent int64
        var failed int64
        db.Model(&model.SendRecord{}).Where("task_id = ?", task.ID).Count(&total)
        db.Model(&model.SendRecord{}).Where("task_id = ? AND status = ?", task.ID, model.RecordStatusSent).Count(&sent)
        db.Model(&model.SendRecord{}).Where("task_id = ? AND status = ?", task.ID, model.RecordStatusFailed).Count(&failed)
        c.JSON(http.StatusOK, gin.H{"status": task.Status, "total": total, "sent": sent, "failed": failed})
    })
}