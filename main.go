package main

import (
	"telex-integration/config"
	"telex-integration/loki"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))
	r.GET("/integration.json", config.GetIntegrationJSON)
	r.POST("/tick", loki.TickHandler)
	r.GET("/tick", loki.TickHandler)

	r.Run(":8080") // Runs on http://localhost:8080
}
