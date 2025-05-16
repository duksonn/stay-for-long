package ports

import "github.com/duksonn/stay-for-long/internal/domain"

// StatsService defines the interface for handling stats business operations
type StatsService interface {
	// CalculateStats computes the average, minimum, and maximum nightly rates for a set of bookings
	CalculateStats(requests domain.Bookings) *domain.StatsResult

	// MaximizeProfit finds the optimal combination of bookings that maximizes total profit
	// while ensuring no booking periods overlap
	MaximizeProfit(requests domain.Bookings) *domain.MaximizeResult
}
