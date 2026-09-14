package timeseries

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFillDays_ZeroFillsGapsAndEndsToday(t *testing.T) {
	today := time.Date(2026, 9, 14, 15, 30, 0, 0, time.UTC)

	got := FillDays(map[string]int{"2026-09-12": 3, "2026-09-14": 1, "2026-09-01": 99}, today, 4)

	require.Len(t, got, 4)
	require.Equal(t, "2026-09-11", got[0].Day.Format(time.DateOnly))
	require.Equal(t, []int{0, 3, 0, 1}, []int{got[0].Count, got[1].Count, got[2].Count, got[3].Count})
}
