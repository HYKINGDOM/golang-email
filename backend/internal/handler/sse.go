package handler

import (
    "strconv"
    "io"

    "github.com/gin-gonic/gin"

    "github.com/example/golang-email/internal/pkg/sse"
)

func RegisterSSE(rg *gin.RouterGroup, broker *sse.Broker) {
    rg.GET("/sse/tasks/:id", func(c *gin.Context) {
        c.Header("Content-Type", "text/event-stream")
        c.Header("Cache-Control", "no-cache")
        c.Header("Connection", "keep-alive")
        id, _ := strconv.Atoi(c.Param("id"))
        ch := broker.Subscribe(uint(id))
        defer broker.Unsubscribe(uint(id), ch)
        c.Stream(func(w io.Writer) bool {
            evt, ok := <-ch
            if !ok { return false }
            c.SSEvent("message", evt)
            return true
        })
    })
}