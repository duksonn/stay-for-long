package domain

import (
	"math"
	"time"
)

// Bookings represents a collection of Booking pointers
type Bookings []*Booking

// Booking represents a hotel booking with its essential information
type Booking struct {
	RequestID   string
	CheckIn     time.Time
	Nights      int
	SellingRate float64
	Margin      float64
}

// StatsResult holds statistical information about booking profits
type StatsResult struct {
	AvgNight float64
	MinNight float64
	MaxNight float64
}

// MaximizeResult contains the optimal booking combination and its statistics
type MaximizeResult struct {
	RequestIDs  []string
	TotalProfit float64
	AvgNight    float64
	MinNight    float64
	MaxNight    float64
}

// ProfitPerNight calculates the profit per night for a booking
func (b *Booking) ProfitPerNight() float64 {
	if b.Nights == 0 {
		return 0
	}

	return roundToTwoDecimals((b.SellingRate * b.Margin / 100) / float64(b.Nights))
}

// ProfitsPerNight returns an array of profit per night for all bookings
func (bb Bookings) ProfitsPerNight() []float64 {
	profits := make([]float64, 0, len(bb))
	for _, b := range bb {
		profits = append(profits, b.ProfitPerNight())
	}

	return profits
}

// CalculateStats computes statistical information about the bookings' profits
func (bb Bookings) CalculateStats() *StatsResult {
	if len(bb) == 0 {
		return &StatsResult{}
	}

	profits := bb.ProfitsPerNight()
	sum, minVal, maxVal := 0.0, profits[0], profits[0]
	for _, p := range profits {
		sum += p
		if p < minVal {
			minVal = p
		}
		if p > maxVal {
			maxVal = p
		}
	}
	avg := sum / float64(len(profits))

	return &StatsResult{
		AvgNight: roundToTwoDecimals(avg),
		MinNight: roundToTwoDecimals(minVal),
		MaxNight: roundToTwoDecimals(maxVal),
	}
}

// OverlapsWith checks if two bookings have overlapping dates
func (b *Booking) OverlapsWith(other *Booking) bool {
	bEnd := b.CheckIn.AddDate(0, 0, b.Nights)
	oEnd := other.CheckIn.AddDate(0, 0, other.Nights)

	return b.CheckIn.Before(oEnd) && other.CheckIn.Before(bEnd)
}

// HasOverlaps checks if any bookings in the collection overlap with each other
func (bb Bookings) HasOverlaps() bool {
	for i := 0; i < len(bb); i++ {
		for j := i + 1; j < len(bb); j++ {
			if bb[i].OverlapsWith(bb[j]) {
				return true
			}
		}
	}

	return false
}

// TotalProfit calculates the total profit for all bookings
func (bb Bookings) TotalProfit() float64 {
	sum := 0.0
	for _, b := range bb {
		sum += b.SellingRate * b.Margin / 100
	}

	return roundToTwoDecimals(sum)
}

// RequestIDs returns an array of all booking request IDs
func (bb Bookings) RequestIDs() []string {
	ids := make([]string, 0, len(bb))
	for _, b := range bb {
		ids = append(ids, b.RequestID)
	}

	return ids
}

// allCombinations generates all possible combinations of bookings
func allCombinations(bookings []*Booking) []Bookings {
	n := len(bookings)
	var result []Bookings
	total := int(math.Pow(2, float64(n)))
	for mask := 1; mask < total; mask++ {
		var combo Bookings
		for j := 0; j < n; j++ {
			if (mask>>j)&1 == 1 {
				combo = append(combo, bookings[j])
			}
		}
		result = append(result, combo)
	}

	return result
}

// MaximizeProfit finds the optimal combination of non-overlapping bookings that maximizes profit
func MaximizeProfit(bookings []*Booking) *MaximizeResult {
	best := findBestCombination(allCombinations(bookings))
	return buildMaximizeResult(best)
}

// findBestCombination finds the combination of bookings with the highest profit and no overlaps
func findBestCombination(combos []Bookings) Bookings {
	var best Bookings
	maxProfit := -1.0
	for _, combo := range combos {
		if combo.HasOverlaps() {
			continue
		}
		profit := combo.TotalProfit()
		if profit > maxProfit {
			maxProfit = profit
			best = combo
		}
	}

	return best
}

// buildMaximizeResult constructs the final result with statistics for the best combination
func buildMaximizeResult(best Bookings) *MaximizeResult {
	if len(best) == 0 {
		return &MaximizeResult{
			RequestIDs:  []string{},
			TotalProfit: 0,
			AvgNight:    0,
			MinNight:    0,
			MaxNight:    0,
		}
	}

	stats := best.CalculateStats()
	return &MaximizeResult{
		RequestIDs:  best.RequestIDs(),
		TotalProfit: best.TotalProfit(),
		AvgNight:    stats.AvgNight,
		MinNight:    stats.MinNight,
		MaxNight:    stats.MaxNight,
	}
}

// roundToTwoDecimals rounds a float64 to two decimal places
func roundToTwoDecimals(n float64) float64 {
	return math.Round(n*100) / 100
}
