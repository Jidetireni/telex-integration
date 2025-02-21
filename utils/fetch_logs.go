package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// LokiResponse represents the structure of the response from Loki
type LokiResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Stream map[string]string `json:"stream"` // Log labels (e.g., job, app)
			Values [][]string        `json:"values"` // 2D array: [timestamp, log]
		} `json:"result"`
	} `json:"data"`
}

// FetchLogs queries Loki and returns log entries
func FetchLogs(lokiURL, query string, start, end time.Time, limit int) ([]string, error) {
	// Build query parameters using url.Values
	params := url.Values{}
	params.Set("query", query) // LogQL query
	params.Set("start", fmt.Sprintf("%d", start.UnixNano()))
	params.Set("end", fmt.Sprintf("%d", end.UnixNano()))
	params.Set("limit", fmt.Sprintf("%d", limit))

	// Construct Loki query URL with time range
	reqURL := fmt.Sprintf("%s/loki/api/v1/query_range?%s", lokiURL, params.Encode())

	// Make request to Loki
	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("loki query failed with status: %s", resp.Status)
	}

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var lokiResponse LokiResponse
	if err := json.Unmarshal(body, &lokiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse Loki response: %v", err)
	}

	// Extract logs
	var logs []string
	for _, result := range lokiResponse.Data.Result {
		for _, value := range result.Values {
			timestamp := value[0]
			logLine := value[1]
			logs = append(logs, fmt.Sprintf("[%s] %s", timestamp, logLine))
		}
	}

	return logs, nil
}

func SendLogsToTelex(returnURL string, logs []string, channelID string) (string, error) {
	// Convert payload to JSON
	data := map[string]interface{}{
		"event_name": "Loki integration",
		"message":    logs[0],
		"status":     "success",
		"username":   "tireni",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// Send POST request to Telex's return_url
	req, err := http.NewRequest("POST", returnURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to build logs request to Telex: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", fmt.Errorf("failed to send logs to Telex: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("telex returned non-OK status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	log.Printf("Logs successfully sent to Telex (%s): %v\n", returnURL, logs)
	return string(body), err

}
