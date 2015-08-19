package main

import "net/http"
import "fmt"
import "github.com/gorilla/mux"

type LogPage struct {
	Title string
	Job   OozieJob
}

func LogHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]
	job := FindJobById(id)
	context := LogPage{Title: fmt.Sprintf("Logs for %s", job.Name), Job: job}
	templates.ExecuteTemplate(response, "logs.html", context)
}
