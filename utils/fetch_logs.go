package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// LokiResponse represents the expected response from Loki
type LokiResponse struct {
	Data struct {
		Result []struct {
			Stream map[string]string `json:"stream"`
			Values [][]string        `json:"values"` // 2D array of [timestamp, log line]
		} `json:"result"`
	} `json:"data"`
}

// FetchLogs queries Loki and returns log entries
func FetchLogs(lokiURL, query string) ([]string, error) {
	// Construct Loki query URL
	lokiQueryURL := fmt.Sprintf("%s/loki/api/v1/query_range?query=%s&limit=5", lokiURL, query)

	// Make request to Loki
	resp, err := http.Get(lokiQueryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse Loki response
	var lokiResponse LokiResponse
	err = json.Unmarshal(body, &lokiResponse)
	if err != nil {
		return nil, err
	}

	// Extract log messages
	var logs []string
	for _, result := range lokiResponse.Data.Result {
		for _, value := range result.Values {
			logs = append(logs, value[1]) // value[1] is the log message
		}
	}

	return logs, nil
}

// SendLogsToTelex sends logs to Telex target_url
func SendLogsToTelex(returnURL string, logs []string, channelID string) {
	// Prepare response payload
	responsePayload := map[string]interface{}{
		"channel_id": channelID,
		"logs":       logs,
	}

	// Convert to JSON
	jsonPayload, _ := json.Marshal(responsePayload)

	// Send response to Telex return_url
	resp, err := http.Post(returnURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Println("Error sending logs to Telex:", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Logs successfully sent to Telex (%s): %v\n", returnURL, logs)
}
