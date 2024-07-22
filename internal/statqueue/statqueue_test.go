package statqueue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	t.Parallel()
	timeNow := time.Now()
	queue := NewStatQueue()
	cnt := 100

	err := queue.CurTail(timeNow)
	require.ErrorIs(t, err, ErrEmptyQueue)

	for i := 0; i < cnt; i++ {
		err := queue.Append(i, timeNow.Add(-time.Duration(cnt-i)*time.Second))
		require.NoError(t, err)
	}
	stat, err := queue.GetLast(time.Duration(cnt)*time.Second, timeNow)
	require.NoError(t, err)
	require.Len(t, stat, cnt)
	value, ok := stat[0].(int)
	require.True(t, ok)
	require.Equal(t, cnt-1, value)
	require.Equal(t, cnt, queue.Len())

	err = queue.CurTail(timeNow.Add(-time.Duration(10) * time.Second))
	require.NoError(t, err)
	stat, err = queue.GetLast(time.Duration(cnt)*time.Second, timeNow)
	require.NoError(t, err)
	require.Len(t, stat, 10)

	require.Equal(t, 10, queue.Len())
}
