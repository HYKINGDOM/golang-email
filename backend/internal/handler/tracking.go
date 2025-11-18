package handler

import (
    "encoding/base64"
    "net/http"
    "net/url"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "github.com/example/golang-email/internal/model"
)

func RegisterTracking(rg *gin.RouterGroup, db *gorm.DB) {
    rg.GET("/t/open/:rid.png", func(c *gin.Context) {
        rid := c.Param("rid")
        var rec model.SendRecord
        if err := db.First(&rec, rid).Error; err == nil {
            var tr model.EmailTracking
            if err := db.Where("record_id = ?", rec.ID).First(&tr).Error; err != nil {
                tr = model.EmailTracking{RecordID: rec.ID, OpenCount: 0}
                db.Create(&tr)
            }
            now := time.Now()
            tr.OpenCount++
            if tr.OpenTime == nil { tr.OpenTime = &now }
            tr.LastOpenTime = &now
            db.Save(&tr)
        }
        c.Header("Content-Type", "image/png")
        c.Header("Cache-Control", "no-store")
        c.Writer.WriteHeader(http.StatusOK)
        c.Writer.Write([]byte{137,80,78,71,13,10,26,10,0,0,0,13,73,72,68,82,0,0,0,1,0,0,0,1,8,6,0,0,0,31,21,196,137,0,0,0,12,73,68,65,84,120,156,99,96,0,0,0,2,0,1,226,33,157,167,0,0,0,0,73,69,78,68,174,66,96,130})
    })

    rg.GET("/t/click", func(c *gin.Context) {
        rid := c.Query("rid")
        raw := c.Query("url")
        uDec, _ := base64.StdEncoding.DecodeString(raw)
        now := time.Now()
        var rec model.SendRecord
        if err := db.First(&rec, rid).Error; err == nil {
            db.Create(&model.LinkClick{RecordID: rec.ID, URL: string(uDec), ClickTime: now})
        }
        if _, err := url.Parse(string(uDec)); err == nil {
            c.Redirect(http.StatusFound, string(uDec))
            return
        }
        c.String(http.StatusBadRequest, "invalid url")
    })
}