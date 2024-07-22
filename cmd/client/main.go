package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/config"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/pb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
)

func main() {
	var configFile string
	pflag.StringVar(&configFile, "config", "./configs/configc.toml", "Path to configuration file")
	pflag.Parse()

	cfg := config.NewConfig()
	if err := cfg.ReadConfig(configFile); err != nil {
		log.Fatal("read config failed: ", err)
		return
	}
	conn, err := grpc.NewClient(cfg.GRPC.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	client := pb.NewServiceStatClient(conn)
	for {
		req, err := getRequest()
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}

		stream, err := client.GetFullStatStream(ctx, req)
		if err != nil {
			log.Error(err)
		}
		stat, err := stream.Recv()
		if err != nil {
			log.Error(err)
		}
		log.Print("Get GRPC response")
		fullStat := stat.FullStat
		if fullStat == nil {
			log.Printf("Full Stat is empty")
			continue
		}
		printStats(*fullStat)
	}
}

func printStats(stat pb.FullStat) {
	log.Printf("LoadAvg %f", stat.LoadStat.LoadAvg)
}

func getRequest() (*pb.ServiceStatRequest, error) {
	interval, err := time.ParseDuration("30s")
	if err != nil {
		return nil, err
	}
	avgWindow, err := time.ParseDuration("1m")
	if err != nil {
		return nil, err
	}
	return &pb.ServiceStatRequest{
		Interval:  durationpb.New(interval),
		Avgwindow: durationpb.New(avgWindow),
	}, nil
}
