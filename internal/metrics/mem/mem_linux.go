//go:build linux

package memstat

import (
	"bufio"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func CurrentStat() (*MemStat, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseStat(file)
}

func parseStat(file io.Reader) (cpuStat *MemStat, err error) {
	memstat := &MemStat{}

	scanner := bufio.NewScanner(file)

	s := reflect.ValueOf(memstat).Elem()

	var i int
	for scanner.Scan() {
		line := scanner.Text()
		if i = strings.IndexRune(line, ':'); i < 0 {
			continue
		}
		key := strings.TrimSpace(line[:i])
		str := strings.TrimSpace(line[i+1:])
		str = strings.ReplaceAll(str, " kB", "")
		value, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, err
		}
		f := s.FieldByName(key)
		if f.IsValid() && f.CanSet() {
			f.SetUint(value)
		}

		// switch key {
		// case "MemTotal":
		// 	memstat.MemTotal = value
		// case "MemFree":
		// 	memstat.MemFree = value
		// case "Buffers":
		// 	memstat.Buffers = value
		// case "Cached":
		// 	memstat.Cached = value
		// case "SwapCached":
		// 	memstat.SwapCached = value
		// case "SwapTotal":
		// 	memstat.SwapTotal = value
		// case "SwapFree":
		// 	memstat.SwapFree = value
		// case "Active":
		// 	memstat.Active = value
		// case "Inactive":
		// 	memstat.Inactive = value
		// case "VmallocUsed":
		// 	memstat.VmallocUsed = value
		// case "Mapped":
		// 	memstat.Mapped = value
		// default:
		// 	continue
		// }
	}

	return memstat, nil
}
