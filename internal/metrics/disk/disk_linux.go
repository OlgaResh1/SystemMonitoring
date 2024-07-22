//go:build linux

package diskstat

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func CurrentStat() (*DiskStat, error) {
	stat := &DiskStat{}

	loadStat, err := statDiskLoad()
	if err != nil {
		return nil, err
	}
	stat.DiskLoad = loadStat
	spaceStat, err := statDiskSpace()
	if err != nil {
		return nil, err
	}
	stat.DiskSpace = spaceStat
	return stat, nil
}

func parceValueToFloat(input string) (result float64, err error) {
	result, err = strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("unexpected format of iostat output, error parce float %s", input)
	}
	return result, nil
}

func parceValueToInt(input string) (result int64, err error) {
	result, err = strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unexpected format of iostat output, error parce int %s", input)
	}
	return result, nil
}

func statDiskLoad() (loadStat map[string]DiskLoad, err error) {
	loadStat = make(map[string]DiskLoad)
	cmd := exec.Command("iostat", "-d", "-k")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if i < 3 {
			continue
		}
		fields := strings.Fields(strings.ReplaceAll(line, ",", "."))
		if len(fields) < 6 {
			continue
		}
		diskLoad := DiskLoad{}
		diskLoad.TransfersPerSec, err = parceValueToFloat(fields[3])
		if err != nil {
			return nil, err
		}
		diskLoad.WritedPerSec, err = parceValueToFloat(fields[2])
		if err != nil {
			return nil, err
		}
		diskLoad.ReadedPerSec, err = parceValueToFloat(fields[1])
		if err != nil {
			return nil, err
		}
		loadStat[fields[0]] = diskLoad
	}
	return loadStat, nil
}

func statDiskSpace() (spaceStat map[string]DiskSpace, err error) {
	spaceStat = make(map[string]DiskSpace)
	cmd := exec.Command("df", "-k")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(strings.ReplaceAll(line, ",", "."))
		if len(fields) < 5 {
			continue
		}
		if strings.Index(fields[0], "/dev") != 0 {
			continue
		}
		diskSpace := DiskSpace{}
		diskSpace.Total, err = parceValueToInt(fields[1])
		if err != nil {
			return nil, err
		}
		diskSpace.Used, err = parceValueToInt(fields[2])
		if err != nil {
			return nil, err
		}
		diskSpace.Available, err = parceValueToInt(fields[3])
		if err != nil {
			return nil, err
		}
		diskSpace.UsedPercent, err = parceValueToInt(strings.ReplaceAll(fields[4], "%", ""))
		if err != nil {
			return nil, err
		}
		spaceStat[fields[0]] = diskSpace
	}
	cmd = exec.Command("df", "-i")
	output, err = cmd.Output()
	if err != nil {
		return nil, err
	}
	lines = strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(strings.ReplaceAll(line, ",", "."))
		if len(fields) < 5 {
			continue
		}
		if strings.Index(fields[0], "/dev") != 0 {
			continue
		}
		diskSpace := spaceStat[fields[0]]
		diskSpace.INodes, err = parceValueToInt(fields[1])
		if err != nil {
			return nil, err
		}
		spaceStat[fields[0]] = diskSpace
	}
	return spaceStat, nil
}
