package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_best_final/pkg/config"
	"go.uber.org/zap/zapcore"
)

func TestReadConfig(t *testing.T) {
	c := config.InitConfig()
	data := `timeout: 2
loglevel: 5
logfile: access.txt
targetfile: target.csv`
	err := c.ReadConfig([]byte(data))
	assert.Nil(t, err)
	assert.Equal(t, time.Duration(2)*time.Second, c.Timeout)
	assert.Equal(t, zapcore.FatalLevel, c.LogLevel)
	assert.Equal(t, "access.txt", c.LogFile)
	assert.Equal(t, "target.csv", c.TargetFile)
}

func TestReadConfigError(t *testing.T) {
	c := config.InitConfig()
	data := `some text`
	err := c.ReadConfig([]byte(data))
	assert.NotNil(t, err)
}
