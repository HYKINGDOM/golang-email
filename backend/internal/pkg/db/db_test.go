package db

import (
    "encoding/json"
    "os"
    "runtime"
    "sync"
    "testing"
    "time"

    "gorm.io/gorm"
)

type Report struct {
    Timestamp string `json:"timestamp"`
    Basic struct { Ok bool `json:"ok"`; DurationMs int64 `json:"duration_ms"` } `json:"basic"`
    Perf struct { Samples int `json:"samples"`; AvgMs int64 `json:"avg_ms"`; MaxMs int64 `json:"max_ms"`; ErrRate float64 `json:"err_rate"` } `json:"perf"`
    Tx struct { Ok bool `json:"ok"`; DurationMs int64 `json:"duration_ms"` } `json:"tx"`
    Concurrency struct { Workers int `json:"workers"`; DurationMs int64 `json:"duration_ms"`; Errors int `json:"errors"` } `json:"concurrency"`
    Resources struct { MemAlloc uint64 `json:"mem_alloc"` } `json:"resources"`
}

func writeReport(t *testing.T, r *Report) {
    _ = os.MkdirAll("test-reports", 0755)
    f, err := os.Create("test-reports/db_report.json")
    if err != nil { t.Logf("report create err: %v", err); return }
    defer f.Close()
    enc := json.NewEncoder(f)
    _ = enc.Encode(r)
}

func TestDatabaseConnectivity(t *testing.T) {
    gdb, sqlDB, err := Connect()
    rep := &Report{Timestamp: time.Now().Format(time.RFC3339)}
    if err != nil {
        rep.Basic.Ok = false
        rep.Basic.DurationMs = 0
        writeReport(t, rep)
        t.Skip("数据库不可达: "+err.Error())
        return
    }
    s := time.Now()
    if er := sqlDB.Ping(); er != nil { t.Fatalf("ping失败: %v", er) }
    rep.Basic.Ok = true
    rep.Basic.DurationMs = time.Since(s).Milliseconds()

    n := 100
    var durs []int64
    errs := 0
    for i := 0; i < n; i++ {
        st := time.Now()
        if er := sqlDB.Ping(); er != nil { errs++ }
        durs = append(durs, time.Since(st).Milliseconds())
    }
    rep.Perf.Samples = n
    rep.Perf.AvgMs = avgInt64(durs)
    rep.Perf.MaxMs = maxInt64(durs)
    rep.Perf.ErrRate = float64(errs)/float64(n)

    st := time.Now()
    er := gdb.Transaction(func(tx *gorm.DB) error { return tx.Exec("SELECT 1").Error })
    rep.Tx.Ok = er == nil
    rep.Tx.DurationMs = time.Since(st).Milliseconds()

    workers := 20
    start := time.Now()
    var wg sync.WaitGroup
    concErrs := 0
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() { defer wg.Done(); for j := 0; j < 50; j++ { if er := sqlDB.Ping(); er != nil { concErrs++ } } }()
    }
    wg.Wait()
    rep.Concurrency.Workers = workers
    rep.Concurrency.DurationMs = time.Since(start).Milliseconds()
    rep.Concurrency.Errors = concErrs

    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    rep.Resources.MemAlloc = m.Alloc
    writeReport(t, rep)
}

func avgInt64(arr []int64) int64 { var s int64; for _, v := range arr { s+=v }; if len(arr)==0 { return 0 }; return s/int64(len(arr)) }
func maxInt64(arr []int64) int64 { var m int64; for _, v := range arr { if v>m { m=v } }; return m }
