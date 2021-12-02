package main

import (
	"fmt"
	"os"
	"spoplaylist/Infrastructure/api/playlist"
	"spoplaylist/Infrastructure/api/source_music"
	"spoplaylist/config"
	"spoplaylist/handlers"
	"spoplaylist/use_cases"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
	// Monitoring
	application, err := newrelic.NewApplication(newrelic.ConfigAppName("spoplaylist"),
		newrelic.ConfigDebugLogger(os.Stdout))
	if err != nil {
		fmt.Errorf("error stating new relic monitoring, details: %w", err)
	}

	// Build dependencies
	dependencies := config.BuildDependencies()
	handler := handlers.Handler{NewRelicApp: application,
		AdminPlaylist: use_cases.NewAdminPlaylist(playlist.NewSpotify(dependencies.ClientSpotify),
			source_music.NewBillboard(),
			dependencies.CacheAdminPlaylist)}

	// Create routes
	router := gin.Default()
	router.GET("/hot-100/:date", handler.PutHot100)
	router.PUT("/hot-100/:date", handler.PutHot100)

	// Start server
	port := os.Getenv("PORT")
	router.Run(fmt.Sprintf("localhost:%s", port))
}
