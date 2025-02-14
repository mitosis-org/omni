package flowgen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnits(t *testing.T) {
	t.Parallel()
	require.Equal(t, "1000000000000000000", ether1.String())
}
