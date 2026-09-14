// Package timeseries holds the tiny per-day aggregation shared by the modules' stats queries.
package timeseries

import (
	"database/sql"
	"time"
)

type DayCount struct {
	Day   time.Time
	Count int
}

// ScanPerDay reads "(date, count)" rows and returns the last `days` days, oldest
// first and ending today (UTC), with days that had no rows filled in as zero.
func ScanPerDay(rows *sql.Rows, days int) ([]DayCount, error) {
	counts := map[string]int{}
	for rows.Next() {
		var day time.Time
		var n int
		if err := rows.Scan(&day, &n); err != nil {
			return nil, err
		}
		counts[day.Format(time.DateOnly)] = n
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return FillDays(counts, time.Now().UTC(), days), nil
}

func FillDays(counts map[string]int, today time.Time, days int) []DayCount {
	series := make([]DayCount, 0, days)
	start := today.Truncate(24*time.Hour).AddDate(0, 0, -(days - 1))
	for i := 0; i < days; i++ {
		day := start.AddDate(0, 0, i)
		series = append(series, DayCount{Day: day, Count: counts[day.Format(time.DateOnly)]})
	}
	return series
}
