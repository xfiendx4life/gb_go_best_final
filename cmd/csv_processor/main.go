package main

import (
	"flag"
	"log"
	"os"

	"github.com/xfiendx4life/gb_go_best_final/pkg/config"
	"github.com/xfiendx4life/gb_go_best_final/pkg/logger"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "../../config.yaml", "use to set config destination")
	flag.Parse()
	conf := config.InitConfig()
	confFile, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("can't read config file %s", err)
	}
	err = conf.ReadConfig(confFile)
	if err != nil {
		log.Fatalf("can't parse config file: %s", err)
	}
	z := logger.InitLogger(&conf.LogLevel, conf.LogFile)
	z.Info("logger initiated")
}
