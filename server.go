package main

import "github.com/gorilla/mux"
import "net/http"
import "io/ioutil"
import "encoding/json"

type Config struct {
	OozieURL string `json:"oozie_url"`
	HueURL   string `json:"hue"`
}

var Conf = new(Config)

func main() {
	ReadConfig()

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/flow/{name}", FlowHandler)
	r.HandleFunc("/all_logs/{id}", LogHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", r)
	log.Info("Running on localhost:8080")
	http.ListenAndServe(":8080", r)
}

func ReadConfig() {
	buf, err := ioutil.ReadFile("config.json")
	check(err)

	err = json.Unmarshal(buf, &Conf)
	check(err)

}
