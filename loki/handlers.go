package loki

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"telex-integration/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// LogRequest represents the data needed to fetch logs
type LogRequest struct {
	LokiURL   string
	Query     string
	StartTime string
	EndTime   string
	ReturnURL string
	ChannelID string
}

// RequestBody represents the JSON structure sent by Telex
type RequestBody struct {
	ChannelID string    `json:"channel_id"`
	ReturnURL string    `json:"return_url"`
	Settings  []Setting `json:"settings"`
}

// Setting represents each setting field
type Setting struct {
	Label    string      `json:"label"`
	Type     string      `json:"type"`
	Required bool        `json:"required"`
	Default  interface{} `json:"default"` // <-- Supports both string and number
}

var LatestReturnURL string

// TickHandler handles POST requests from Telex
func TickHandler(c *gin.Context) {
	var reqBody RequestBody

	// Parse incoming JSON request
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "error_msg": err.Error()})

	}

	LatestReturnURL = reqBody.ReturnURL

	// Extract settings
	var lokiURL, query string
	for _, setting := range reqBody.Settings {
		switch setting.Label {
		case "Loki Server URL":
			if url, ok := setting.Default.(string); ok {
				lokiURL = url
			}
		case "Loki Query":
			if q, ok := setting.Default.(string); ok {
				query = q
			}
		}
	}

	// Validate required settings
	if lokiURL == "" || query == "" {
		log.Println("Missing required settings (Loki URL, Query)")
		return
	}

	// Get time range (last 5 minutes)
	endTime := time.Now().UTC()
	startTime := endTime.Add(-5 * time.Minute)

	// Send log request to the LogsEndpointHandler via the channel
	logs, err := utils.FetchLogs(lokiURL, query, startTime, endTime, 10)
	if err != nil {
		log.Printf("Error fetching logs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs", "error_msg": err.Error()})
		return
	}

	_ = os.Getenv("WEBHOOK_URL")
	telex_response, err := utils.SendLogsToTelex(reqBody.ReturnURL, logs, reqBody.ChannelID)
	if err != nil {
		log.Printf("Error sending logs to telex: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending logs to telex:", "error_msg": err.Error()})
		return
	}
	fmt.Println(telex_response)

	c.JSON(http.StatusOK, gin.H{
		"channel_id":    reqBody.ChannelID,
		"return_url":    reqBody.ReturnURL,
		"logs":          logs,
		"telex_reponse": telex_response,
	})

}

// LogsEndpointHandler fetches logs from Loki when triggered
// func LogsEndpointHandler(c *gin.Context) {
// 	// Wait for a log request from the channel
// 	req := <-logChan

// 	// Fetch logs from Loki
// 	logs, err := utils.FetchLogs(req.LokiURL, req.Query, req.StartTime, req.EndTime)
// 	if err != nil {
// 		log.Println("Failed to fetch logs from Loki:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
// 		return
// 	}

// 	// Send logs to Telex return_url
// 	utils.SendLogsToTelex(req.ReturnURL, logs, req.ChannelID)

// 	// Respond to the client
// 	c.JSON(http.StatusOK, gin.H{"message": "Logs fetched and sent successfully"})
// }
