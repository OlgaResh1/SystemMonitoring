package aggregator

import (
	"context"
	"fmt"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/statqueue"
)

type StatQueueType int

type statQueue interface {
	GetLast(interval time.Duration, time time.Time) (stat []any, err error)
	Append(stat any, time time.Time) error
}

type metric interface {
	CurrentStat(queueType int) (any, error)
	AggregatedStat(queueType int, avgWindow time.Duration) (any, error)
}

type logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Aggregator struct {
	logger   logger
	interval time.Duration
	queues   map[StatQueueType]statQueue
	metrics  metric
}

func New(logger logger, interval time.Duration) Aggregator {
	a := Aggregator{
		logger:   logger,
		interval: interval,
	}
	a.queues = make(map[StatQueueType]statQueue)
	return a
}

func (a *Aggregator) RegisterMetric(metrics metric) {
	a.metrics = metrics
}

func (a Aggregator) AddQueue(queueType StatQueueType) {
	a.queues[queueType] = statqueue.NewStatQueue()
}

func (a Aggregator) Run(ctx context.Context) error {
	if a.interval == 0 {
		err := fmt.Errorf("interval cannot be zero")
		a.logger.Error("error fetching statistics", "error", err)
		return err
	}
	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			for queueType, queue := range a.queues {
				stat, err := a.metrics.CurrentStat(int(queueType))
				if err != nil {
					a.logger.Error("error fetching statistics", "error", err)
					continue
				}
				queue.Append(stat, time.Now())
			}
		}
	}
}

func (a Aggregator) StatByInterval(queueType StatQueueType, avgWindow time.Duration) ([]any, error) {
	return a.queues[queueType].GetLast(avgWindow, time.Now())
}
