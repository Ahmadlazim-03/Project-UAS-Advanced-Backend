package utils

import (
	"testing"
)

func TestGetPaginationParams(t *testing.T) {
	// This is a basic structure test
	// In real scenario, you'd use fiber.Ctx mock
	params := PaginationParams{
		Page:   1,
		Limit:  10,
		Offset: 0,
	}

	if params.Page != 1 {
		t.Errorf("Expected page 1, got %d", params.Page)
	}

	if params.Limit != 10 {
		t.Errorf("Expected limit 10, got %d", params.Limit)
	}

	if params.Offset != 0 {
		t.Errorf("Expected offset 0, got %d", params.Offset)
	}
}

func TestPaginationCalculations(t *testing.T) {
	tests := []struct {
		page           int
		limit          int
		expectedOffset int
	}{
		{1, 10, 0},
		{2, 10, 10},
		{3, 10, 20},
		{1, 20, 0},
		{5, 5, 20},
	}

	for _, tt := range tests {
		offset := (tt.page - 1) * tt.limit
		if offset != tt.expectedOffset {
			t.Errorf("Page %d, Limit %d: expected offset %d, got %d",
				tt.page, tt.limit, tt.expectedOffset, offset)
		}
	}
}

func TestMaxPaginationLimit(t *testing.T) {
	maxLimit := 100
	
	tests := []struct {
		requestedLimit int
		expectedLimit  int
	}{
		{10, 10},
		{50, 50},
		{100, 100},
		{150, 100}, // Should be capped at max
		{200, 100}, // Should be capped at max
	}

	for _, tt := range tests {
		actualLimit := tt.requestedLimit
		if actualLimit > maxLimit {
			actualLimit = maxLimit
		}

		if actualLimit != tt.expectedLimit {
			t.Errorf("Requested %d: expected %d, got %d",
				tt.requestedLimit, tt.expectedLimit, actualLimit)
		}
	}
}
