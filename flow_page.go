package main

import "net/http"
import "github.com/gorilla/mux"

type FlowPage struct {
	Title string
	Flows []OozieJob
	Dag   WorkflowDAG
	Conf  Config
}

func FlowHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	flow := vars["name"]
	response.Header().Set("Content-type", "text/html")

	flow_history := FlowHistory(flow)
	dag := FlowDefinition(flow_history[0].Id)
	context := FlowPage{Conf: *Conf, Title: flow, Flows: flow_history, Dag: dag}
	templates.ExecuteTemplate(response, "flow.html", context)
}
