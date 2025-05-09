package config

import (
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Log    LogConfig    `json:"log"`
	Server ServerConfig `json:"server"`
}

// LogConfig 日志配置
type LogConfig struct {
	Dir   string `json:"dir"`
	File  string `json:"file"`
	Level string `json:"level"`
}

func (log *LogConfig) GetLogDir() string {
	// 获取log dir的绝对路径
	if path.IsAbs(log.Dir) {
		return log.Dir
	}
	pwd, _ := os.Getwd()
	return path.Join(pwd, log.Dir)
}

func (log *LogConfig) GetLogPath() string {
	return path.Join(log.GetLogDir(), log.File)
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host           string        `json:"host"`
	Port           int           `json:"port"`
	ReadTimeout    time.Duration `json:"read_timeout"`
	WriteTimeout   time.Duration `json:"write_timeout"`
	MaxRequestSize int64         `json:"max_request_size"`
	GinMode        string        `json:"gin_mode"`
}

func DefaultConfig() *Config {
	return &Config{
		Log: LogConfig{
			Dir:   "logs",
			File:  "container-commit.log",
			Level: "info",
		},
		Server: ServerConfig{
			Host:           "0.0.0.0",
			Port:           8080,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxRequestSize: 10 * 1024 * 1024, // 10 MB
			GinMode:        gin.ReleaseMode,
		},
	}
}

func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		Dir:   "logs",
		File:  "container-commit.log",
		Level: "info",
	}
}
