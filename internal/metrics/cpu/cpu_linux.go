//go:build linux

package cpustat

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/metrics/common"
)

func CurrentStat() (common.Metric, error) {
	cmd := exec.Command("iostat", "-c")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	cpu := &CPUStat{}

	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if !strings.Contains(line, "avg-cpu") {
			continue
		}
		if len(lines) <= i+1 {
			return nil, fmt.Errorf("unexpected format of iostat output: %s", line)
		}
		fields := strings.Fields(strings.ReplaceAll(lines[i+1], ",", "."))
		if len(fields) != 6 {
			return nil, fmt.Errorf("unexpected format of iostat output, count columns %d", len(fields))
		}
		if cpu.User, err = parseValueToFloat(fields[0]); err != nil {
			return nil, err
		}
		if cpu.Nice, err = parseValueToFloat(fields[1]); err != nil {
			return nil, err
		}
		if cpu.System, err = parseValueToFloat(fields[2]); err != nil {
			return nil, err
		}
		if cpu.Iowait, err = parseValueToFloat(fields[3]); err != nil {
			return nil, err
		}
		if cpu.Steal, err = parseValueToFloat(fields[4]); err != nil {
			return nil, err
		}
		if cpu.Idle, err = parseValueToFloat(fields[5]); err != nil {
			return nil, err
		}
		return cpu, nil
	}
	return nil, fmt.Errorf("unexpected format of iostat output")
}

func parseValueToFloat(input string) (result float32, err error) {
	res, err := strconv.ParseFloat(input, 32)
	if err != nil {
		return 0, fmt.Errorf("unexpected format of iostat output, error parse float %s", input)
	}
	result = float32(res)
	return result, nil
}
