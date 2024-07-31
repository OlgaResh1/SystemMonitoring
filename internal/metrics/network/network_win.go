//go:build windows

package networkstat

import "github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/common"

type ListenSocket struct {
	Protocol     string
	LocalAddress string
	PeerAddress  string
	Process      string
}

type NetworkStat struct {
	common.Metric
	ListenSockets []ListenSocket
	SocketStates  map[string]int32
}

func CurrentStat() (*NetworkStat, error) {
	return &NetworkStat{}, nil
}

func AggregatedStats(stat []common.Metric) (common.Metric, error) {
	return &NetworkStat{}, nil
}
