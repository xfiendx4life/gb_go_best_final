package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/xfiendx4life/gb_go_best_final/pkg/config"
	"github.com/xfiendx4life/gb_go_best_final/pkg/csvreader"
	"github.com/xfiendx4life/gb_go_best_final/pkg/logger"
	"go.uber.org/zap"
)

func openSource(path string, z *zap.SugaredLogger) (io.Reader, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		z.Fatalf("can't open source file: %s", err)
		return nil, fmt.Errorf("can't open source file: %s", err)
	}
	return bytes.NewReader(file), nil
}

func main() {
	var configFile string
	var query string
	flag.StringVar(&configFile, "config", "../../config.yaml", "use to set config destination")
	flag.StringVar(&query, "query", "", "use flag to set query")
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
	table := csvreader.NewData()
	ctx, cancel := context.WithTimeout(context.Background(), conf.Timeout)
	_ = cancel
	data, err := openSource(conf.TargetFile, z)
	if err != nil {
		return
	}
	if query == "" {
		z.Fatal("query not set")
	}
	table, err = table.ProceedFullTable(ctx, data, query, z)
	if err != nil {
		z.Fatalf("can't proceed query: %s", err)
	}
	fmt.Printf("%#v\n", table.GetTable())
}
