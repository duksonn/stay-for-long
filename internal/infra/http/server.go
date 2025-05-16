package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/duksonn/stay-for-long/cmd/di"
	"github.com/duksonn/stay-for-long/internal/infra/http/handler"
)

func Routes(deps *di.Dependencies) (*mux.Router, error) {
	router := mux.NewRouter()

	// Stats endpoints
	statsHandler, err := handler.NewStatsHandler(deps.StatsSvc)
	if err != nil {
		return nil, err
	}

	router.HandleFunc("/stats", statsHandler.HandlerCalculateStats).Methods(http.MethodPost)
	router.HandleFunc("/maximize", statsHandler.HandlerMaximizeProfit).Methods(http.MethodPost)

	return router, nil
}
