package service

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"regexp"

	"github.com/example/golang-email/internal/model"
	"github.com/example/golang-email/internal/pkg/crypto"
	mail "github.com/example/golang-email/internal/pkg/email"
	"github.com/example/golang-email/internal/pkg/sse"
	"github.com/spf13/viper"
)

type RateLimiter struct {
	mu       sync.Mutex
	interval time.Duration
	last     time.Time
}

func NewRateLimiter(interval time.Duration) *RateLimiter {
	return &RateLimiter{interval: interval}
}

func (r *RateLimiter) Wait(jitter time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	wait := r.last.Add(r.interval + time.Duration(rand.Int63n(int64(jitter)))).Sub(now)
	if wait > 0 {
		time.Sleep(wait)
	}
	r.last = time.Now()
}

type AccountBreaker struct {
	mu        sync.Mutex
	failures  int
	threshold int
	openUntil time.Time
	coolDown  time.Duration
}

func NewAccountBreaker(threshold int, coolDown time.Duration) *AccountBreaker {
	return &AccountBreaker{threshold: threshold, coolDown: coolDown}
}

func (b *AccountBreaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return time.Now().After(b.openUntil)
}

func (b *AccountBreaker) Success() {
	b.mu.Lock()
	b.failures = 0
	b.mu.Unlock()
}

func (b *AccountBreaker) Fail() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.threshold {
		b.openUntil = time.Now().Add(b.coolDown)
		b.failures = 0
	}
}

type SendEngine struct {
	db         *gorm.DB
	workerSize int
	accountRL  map[uint]*RateLimiter
	accountCB  map[uint]*AccountBreaker
	jitter     time.Duration
	broker     *sse.Broker
}

func NewSendEngine(db *gorm.DB, workerSize int, broker *sse.Broker) *SendEngine {
	return &SendEngine{
		db:         db,
		workerSize: workerSize,
		accountRL:  make(map[uint]*RateLimiter),
		accountCB:  make(map[uint]*AccountBreaker),
		jitter:     time.Duration(viper.GetInt("email.jitter_ms")) * time.Millisecond,
		broker:     broker,
	}
}

func (e *SendEngine) ensureAccountControl(id uint) {
	if _, ok := e.accountRL[id]; !ok {
		e.accountRL[id] = NewRateLimiter(time.Duration(viper.GetInt("email.rate_limit_interval_ms")) * time.Millisecond)
	}
	if _, ok := e.accountCB[id]; !ok {
		e.accountCB[id] = NewAccountBreaker(viper.GetInt("email.breaker_threshold"), time.Duration(viper.GetInt("email.breaker_cooldown_minutes"))*time.Minute)
	}
}

func (e *SendEngine) StartTask(taskID uint) {
	var task model.SendTask
	if err := e.db.First(&task, taskID).Error; err != nil {
		logrus.WithError(err).Error("任务加载失败")
		return
	}
	var tpl model.EmailTemplate
	if err := e.db.First(&tpl, task.TemplateID).Error; err != nil {
		logrus.WithError(err).Error("模板加载失败")
		return
	}
	var senders []model.EmailConfig
	if err := e.db.Where("is_active = ?", true).Find(&senders).Error; err != nil {
		logrus.WithError(err).Error("发件账号加载失败")
		return
	}
	recipients := parseCSV(task.RecipientList)

	ch := make(chan string)
	wg := &sync.WaitGroup{}
	for i := 0; i < e.workerSize; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			idx := 0
			for to := range ch {
				sender := senders[idx%len(senders)]
				idx++
				e.ensureAccountControl(sender.ID)
				if !e.accountCB[sender.ID].Allow() {
					continue
				}
				e.accountRL[sender.ID].Wait(e.jitter)
				pass, err := crypto.DecryptString(sender.Password)
				if err != nil {
					logrus.WithError(err).Error("密码解密失败")
					e.accountCB[sender.ID].Fail()
					continue
				}
				preID := e.record(task.ID, sender.Username, to, string(model.RecordStatusPending), "", 0)
				body := tpl.Content
				if tpl.TrackingEnabled && tpl.IsRichText {
					body = e.injectTracking(body, preID)
				}
				req := &mail.SendRequest{From: sender.Username, To: to, Subject: tpl.Subject, Body: body, IsHTML: tpl.IsRichText}
				cfg := mail.SMTPConfig{Host: sender.Host, Port: sender.Port, Username: sender.Username, Password: pass}
				var sendErr error
				for retry := 0; retry < viper.GetInt("email.retry_times"); retry++ {
					sendErr = mail.SendSMTP(cfg, req)
					if sendErr == nil {
						rid := e.record(task.ID, sender.Username, to, "sent", "", retry)
						e.publish(task.ID, "progress", map[string]any{"id": rid, "recipient": to, "status": "sent"})
						e.accountCB[sender.ID].Success()
						break
					}
					rid := e.record(task.ID, sender.Username, to, "failed", sendErr.Error(), retry+1)
					e.publish(task.ID, "progress", map[string]any{"id": rid, "recipient": to, "status": "failed", "error": sendErr.Error()})
					time.Sleep(time.Duration(500+rand.Intn(500)) * time.Millisecond)
				}
				if sendErr != nil {
					e.accountCB[sender.ID].Fail()
				}
			}
		}(i)
	}

	for _, r := range recipients {
		ch <- r
	}
	close(ch)
	wg.Wait()
	e.db.Model(&task).Update("status", model.TaskStatusFinished)
	e.publish(task.ID, "finished", nil)
}

