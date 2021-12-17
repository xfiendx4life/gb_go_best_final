package config

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	// _ "gopkg.in/yaml.v3"
)

type Config interface {
	// read config from file with path
	ReadConfig(data []byte, z *zap.SugaredLogger) (err error)
}

type ConfYML struct {
	Timeout    time.Duration `yaml:"timeout"`
	LogLevel   zapcore.Level `yaml:"loglevel"`
	LogFile    string        `yaml:"logfile"`
	TargetFile string        `yaml:"targetfile"`
}
