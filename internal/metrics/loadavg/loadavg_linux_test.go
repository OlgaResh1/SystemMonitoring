package loadavg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadAvg(t *testing.T) {
	loadavg, err := CurrentStat()
	require.NoError(t, err)
	require.NotNil(t, loadavg)
}
