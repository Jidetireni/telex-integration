package main

import (
	"net/http"
	"telex-integration/loki"

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
				"app_logo":         "https://grafana.com/media/docs/loki/logo-grafana-loki.png",
				"app_name":         "Grafana-Loki Integration",
				"app_url":          "https://telex-integration.onrender.com/",
				"background_color": "#fff",
			},
			"is_active":            true,
			"integration_category": "Monitoring & Logging",
			"integration_type":     "interval",
			"output": []map[string]interface{}{
				{"label": "output_channel_1", "value": true},
				{"label": "output_channel_2", "value": false},
			},
			"key_features": []string{
				"Fetch logs from Loki at regular intervals",
				"Filter logs based on custom Loki queries",
				"Send logs to a designated Telex channel",
				"Monitor specific applications or services",
			},
			"permissions": map[string]map[string]interface{}{
				"monitoring_user": {
					"always_online": true,
					"display_name":  "Loki Monitor",
				},
			},
			"settings": []map[string]interface{}{
				{"label": "Loki Server URL", "type": "text", "required": true, "default": "http://34.229.14.54:3100"},
				{"label": "Loki Query", "type": "text", "required": true, "default": "{job=\"varlogs\"}"},
				{"label": "Interval", "type": "text", "required": true, "default": "* * * * *"},
				// {"label": "Key", "type": "text", "required": true, "default": "1234567890"},
				// {"label": "Do you want to continue", "type": "checkbox", "required": true, "default": "Yes"},
				// {"label": "Provide Speed", "type": "number", "required": true, "default": "1000"},
				// {"label": "Sensitivity Level", "type": "dropdown", "required": true, "default": "Low", "options": []string{"High", "Low"}},
				// {"label": "Alert Admin", "type": "multi-checkbox", "required": true, "default": "Super-Admin", "options": []string{"Super-Admin", "Admin", "Manager", "Developer"}},
			},
			"tick_url":   "https://telex-integration.onrender.com/tick/",
			"target_url": loki.LatestReturnURL,
		},
	}

	c.JSON(http.StatusOK, integrationJSON)
}
