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
				"app_description":  "Fetches logs from a Loki server and sends them to a Telex channel at defined intervals.",
				"app_logo":         "https://miro.medium.com/v2/resize:fit:1400/1*k-hdOAQjRXKoyguzKuoeKg.png",
				"app_name":         "Grafana-Loki Integration",
				"app_url":          "https://telex-integration-production.up.railway.app/",
				"background_color": "#1F2937",
			},
			"integration_category": "Monitoring & Logging",
			"integration_type":     "interval",
			"is_active":            true,
			"output": []map[string]interface{}{
				{"label": "log_channel", "value": true},
			},
			"key_features": []string{
				"Fetch logs from Loki at regular intervals",
				"Filter logs based on custom Loki queries",
				"Send logs to a designated Telex channel",
				"Monitor specific applications or services",
			},
			"permissions": map[string]interface{}{
				"monitoring_user": map[string]interface{}{
					"always_online": true,
					"display_name":  "Loki Monitor Bot",
				},
			},
			"settings": []map[string]interface{}{
				{"label": "Loki Server URL", "type": "text", "required": true, "default": "http://localhost:3100"},
				{"label": "Loki Query", "type": "text", "required": true, "default": "{job='nginx'}"},
				{"label": "Interval", "type": "text", "required": true, "default": "*/5 * * * *", "description": "Cron expression defining how often logs are fetched"},
				{"label": "Maximum Logs", "type": "number", "required": false, "default": 10, "description": "Maximum number of logs to fetch per interval"},
				{"label": "Alert Admins", "type": "multi-checkbox", "required": false, "default": []string{"Super-Admin"}, "options": []string{"Super-Admin", "Admin", "Manager", "Developer"}},
			},
			"tick_url": "https://telex-integration-production.up.railway.app/fetch-logs",
		},
	}

	c.JSON(http.StatusOK, integrationJSON)
}
