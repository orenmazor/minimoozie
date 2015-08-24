package main

import "io/ioutil"
import "fmt"
import "net/http"
import "encoding/json"
import "encoding/xml"

type OozieResultSet struct {
	Total     int           `json:"total"`
	Workflows []OozieJob    `json:"workflows"`
	Bundles   []OozieBundle `json:"bundlejobs"`
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

type JobDefinition struct {
	Start        Node   `xml:"start"`
	Actions      []Node `xml:"action"`
	Coordinators []Node `xml:"coordinator"`
}

func RunningJobs() []OozieJob {
	return getJobs("status%3DRUNNING")
}

func SuccessfulJobs() []OozieJob {
	return getJobs("status%3DSUCCEEDED")
}

func FailedJobs() []OozieJob {
	//not sure about this...submits too many requests for the front page, maybe?
	jobs := getJobs("status%3DKILLED&len=10")
	detailedJobs := make([]OozieJob, len(jobs))
	for index, job := range jobs {
		detailedJobs[index] = FindJobById(job.Id)
	}
	return detailedJobs
}

func FlowHistory(flowName string) []OozieJob {
	return getJobs(fmt.Sprintf("name%%3D%s", flowName))
}

func FlowDefinition(flowId string) JobDefinition {
	oozieURL := Conf.OozieURL
	fullURL := fmt.Sprintf("%s/oozie/v1/job/%s?show=definition", oozieURL, flowId)
	log.Info(fullURL)
	resp, err := http.Get(fullURL)
	check(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	var dag JobDefinition
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

func RunningBundles() []OozieBundle {
	oozieURL := Conf.OozieURL
	fullURL := fmt.Sprintf("%s/oozie/v1/jobs?jobtype=bundle", oozieURL)
	resp, err := http.Get(fullURL)
	check(err)
	defer resp.Body.Close()

	results := new(OozieResultSet)
	err = json.NewDecoder(resp.Body).Decode(results)
	check(err)

	var runningBundles []OozieBundle

	for _, bundle := range results.Bundles {
		if bundle.Status == "RUNNING" {
			runningBundles = append(runningBundles, bundle)
		}
	}

	return runningBundles
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
