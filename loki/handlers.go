package loki

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
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
		log.Println("❌ Invalid request format:", err)

	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("unable to read request body")
	}
	log.Printf("Request Body: %s", string(body))

	log.Println(reqBody)

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
		log.Println("❌ Missing required settings (Loki URL, Query)")
	}

	// Using WaitGroup to manage goroutine
	var wg sync.WaitGroup
	var logs []string
	var mu sync.Mutex // Mutex for safe concurrent access to logs slice

	wg.Add(1)
	go func() {
		defer wg.Done()

		endTime := time.Now().UTC()
		startTime := endTime.Add(-5 * time.Minute)

		// Fetch logs
		lokiUrl := "http://100.27.210.53:3100"
		queryy := "{job=\"varlogs\"}"
		fetchedLogs, err := utils.FetchLogs(lokiUrl, queryy, startTime, endTime, 10)
		if err != nil {
			log.Printf("❌ Error fetching logs: %v", err)

		}

		// Protect shared slice with a mutex
		mu.Lock()
		logs = append(logs, fetchedLogs...)
		mu.Unlock()
	}()

	// Wait for goroutine to finish
	wg.Wait()

	// Send logs to Telex
	logMessage := strings.Join(logs, "\n")
	data := map[string]interface{}{
		"event_name": "loki_integration_event",
		"message":    logMessage,
		"status":     "success",
		"username":   "Loki integration",
	}

	telex_url := os.Getenv("WEBHOOK_URL")
	telexResponse, err := utils.SendLogsToTelex(telex_url, data)
	if err != nil {
		log.Println("❌ Error sending logs to Telex:", err)

	}

	// Print successful response for debugging
	fmt.Println("✅ Logs sent to Telex:", telexResponse)
	c.JSON(http.StatusOK, data)
}
