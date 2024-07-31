package cpustat

import (
	"fmt"

	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/common"
)

type CPUStat struct {
	common.Metric
	User, Nice, System, Idle, Iowait, Steal float32
}

func (s *CPUStat) MetricType() int {
	return common.CPUStatType
}

func AggregatedStats(stat []common.Metric) (common.Metric, error) {
	cpu := &CPUStat{}
	if len(stat) == 0 {
		return cpu, nil
	}
	for _, s := range stat {
		if s, ok := s.(*CPUStat); ok {
			cpu.User += s.User
			cpu.Nice += s.Nice
			cpu.System += s.System
			cpu.Idle += s.Idle
			cpu.Iowait += s.Iowait
			cpu.Steal += s.Steal
		} else {
			return nil, fmt.Errorf("error aggregate, expected *CPUStat")
		}
	}
	cnt := float32(len(stat))

	return &CPUStat{
		User:   cpu.User / cnt,
		Nice:   cpu.Nice / cnt,
		System: cpu.System / cnt,
		Idle:   cpu.Idle / cnt,
		Iowait: cpu.Iowait / cnt,
		Steal:  cpu.Steal / cnt,
	}, nil
}
