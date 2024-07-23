package metrics

import (
	"fmt"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/aggregator"
	cpustat "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/cpu"
	diskstat "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/disk"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/loadavg"
	memstat "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/mem"
	networkstat "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/network"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/pb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (m Metrics) AggregatedStatPB(queueType int, avgWindow time.Duration) (any, error) {
	stat, err := m.AggregatedStat(queueType, avgWindow)
	if err != nil {
		return nil, err
	}
	switch aggregator.StatQueueType(queueType) {
	case LoadAvgStatType:
		s, ok := stat.(*loadavg.LoadAvg)
		if !ok {
			return nil, fmt.Errorf("invalid response type: expected *loadavg.LoadAvg, got %T", stat)
		}
		return &pb.LoadAvgStat{
			LoadAvg: s.LoadAvg1,
		}, nil
	case CPUStatType:
		s, ok := stat.(*cpustat.CPUStat)
		if !ok {
			return nil, fmt.Errorf("invalid response type: expected *cpustat.CPUStat, got %T", stat)
		}
		return &pb.CpuStat{
			UserPr:   s.User,
			NicePr:   s.Nice,
			SystemPr: s.System,
			IdlePr:   s.Idle,
			IowaitPr: s.Iowait,
			StealPr:  s.Steal,
		}, nil

	case MemStatType:
		s, ok := stat.(*memstat.MemStat)
		if !ok {
			return nil, fmt.Errorf("invalid response type: expected *memstat.MemStat, got %T", stat)
		}
		return &pb.MemStat{
			MemTotal:    s.MemTotal,
			MemFree:     s.MemFree,
			Buffers:     s.Buffers,
			Cached:      s.Cached,
			SwapCached:  s.SwapCached,
			SwapTotal:   s.SwapTotal,
			SwapFree:    s.SwapFree,
			Active:      s.Active,
			Inactive:    s.Inactive,
			VmallocUsed: s.VmallocUsed,
			Mapped:      s.Mapped,
		}, nil
	case DiskStatType:
		s, ok := stat.(*diskstat.DiskStat)
		if !ok {
			return nil, fmt.Errorf("invalid response type: expected *diskstat.DiskStat, got %T", stat)
		}
		pbStat := &pb.DiskStat{
			DiskLoad:  make(map[string]*pb.DiskStat_DiskLoad),
			DiskSpace: make(map[string]*pb.DiskStat_DiskSpace),
		}
		for disk, load := range s.DiskLoad {
			pbStat.DiskLoad[disk] = &pb.DiskStat_DiskLoad{
				Tps: load.TransfersPerSec,
				Wps: load.WritedPerSec,
				Rps: load.ReadedPerSec,
			}
		}
		for disk, load := range s.DiskSpace {
			pbStat.DiskSpace[disk] = &pb.DiskStat_DiskSpace{
				Total:       load.Total,
				Used:        load.Used,
				Available:   load.Available,
				Inodes:      load.INodes,
				PercentUsed: load.UsedPercent,
			}
		}
		return pbStat, nil
	case NetworkStatType:
		s, ok := stat.(*networkstat.NetworkStat)
		if !ok {
			return nil, fmt.Errorf("invalid response type: expected *networkstat.NetworkStat, got %T", stat)
		}
		pbStat := &pb.NetworkStat{
			SocketStates: make(map[string]int32),
		}
		for _, sockets := range s.ListenSockets {
			pbStat.ListenSockets = append(pbStat.ListenSockets, &pb.NetworkStat_ListenSocket{
				LocalAddress: sockets.LocalAddress,
				Process:      sockets.Process,
				Protocol:     sockets.Protocol,
				PeerAddress:  sockets.PeerAddress,
			})
		}
		for state, cnt := range s.SocketStates {
			pbStat.SocketStates[state] += cnt
		}
		return pbStat, nil
	}

	return nil, nil
}

func (m Metrics) AggregatedFullStatPB(avgWindow time.Duration) (any, error) {
	load, err := m.AggregatedStatPB(int(LoadAvgStatType), avgWindow)
	if err != nil {
		return nil, err
	}
	cpu, err := m.AggregatedStatPB(int(CPUStatType), avgWindow)
	if err != nil {
		return nil, err
	}
	mem, err := m.AggregatedStatPB(int(MemStatType), avgWindow)
	if err != nil {
		return nil, err
	}
	disk, err := m.AggregatedStatPB(int(DiskStatType), avgWindow)
	if err != nil {
		return nil, err
	}
	network, err := m.AggregatedStatPB(int(NetworkStatType), avgWindow)
	if err != nil {
		return nil, err
	}
	return &pb.FullStat{
		Time:        timestamppb.Now(),
		Interval:    durationpb.New(avgWindow),
		LoadStat:    load.(*pb.LoadAvgStat),
		CpuStat:     cpu.(*pb.CpuStat),
		MemStat:     mem.(*pb.MemStat),
		DiskStat:    disk.(*pb.DiskStat),
		NetworkStat: network.(*pb.NetworkStat),
	}, nil
}
