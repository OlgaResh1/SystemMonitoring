package metrics

import (
	"fmt"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/aggregator"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/config"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/logger"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/common"
	cpu "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/cpu"
	disk "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/disk"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/loadavg"
	mem "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/mem"
	network "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/network"
)

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
	switch queueType {
	case common.LoadAvgStatType:
		return loadavg.CurrentStat()
	case common.CPUStatType:
		return cpu.CurrentStat()
	case common.MemStatType:
		return mem.CurrentStat()
	case common.DiskStatType:
		return disk.CurrentStat()
	case common.NetworkStatType:
		return network.CurrentStat()
	}
	return nil, common.ErrUnknownQueueType
}

func (m Metrics) AggregatedStat(queueType int, avgWindow time.Duration) (any, error) {
	stat, err := m.avg.StatByInterval(aggregator.StatQueueType(queueType), avgWindow)
	if err != nil {
		return nil, err
	}
	metrics := make([]common.Metric, len(stat))
	for i, v := range stat {
		metrics[i] = v.(common.Metric)
	}

	switch queueType {
	case common.LoadAvgStatType:
		return loadavg.AggregatedStats(metrics)
	case common.CPUStatType:
		return cpu.AggregatedStats(metrics)
	case common.MemStatType:
		return mem.AggregatedStats(metrics)
	case common.DiskStatType:
		return disk.AggregatedStats(metrics)
	case common.NetworkStatType:
		return network.AggregatedStats(metrics)
	}
	return nil, common.ErrUnknownQueueType
}

func (m Metrics) AggregatedFullStat(avgWindow time.Duration) (any, error) {
	fullstat := &FullStat{}
	fullstat.timestamp = time.Now()
	fullstat.interval = avgWindow

	if m.loadavgEnabled {
		stat, err := m.AggregatedStat(common.LoadAvgStatType, avgWindow)
		if err != nil {
			m.logger.Error(fmt.Sprintf("error aggregate loadavg statistic: %v", err))
			return nil, err
		}
		fullstat.loadStat = stat.(loadavg.LoadAvg)
	}
	if m.cpuEnabled {
		stat, err := m.AggregatedStat(common.CPUStatType, avgWindow)
		if err != nil {
			m.logger.Error(fmt.Sprintf("error aggregate cpu statistic: %v", err))
			return nil, err
		}
		fullstat.cpuStat = stat.(cpu.CPUStat)
	}
	if m.memoryEnabled {
		stat, err := m.AggregatedStat(common.MemStatType, avgWindow)
		if err != nil {
			m.logger.Error(fmt.Sprintf("error aggregate mem statistic: %v", err))
			return nil, err
		}
		fullstat.memStat = stat.(mem.MemStat)
	}
	if m.diskEnabled {
		stat, err := m.AggregatedStat(common.DiskStatType, avgWindow)
		if err != nil {
			m.logger.Error(fmt.Sprintf("error aggregate disk statistic: %v", err))
			return nil, err
		}
		fullstat.diskStat = stat.(disk.DiskStat)
	}
	if m.networkEnabled {
		stat, err := m.AggregatedStat(common.NetworkStatType, avgWindow)
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
		m.avg.AddQueue(common.LoadAvgStatType)
	}
	if m.cpuEnabled {
		m.avg.AddQueue(common.CPUStatType)
	}
	if m.memoryEnabled {
		m.avg.AddQueue(common.MemStatType)
	}
	if m.diskEnabled {
		m.avg.AddQueue(common.DiskStatType)
	}
	if m.networkEnabled {
		m.avg.AddQueue(common.NetworkStatType)
	}
}
