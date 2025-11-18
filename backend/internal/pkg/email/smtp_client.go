package email

import (
    "crypto/tls"
    "fmt"
    "net"
    "net/smtp"
    "strings"
    "time"

    "github.com/sirupsen/logrus"
)

type SMTPConfig struct {
    Host     string
    Port     int
    Username string
    Password string
}

type SendRequest struct {
    From    string
    To      string
    Subject string
    Body    string
    IsHTML  bool
}

func buildMessage(req *SendRequest) []byte {
    headers := []string{
        fmt.Sprintf("From: %s", req.From),
        fmt.Sprintf("To: %s", req.To),
        fmt.Sprintf("Subject: %s", req.Subject),
        "MIME-Version: 1.0",
    }
    if req.IsHTML {
        headers = append(headers, "Content-Type: text/html; charset=\"UTF-8\"")
    } else {
        headers = append(headers, "Content-Type: text/plain; charset=\"UTF-8\"")
    }
    msg := strings.Join(headers, "\r\n") + "\r\n\r\n" + req.Body
    return []byte(msg)
}

func SendSMTP(cfg SMTPConfig, req *SendRequest) error {
    addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
    msg := buildMessage(req)
    auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

    timeout := 10 * time.Second
    dialer := &net.Dialer{Timeout: timeout}
    conn, err := dialer.Dial("tcp", addr)
    if err != nil {
        return err
    }
    tlsConn := tls.Client(conn, &tls.Config{ServerName: cfg.Host})
    c, err := smtp.NewClient(tlsConn, cfg.Host)
    if err != nil {
        return err
    }
    defer func() {
        if err := c.Quit(); err != nil {
            logrus.WithError(err).Warn("SMTP连接关闭失败")
        }
    }()
    if err := c.Auth(auth); err != nil {
        return err
    }
    if err := c.Mail(req.From); err != nil {
        return err
    }
    if err := c.Rcpt(req.To); err != nil {
        return err
    }
    w, err := c.Data()
    if err != nil {
        return err
    }
    if _, err := w.Write(msg); err != nil {
        return err
    }
    if err := w.Close(); err != nil {
        return err
    }
    return nil
}