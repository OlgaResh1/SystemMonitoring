package cpustat

import (
	"testing"

	require "github.com/stretchr/testify/require"
)

func TestCpuStat(t *testing.T) {
	cpustat, err := CurrentStat()
	require.NoError(t, err)
	require.NotNil(t, cpustat)
	// require.NotEqual(t, cpustat.User, 0)
}
