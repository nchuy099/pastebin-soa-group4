package router

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"get-paste-service/handler"
)

func SetupRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET("/paste/:id", handler.GetPasteHandler)

	// Health check
	router.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return router
}
