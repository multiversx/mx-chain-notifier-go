package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBloom_IsInSetFromStrings(t *testing.T) {
	t.Parallel()

	n := 4
	bloom := NewBloom(uint(n))

	d1 := []byte("This")
	d2 := []byte("Should")
	d3 := []byte("Be")
	d4 := []byte("In set")

	err := bloom.SetMany([][]byte{d1, d2, d3, d4})
	require.Nil(t, err)

	require.True(t, bloom.IsInSet(d2))
	require.False(t, bloom.IsInSet([]byte("Not in set")))
}

func TestBloom_IsInSetFromHexStrings(t *testing.T) {
	t.Parallel()

	n := 4
	bloom := NewBloom(uint(n))

	d1 := []byte("425553442d313233343536")
	d2 := []byte("3135")

	err := bloom.SetMany([][]byte{d1, d2})
	require.Nil(t, err)

	require.True(t, bloom.IsInSet(d1))
	require.True(t, bloom.IsInSet(d2))
	require.False(t, bloom.IsInSet([]byte("3135")))
}
