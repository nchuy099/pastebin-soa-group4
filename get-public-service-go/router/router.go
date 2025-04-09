package router

import (
	"net/http"

	"get-public-service/handler"

	"github.com/julienschmidt/httprouter"
)

func SetupRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/paste", handler.GetPublicPastes)

	// Health check
	router.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return router
}
