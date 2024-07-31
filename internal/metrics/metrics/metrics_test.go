package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/aggregator"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/config"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/logger"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/common"
	cpustat "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/cpu"
	diskstat "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/disk"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/loadavg"
	memstat "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/mem"
	networkstat "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/network"
	"github.com/stretchr/testify/require"
)

func TestMetrics(t *testing.T) {
	logg := logger.New("info", "json", false)

	cfg := config.NewConfig()
	cfg.Metrics.LoadavgEnabled = true
	cfg.Metrics.CPUEnabled = true
	cfg.Metrics.MemoryEnabled = true
	cfg.Metrics.DiskEnabled = true
	cfg.Metrics.NetworkEnabled = true

	avg := aggregator.New(logg, 1*time.Second)
	metrics := New(cfg, avg, *logg)
	avg.RegisterMetric(metrics)

	ctx, cancel := context.WithCancel(context.Background())
	go avg.Run(ctx)
	defer cancel()

	time.Sleep(3 * time.Second)

	stat, err := metrics.AggregatedStat(int(common.LoadAvgStatType), 3*time.Second)
	require.NoError(t, err)
	require.NotNil(t, stat)
	load, ok := stat.(*loadavg.LoadAvg)
	require.True(t, ok)
	require.False(t, load.LoadAvg1 == 0 && load.LoadAvg5 == 0 && load.LoadAvg15 == 0)

	stat, err = metrics.AggregatedStat(int(common.CPUStatType), 3*time.Second)
	require.NoError(t, err)
	require.NotNil(t, stat)
	cpuStat, ok := stat.(*cpustat.CPUStat)
	require.True(t, ok)
	require.True(t, cpuStat.User+cpuStat.System+cpuStat.Idle > 0)

	stat, err = metrics.AggregatedStat(int(common.MemStatType), 3*time.Second)
	require.NoError(t, err)
	require.NotNil(t, stat)
	memStat, ok := stat.(*memstat.MemStat)
	require.True(t, ok)
	require.True(t, memStat.MemTotal > 0)

	stat, err = metrics.AggregatedStat(int(common.DiskStatType), 3*time.Second)
	require.NoError(t, err)
	require.NotNil(t, stat)
	diskStat, ok := stat.(*diskstat.DiskStat)
	require.True(t, ok)
	require.True(t, len(diskStat.DiskSpace) > 0)
	require.True(t, len(diskStat.DiskLoad) > 0)

	stat, err = metrics.AggregatedStat(int(common.NetworkStatType), 3*time.Second)
	require.NoError(t, err)
	require.NotNil(t, stat)
	networkStat, ok := stat.(*networkstat.NetworkStat)
	require.True(t, ok)
	require.True(t, len(networkStat.SocketStates) > 0)

	filStatPB, err := metrics.AggregatedFullStatPB(3 * time.Second)
	require.NoError(t, err)
	require.NotNil(t, filStatPB)

	fullStat, err := metrics.AggregatedFullStat(3 * time.Second)
	require.NoError(t, err)
	require.NotNil(t, fullStat)
	full, ok := fullStat.(*FullStat)
	require.Equal(t, true, ok)
	require.NotNil(t, full)
}
