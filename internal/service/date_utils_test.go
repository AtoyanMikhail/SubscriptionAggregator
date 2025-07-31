package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseMonthYear_Success(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "January 2025",
			input:    "01-2025",
			expected: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "December 2024",
			input:    "12-2024",
			expected: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "June 2023",
			input:    "06-2023",
			expected: time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseMonthYear(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseMonthYear_InvalidFormat(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Empty string",
			input: "",
		},
		{
			name:  "Wrong separator",
			input: "01/2025",
		},
		{
			name:  "Missing month",
			input: "-2025",
		},
		{
			name:  "Missing year",
			input: "01-",
		},
		{
			name:  "Invalid month",
			input: "13-2025",
		},
		{
			name:  "Invalid month zero",
			input: "00-2025",
		},
		{
			name:  "Invalid year",
			input: "01-abc",
		},
		{
			name:  "Year too low",
			input: "01-1800",
		},
		{
			name:  "Year too high",
			input: "01-3500",
		},
		{
			name:  "Too many parts",
			input: "01-02-2025",
		},
		{
			name:  "Non-numeric month",
			input: "ab-2025",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseMonthYear(tt.input)
			assert.Error(t, err)
			assert.Equal(t, ErrInvalidDateFormat, err)
			assert.True(t, result.IsZero())
		})
	}
}

func TestFormatMonthYear(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "January 2025",
			input:    time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC),
			expected: "01-2025",
		},
		{
			name:     "December 2024",
			input:    time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			expected: "12-2024",
		},
		{
			name:     "June 2023",
			input:    time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
			expected: "06-2023",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatMonthYear(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateMonthsInPeriod(t *testing.T) {
	tests := []struct {
		name      string
		startDate time.Time
		endDate   time.Time
		expected  int
	}{
		{
			name:      "Same month",
			startDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
			expected:  1,
		},
		{
			name:      "Full year",
			startDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
			expected:  12,
		},
		{
			name:      "Partial months",
			startDate: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC),
			expected:  2,
		},
		{
			name:      "Same date",
			startDate: time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
			expected:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateMonthsInPeriod(tt.startDate, tt.endDate)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function to create a time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}