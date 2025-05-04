package router

import (
	"net/http"

	"get-stats-service/handler"

	"github.com/julienschmidt/httprouter"
)

// SetupRouter configures all the routes for the application
func SetupRouter() *httprouter.Router {
	router := httprouter.New()

	// Stats endpoint
	router.GET("/api/paste/stats", handler.GetMonthlyStats)

	// Health check
	router.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return router
}
