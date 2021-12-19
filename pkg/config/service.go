package config

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

func InitConfig() (conf *ConfYML) {
	return &ConfYML{}
}

func ReadFromFile(path string) (data []byte, err error) {
	data, err = os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("can't open config file: %s", err)
	}
	return data, err
}

func (conf *ConfYML) ReadConfig(data []byte) (err error) {
	fakestruct := struct {
		Timeout    int    `yaml:"timeout"`
		LogLevel   int8   `yaml:"loglevel"`
		LogFile    string `yaml:"logfile"`
		TargetFile string `yaml:"targetfile"`
		Separator  string `yaml:"separator"`
	}{}
	err = yaml.Unmarshal(data, &fakestruct)
	if err != nil {
		return fmt.Errorf("can't unmarshall data: %s", err)
	}
	conf.LogFile = fakestruct.LogFile
	conf.TargetFile = fakestruct.TargetFile
	conf.Timeout = time.Duration(fakestruct.Timeout) * time.Second
	conf.LogLevel = zapcore.Level(fakestruct.LogLevel)
	conf.Separator = fakestruct.Separator
	return nil
}

func (conf *ConfYML) GetTimeout() (timeout time.Duration) {
	return conf.Timeout
}

func (conf *ConfYML) GetLogLevel() (LogLevel zapcore.Level) {
	return conf.LogLevel
}

func (conf *ConfYML) GetTargetFile() string {
	return conf.TargetFile
}

func (conf *ConfYML) GetSeparator() rune {
	return rune(conf.Separator[0])
}
