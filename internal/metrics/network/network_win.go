//go:build windows

package networkstat

type ListenSocket struct {
	Protocol     string
	LocalAddress string
	PeerAddress  string
	Process      string
}

type NetworkStat struct {
	ListenSockets []ListenSocket
	SocketStates  map[string]int32
}

func CurrentStat() (*NetworkStat, error) {
	return &NetworkStat{}, nil
}

func AggregatedStats(stat []any) (*NetworkStat, error) {
	return &NetworkStat{}, nil
}
