package main

import "net/http"
import "fmt"
import "github.com/gorilla/mux"

type LogPage struct {
	Title string
	Job   OozieJob
	Conf  Config
}

func LogHandler(response http.ResponseWriter, request *http.Request) {
	Authorize(response, request)

	vars := mux.Vars(request)
	id := vars["id"]
	job := FindJobById(id)
	context := LogPage{Conf: *Conf, Title: fmt.Sprintf("Logs for %s", job.Name), Job: job}
	templates.ExecuteTemplate(response, "logs.html", context)
}
