package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
