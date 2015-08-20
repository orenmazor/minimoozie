package main

import "io/ioutil"
import "fmt"
import "net/http"
import "encoding/json"
import "encoding/xml"

type OozieAction struct {
	Id         string `json:"id"`
	ExternalId string `json:"externalid"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	StartTime  string `json:"starttime"`
	EndTime    string `json:"endtime"`
	Status     string `json:"externalstatus"`
	TrackerURI string `json:"trackeruri"`
	ConsoleURI string `json:"consoleuri"`
	Config     string `json:"conf"`
}

type OozieJob struct {
	Coordinator string        `json:"parentId"`
	Name        string        `json:"appName"`
	Id          string        `json:"id"`
	Status      string        `json:"status"`
	StartTime   string        `json:"startTime"`
	EndTime     string        `json:"endTime"`
	ConsoleURL  string        `json:"consoleurl"`
	Actions     []OozieAction `json:"actions"`
}

type OozieResultSet struct {
	Total     int        `json:"total"`
	Workflows []OozieJob `json:"workflows"`
}

type Edge struct {
	To string `xml:"to,attr"`
}

type Node struct {
	To    string `xml:"to,attr"`
	Name  string `xml:"name,attr"`
	Ok    Edge   `xml:"ok"`
	Error Edge   `xml:"error"`
}

type WorkflowDAG struct {
	Start   Node   `xml:"start"`
	Actions []Node `xml:"action"`
}

func RunningJobs() []OozieJob {
	return getJobs("status%3DRUNNING")
}

func SuccessfulJobs() []OozieJob {
	return getJobs("status%3DSUCCEEDED")
}

func FailedJobs() []OozieJob {
	return getJobs("status%3DKILLED")
}

func FlowHistory(flowName string) []OozieJob {
	return getJobs(fmt.Sprintf("name%%3D%s", flowName))
}

func FlowDefinition(flowId string) WorkflowDAG {
	oozieURL := Conf.OozieURL
	fullURL := fmt.Sprintf("%s/oozie/v1/job/%s?show=definition", oozieURL, flowId)
	log.Info(fullURL)
	resp, err := http.Get(fullURL)
	check(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	var dag WorkflowDAG
	err = xml.Unmarshal(body, &dag)
	check(err)

	return dag
}

func FindJobById(flowId string) OozieJob {
	oozieURL := Conf.OozieURL
	fullURL := fmt.Sprintf("%s/oozie/v2/job/%s?show=info&timezone=GMT", oozieURL, flowId)
	log.Info(fullURL)
	resp, err := http.Get(fullURL)
	check(err)
	defer resp.Body.Close()

	result := new(OozieJob)
	err = json.NewDecoder(resp.Body).Decode(result)
	check(err)

	return *result

}

func getJobs(filter string) []OozieJob {
	oozieURL := Conf.OozieURL
	fullURL := fmt.Sprintf("%s/oozie/v1/jobs?filter=%s", oozieURL, filter)
	log.Info(fullURL)
	resp, err := http.Get(fullURL)
	check(err)

	defer resp.Body.Close()

	results := new(OozieResultSet)
	err = json.NewDecoder(resp.Body).Decode(results)
	check(err)

	log.Info(fmt.Sprintf("received %d workflows", results.Total))
	return results.Workflows

}
