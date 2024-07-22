package diskstat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiskStat(t *testing.T) {
	stat, err := CurrentStat()
	require.NoError(t, err)
	require.NotNil(t, stat)
}
