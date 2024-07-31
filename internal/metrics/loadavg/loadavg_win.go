//go:build windows

package loadavg

import "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/common"

type LoadAvg struct {
	common.Metric
	LoadAvg1  float32
	LoadAvg5  float32
	LoadAvg15 float32
	// CntRunnableEntities int
	// CntEntities         int
	// LastPID             int
}

func CurrentStat() (*LoadAvg, error) {
	return &LoadAvg{}, nil
}

func AggregatedStats(stat []common.Metric) (common.Metric, error) {
	return &LoadAvg{}, nil
}
