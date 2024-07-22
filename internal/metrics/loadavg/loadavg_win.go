//go:build windows

package loadavg

type LoadAvg struct {
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

func AggregatedStats(stat []any) (*LoadAvg, error) {
	return &LoadAvg{}, nil
}
