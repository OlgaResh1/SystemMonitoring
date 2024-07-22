package statqueue

import (
	"errors"
	"sync"
	"time"
)

type statElement struct {
	stat      any
	timestamp time.Time
}

type StatQueue struct {
	stats []statElement
	mutex sync.Mutex
}

var ErrEmptyQueue = errors.New("empty queue")

func NewStatQueue() *StatQueue {
	return &StatQueue{
		stats: make([]statElement, 0, 10),
	}
}

func (s *StatQueue) Append(stat any, time time.Time) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.stats = append(s.stats, statElement{stat: stat, timestamp: time})
	return nil
}

func (s *StatQueue) GetLast(interval time.Duration, timeNow time.Time) (stat []any, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	timeend := timeNow.Add(-interval)
	for i := len(s.stats) - 1; i >= 0; i-- {
		if s.stats[i].timestamp.Compare(timeend) >= 0 {
			stat = append(stat, s.stats[i].stat)
		} else {
			break
		}
	}
	return stat, err
}

func (s *StatQueue) Len() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return len(s.stats)
}

func (s *StatQueue) CurTail(timecut time.Time) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.stats) == 0 {
		return ErrEmptyQueue
	}
	for len(s.stats) > 0 && s.stats[0].timestamp.Before(timecut) {
		s.stats = s.stats[1:]
	}

	return nil
}
