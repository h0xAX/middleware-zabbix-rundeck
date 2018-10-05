package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	apiToken   = "3T7WByHHv4EaphbFib2L5kYHG617f098"
	rundeckURL = "http://172.23.215.206:4440"
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

func getJobID(jobFilter string) string {
	data := map[string]string{"jobFilter": jobFilter}
	jData, err := json.Marshal(data)
	checkErr(err)
	req, err := http.NewRequest(http.MethodPost, rundeckURL+"/api/27/project/Test/jobs", bytes.NewBuffer(jData))
	checkErr(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Rundeck-Auth-Token", apiToken)
	resp, err := http.DefaultClient.Do(req)
	checkErr(err)
	body, _ := ioutil.ReadAll(resp.Body)
	var j jobs
	xml.Unmarshal(body, &j)
	if len(j.Jobs) > 0 {
		return j.Jobs[0].ID
	}
	return ""
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
	}
}

func runJobs(jobID, host string) {
	content, err := json.Marshal(map[string]string{"loglevel": "verbose", "filter": host})
	checkErr(err)
	req, err := http.NewRequest(http.MethodGet, rundeckURL+"/api/2/job/"+jobID+"/run", bytes.NewBuffer(content))
	checkErr(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Rundeck-Auth-Token", apiToken)
	resp, err := http.DefaultClient.Do(req)
	checkErr(err)
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	fmt.Println(string(body))
}

func main() {
	jID := getJobID("/tmp")
	fmt.Println("JOBID:", jID)
	if len(jID) > 2 {
		runJobs(jID, "")
	}
}
