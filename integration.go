package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getIntegrationJSON(c *gin.Context) {
	integrationJSON := map[string]interface{}{
		"data": map[string]interface{}{
			"date": map[string]string{
				"created_at": "2025-02-18",
				"updated_at": "2025-02-18",
			},
			"descriptions": map[string]string{
				"app_name":         "Grafana-Loki Integration",
				"app_description":  "Fetches logs from a Loki server and sends them to a Telex channel at defined intervals.",
				"app_logo":         "https://miro.medium.com/v2/resize:fit:1400/1*k-hdOAQjRXKoyguzKuoeKg.png",
				"app_url":          "https://telex-integration-production.up.railway.app",
				"background_color": "#fff",
			},
			"is_active":        true,
			"integration_type": "interval",
			"author":           "Tireni",
			"key_features": []string{
				"Fetch logs from Loki at regular intervals",
				"Filter logs based on custom Loki queries",
				"Send logs to a designated Telex channel",
				"Monitor specific applications or services",
			},
			"settings": []map[string]interface{}{
				{"label": "Loki Server URL", "type": "text", "required": true, "default": "http://localhost:3100"},
				{"label": "Loki Query", "type": "text", "required": true, "default": "{job='varlogs'}"},
				{"label": "Interval", "type": "text", "required": true, "default": "*/5 * * * *", "description": "Cron expression defining how often logs are fetched"},
			},
			"target_url": "https://telex-integration-production.up.railway.app/fetch-logs",
			"tick_url":   "",
		},
	}

	c.JSON(http.StatusOK, integrationJSON)
}
