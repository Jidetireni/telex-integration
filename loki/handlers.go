package loki

import (
	"log"
	"net/http"
	"telex-integration/utils"

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
	Label    string `json:"label"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Default  string `json:"default"`
}

// TickHandler handles POST requests from Telex
func TickHandler(c *gin.Context) {
	var reqBody RequestBody

	// Parse incoming JSON request
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Immediately respond with 202 Accepted
	c.JSON(http.StatusAccepted, gin.H{"message": "Request received, processing in background"})

	go func() {
		// Extract settings
		var lokiURL, query string
		for _, setting := range reqBody.Settings {
			if setting.Label == "Loki Server URL" {
				lokiURL = setting.Default
			} else if setting.Label == "Loki Query" {
				query = setting.Default
			}
		}

		// Validate settings
		if lokiURL == "" || query == "" {
			log.Println("Loki URL or query missing in settings")
			return
		}

		// Fetch logs from Loki
		logs, err := utils.FetchLogs(lokiURL, query)
		if err != nil {
			log.Println("Failed to fetch logs from Loki:", err)
			return
		}

		// Send logs to Telex return_url
		utils.SendLogsToTelex(reqBody.ReturnURL, logs, reqBody.ChannelID)
	}()

	// Respond to Telex
	c.JSON(http.StatusOK, gin.H{"message": "Logs sent successfully"})
}
