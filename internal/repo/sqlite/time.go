package sqlite

import "time"

// All timestamps in SQLite are stored as UTC RFC3339 (e.g. 2006-01-02T15:04:05Z).
const dbTimeLayout = time.RFC3339

func formatTime(t time.Time) string {
	if t.IsZero() {
		return time.Now().UTC().Format(dbTimeLayout)
	}
	return t.UTC().Format(dbTimeLayout)
}

func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(dbTimeLayout, s)
	if err != nil {
		return time.Time{}
	}
	return t
}
