package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort int
	ServerMode string
	DBDSN      string
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	_ = viper.ReadInConfig()
	viper.SetEnvPrefix("EMAIL")
	viper.AutomaticEnv()

	viper.SetDefault("redis.address", "127.0.0.1:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("sse.base_path", "/sse")
	viper.SetDefault("sse.retry_ms", 2000)
	viper.SetDefault("tracking.pixel_path", "/t/open")
	viper.SetDefault("tracking.click_path", "/t/click")
	viper.SetDefault("tracking.enabled_default", true)
	viper.SetDefault("email.rate_limit_interval_ms", 2000)
	viper.SetDefault("email.jitter_ms", 1500)
	viper.SetDefault("email.retry_times", 3)
	viper.SetDefault("email.breaker_threshold", 3)
	viper.SetDefault("email.breaker_cooldown_minutes", 2)

	return &Config{
		ServerPort: viper.GetInt("server.port"),
		ServerMode: viper.GetString("server.mode"),
		DBDSN:      viper.GetString("db.dsn"),
	}
}

func (c *Config) ServerAddr() string {
	return fmt.Sprintf(":%d", c.ServerPort)
}
