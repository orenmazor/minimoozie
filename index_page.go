package main

import "html/template"
import "net/http"

var templates = template.Must(template.ParseGlob("templates/*.html"))

type Page struct {
	Title string
	Jobs  map[string][]OozieJob
}

func IndexHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "text/html")
	jobs := make(map[string][]OozieJob)
	jobs["running"] = RunningJobs()
	jobs["failed"] = FailedJobs()
	jobs["successful"] = SuccessfulJobs()
	context := Page{Title: "Home", Jobs: jobs}
	templates.ExecuteTemplate(response, "index.html", context)
}
