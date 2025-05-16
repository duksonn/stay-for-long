package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/duksonn/stay-for-long/internal/domain"
	"github.com/duksonn/stay-for-long/internal/ports"
)

var (
	// ErrNilStatsService is returned when the stats service is nil
	ErrNilStatsService = errors.New("stats service cannot be nil")
	// ErrInvalidRequest is returned when the request body is invalid
	ErrInvalidRequest = errors.New("invalid request body")
	// ErrInvalidJSON is returned when the json is invalid
	ErrInvalidJSON = errors.New("invalid request json")
	// ErrInvalidDateFormat is returned when the date has invalid format
	ErrInvalidDateFormat = errors.New("invalid date format")
)

// StatsHandler handles HTTP requests for stats-related operations
// It provides endpoints for calculating booking statistics and maximizing profit
type StatsHandler struct {
	statsService ports.StatsService
}

// NewStatsHandler creates a new instance of StatsHandler
// Returns ErrNilStatsService if the stats service is nil
func NewStatsHandler(statsSvc ports.StatsService) (*StatsHandler, error) {
	if statsSvc == nil {
		return nil, ErrNilStatsService
	}

	return &StatsHandler{statsService: statsSvc}, nil
}

// HandlerCalculateStats processes HTTP requests to calculate booking statistics
// It accepts a list of booking requests and returns average, minimum, and maximum nightly rates
func (h *StatsHandler) HandlerCalculateStats(w http.ResponseWriter, r *http.Request) {
	var dtos []bookingRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrInvalidRequest)
		return
	}
	if err := json.Unmarshal(body, &dtos); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrInvalidJSON)
		return
	}

	requests, err := parseBookingRequests(dtos)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	stats := h.statsService.CalculateStats(requests)
	response := statsResultResponse{
		AvgNight: stats.AvgNight,
		MinNight: stats.MinNight,
		MaxNight: stats.MaxNight,
	}
	writeJSONResponse(w, http.StatusOK, response)
}

// HandlerMaximizeProfit processes HTTP requests to find the optimal booking combination
// that maximizes profit while avoiding booking overlaps
func (h *StatsHandler) HandlerMaximizeProfit(w http.ResponseWriter, r *http.Request) {
	var dtos []bookingRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrInvalidRequest)
		return
	}
	if err := json.Unmarshal(body, &dtos); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrInvalidJSON)
		return
	}

	requests, err := parseBookingRequests(dtos)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	result := h.statsService.MaximizeProfit(requests)
	response := maximizeResultResponse{
		RequestIDs:  result.RequestIDs,
		TotalProfit: result.TotalProfit,
		AvgNight:    result.AvgNight,
		MinNight:    result.MinNight,
		MaxNight:    result.MaxNight,
	}
	writeJSONResponse(w, http.StatusOK, response)
}

// parseBookingRequests converts a slice of bookingRequest DTOs to domain.Booking objects
// It handles date parsing and validation of the input data
func parseBookingRequests(dtos []bookingRequest) ([]*domain.Booking, error) {
	requests := make([]*domain.Booking, 0, len(dtos))
	for _, dto := range dtos {
		checkIn, err := time.Parse(time.DateOnly, dto.CheckIn)
		if err != nil {
			return nil, ErrInvalidDateFormat
		}
		requests = append(requests, &domain.Booking{
			RequestID:   dto.RequestID,
			CheckIn:     checkIn,
			Nights:      dto.Nights,
			SellingRate: dto.SellingRate,
			Margin:      dto.Margin,
		})
	}

	return requests, nil
}

// writeJSONResponse is a helper function to write JSON responses
// It sets the appropriate headers and handles JSON encoding errors
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
