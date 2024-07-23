package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/config"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/pb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
)

func main() {
	var configFile string
	pflag.StringVar(&configFile, "config", "./configs/configc.toml", "Path to configuration file")
	pflag.Parse()

	formater := new(logrus.TextFormatter)
	formater.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(formater)

	cfg := config.NewConfig()
	if err := cfg.ReadConfig(configFile); err != nil {
		logrus.Fatal("read config failed: ", err)
		return
	}
	conn, err := grpc.NewClient(cfg.GRPC.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatal(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	client := pb.NewServiceStatClient(conn)
	for {
		req := &pb.ServiceStatRequest{
			Interval:  durationpb.New(cfg.Stat.IntervalStat),
			Avgwindow: durationpb.New(cfg.Stat.AvgWindow),
		}
		stream, err := client.GetFullStatStream(ctx, req)
		if err != nil {
			logrus.Error(err)
			break
		}
		stat, err := stream.Recv()
		if err != nil {
			logrus.Error(err)
			break
		}
		if stat == nil {
			logrus.Printf("Error, full Stat is empty")
			continue
		}
		printStats(*stat.FullStat)
	}
}

func Format(n uint64) string {
	in := strconv.FormatUint(n, 10)
	numOfDigits := len(in)
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = '\''
		}
	}
}

func printStats(stat pb.FullStat) {
	if stat.GetTime() != nil {
		log.Printf("-----SYSTEM STATISTIC %s-----", stat.Time.AsTime().Format("2006-01-02 15:04:05"))
	}
	if stat.GetLoadStat() != nil {
		log.Printf(" - LoadAvg\t\t\t\t\t%f", stat.LoadStat.LoadAvg)
	}
	if stat.GetCpuStat() != nil {
		log.Printf(" - CPU")
		log.Printf(" 	- UserPr\t\t\t\t%f", stat.CpuStat.UserPr)
		log.Printf(" 	- NicePr\t\t\t\t%f", stat.CpuStat.NicePr)
		log.Printf(" 	- SystemPr\t\t\t\t%f", stat.CpuStat.SystemPr)
		log.Printf(" 	- IdlePr\t\t\t\t%f", stat.CpuStat.IdlePr)
		log.Printf(" 	- IowaitPr\t\t\t\t%f", stat.CpuStat.IowaitPr)
		log.Printf(" 	- StealPr\t\t\t\t%f", stat.CpuStat.StealPr)
	}
	if stat.GetMemStat() != nil {
		log.Printf(" - Memory")
		log.Printf("     - MemTotal\t\t\t\t%s kb", Format(stat.MemStat.MemTotal))
		log.Printf("     - MemFree\t\t\t\t%s kb", Format(stat.MemStat.MemFree))
		log.Printf("     - Buffers\t\t\t\t%s kb", Format(stat.MemStat.Buffers))
		log.Printf("     - Cached\t\t\t\t%s kb", Format(stat.MemStat.Cached))
		log.Printf("     - SwapCached\t\t\t\t%s kb", Format(stat.MemStat.SwapCached))
		log.Printf("     - SwapTotal\t\t\t\t%s kb", Format(stat.MemStat.SwapTotal))
		log.Printf("     - SwapFree\t\t\t\t%s kb", Format(stat.MemStat.SwapFree))
		log.Printf("     - Active\t\t\t\t%s kb", Format(stat.MemStat.Active))
		log.Printf("     - Inactive\t\t\t\t%s kb", Format(stat.MemStat.Inactive))
		log.Printf("     - VmallocUsed\t\t\t\t%s kb", Format(stat.MemStat.VmallocUsed))
	}
	if stat.GetDiskStat() != nil {
		if stat.GetDiskStat().GetDiskLoad() != nil {
			log.Printf(" - DiskLoad")
			for disk, diskload := range stat.DiskStat.DiskLoad {
				log.Printf("     - %s: ", disk)
				log.Printf("     	- TransfersPerSec\t\t%f", diskload.Tps)
				log.Printf("     	- WritedPerSec\t\t\t%f", diskload.Rps)
				log.Printf("     	- ReadedPerSec\t\t\t%f", diskload.Wps)
			}
		}
		if stat.GetDiskStat().GetDiskSpace() != nil {
			log.Printf(" - DiskSpace")
			for disk, diskspace := range stat.DiskStat.DiskSpace {
				log.Printf("     - %s: ", disk)
				log.Printf("     	- Total\t\t\t\t%d", diskspace.Total)
				log.Printf("     	- Used\t\t\t\t%d", diskspace.Used)
				log.Printf("     	- Available\t\t\t%d", diskspace.Available)
				log.Printf("     	- Inodes\t\t\t%d", diskspace.Inodes)
				log.Printf("     	- PercentUsed\t\t\t%d", diskspace.PercentUsed)
			}
		}
	}
	if stat.GetNetworkStat() != nil {
		log.Printf(" - Network Sockets States")
		if stat.GetNetworkStat().SocketStates != nil {
			for state, cnt := range stat.NetworkStat.SocketStates {
				log.Printf("     - %s:\t\t\t%d", state, cnt)
			}
		}
		if stat.GetNetworkStat().ListenSockets != nil {
			log.Printf(" - Network Listen Sockets List")
			for _, sock := range stat.NetworkStat.ListenSockets {
				log.Printf("     - %s\t\t%s\t\t%s\t\t%s", sock.Protocol, sock.LocalAddress, sock.PeerAddress, sock.Process)
			}
		}
	}
	log.Printf("-----END OF STATISTIC-----")
}
