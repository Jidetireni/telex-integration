package main

import (
	"telex-integration/loki"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/integration", getIntegrationJSON)
	r.POST("/fetch-logs", loki.TickHandler)
	r.GET("/logs-endpoint", loki.LogsEndpointHandler)

	r.Run(":8080") // Runs on http://localhost:8080
}
