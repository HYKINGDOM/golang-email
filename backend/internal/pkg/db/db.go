package db

import (
    "database/sql"
    "fmt"
    "time"

    "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Options struct {
    Host string
    Port int
    User string
    Password string
    Name string
    SSLMode string
    TimeZone string
    MaxOpen int
    MaxIdle int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime time.Duration
}

func LoadOptions() *Options {
    viper.SetDefault("db.host", "192.168.5.8")
    viper.SetDefault("db.port", 5432)
    viper.SetDefault("db.user", "user_S5YD4z")
    viper.SetDefault("db.password", "password_bFkZC5")
    viper.SetDefault("db.name", "email_service")
    viper.SetDefault("db.sslmode", "disable")
    viper.SetDefault("db.timezone", "Asia/Shanghai")
    viper.SetDefault("db.pool.max_open", 50)
    viper.SetDefault("db.pool.max_idle", 25)
    viper.SetDefault("db.pool.conn_max_lifetime_minutes", 30)
    viper.SetDefault("db.pool.conn_max_idle_minutes", 10)
    return &Options{
        Host:     viper.GetString("db.host"),
        Port:     viper.GetInt("db.port"),
        User:     viper.GetString("db.user"),
        Password: viper.GetString("db.password"),
        Name:     viper.GetString("db.name"),
        SSLMode:  viper.GetString("db.sslmode"),
        TimeZone: viper.GetString("db.timezone"),
        MaxOpen:  viper.GetInt("db.pool.max_open"),
        MaxIdle:  viper.GetInt("db.pool.max_idle"),
        ConnMaxLifetime: time.Minute * time.Duration(viper.GetInt("db.pool.conn_max_lifetime_minutes")),
        ConnMaxIdleTime: time.Minute * time.Duration(viper.GetInt("db.pool.conn_max_idle_minutes")),
    }
}

func DSN(o *Options) string {
    return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", o.Host, o.User, o.Password, o.Name, o.Port, o.SSLMode, o.TimeZone)
}

func Connect() (*gorm.DB, *sql.DB, error) {
    opts := LoadOptions()
    dsn := DSN(opts)
    gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil { return nil, nil, err }
    sqlDB, err := gdb.DB()
    if err != nil { return nil, nil, err }
    sqlDB.SetMaxOpenConns(opts.MaxOpen)
    sqlDB.SetMaxIdleConns(opts.MaxIdle)
    sqlDB.SetConnMaxLifetime(opts.ConnMaxLifetime)
    sqlDB.SetConnMaxIdleTime(opts.ConnMaxIdleTime)
    if err := sqlDB.Ping(); err != nil { return nil, nil, err }
    logrus.WithFields(logrus.Fields{"host": opts.Host, "db": opts.Name}).Info("数据库连接成功")
    return gdb, sqlDB, nil
}
