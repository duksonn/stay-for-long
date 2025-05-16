package di

import "github.com/duksonn/stay-for-long/internal/application"

// Dependencies list the use cases application services of the system
type Dependencies struct {
	StatsSvc *application.StatsService
}

// Init return the initialized dependencies of the system
func Init() *Dependencies {
	// Services
	statsSvc := application.NewStatsService()

	return &Dependencies{StatsSvc: statsSvc}
}
