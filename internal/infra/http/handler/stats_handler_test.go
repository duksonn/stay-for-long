package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/duksonn/stay-for-long/internal/domain"
	"github.com/duksonn/stay-for-long/internal/infra/http/handler"
	"github.com/duksonn/stay-for-long/internal/mocks"
	"github.com/duksonn/stay-for-long/internal/ports"
)

func TestNewStatsHandler(t *testing.T) {
	tests := []struct {
		name        string
		statsSvc    ports.StatsService
		wantErr     error
		description string
	}{
		{
			name:        "successful creation",
			statsSvc:    mocks.NewMockStatsService(gomock.NewController(t)),
			wantErr:     nil,
			description: "should create handler successfully with valid service",
		},
		{
			name:        "nil service",
			statsSvc:    nil,
			wantErr:     handler.ErrNilStatsService,
			description: "should return error when service is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := handler.NewStatsHandler(tt.statsSvc)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, h)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, h)
			}
		})
	}
}

func TestStatsHandler_HandlerCalculateStats(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    []map[string]interface{}
		mock           func(*mocks.MockStatsService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful calculation",
			requestBody: []map[string]interface{}{
				{
					"request_id":   "bookata_XY123",
					"check_in":     "2020-01-01",
					"nights":       5,
					"selling_rate": 200,
					"margin":       20,
				},
				{
					"request_id":   "kayete_PP234",
					"check_in":     "2020-01-04",
					"nights":       4,
					"selling_rate": 156,
					"margin":       22,
				},
			},
			mock: func(m *mocks.MockStatsService) {
				m.EXPECT().
					CalculateStats(gomock.Any()).
					Return(&domain.StatsResult{
						AvgNight: 178,
						MinNight: 156,
						MaxNight: 200,
					})
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"avg_night": float64(178),
				"min_night": float64(156),
				"max_night": float64(200),
			},
		},
		{
			name:           "invalid json",
			requestBody:    nil,
			mock:           func(m *mocks.MockStatsService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "invalid date format",
			requestBody: []map[string]interface{}{
				{
					"request_id":   "bookata_XY123",
					"check_in":     "invalid-date",
					"nights":       5,
					"selling_rate": 200,
					"margin":       20,
				},
			},
			mock:           func(m *mocks.MockStatsService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStatsService := mocks.NewMockStatsService(ctrl)
			h, _ := handler.NewStatsHandler(mockStatsService)

			tt.mock(mockStatsService)

			var body []byte
			var err error
			if tt.requestBody != nil {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			} else {
				body = []byte("invalid json")
			}

			req := httptest.NewRequest(http.MethodPost, "/stats", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			h.HandlerCalculateStats(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err = json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)
			for k, v := range tt.expectedBody {
				assert.Equal(t, v, response[k])
			}
		})
	}
}

func TestStatsHandler_HandlerMaximizeProfit(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    []map[string]interface{}
		mock           func(*mocks.MockStatsService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful maximization",
			requestBody: []map[string]interface{}{
				{
					"request_id":   "bookata_XY123",
					"check_in":     "2020-01-01",
					"nights":       5,
					"selling_rate": 200,
					"margin":       20,
				},
				{
					"request_id":   "kayete_PP234",
					"check_in":     "2020-01-04",
					"nights":       4,
					"selling_rate": 156,
					"margin":       22,
				},
			},
			mock: func(m *mocks.MockStatsService) {
				m.EXPECT().
					MaximizeProfit(gomock.Any()).
					Return(&domain.MaximizeResult{
						RequestIDs:  []string{"bookata_XY123"},
						TotalProfit: 200,
						AvgNight:    200,
						MinNight:    200,
						MaxNight:    200,
					})
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"request_ids":  []interface{}{"bookata_XY123"},
				"total_profit": float64(200),
				"avg_night":    float64(200),
				"min_night":    float64(200),
				"max_night":    float64(200),
			},
		},
		{
			name:           "invalid json",
			requestBody:    nil,
			mock:           func(m *mocks.MockStatsService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "invalid date format",
			requestBody: []map[string]interface{}{
				{
					"request_id":   "bookata_XY123",
					"check_in":     "invalid-date",
					"nights":       5,
					"selling_rate": 200,
					"margin":       20,
				},
			},
			mock:           func(m *mocks.MockStatsService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStatsService := mocks.NewMockStatsService(ctrl)
			h, _ := handler.NewStatsHandler(mockStatsService)

			tt.mock(mockStatsService)

			var body []byte
			var err error
			if tt.requestBody != nil {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			} else {
				body = []byte("invalid json")
			}

			req := httptest.NewRequest(http.MethodPost, "/maximize", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			h.HandlerMaximizeProfit(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err = json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)
			for k, v := range tt.expectedBody {
				assert.Equal(t, v, response[k])
			}
		})
	}
}
