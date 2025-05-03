package router

import (
	"get-paste-service/handler"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func SetupRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/paste/:id", handler.GetPasteHandler)

	// Health check
	router.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return router
}
