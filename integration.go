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
				"app_logo":         "https://grafana.com/media/docs/loki/logo-grafana-loki.png",
				"app_name":         "Grafana-Loki Integration",
				"app_url":          "https://telex-integration.onrender.com",
				"background_color": "#fff",
			},
			"is_active":            true,
			"integration_category": "Monitoring & Logging",
			"integration_type":     "interval",
			"key_features": []string{
				"Fetch logs from Loki at regular intervals",
				"Filter logs based on custom Loki queries",
				"Send logs to a designated Telex channel",
				"Monitor specific applications or services",
			},
			"permissions": map[string]map[string]interface{}{
				"monitoring_user": {
					"always_online": true,
					"display_name":  "Performance Monitor",
				},
			},

			"author": "Tireni",
			"settings": []map[string]interface{}{
				{"label": "Loki Server URL", "type": "text", "required": true, "default": "http://localhost:3100"},
				{"label": "Loki Query", "type": "text", "required": true, "default": "{job='varlogs'}"},
				{"label": "Interval", "type": "text", "required": true, "default": "*/5 * * * *", "description": "Cron expression defining how often logs are fetched"},
			},
			"tick_url":   "https://telex-integration.onrender.com/tick",
			"target_url": "",
		},
	}

	c.JSON(http.StatusOK, integrationJSON)
}
