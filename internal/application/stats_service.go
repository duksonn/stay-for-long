package application

import (
	"github.com/duksonn/stay-for-long/internal/domain"
	"github.com/duksonn/stay-for-long/internal/ports"
)

// Ensure StatsService implements the ports.StatsService interface
var _ ports.StatsService = (*StatsService)(nil)

// StatsService implements the ports.StatsService interface and provides stats management functionality
// It handles the business logic for calculating booking statistics and maximizing profit
type StatsService struct{}

// NewStatsService creates and returns a new instance of StatsService
func NewStatsService() *StatsService {
	return &StatsService{}
}

// CalculateStats computes the average, minimum, and maximum nightly rates for a set of bookings
func (s StatsService) CalculateStats(requests domain.Bookings) *domain.StatsResult {
	return requests.CalculateStats()
}

// MaximizeProfit finds the optimal combination of bookings that maximizes total profit
// while ensuring no booking periods overlap
func (s StatsService) MaximizeProfit(requests domain.Bookings) *domain.MaximizeResult {
	return domain.MaximizeProfit(requests)
}
