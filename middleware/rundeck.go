package middleware

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

type jobs struct {
	XMLName xml.Name `xml:"jobs"`
	Count   int      `xml:"count,attr"`
	Jobs    []job    `xml:"job"`
}

type job struct {
	XMLName         xml.Name `xml:"job"`
	Href            string   `xml:"href,attr"`
	ID              string   `xml:"id,attr"`
	ScheduleEnabled bool     `xml:"scheduleEnabled,attr"`
	Scheduled       bool     `xml:"sheduled,attr"`
	Enabled         bool     `xml:"enabled,attr"`
	Permalink       string   `xml:"permalink,attr"`
	Group           string   `xml:"group"`
	Description     string   `xml:"description"`
	Project         string   `xml:"project"`
	Name            string   `xml:"name"`
}

type result struct {
	XMLName    xml.Name   `xml:"result"`
	Success    string     `xml:"success,attr"`
	Executions executions `xml:"executions"`
}

type executions struct {
	XMLName   xml.Name `xml:"executions"`
	Count     string   `xml:"count,atrr"`
	Execution []exec   `xml:"execution"`
}

type exec struct {
	XMLName xml.Name `xml:"execution"`
	ExecID  string   `xml:"id,attr"`
}

// Rundeck : Rundeck object
type Rundeck struct {
	APIToken   string
	RundeckURL string
}

//NewRundeck : Return a rundeck object
func NewRundeck(apiToken, rundeckURL string) *Rundeck {
	return &Rundeck{
		APIToken:   apiToken,
		RundeckURL: rundeckURL,
	}
}

// GetJobID : Return rundeck job ID based on the string
func (r *Rundeck) GetJobID(jobFilter string) string {
	data := map[string]string{"jobFilter": jobFilter}
	jData, err := json.Marshal(data)
	CheckErr(err)
	req, err := http.NewRequest(http.MethodPost, r.RundeckURL+"/api/27/project/Test/jobs", bytes.NewBuffer(jData))
	CheckErr(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Rundeck-Auth-Token", r.APIToken)
	resp, err := http.DefaultClient.Do(req)
	CheckErr(err)
	body, _ := ioutil.ReadAll(resp.Body)
	var j jobs
	xml.Unmarshal(body, &j)
	if len(j.Jobs) > 0 {
		return j.Jobs[0].ID
	}
	return ""
}

// RunJobs : Run the job specified by jobID on the host
func (r *Rundeck) RunJobs(jobID, host string) (string, string) {
	content, err := json.Marshal(map[string]string{"loglevel": "verbose", "filter": host})
	CheckErr(err)
	req, err := http.NewRequest(http.MethodGet, r.RundeckURL+"/api/2/job/"+jobID+"/run", bytes.NewBuffer(content))
	CheckErr(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Rundeck-Auth-Token", r.APIToken)
	resp, err := http.DefaultClient.Do(req)
	CheckErr(err)
	body, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	var res result
	err = xml.Unmarshal(body, &res)
	CheckErr(err)
	if len(res.Executions.Execution) > 0 {
		log.Println("PARSED:", res.Success, res.Executions.Execution[0].ExecID)
		return res.Success, res.Executions.Execution[0].ExecID
	}
	return res.Success, ""
}

// CheckErr : Error checking
func CheckErr(err error) {
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
}
