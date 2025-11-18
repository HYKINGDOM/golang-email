package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

func RequestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        method := c.Request.Method

        c.Next()

        latency := time.Since(start)
        status := c.Writer.Status()
        logrus.WithFields(logrus.Fields{
            "path":    path,
            "method":  method,
            "status":  status,
            "latency": latency.String(),
        }).Info("请求追踪")
    }
}