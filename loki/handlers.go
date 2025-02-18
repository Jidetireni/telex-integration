package loki

import (
	"log"
	"net/http"
	"telex-integration/utils"
	"time"

	"github.com/gin-gonic/gin"
)

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
		return
	}

	// Immediately respond with 202 Accepted
	c.JSON(http.StatusAccepted, gin.H{"message": "Request received, processing in background"})

	go func() {
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

		// Format time in RFC3339 (ISO 8601)
		start := startTime.Format(time.RFC3339)
		end := endTime.Format(time.RFC3339)

		// Fetch logs from Loki
		logs, err := utils.FetchLogs(lokiURL, query, start, end)
		if err != nil {
			log.Println("Failed to fetch logs from Loki:", err)
			return
		}

		// Send logs to Telex return_url
		utils.SendLogsToTelex(reqBody.ReturnURL, logs, reqBody.ChannelID)
	}()
}
