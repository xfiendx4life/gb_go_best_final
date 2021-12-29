package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/xfiendx4life/gb_go_best_final/internal/config"
	"github.com/xfiendx4life/gb_go_best_final/internal/logger"
	"github.com/xfiendx4life/gb_go_best_final/pkg/consolewriter"
	"github.com/xfiendx4life/gb_go_best_final/pkg/csvreader"
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

func readQuery() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "../../config.yaml", "use to set config destination")
	flag.Parse()
	fmt.Println("Write your query")
	query := readQuery()
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
	ctx, timeoutCancel := context.WithTimeout(context.Background(), conf.Timeout)
	_ = timeoutCancel
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	_ = cancel
	data, err := openSource(conf.TargetFile, z)
	if err != nil {
		return
	}
	if query == "" {
		z.Fatal("query not set")
	}
	resChan := make(chan csvreader.Table)
	errChan := make(chan error)
	go table.ProceedFullTable(ctx, data, query, z, resChan, errChan)
	select {
	case <-ctx.Done():
		z.Error("Done with context")
		cancel()
	case err := <-errChan:
		z.Errorf("error while processing table: %s", err)
	case res := <-resChan:
		consolewriter.NewConsoleWriter().Write(os.Stdout, res)
	}

}
