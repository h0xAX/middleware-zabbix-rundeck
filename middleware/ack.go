package middleware

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type redisResp struct {
	Status string `json:"status"`
	Text   string `json:"text"`
}

//AckEvent : send ack data to redis
func AckEvent(jobID, host, trigger, itemKey, executionID, redisAPI string) bool {
	ackData := map[string]string{
		jobID:       jobID,
		host:        host,
		trigger:     trigger,
		itemKey:     itemKey,
		executionID: executionID,
	}
	data, err := json.Marshal(ackData)
	checkErr(err)
	req, err := http.NewRequest(http.MethodPost, redisAPI, bytes.NewBuffer(data))
	checkErr(err)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Do(req)
	checkErr(err)
	var respData redisResp
	err = json.NewDecoder(resp.Body).Decode(&respData)
	checkErr(err)
	if respData.Status == "Accept" {
		return true
	}
	return false
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
