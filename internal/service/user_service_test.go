package service_test

import (
	"testing"
	"time"

	"github.com/example/go-user-api/internal/service"
)

func TestCalculateAge(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		dob      time.Time
		wantAge  int
	}{
		{
			name:    "birthday already passed this year",
			dob:     time.Date(now.Year()-30, now.Month()-1, now.Day(), 0, 0, 0, 0, time.UTC),
			wantAge: 30,
		},
		{
			name:    "birthday not yet this year",
			dob:     time.Date(now.Year()-25, now.Month()+1, now.Day(), 0, 0, 0, 0, time.UTC),
			wantAge: 24,
		},
		{
			name:    "birthday is today",
			dob:     time.Date(now.Year()-20, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			wantAge: 20,
		},
		{
			name:    "newborn (dob = today)",
			dob:     now,
			wantAge: 0,
		},
		{
			name:    "classic date: born 1990-05-10",
			dob:     time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC),
			wantAge: calculateExpected(time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC)),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := service.CalculateAge(tc.dob)
			if got != tc.wantAge {
				t.Errorf("CalculateAge(%v) = %d, want %d", tc.dob.Format("2006-01-02"), got, tc.wantAge)
			}
		})
	}
}

// calculateExpected mirrors CalculateAge so the "1990-05-10" test stays correct
// regardless of when the test suite is run.
func calculateExpected(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		years--
	}
	return years
}
