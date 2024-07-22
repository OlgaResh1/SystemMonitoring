//go:build linux

package networkstat

import (
	"os/exec"
	"strings"
)

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
	listenSockets, err := liseningList()
	if err != nil {
		return nil, err
	}
	socketStates, err := socketStates()
	if err != nil {
		return nil, err
	}
	return &NetworkStat{
		ListenSockets: listenSockets,
		SocketStates:  socketStates,
	}, nil
}

func AggregatedStats(stat []any) (*NetworkStat, error) {
	networkStat := &NetworkStat{
		SocketStates: make(map[string]int32),
	}
	if len(stat) == 0 {
		return networkStat, nil
	}
	listenMap := make(map[string]ListenSocket)
	cntByState := make(map[string]int32)
	for _, s := range stat {
		if s, ok := s.(*NetworkStat); ok {
			for _, ls := range s.ListenSockets {
				listenMap[ls.LocalAddress+ls.Process] = ls
			}
			for state, cnt := range s.SocketStates {
				networkStat.SocketStates[state] += cnt
				cntByState[state]++
			}
		}
	}
	for _, listener := range listenMap {
		networkStat.ListenSockets = append(networkStat.ListenSockets, listener)
	}
	for state, cntAll := range networkStat.SocketStates {
		networkStat.SocketStates[state] = cntAll / cntByState[state]
	}
	return networkStat, nil
}

func liseningList() ([]ListenSocket, error) {
	cmd := exec.Command("ss", "-lntup")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n")
	listenSockets := make([]ListenSocket, 0, len(lines)-1)
	for i, line := range lines {
		if i == 0 {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		var ProcessName string
		if len(fields) > 6 {
			ProcessName = fields[6]
		}
		listener := ListenSocket{
			Protocol:     fields[0],
			LocalAddress: fields[4],
			PeerAddress:  fields[5],
			Process:      ProcessName,
		}
		listenSockets = append(listenSockets, listener)
	}
	return listenSockets, nil
}

func socketStates() (map[string]int32, error) {
	cmd := exec.Command("ss", "-ta")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	socketStates := make(map[string]int32)
	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		socketStates[fields[0]]++
	}
	return socketStates, nil
}
