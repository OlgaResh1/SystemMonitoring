//go:build windows

package cpustat

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

var kernel32 = syscall.MustLoadDLL("kernel32.dll")

type fileTime struct {
	dwLowDateTime  uint32
	dwHighDateTime uint32
}

func CurrentStat() (*CPUStat, error) {
	var idle, kernel, user fileTime

	sysTimes, err := kernel32.FindProc("GetSystemTimes")
	if err != nil {
		return nil, fmt.Errorf("error find GetSystemTimes: %v", err)
	}
	r0, _, e1 := sysTimes.Call(uintptr(unsafe.Pointer(&idle)), uintptr(unsafe.Pointer(&kernel)), uintptr(unsafe.Pointer(&user)))
	if uint32(r0) == 0 {
		return nil, fmt.Errorf("error call GetSystemTimes: %v", e1)
	}

	idleFirst := int64(idle.dwLowDateTime) | (int64(idle.dwHighDateTime) << 32)
	kernelFirst := int64(kernel.dwLowDateTime) | (int64(kernel.dwHighDateTime) << 32)
	userFirst := int64(user.dwLowDateTime) | (int64(user.dwHighDateTime) << 32)

	time.Sleep(500 * time.Millisecond)

	r0, _, e1 = sysTimes.Call(uintptr(unsafe.Pointer(&idle)), uintptr(unsafe.Pointer(&kernel)), uintptr(unsafe.Pointer(&user)))
	if uint32(r0) == 0 {
		return nil, fmt.Errorf("error call GetSystemTimes: %v", e1)
	}

	idleSecond := int64(idle.dwLowDateTime) | (int64(idle.dwHighDateTime) << 32)
	kernelSecond := int64(kernel.dwLowDateTime) | (int64(kernel.dwHighDateTime) << 32)
	userSecond := int64(user.dwLowDateTime) | (int64(user.dwHighDateTime) << 32)

	totalIdle := float32(idleSecond - idleFirst)
	totalKernel := float32(kernelSecond - kernelFirst)
	totalUser := float32(userSecond - userFirst)
	totalSys := float32(totalKernel + totalUser)

	// fmt.Printf("Idle: %f%%\nKernel: %f%%\nUser: %f%%\n", (totalIdle/totalSys)*100, (totalKernel/totalSys)*100, (totalUser/totalSys)*100)
	// fmt.Printf("\nTotal: %f%%\n", (totalSys-totalIdle)*100/totalSys)

	return &CPUStat{
		User:   (totalUser / totalSys) * 100,
		System: (totalKernel / totalSys) * 100,
		Idle:   (totalIdle / totalSys) * 100,
	}, nil
}
