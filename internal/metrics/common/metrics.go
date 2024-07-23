package metrics

import (
	"errors"
	"fmt"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/aggregator"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/config"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/logger"
	cpu "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/cpu"
	disk "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/disk"
	loadavg "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/loadavg"
	mem "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/mem"
	network "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/network"
)

const (
	LoadAvgStatType aggregator.StatQueueType = iota
	CPUStatType
	MemStatType
	DiskStatType
	NetworkStatType
)

var ErrUnknownQueueType = errors.New("unknown queue type")

type Metrics struct {
	avg            aggregator.Aggregator
	logger         logger.Logger
	loadavgEnabled bool
	cpuEnabled     bool
	memoryEnabled  bool
	diskEnabled    bool
	networkEnabled bool
}
type FullStat struct {
	timestamp   time.Time
	interval    time.Duration
	loadStat    loadavg.LoadAvg
	cpuStat     cpu.CPUStat
	memStat     mem.MemStat
	diskStat    disk.DiskStat
	networkStat network.NetworkStat
}

func New(cfg config.Config, avg aggregator.Aggregator, logger logger.Logger) *Metrics {
	m := &Metrics{
		avg:    avg,
		logger: logger,
	}
	m.setup(cfg)
	return m
}

func (m Metrics) CurrentStat(queueType int) (any, error) {
	switch aggregator.StatQueueType(queueType) {
	case LoadAvgStatType:
		return loadavg.CurrentStat()
	case CPUStatType:
		return cpu.CurrentStat()
	case MemStatType:
		return mem.CurrentStat()
	case DiskStatType:
		return disk.CurrentStat()
	case NetworkStatType:
		return network.CurrentStat()
	}
	return nil, ErrUnknownQueueType
}

func (m Metrics) AggregatedStat(queueType int, avgWindow time.Duration) (any, error) {
	stat, err := m.avg.StatByInterval(aggregator.StatQueueType(queueType), avgWindow)
	if err != nil {
		return nil, err
	}

	switch aggregator.StatQueueType(queueType) {
	case LoadAvgStatType:
		return loadavg.AggregatedStats(stat)
	case CPUStatType:
		return cpu.AggregatedStats(stat)
	case MemStatType:
		return mem.AggregatedStats(stat)
	case DiskStatType:
		return disk.AggregatedStats(stat)
	case NetworkStatType:
		return network.AggregatedStats(stat)
	}
	return nil, ErrUnknownQueueType
}

func (m Metrics) AggregatedFullStat(avgWindow time.Duration) (any, error) {
	fullstat := &FullStat{}
	fullstat.timestamp = time.Now()
	fullstat.interval = avgWindow

	if m.loadavgEnabled {
		stat, err := m.AggregatedStat(int(LoadAvgStatType), avgWindow)
		if err != nil {
			m.logger.Error(fmt.Sprintf("error aggregate loadavg statistic: %v", err))
			return nil, err
		}
		fullstat.loadStat = stat.(loadavg.LoadAvg)
	}
	if m.cpuEnabled {
		stat, err := m.AggregatedStat(int(CPUStatType), avgWindow)
		if err != nil {
			m.logger.Error(fmt.Sprintf("error aggregate cpu statistic: %v", err))
			return nil, err
		}
		fullstat.cpuStat = stat.(cpu.CPUStat)
	}
	if m.memoryEnabled {
		stat, err := m.AggregatedStat(int(MemStatType), avgWindow)
		if err != nil {
			m.logger.Error(fmt.Sprintf("error aggregate mem statistic: %v", err))
			return nil, err
		}
		fullstat.memStat = stat.(mem.MemStat)
	}
	if m.diskEnabled {
		stat, err := m.AggregatedStat(int(DiskStatType), avgWindow)
		if err != nil {
			m.logger.Error(fmt.Sprintf("error aggregate disk statistic: %v", err))
			return nil, err
		}
		fullstat.diskStat = stat.(disk.DiskStat)
	}
	if m.networkEnabled {
		stat, err := m.AggregatedStat(int(NetworkStatType), avgWindow)
		if err != nil {
			m.logger.Error(fmt.Sprintf("error aggregate network statistic: %v", err))
			return nil, err
		}
		fullstat.networkStat = stat.(network.NetworkStat)
	}

	return fullstat, nil
}

func (m Metrics) setup(cfg config.Config) {
	m.loadavgEnabled = cfg.Metrics.LoadavgEnabled
	m.cpuEnabled = cfg.Metrics.CPUEnabled
	m.memoryEnabled = cfg.Metrics.MemoryEnabled
	m.diskEnabled = cfg.Metrics.DiskEnabled
	m.networkEnabled = cfg.Metrics.NetworkEnabled

	if m.loadavgEnabled {
		m.avg.AddQueue(LoadAvgStatType)
	}
	if m.cpuEnabled {
		m.avg.AddQueue(CPUStatType)
	}
	if m.memoryEnabled {
		m.avg.AddQueue(MemStatType)
	}
	if m.diskEnabled {
		m.avg.AddQueue(DiskStatType)
	}
	if m.networkEnabled {
		m.avg.AddQueue(NetworkStatType)
	}
}