func (e *SendEngine) record(taskID uint, sender string, recipient string, status string, errMsg string, retry int) uint {
	now := time.Now()
	rec := &model.SendRecord{TaskID: taskID, SenderEmail: sender, RecipientEmail: recipient, SendTime: &now, Status: model.SendRecordStatus(status), ErrorMessage: errMsg, RetryCount: retry}
	if er := e.db.Create(rec).Error; er != nil {
		logrus.WithError(er).Error("发送记录写入失败")
	}
	return rec.ID
}

func (e *SendEngine) publish(taskID uint, typ string, data any) {
	if e.broker != nil {
		e.broker.Publish(taskID, sse.Event{TaskID: taskID, Type: typ, Data: data})
	}
}

func (e *SendEngine) RetryRecord(recordID uint) error {
	var rec model.SendRecord
	if err := e.db.First(&rec, recordID).Error; err != nil {
		return err
	}
	var task model.SendTask
	if err := e.db.First(&task, rec.TaskID).Error; err != nil {
		return err
	}
	var tpl model.EmailTemplate
	if err := e.db.First(&tpl, task.TemplateID).Error; err != nil {
		return err
	}
	var senders []model.EmailConfig
	if err := e.db.Where("is_active = ?", true).Find(&senders).Error; err != nil {
		return err
	}
	if len(senders) == 0 {
		return fmt.Errorf("无可用发件人")
	}
	sender := senders[int(time.Now().UnixNano())%len(senders)]
	e.ensureAccountControl(sender.ID)
	if !e.accountCB[sender.ID].Allow() {
		return fmt.Errorf("账号熔断中")
	}
	e.accountRL[sender.ID].Wait(e.jitter)
	pass, err := crypto.DecryptString(sender.Password)
	if err != nil {
		return err
	}
	req := &mail.SendRequest{From: sender.Username, To: rec.RecipientEmail, Subject: tpl.Subject, Body: tpl.Content, IsHTML: tpl.IsRichText}
	cfg := mail.SMTPConfig{Host: sender.Host, Port: sender.Port, Username: sender.Username, Password: pass}
	var sendErr error
	for retry := 0; retry < 3; retry++ {
		sendErr = mail.SendSMTP(cfg, req)
		if sendErr == nil {
			e.db.Model(&rec).Updates(map[string]any{"status": model.RecordStatusSent, "error_message": "", "retry_count": rec.RetryCount + retry + 1})
			e.publish(task.ID, "progress", map[string]any{"recipient": rec.RecipientEmail, "status": "sent"})
			e.accountCB[sender.ID].Success()
			return nil
		}
		e.db.Model(&rec).Updates(map[string]any{"status": model.RecordStatusFailed, "error_message": sendErr.Error(), "retry_count": rec.RetryCount + retry + 1})
		e.publish(task.ID, "progress", map[string]any{"recipient": rec.RecipientEmail, "status": "failed", "error": sendErr.Error()})
		time.Sleep(time.Duration(500+rand.Intn(500)) * time.Millisecond)
	}
	e.accountCB[sender.ID].Fail()
	return sendErr
}

func parseCSV(s string) []string {
	parts := strings.Split(s, ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			res = append(res, t)
		}
	}
	return res
}

func (e *SendEngine) injectTracking(html string, recordID uint) string {
	base := viper.GetString("server.base_url")
	if base == "" {
		base = "http://localhost:8080"
	}
	pixel := fmt.Sprintf("<img src=\"%s/t/open/%d.png\" style=\"display:none\" />", base, recordID)
	out := strings.ReplaceAll(html, "</body>", pixel+"</body>")
	out = rewriteLinks(out, base, recordID)
	return out
}

func rewriteLinks(html string, base string, rid uint) string {
	re := regexp.MustCompile(`href=\"(https?://[^\"]+)\"`)
	return re.ReplaceAllStringFunc(html, func(m string) string {
		u := re.FindStringSubmatch(m)
		if len(u) > 1 {
			enc := base64.StdEncoding.EncodeToString([]byte(u[1]))
			return fmt.Sprintf("href=\"%s/t/click?rid=%d&url=%s\"", base, rid, enc)
		}
		return m
	})
}
