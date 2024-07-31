package common

import (
	"errors"
)

type Metric interface {
	MetricType() int
}

const (
	LoadAvgStatType int = iota
	CPUStatType
	MemStatType
	DiskStatType
	NetworkStatType
)

var ErrUnknownQueueType = errors.New("unknown queue type")
