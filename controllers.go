package main

import "strings"
import "html/template"
import "fmt"
import "net/http"
import "github.com/gorilla/mux"

var templates = template.Must(template.ParseGlob("templates/*.html"))

func SearchHandler(response http.ResponseWriter, request *http.Request) {
	if Authorize(response, request) {
		type SearchPage struct {
			Title string
			Conf  Config
			Jobs  []OozieJob
		}

		query := strings.ToLower(request.PostFormValue("query"))

		response.Header().Set("Content-type", "text/html")
		var matching_jobs []OozieJob
		//oh god don't look at me IM UGLY
		for _, bundle := range RunningBundles() {
			for _, coord := range FlowDefinition(bundle.Id).Coordinators {
				if strings.Contains(strings.ToLower(coord.Name), query) {
					matching_jobs = append(matching_jobs, OozieJob{Name: coord.Name, Bundle: bundle.Name})
				}
			}
		}
		context := SearchPage{Conf: *Conf, Title: query, Jobs: matching_jobs}
		templates.ExecuteTemplate(response, "search.html", context)

	}
}

func FlowHandler(response http.ResponseWriter, request *http.Request) {
	if Authorize(response, request) {

		type FlowPage struct {
			Title string
			Flows []OozieJob
			Dag   JobDefinition
			Conf  Config
		}

		vars := mux.Vars(request)
		flow := vars["name"]
		response.Header().Set("Content-type", "text/html")

		flow_history := FlowHistory(flow)
		dag := FlowDefinition(flow_history[0].Id)
		context := FlowPage{Conf: *Conf, Title: flow, Flows: flow_history, Dag: dag}
		templates.ExecuteTemplate(response, "flow.html", context)
	}
}

func IndexHandler(response http.ResponseWriter, request *http.Request) {
	if Authorize(response, request) {

		type Page struct {
			Title string
			Jobs  map[string][]OozieJob
			Conf  Config
		}

		response.Header().Set("Content-type", "text/html")
		jobs := make(map[string][]OozieJob)
		jobs["running"] = RunningJobs()
		jobs["failed"] = FailedJobs()
		jobs["successful"] = SuccessfulJobs()
		context := Page{Conf: *Conf, Title: "Home", Jobs: jobs}

		templates.ExecuteTemplate(response, "index.html", context)
	}
}

func LogHandler(response http.ResponseWriter, request *http.Request) {
	if Authorize(response, request) {

		type LogPage struct {
			Title string
			Job   OozieJob
			Conf  Config
		}

		vars := mux.Vars(request)
		id := vars["id"]
		job := FindJobById(id)
		context := LogPage{Conf: *Conf, Title: fmt.Sprintf("Logs for %s", job.Name), Job: job}
		templates.ExecuteTemplate(response, "logs.html", context)
	}
}
