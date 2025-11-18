package handler

import (
    "net/http"
    "runtime"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterDBDiagnostics(rg *gin.RouterGroup, db *gorm.DB) {
    rg.GET("/db/ping", func(c *gin.Context) {
        sqlDB, err := db.DB(); if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"message": "DB对象获取失败"}); return }
        start := time.Now()
        er := sqlDB.Ping()
        dur := time.Since(start)
        stats := sqlDB.Stats()
        c.JSON(http.StatusOK, gin.H{"ok": er==nil, "duration_ms": dur.Milliseconds(), "stats": stats})
    })

    rg.GET("/db/perf", func(c *gin.Context) {
        sqlDB, _ := db.DB()
        n := 100
        var durations []int64
        var errs int
        for i := 0; i < n; i++ {
            s := time.Now()
            er := sqlDB.Ping()
            if er != nil { errs++ }
            durations = append(durations, time.Since(s).Milliseconds())
        }
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        c.JSON(http.StatusOK, gin.H{"samples": n, "avg_ms": avg(durations), "max_ms": max(durations), "err_rate": float64(errs)/float64(n), "mem_alloc": m.Alloc})
    })

    rg.POST("/db/transaction", func(c *gin.Context) {
        start := time.Now()
        er := db.Transaction(func(tx *gorm.DB) error {
            type Tmp struct{ ID int }
            if er := tx.Exec("SELECT 1").Error; er != nil { return er }
            return nil
        })
        dur := time.Since(start)
        c.JSON(http.StatusOK, gin.H{"ok": er==nil, "duration_ms": dur.Milliseconds()})
    })

    rg.GET("/db/concurrency", func(c *gin.Context) {
        sqlDB, _ := db.DB()
        workers := 20
        var wg sync.WaitGroup
        errs := 0
        start := time.Now()
        for i := 0; i < workers; i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                for j := 0; j < 50; j++ {
                    if er := sqlDB.Ping(); er != nil { errs++ }
                }
            }()
        }
        wg.Wait()
        dur := time.Since(start)
        stats := sqlDB.Stats()
        c.JSON(http.StatusOK, gin.H{"workers": workers, "duration_ms": dur.Milliseconds(), "errors": errs, "stats": stats})
    })
}

func avg(arr []int64) int64 { var s int64; for _, v := range arr { s+=v }; if len(arr)==0 { return 0 }; return s/int64(len(arr)) }
func max(arr []int64) int64 { var m int64; for _, v := range arr { if v>m { m=v } }; return m }
