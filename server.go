package main

import "github.com/gorilla/mux"
import "net/http"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/flow/{name}", FlowHandler)
	r.HandleFunc("/all_logs/{id}", LogHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", r)
	log.Info("Running on localhost:8080")
	http.ListenAndServe(":8080", r)
}
