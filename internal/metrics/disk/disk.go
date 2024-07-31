package diskstat

import (
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/common"
)

type DiskLoad struct {
	TransfersPerSec, WritedPerSec, ReadedPerSec float64
}

type DiskSpace struct {
	Total, Used, Available, INodes, UsedPercent int64
}
type DiskStat struct {
	common.Metric
	DiskLoad  map[string]DiskLoad
	DiskSpace map[string]DiskSpace
}

func (s *DiskStat) MetricType() int {
	return common.DiskStatType
}

func AggregatedStats(stat []common.Metric) (common.Metric, error) {
	diskStat := &DiskStat{
		DiskLoad:  make(map[string]DiskLoad),
		DiskSpace: make(map[string]DiskSpace),
	}
	if len(stat) == 0 {
		return diskStat, nil
	}
	cntByDiskLoad := make(map[string]float64)
	cntByDiskSpace := make(map[string]int64)
	for _, s := range stat {
		if s, ok := s.(*DiskStat); ok {
			for diskname, loadstat := range s.DiskLoad {
				diskStat.DiskLoad[diskname] = DiskLoad{
					TransfersPerSec: loadstat.TransfersPerSec + diskStat.DiskLoad[diskname].TransfersPerSec,
					WritedPerSec:    loadstat.WritedPerSec + diskStat.DiskLoad[diskname].WritedPerSec,
					ReadedPerSec:    loadstat.ReadedPerSec + diskStat.DiskLoad[diskname].ReadedPerSec,
				}
				cntByDiskLoad[diskname]++
			}
			for diskname, loadspace := range s.DiskSpace {
				diskStat.DiskSpace[diskname] = DiskSpace{
					Total:       loadspace.Total + diskStat.DiskSpace[diskname].Total,
					Used:        loadspace.Used + diskStat.DiskSpace[diskname].Used,
					Available:   loadspace.Available + diskStat.DiskSpace[diskname].Available,
					INodes:      loadspace.INodes + diskStat.DiskSpace[diskname].INodes,
					UsedPercent: loadspace.UsedPercent + diskStat.DiskSpace[diskname].UsedPercent,
				}
				cntByDiskSpace[diskname]++
			}
		}
	}
	for diskname, loadstat := range diskStat.DiskLoad {
		diskStat.DiskLoad[diskname] = DiskLoad{
			TransfersPerSec: loadstat.TransfersPerSec / cntByDiskLoad[diskname],
			WritedPerSec:    loadstat.WritedPerSec / cntByDiskLoad[diskname],
			ReadedPerSec:    loadstat.ReadedPerSec / cntByDiskLoad[diskname],
		}
	}
	for diskname, loadspace := range diskStat.DiskSpace {
		diskStat.DiskSpace[diskname] = DiskSpace{
			Total:       loadspace.Total / cntByDiskSpace[diskname],
			Used:        loadspace.Used / cntByDiskSpace[diskname],
			Available:   loadspace.Available / cntByDiskSpace[diskname],
			INodes:      loadspace.INodes / cntByDiskSpace[diskname],
			UsedPercent: loadspace.UsedPercent / cntByDiskSpace[diskname],
		}
	}

	return diskStat, nil
}
