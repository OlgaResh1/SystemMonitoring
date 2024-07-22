//go:build linux

package loadavg

import (
	"fmt"
	"io"
	"os"
)

type LoadAvg struct {
	LoadAvg1  float32
	LoadAvg5  float32
	LoadAvg15 float32
	// CntRunnableEntities int
	// CntEntities         int
	// LastPID             int
}

func CurrentStat() (*LoadAvg, error) {
	file, err := os.Open("/proc/loadavg")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parceStat(file)
}

func parceStat(file io.Reader) (*LoadAvg, error) {
	loadavg := &LoadAvg{}
	ret, err := fmt.Fscanf(file, "%f %f %f", &loadavg.LoadAvg1,
		&loadavg.LoadAvg5, &loadavg.LoadAvg15)
	if err != nil || ret != 3 {
		return nil, fmt.Errorf("unexpected format of /proc/loadavg")
	}
	return loadavg, nil
}

func AggregatedStats(stat []any) (*LoadAvg, error) {
	loadavg := &LoadAvg{}
	if len(stat) == 0 {
		return loadavg, nil
	}
	for _, s := range stat {
		if s, ok := s.(*LoadAvg); ok {
			loadavg.LoadAvg1 += s.LoadAvg1
			loadavg.LoadAvg5 += s.LoadAvg5
			loadavg.LoadAvg15 += s.LoadAvg15
		} else {
			return nil, fmt.Errorf("error aggregate, expected *LoadAvgStat")
		}
	}
	cnt := float32(len(stat))
	return &LoadAvg{
		LoadAvg1:  loadavg.LoadAvg1 / cnt,
		LoadAvg5:  loadavg.LoadAvg5 / cnt,
		LoadAvg15: loadavg.LoadAvg15 / cnt,
	}, nil
}
