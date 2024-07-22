package app

import (
	"time"
)

type App struct {
	logger  Logger
	metrics Metrics
}

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Metrics interface {
	CurrentStat(queueType int) (any, error)
	AggregatedStat(queueType int, avgWindow time.Duration) (any, error)
	AggregatedFullStat(avgWindow time.Duration) (any, error)
	AggregatedStatPB(queueType int, avgWindow time.Duration) (any, error)
	AggregatedFullStatPB(avgWindow time.Duration) (any, error)
}

func New(logger Logger, metrics Metrics) *App {
	return &App{logger: logger, metrics: metrics}
}

func (a App) StatByInterval(queueType int, avgWindow time.Duration) (any, error) {
	return a.metrics.AggregatedStat(queueType, avgWindow)
}

func (a App) FullStatByInterval(avgWindow time.Duration) (any, error) {
	return a.metrics.AggregatedFullStat(avgWindow)
}

func (a App) StatByIntervalPB(queueType int, avgWindow time.Duration) (any, error) {
	return a.metrics.AggregatedStatPB(queueType, avgWindow)
}

func (a App) FullStatByIntervalPB(avgWindow time.Duration) (any, error) {
	return a.metrics.AggregatedFullStatPB(avgWindow)
}
