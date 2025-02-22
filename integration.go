package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Date struct for created_at and updated_at fields
type Date struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Description struct for app details
type Description struct {
	AppDescription  string `json:"app_description"`
	AppLogo         string `json:"app_logo"`
	AppName         string `json:"app_name"`
	AppURL          string `json:"app_url"`
	BackgroundColor string `json:"background_color"`
}

// Permissions struct for monitoring user
type Permissions struct {
	MonitoringUser struct {
		AlwaysOnline bool   `json:"always_online"`
		DisplayName  string `json:"display_name"`
	} `json:"monitoring_user"`
}

// Setting struct for integration settings
type Setting struct {
	Label    string      `json:"label"`
	Type     string      `json:"type"`
	Required bool        `json:"required"`
	Default  interface{} `json:"default"` // Supports string, number, boolean
}

// Data struct for the full integration details
type Data struct {
	Date                Date        `json:"date"`
	Descriptions        Description `json:"descriptions"`
	IntegrationCategory string      `json:"integration_category"`
	IntegrationType     string      `json:"integration_type"`
	IsActive            bool        `json:"is_active"`
	KeyFeatures         []string    `json:"key_features"`
	Author              string      `json:"author"`
	Settings            []Setting   `json:"settings"`
	TickURL             string      `json:"tick_url"`
	TargetURL           string      `json:"target_url"`
}

// Response struct to wrap the JSON response
type Response struct {
	Data Data `json:"data"`
}

// getIntegrationJSON handles the request and sends the structured JSON response
func getIntegrationJSON(c *gin.Context) {
	response := Response{
		Data: Data{
			Date: Date{
				CreatedAt: time.Now().Format("2006-01-02"),
				UpdatedAt: time.Now().Format("2006-01-02"),
			},
			Descriptions: Description{
				AppDescription:  "Fetches logs from a Loki server and sends them to a Telex channel at defined intervals.",
				AppLogo:         "https://grafana.com/media/docs/loki/logo-grafana-loki.png",
				AppName:         "Grafana-Loki Integration",
				AppURL:          "https://telex-integration.onrender.com/",
				BackgroundColor: "#fff",
			},
			IntegrationCategory: "Monitoring & Logging",
			IntegrationType:     "interval",
			IsActive:            true,

			KeyFeatures: []string{
				"Fetch logs from Loki at regular intervals",
				"Filter logs based on custom Loki queries",
				"Send logs to a designated Telex channel",
				"Monitor specific applications or services",
			},

			Author: "Tireni",
			Settings: []Setting{
				{"loki Server URL", "text", true, "http://100.27.210.53:3100"},
				{"loki Query", "text", true, "{job=\"varlogs\"}"},
				{"interval", "text", true, "2 * * * *"},
			},
			TickURL:   "https://telex-integration.onrender.com/tick",
			TargetURL: "",
		},
	}

	c.JSON(http.StatusOK, response)
}
