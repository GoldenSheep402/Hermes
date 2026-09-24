package settle

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSeedingBonusPoints(t *testing.T) {
	t.Parallel()

	// Brand-new seed (T≈0) earns near-zero regardless of size.
	require.InDelta(t, 0, seedingBonusPoints(0, 10, 1, 1), 1e-9)

	// After many weeks, time factor ≈ 1; single seeder gets 1+sqrt(2) multiplier.
	got := seedingBonusPoints(10, 1, 1, 1)
	want := 1 * (1 + math.Sqrt2)
	require.InDelta(t, want, got, 1e-6)

	// More seeders reduce the seed factor.
	lonely := seedingBonusPoints(10, 1, 1, 1)
	crowded := seedingBonusPoints(10, 1, 20, 1)
	require.Greater(t, lonely, crowded)
}

func TestParsePairMember(t *testing.T) {
	t.Parallel()
	u, tor, ok := parsePairMember("user1:torrent2")
	require.True(t, ok)
	require.Equal(t, "user1", u)
	require.Equal(t, "torrent2", tor)

	_, _, ok = parsePairMember("bad")
	require.False(t, ok)
}
