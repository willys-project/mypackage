package daterange

import (
	"time"
)

// GenerateMonthlyInterval generates the start and end dates for the month closest to the current date.
// Logic fungsi ini sama persis dengan kode asli.
func GenerateMonthlyInterval(startTime time.Time) map[string]string {
	var closestInterval map[string]string
	now := time.Now()

	for i := 0; i < 12; i++ {
		startOfMonth := startTime.AddDate(0, i, 0)
		endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

		if startOfMonth.After(now) {
			formattedStartOfMonth := startOfMonth.Format("2006-01-02")
			formattedEndOfMonth := endOfMonth.Format("2006-01-02")

			closestInterval = map[string]string{
				"start": formattedStartOfMonth,
				"end":   formattedEndOfMonth,
			}
			break
		}
	}

	return closestInterval
}
