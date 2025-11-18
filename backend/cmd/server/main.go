package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	cfgpkg "github.com/example/golang-email/internal/config"
	handler "github.com/example/golang-email/internal/handler"
	mw "github.com/example/golang-email/internal/middleware"
	"github.com/example/golang-email/internal/model"
	dbpkg "github.com/example/golang-email/internal/pkg/db"
	"github.com/example/golang-email/internal/pkg/sse"
	"github.com/example/golang-email/internal/repository"
	"github.com/example/golang-email/internal/service"
)

func main() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "release")
	viper.SetDefault("db.dsn", "host=localhost user=postgres password=postgres dbname=email_service port=5432 sslmode=disable TimeZone=Asia/Shanghai")

	c := cfgpkg.Load()

	logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
	logrus.SetLevel(logrus.InfoLevel)

	if c.ServerMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	gdb, sqlDB, err := dbpkg.Connect()
	if err != nil {
		logrus.WithError(err).Fatal("数据库连接失败")
	}
	db := gdb
	_ = sqlDB
	if err := db.AutoMigrate(&model.EmailConfig{}, &model.EmailTemplate{}, &model.SendTask{}, &model.SendRecord{}, &model.EmailTracking{}, &model.LinkClick{}); err != nil {
		logrus.WithError(err).Fatal("数据库迁移失败")
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(mw.RequestLogger())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
		MaxAge:          12 * time.Hour,
	}))

	api := r.Group("/")
	handler.RegisterHealth(api, db)

	emailRepo := repository.NewEmailConfigRepository(db)
	emailSvc := service.NewEmailConfigService(emailRepo)
	handler.RegisterEmailConfig(api, emailSvc)

	broker := sse.NewBroker()
	engine := service.NewSendEngine(db, 4, broker)
	handler.RegisterTask(api, db, engine)

	tplRepo := repository.NewEmailTemplateRepository(db)
	handler.RegisterTemplate(api, tplRepo)
	handler.RegisterStats(api, db)
	handler.RegisterSSE(api, broker)
	handler.RegisterRecord(api, engine)
	handler.RegisterTracking(api, db)
	handler.RegisterTrackingStats(api, db)
	handler.RegisterDBDiagnostics(api, db)

	srv := &http.Server{
		Addr:              c.ServerAddr(),
		Handler:           r,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	logrus.WithFields(logrus.Fields{"addr": c.ServerAddr()}).Info("HTTP 服务启动")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.WithError(err).Fatal("HTTP 服务启动失败")
	}
}
