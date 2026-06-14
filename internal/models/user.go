package models

import "time"

// CreateUserRequest is the request body for POST /users.
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	DOB  string `json:"dob"  validate:"required,datetime=2006-01-02"`
}

// UpdateUserRequest is the request body for PUT /users/:id.
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	DOB  string `json:"dob"  validate:"required,datetime=2006-01-02"`
}

// UserResponse is the standard response for create / update.
type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
}

// UserWithAgeResponse is the response for get / list – includes computed age.
type UserWithAgeResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Age  int    `json:"age"`
}

// PaginatedUsersResponse wraps a page of users with metadata.
type PaginatedUsersResponse struct {
	Data       []UserWithAgeResponse `json:"data"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"total_pages"`
}

// CalculateAge returns the completed years between dob and now.
func CalculateAge(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()
	// subtract one year if the birthday hasn't occurred yet this calendar year
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		years--
	}
	return years
}
