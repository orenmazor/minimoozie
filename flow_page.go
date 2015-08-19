package main

import "net/http"
import "github.com/gorilla/mux"

type FlowPage struct {
	Title string
	Flows []OozieJob
}

func FlowHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	flow := vars["name"]
	response.Header().Set("Content-type", "text/html")

	flow_history := FlowHistory(flow)
	context := FlowPage{Title: flow, Flows: flow_history}
	templates.ExecuteTemplate(response, "flow.html", context)
}
