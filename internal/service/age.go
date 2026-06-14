package service

import (
	"time"

	"github.com/example/go-user-api/internal/models"
)

// CalculateAge is exported so tests and external callers can use it directly.
// It delegates to models.CalculateAge which owns the implementation.
func CalculateAge(dob time.Time) int {
	return models.CalculateAge(dob)
}