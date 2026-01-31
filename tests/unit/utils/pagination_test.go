package utils_test

import (
	"golang-backend/utils"
	"testing"
)

func TestCalculatePagination(t *testing.T) {
	tests := []struct {
		name        string
		total       int64
		page        int
		perPage     int
		wantLast    int
		wantFrom    int
		wantTo      int
		wantHasMore bool
	}{
		{
			name:        "First page, full",
			total:       50,
			page:        1,
			perPage:     10,
			wantLast:    5,
			wantFrom:    1,
			wantTo:      10,
			wantHasMore: true,
		},
		{
			name:        "Last page, partial",
			total:       55,
			page:        6,
			perPage:     10,
			wantLast:    6,
			wantFrom:    51,
			wantTo:      55,
			wantHasMore: false,
		},
		{
			name:        "Empty result",
			total:       0,
			page:        1,
			perPage:     10,
			wantLast:    1,
			wantFrom:    0,
			wantTo:      0,
			wantHasMore: false,
		},
		{
			name:        "Page out of range",
			total:       20,
			page:        3,
			perPage:     10,
			wantLast:    2,
			wantFrom:    0,
			wantTo:      0,
			wantHasMore: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.CalculatePagination(tt.total, tt.page, tt.perPage)
			if got.LastPage != tt.wantLast {
				t.Errorf("LastPage = %v, want %v", got.LastPage, tt.wantLast)
			}
			if got.From != tt.wantFrom {
				t.Errorf("From = %v, want %v", got.From, tt.wantFrom)
			}
			if got.To != tt.wantTo {
				t.Errorf("To = %v, want %v", got.To, tt.wantTo)
			}
			if got.HasMorePages != tt.wantHasMore {
				t.Errorf("HasMorePages = %v, want %v", got.HasMorePages, tt.wantHasMore)
			}
		})
	}
}
