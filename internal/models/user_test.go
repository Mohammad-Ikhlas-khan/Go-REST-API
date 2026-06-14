package models

import (
	"testing"
	"time"
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateAge(tt.dob)
			if got != tt.wantAge {
				t.Errorf("CalculateAge(%v) = %d, want %d", tt.dob, got, tt.wantAge)
			}
		})
	}
}
