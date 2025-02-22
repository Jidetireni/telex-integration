package loki

import (
	"log"
	"net/http"
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

// TickHandler handles POST requests from Telex
func TickHandler(c *gin.Context) {
	var reqBody RequestBody

	// Parse incoming JSON request
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "error_msg": err.Error()})

	}

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required settings"})
		return
	}

	// **Send HTTP 202 Accepted BEFORE starting Goroutine**
	c.JSON(http.StatusAccepted, gin.H{"message": "Processing in background", "channel_id": reqBody.ChannelID})

	// Process logs in the background
	go func() {
		// Get time range (last 5 minutes)
		endTime := time.Now().UTC()
		startTime := endTime.Add(-5 * time.Minute)

		// Fetch logs from Loki
		logs, err := utils.FetchLogs(lokiURL, query, startTime, endTime, 10)
		if err != nil {
			log.Printf("Error fetching logs: %v", err)
			return
		}

		// Send logs to Telex
		telexResponse, err := utils.SendLogsToTelex(reqBody.ReturnURL, logs, reqBody.ChannelID)
		if err != nil {
			log.Printf("Error sending logs to Telex: %v", err)
			return
		}

		// Log success (No `c.JSON()` here since response is already sent)
		log.Printf("Successfully processed logs for Channel ID %s. Telex Response: %v", reqBody.ChannelID, telexResponse)
	}()
}
