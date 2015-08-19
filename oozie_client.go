package main

import "os"
import "fmt"
import "net/http"
import "encoding/json"

type OozieJob struct {
	Coordinator string `json:"parentId"`
	Name        string `json:"appName"`
	Id          string `json:"id"`
	Status      string `json:"status"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	ConsoleURL  string `json:"consoleurl"`
}

type OozieResultSet struct {
	Total     int        `json:"total"`
	Workflows []OozieJob `json:"workflows"`
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

func getJobs(filter string) []OozieJob {
	oozieURL := os.Getenv("OOZIE_URL")
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
