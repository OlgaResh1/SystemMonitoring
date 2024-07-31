package memstat

import (
	"fmt"

	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/common"
)

type MemStat struct {
	MemTotal, MemFree, Buffers, Cached, SwapCached, SwapTotal, SwapFree,
	Active, Inactive, VmallocUsed, Mapped uint64
}

func (s *MemStat) MetricType() int {
	return common.DiskStatType
}

func AggregatedStats(stat []common.Metric) (common.Metric, error) {
	if len(stat) == 0 {
		return &MemStat{}, nil
	}
	memstat := MemStat{}
	for _, s := range stat {
		if s, ok := s.(*MemStat); ok {
			memstat.MemTotal += s.MemTotal
			memstat.MemFree += s.MemFree
			memstat.Buffers += s.Buffers
			memstat.Cached += s.Cached
			memstat.SwapCached += s.SwapCached
			memstat.SwapTotal += s.SwapTotal
			memstat.SwapFree += s.SwapFree
			memstat.Active += s.Active
			memstat.Inactive += s.Inactive
			memstat.VmallocUsed += s.VmallocUsed
			memstat.Mapped += s.Mapped
		} else {
			return nil, fmt.Errorf("error aggregate, expected *CPUStat")
		}
	}
	cnt := uint64(len(stat))

	return &MemStat{
		MemTotal:    memstat.MemTotal / cnt,
		MemFree:     memstat.MemFree / cnt,
		Buffers:     memstat.Buffers / cnt,
		Cached:      memstat.Cached / cnt,
		SwapCached:  memstat.SwapCached / cnt,
		SwapTotal:   memstat.SwapTotal / cnt,
		SwapFree:    memstat.SwapFree / cnt,
		Active:      memstat.Active / cnt,
		Inactive:    memstat.Inactive / cnt,
		VmallocUsed: memstat.VmallocUsed / cnt,
		Mapped:      memstat.Mapped / cnt,
	}, nil
}
