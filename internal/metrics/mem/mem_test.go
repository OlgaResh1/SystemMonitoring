package memstat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemStat(t *testing.T) {
	memstat, err := CurrentStat()
	require.NoError(t, err)
	require.NotNil(t, memstat)
	require.NotEqual(t, memstat.MemTotal, 0)
}
