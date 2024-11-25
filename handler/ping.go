package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func CreatePingHandlers(router *mux.Router) {
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		defer handleError(w)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}).Methods("GET")
}
