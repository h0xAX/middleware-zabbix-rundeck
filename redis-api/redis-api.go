package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mediocregopher/radix.v2/redis"
)

type ackData struct {
	JobID       string `json:"jobID"`
	Host        string `json:"host"`
	Trigger     string `json:"trigger"`
	ItemKey     string `json:"itemKey"`
	ExecutionID string `json:"executionID"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", addToRedis)
	http.ListenAndServe(":8070", mux)
}

func addToRedis(w http.ResponseWriter, r *http.Request) {
	result := make(map[string]string)
	if r.Method == http.MethodPost {
		log.Println("Received POST request from", r.RemoteAddr)
		client, err := redis.Dial("tcp", "172.23.215.195:6379")
		checkErr(err)
		if r.Header.Get("Content-Type") == "application/json" {
			var ack ackData
			err = json.NewDecoder(r.Body).Decode(&ack)
			checkErr(err)
			log.Println("Host:", ack.Host, "Trigger:", ack.Trigger, "JobID:", ack.JobID, "ExecutionID:", ack.ExecutionID)
			if err = client.Cmd("HMSET", ack.ExecutionID, "time", time.Now().UnixNano(), "host", ack.Host, "trigger", ack.Trigger, "items", ack.ItemKey, "jobID", ack.JobID, "src", r.RemoteAddr, "status", "", "ack", "").Err; err != nil {
				log.Println(err)
				result["Status"] = "Reject"
			} else {
				log.Println("Successfully added to redis")
				result["Status"] = "Accept"
			}
		} else {
			http.Error(w, "Request is not in JSON format", http.StatusNotAcceptable)
			result["Status"] = "Reject"
		}
	} else {
		result["Status"] = "Running" // Show this for other methods.
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(result)
		checkErr(err)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
