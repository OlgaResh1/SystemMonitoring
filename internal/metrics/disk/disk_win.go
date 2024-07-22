//go:build windows

package diskstat

import (
	"fmt"
	"syscall"
	"unsafe"
)

var kernel32 = syscall.MustLoadDLL("kernel32.dll")

func CurrentStat() (*DiskStat, error) {
	stat := &DiskStat{
		DiskLoad:  make(map[string]DiskLoad),
		DiskSpace: make(map[string]DiskSpace),
	}

	drivesProc, err := kernel32.FindProc("GetLogicalDrives")
	if err != nil {
		return nil, fmt.Errorf("error find GetLogicalDrive: %v", err)
	}

	var bitMask uint32
	r0, _, e1 := drivesProc.Call()
	bitMask = uint32(r0)
	if bitMask == 0 {
		return nil, fmt.Errorf("error call GetLogicalDriveStrings: %v", e1)
	}
	var drivers []byte
	for i := 0; i < 26; i++ {
		if (bitMask & (1 << i)) != 0 {
			drivers = append(drivers, byte('A'+i))
		}
	}
	for i := 0; i < len(drivers); i++ {
		d := fmt.Sprintf("%c:\\", drivers[i])
		dstat, err := getDiskSpace(d)
		if err != nil {
			return nil, fmt.Errorf("error get disk space for %s: %v", d, err)
		}
		stat.DiskSpace[d] = dstat
	}

	return stat, nil
}

func getDiskSpace(drive string) (DiskSpace, error) {
	diskSpace := DiskSpace{}

	kernel32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return diskSpace, fmt.Errorf("error find kernel32.dll: %v", err)
	}
	freeSpaceProc, err := kernel32.FindProc("GetDiskFreeSpaceExW")
	if err != nil {
		return diskSpace, fmt.Errorf("error find GetDiskFreeSpaceExW: %v", err)
	}
	drivebytes, err := syscall.UTF16FromString(drive)
	if err != nil {
		return diskSpace, err
	}

	var lpFreeBytesAvailable, lpTotalBytes, lpTotalFreeBytes int64
	r1, _, e1 := freeSpaceProc.Call(uintptr(unsafe.Pointer(&drivebytes[0])), uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalBytes)), uintptr(unsafe.Pointer(&lpTotalFreeBytes)))

	if r1 != 0 {
		diskSpace.Available = lpFreeBytesAvailable / 1024 / 1024 // KB
		diskSpace.Total = lpTotalBytes / 1024 / 1024             // KB
	} else {
		return diskSpace, fmt.Errorf("error call GetDiskFreeSpaceExW: %v", e1)
	}
	return diskSpace, nil
}
