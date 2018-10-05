package main

import (
	"encoding/json"
	"log"
	"middleware-zabbix-rundeck/middleware"
	"net/http"
)

const (
	apiToken       = "3T7WByHHv4EaphbFib2L5kYHG617f098"
	rundeckURL     = "http://172.23.215.206:4440"
	serverPublicIP = "172.16.128.39"
	redisAPI       = "http://172.23.215.206:8080"
)

type zabbixRequest struct {
	EventID     string `json:"eventID"`
	ZabbixHost  string `json:"zabbix_host"`
	Host        string `json:"host"`
	TriggerDesc string `json:"triggerDescription"`
	TriggerID   string `json:"triggerId"`
	ItemKey     string `json:"itemKey"`
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)

	}
}

func main() {
	http.HandleFunc("/submit", rundeckJob)
	http.Handle("favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(serverPublicIP+":8090", nil)
}

func rundeckJob(w http.ResponseWriter, r *http.Request) {
	if cType := r.Header.Get("Content-Type"); cType != "application/json" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var z zabbixRequest
	err := json.NewDecoder(r.Body).Decode(&z)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Recieved:", z.EventID, z.Host, z.ItemKey, z.TriggerDesc, z.TriggerID, z.ZabbixHost)

	rundeck := middleware.NewRundeck(apiToken, rundeckURL)
	if jobID := rundeck.GetJobID(z.TriggerDesc); jobID != "" {
		_, execID := rundeck.RunJobs(jobID, z.Host)
		if len(execID) > 1 {
			if middleware.AckEvent(z.EventID, z.Host, z.TriggerDesc, z.ItemKey, jobID, redisAPI) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Jobs executed, Acknowldge completed."))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to acknowledge"))
			return
		}
	}

}
