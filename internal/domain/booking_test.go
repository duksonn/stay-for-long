package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/duksonn/stay-for-long/internal/domain"
)

func TestBooking_ProfitPerNight(t *testing.T) {
	tests := []struct {
		name     string
		booking  *domain.Booking
		expected float64
	}{
		{
			name: "normal booking",
			booking: &domain.Booking{
				SellingRate: 1000,
				Margin:      20,
				Nights:      5,
			},
			expected: 40, // (1000 * 20%) / 5
		},
		{
			name: "zero nights",
			booking: &domain.Booking{
				SellingRate: 1000,
				Margin:      20,
				Nights:      0,
			},
			expected: 0,
		},
		{
			name: "zero margin",
			booking: &domain.Booking{
				SellingRate: 1000,
				Margin:      0,
				Nights:      5,
			},
			expected: 0,
		},
		{
			name: "zero selling rate",
			booking: &domain.Booking{
				SellingRate: 0,
				Margin:      20,
				Nights:      5,
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.booking.ProfitPerNight())
		})
	}
}

func TestBookings_ProfitsPerNight(t *testing.T) {
	bookings := domain.Bookings{
		{
			SellingRate: 1000,
			Margin:      20,
			Nights:      5,
		},
		{
			SellingRate: 2000,
			Margin:      15,
			Nights:      4,
		},
	}

	expected := []float64{40, 75} // (1000 * 20%) / 5, (2000 * 15%) / 4
	result := bookings.ProfitsPerNight()

	assert.Equal(t, expected, result)
}

func TestBookings_CalculateStats(t *testing.T) {
	tests := []struct {
		name     string
		bookings domain.Bookings
		expected *domain.StatsResult
	}{
		{
			name: "multiple bookings",
			bookings: domain.Bookings{
				{SellingRate: 1000, Margin: 20, Nights: 5},
				{SellingRate: 2000, Margin: 15, Nights: 4},
				{SellingRate: 3000, Margin: 25, Nights: 6},
			},
			expected: &domain.StatsResult{
				AvgNight: 80, // (40 + 75 + 125) / 3
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
				{SellingRate: 1000, Margin: 20, Nights: 5},
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
			result := tt.bookings.CalculateStats()
			assert.Equal(t, tt.expected.AvgNight, result.AvgNight)
			assert.Equal(t, tt.expected.MinNight, result.MinNight)
			assert.Equal(t, tt.expected.MaxNight, result.MaxNight)
		})
	}
}

func TestBooking_OverlapsWith(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		booking1 *domain.Booking
		booking2 *domain.Booking
		expected bool
	}{
		{
			name: "no overlap - sequential",
			booking1: &domain.Booking{
				CheckIn: baseTime,
				Nights:  3,
			},
			booking2: &domain.Booking{
				CheckIn: baseTime.AddDate(0, 0, 4),
				Nights:  3,
			},
			expected: false,
		},
		{
			name: "overlap - same dates",
			booking1: &domain.Booking{
				CheckIn: baseTime,
				Nights:  3,
			},
			booking2: &domain.Booking{
				CheckIn: baseTime,
				Nights:  3,
			},
			expected: true,
		},
		{
			name: "overlap - partial",
			booking1: &domain.Booking{
				CheckIn: baseTime,
				Nights:  5,
			},
			booking2: &domain.Booking{
				CheckIn: baseTime.AddDate(0, 0, 3),
				Nights:  5,
			},
			expected: true,
		},
		{
			name: "overlap - contained",
			booking1: &domain.Booking{
				CheckIn: baseTime,
				Nights:  7,
			},
			booking2: &domain.Booking{
				CheckIn: baseTime.AddDate(0, 0, 2),
				Nights:  3,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.booking1.OverlapsWith(tt.booking2))
			assert.Equal(t, tt.expected, tt.booking2.OverlapsWith(tt.booking1))
		})
	}
}

func TestBookings_HasOverlaps(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		bookings domain.Bookings
		expected bool
	}{
		{
			name: "no overlaps",
			bookings: domain.Bookings{
				{CheckIn: baseTime, Nights: 3},
				{CheckIn: baseTime.AddDate(0, 0, 4), Nights: 3},
				{CheckIn: baseTime.AddDate(0, 0, 8), Nights: 3},
			},
			expected: false,
		},
		{
			name: "has overlaps",
			bookings: domain.Bookings{
				{CheckIn: baseTime, Nights: 3},
				{CheckIn: baseTime.AddDate(0, 0, 2), Nights: 3},
				{CheckIn: baseTime.AddDate(0, 0, 8), Nights: 3},
			},
			expected: true,
		},
		{
			name:     "empty bookings",
			bookings: domain.Bookings{},
			expected: false,
		},
		{
			name: "single booking",
			bookings: domain.Bookings{
				{CheckIn: baseTime, Nights: 3},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.bookings.HasOverlaps())
		})
	}
}

func TestBookings_TotalProfit(t *testing.T) {
	tests := []struct {
		name     string
		bookings domain.Bookings
		expected float64
	}{
		{
			name: "multiple bookings",
			bookings: domain.Bookings{
				{SellingRate: 1000, Margin: 20},
				{SellingRate: 2000, Margin: 15},
				{SellingRate: 3000, Margin: 25},
			},
			expected: 1250, // (1000 * 20%) + (2000 * 15%) + (3000 * 25%)
		},
		{
			name:     "empty bookings",
			bookings: domain.Bookings{},
			expected: 0,
		},
		{
			name: "zero margin",
			bookings: domain.Bookings{
				{SellingRate: 1000, Margin: 0},
				{SellingRate: 2000, Margin: 0},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.bookings.TotalProfit())
		})
	}
}

func TestBookings_RequestIDs(t *testing.T) {
	tests := []struct {
		name     string
		bookings domain.Bookings
		expected []string
	}{
		{
			name: "multiple bookings",
			bookings: domain.Bookings{
				{RequestID: "req1"},
				{RequestID: "req2"},
				{RequestID: "req3"},
			},
			expected: []string{"req1", "req2", "req3"},
		},
		{
			name:     "empty bookings",
			bookings: domain.Bookings{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.bookings.RequestIDs())
		})
	}
}

func TestMaximizeProfit(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		bookings []*domain.Booking
		expected *domain.MaximizeResult
	}{
		{
			name: "optimal combination with overlaps",
			bookings: []*domain.Booking{
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
				TotalProfit: 950, // (2000 * 25%) + (1500 * 30%)
				AvgNight:    158.33,
				MinNight:    150,
				MaxNight:    166.67,
			},
		},
		{
			name:     "empty bookings",
			bookings: []*domain.Booking{},
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
			bookings: []*domain.Booking{
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
				TotalProfit: 700, // (1000 * 20%) + (2000 * 25%)
				AvgNight:    116.67,
				MinNight:    66.67,
				MaxNight:    166.67,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := domain.MaximizeProfit(tt.bookings)
			require.NotNil(t, result)

			assert.Equal(t, tt.expected.RequestIDs, result.RequestIDs)
			assert.Equal(t, tt.expected.TotalProfit, result.TotalProfit)
			assert.Equal(t, tt.expected.AvgNight, result.AvgNight)
			assert.Equal(t, tt.expected.MinNight, result.MinNight)
			assert.Equal(t, tt.expected.MaxNight, result.MaxNight)
		})
	}
}
