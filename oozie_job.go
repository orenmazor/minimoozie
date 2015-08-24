package main

type OozieBundle struct {
	Status       string       `json:"status"`
	Name         string       `json:"bundlejobname"`
	Id           string       `json:"bundlejobid"`
	Coordinators []OozieCoord `string:"coordinator"`
}

type OozieCoord struct {
	Name string `xml:"attr,name"`
}

type OozieAction struct {
	Id           string `json:"id"`
	ExternalId   string `json:"externalid"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	StartTime    string `json:"starttime"`
	EndTime      string `json:"endtime"`
	Status       string `json:"externalstatus"`
	TrackerURI   string `json:"trackeruri"`
	ConsoleUrl   string `json:"consoleurl"`
	Config       string `json:"conf"`
	ErrorMessage string `json:"errormessage"`
}

type OozieJob struct {
	Coordinator  string `json:"parentId"`
	Bundle       string
	Name         string        `json:"appName"`
	Id           string        `json:"id"`
	Status       string        `json:"status"`
	StartTime    string        `json:"startTime"`
	EndTime      string        `json:"endTime"`
	ConsoleURL   string        `json:"consoleurl"`
	Actions      []OozieAction `json:"actions"`
	ErrorMessage string        `json:"errormessage"`
	User         string        `json:"user"`
}

func (job *OozieJob) Errors() []OozieAction {
	var erroredActions []OozieAction

	for _, action := range job.Actions {
		if action.Status == "ERROR" || action.Status == "KILLED" || action.Status == "FAILED/KILLED" {
			erroredActions = append(erroredActions, action)
		}
	}

	return erroredActions
}
