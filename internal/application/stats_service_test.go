package application_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/duksonn/stay-for-long/internal/application"
	"github.com/duksonn/stay-for-long/internal/domain"
)

func TestStatsService_CalculateStats(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	service := application.NewStatsService()

	tests := []struct {
		name     string
		bookings domain.Bookings
		expected *domain.StatsResult
	}{
		{
			name: "multiple bookings",
			bookings: domain.Bookings{
				{
					CheckIn:     baseTime,
					Nights:      5,
					SellingRate: 1000,
					Margin:      20,
				},
				{
					CheckIn:     baseTime.AddDate(0, 0, 6),
					Nights:      4,
					SellingRate: 2000,
					Margin:      15,
				},
				{
					CheckIn:     baseTime.AddDate(0, 0, 11),
					Nights:      6,
					SellingRate: 3000,
					Margin:      25,
				},
			},
			expected: &domain.StatsResult{
				AvgNight: 80,
				MinNight: 40,
				MaxNight: 125,
			},
		},
		{
			name:     "empty bookings",
			bookings: domain.Bookings{},
			expected: &domain.StatsResult{},
		},
		{
			name: "single booking",
			bookings: domain.Bookings{
				{
					CheckIn:     baseTime,
					Nights:      5,
					SellingRate: 1000,
					Margin:      20,
				},
			},
			expected: &domain.StatsResult{
				AvgNight: 40,
				MinNight: 40,
				MaxNight: 40,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CalculateStats(tt.bookings)
			require.NotNil(t, result)

			assert.Equal(t, tt.expected.AvgNight, result.AvgNight)
			assert.Equal(t, tt.expected.MinNight, result.MinNight)
			assert.Equal(t, tt.expected.MaxNight, result.MaxNight)
		})
	}
}

func TestStatsService_MaximizeProfit(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	service := application.NewStatsService()

	tests := []struct {
		name     string
		bookings domain.Bookings
		expected *domain.MaximizeResult
	}{
		{
			name: "optimal combination with overlaps",
			bookings: domain.Bookings{
				{
					RequestID:   "req1",
					CheckIn:     baseTime,
					Nights:      3,
					SellingRate: 1000,
					Margin:      20,
				},
				{
					RequestID:   "req2",
					CheckIn:     baseTime.AddDate(0, 0, 2),
					Nights:      3,
					SellingRate: 2000,
					Margin:      25,
				},
				{
					RequestID:   "req3",
					CheckIn:     baseTime.AddDate(0, 0, 6),
					Nights:      3,
					SellingRate: 1500,
					Margin:      30,
				},
			},
			expected: &domain.MaximizeResult{
				RequestIDs:  []string{"req2", "req3"},
				TotalProfit: 950,
				AvgNight:    158.33,
				MinNight:    150,
				MaxNight:    166.67,
			},
		},
		{
			name:     "empty bookings",
			bookings: domain.Bookings{},
			expected: &domain.MaximizeResult{
				RequestIDs:  []string{},
				TotalProfit: 0,
				AvgNight:    0,
				MinNight:    0,
				MaxNight:    0,
			},
		},
		{
			name: "no overlaps possible",
			bookings: domain.Bookings{
				{
					RequestID:   "req1",
					CheckIn:     baseTime,
					Nights:      3,
					SellingRate: 1000,
					Margin:      20,
				},
				{
					RequestID:   "req2",
					CheckIn:     baseTime.AddDate(0, 0, 4),
					Nights:      3,
					SellingRate: 2000,
					Margin:      25,
				},
			},
			expected: &domain.MaximizeResult{
				RequestIDs:  []string{"req1", "req2"},
				TotalProfit: 700,
				AvgNight:    116.67,
				MinNight:    66.67,
				MaxNight:    166.67,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.MaximizeProfit(tt.bookings)
			require.NotNil(t, result)

			assert.Equal(t, tt.expected.RequestIDs, result.RequestIDs)
			assert.Equal(t, tt.expected.TotalProfit, result.TotalProfit)
			assert.Equal(t, tt.expected.AvgNight, result.AvgNight)
			assert.Equal(t, tt.expected.MinNight, result.MinNight)
			assert.Equal(t, tt.expected.MaxNight, result.MaxNight)
		})
	}
}
