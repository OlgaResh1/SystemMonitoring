package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/aggregator"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/app"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/config"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/logger"
	metrics "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/metrics"
	internalgrpc "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/server/grpc"
	"github.com/spf13/pflag"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

func printVersion() {
	if err := json.NewEncoder(os.Stdout).Encode(struct {
		Release   string
		BuildDate string
		GitHash   string
	}{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
	}
}

func main() {
	var configFile string
	pflag.StringVar(&configFile, "config", "./configs/configd.toml", "Path to configuration file")
	pflag.Parse()
	if pflag.Arg(0) == "version" {
		printVersion()
		return
	}
	cfg := config.NewConfig()
	if err := cfg.ReadConfig(configFile); err != nil {
		log.Fatal("read config failed: ", err)
		return
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	logg := logger.New(cfg.Logger.Level, cfg.Logger.Format, cfg.Logger.AddSource)

	avg := aggregator.New(logg, cfg.Stat.IntervalStat)
	metrics := metrics.New(cfg, avg, *logg)
	avg.RegisterMetric(metrics)
	statapp := app.New(logg, metrics)

	serverGrpc := internalgrpc.NewServer(cfg, logg, statapp)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := avg.Run(ctx); err != nil {
			logg.Error("failed to start aggregator: "+err.Error(), "source", "aggregator")
			cancel()
		} else {
			logg.Info("Aggregator stopped ok", "source", "aggregator")
		}
	}()
	go func() {
		if err := serverGrpc.Start(); err != nil {
			logg.Error("failed to start GRPC-server: "+err.Error(), "source", "grpc")
			cancel()
		}
	}()
	go func() {
		defer wg.Done()
		<-ctx.Done()

		if err := serverGrpc.Stop(); err != nil {
			logg.Error("failed to stop GRPC-server: "+err.Error(), "source", "grpc")
		} else {
			logg.Info("GRPC-server stopped ok", "source", "grpc")
		}
	}()
	wg.Wait()
}
