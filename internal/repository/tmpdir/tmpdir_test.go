package tmpdir

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandomDirname(t *testing.T) {
	actual := randomDirName()
	require.NotEmpty(t, actual)
	require.GreaterOrEqual(t, len(actual), 10)
}

func TestRandomDirnameFallback(t *testing.T) {
	actual := randomDirNameFallback()
	require.NotEmpty(t, actual)
	require.GreaterOrEqual(t, len(actual), 10)
}
