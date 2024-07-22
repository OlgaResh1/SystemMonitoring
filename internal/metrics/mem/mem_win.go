//go:build windows

package memstat

import (
	"fmt"
	"syscall"
	"unsafe"
)

type MEMORYSTATUSEX struct {
	dwLength                uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

var kernel32 = syscall.MustLoadDLL("kernel32.dll")

func CurrentStat() (*MemStat, error) {
	mem := &MemStat{}
	memStatus := MEMORYSTATUSEX{}
	memStatus.dwLength = uint32(unsafe.Sizeof(memStatus))

	proc, err := kernel32.FindProc("GlobalMemoryStatusEx")
	if err != nil {
		return nil, fmt.Errorf("error find proc GlobalMemoryStatusEx: %v", err)
	}
	r1, _, err := proc.Call(uintptr(unsafe.Pointer(&memStatus)))
	if r1 == 0 {
		return nil, fmt.Errorf("error call GlobalMemoryStatusEx: %v", err)
	}
	mem.MemTotal = memStatus.ullTotalPhys
	mem.MemFree = memStatus.ullAvailPhys
	mem.Buffers = memStatus.ullTotalPageFile - memStatus.ullAvailPageFile
	mem.Cached = memStatus.ullTotalPageFile - memStatus.ullAvailPageFile - memStatus.ullTotalPhys + memStatus.ullAvailPhys

	return mem, nil
}
