package networkstat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNetwork(t *testing.T) {
	stat, err := CurrentStat()
	require.NoError(t, err)
	require.NotNil(t, stat)
}
