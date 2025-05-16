package handler

// bookingRequest represents the structure of a booking request as received from the HTTP API
// It contains all necessary information to create a domain.Booking object
type bookingRequest struct {
	RequestID   string  `json:"request_id"`   // Unique identifier for the booking request
	CheckIn     string  `json:"check_in"`     // Check-in date in YYYY-MM-DD format
	Nights      int     `json:"nights"`       // Number of nights for the stay
	SellingRate float64 `json:"selling_rate"` // Total selling rate for the entire stay
	Margin      float64 `json:"margin"`       // Profit margin percentage
}

// statsResultResponse represents the structure of the stats calculation response
// It contains the calculated statistics for a set of bookings
type statsResultResponse struct {
	AvgNight float64 `json:"avg_night"` // Average nightly rate
	MinNight float64 `json:"min_night"` // Minimum nightly rate
	MaxNight float64 `json:"max_night"` // Maximum nightly rate
}

// maximizeResultResponse represents the structure of the profit maximization response
// It contains the optimal booking combination and its associated statistics
type maximizeResultResponse struct {
	RequestIDs  []string `json:"request_ids"`  // List of request IDs that maximize profit
	TotalProfit float64  `json:"total_profit"` // Total profit for the selected bookings
	AvgNight    float64  `json:"avg_night"`    // Average nightly rate for selected bookings
	MinNight    float64  `json:"min_night"`    // Minimum nightly rate for selected bookings
	MaxNight    float64  `json:"max_night"`    // Maximum nightly rate for selected bookings
}
