package api

import "net/http"

func RegisterRoutes(handler *Handler) {
	http.HandleFunc("/api/entries", handler.CreateEntry)
	http.HandleFunc("/api/state", handler.GetState)
	http.HandleFunc("/api/results", handler.GetResults)
	fileServer := http.FileServer(http.Dir("./web"))

	http.Handle("/", fileServer)
}
