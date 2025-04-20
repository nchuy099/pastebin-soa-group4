package router

import (
	"net/http"

	"create-paste-service-go/handler"

	"github.com/julienschmidt/httprouter"
)

// SetupRouter configures all the routes for the application
func SetupRouter() *httprouter.Router {
	router := httprouter.New()

	// Paste creation endpoint
	router.POST("/api/paste", handler.CreatePaste)

	// Health check
	router.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return router
}
