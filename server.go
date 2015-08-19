package main

import "github.com/gorilla/mux"
import "net/http"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/flow/{name}", FlowHandler)
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
