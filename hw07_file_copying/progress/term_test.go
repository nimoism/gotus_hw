package progress

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTerm(t *testing.T) {
	t.Run("term parse size", func(t *testing.T) {
		term := &Term{}
		width, height, err := term.size([]byte("16 22"))
		require.NoError(t, err)
		require.Equal(t, 22, width)
		require.Equal(t, 16, height)
	})

	t.Run("term parse width wrong input", func(t *testing.T) {
		term := Term{}
		_, _, err := term.size([]byte("NAN NAN"))
		require.Error(t, err)
	})
}
