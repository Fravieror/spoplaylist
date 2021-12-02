package main

import (
	"fmt"
	"os"
	"spoplaylist/Infrastructure/api/playlist"
	"spoplaylist/Infrastructure/api/source_music"
	"spoplaylist/handlers"
	"spoplaylist/use_cases"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	spotifyauth "github.com/zmb3/spotify/v2"
)

func main() {	
	// port := os.Getenv("PORT")

	// dependencies := buildDependencies()

	// if err := server.run(port, dependencies); err != nil {
	// 	log.Panic("error running server", err)
	// }
		 
	application, err := newrelic.NewApplication(newrelic.ConfigAppName("spoplaylist"), 
												newrelic.ConfigDebugLogger(os.Stdout))	
	if err != nil {
		fmt.Errorf("error stating new relic monitoring, details: %w", err)
	}

	handler := handlers.Handler{NewRelicApp: application, 
						AdminPlaylist: use_cases.NewAdminPlaylist(playlist.NewSpotify(),
																 source_music.NewBillboard())}

	router := gin.Default()
    router.GET("/hot-100/:date", handler.PutHot100)

    router.Run("localhost:8080")
}

// func buildDependencies() {

// }


func addSongsPlaylist(c *gin.Context, txn *newrelic.Transaction, auth *spotifyauth.Client, songs []string) error {			
	playLists, err := auth.GetPlaylistsForUser(c, UserID)
	if err != nil {
		fmt.Errorf("error getting playlist: %w", err)
		return fmt.Errorf("error consuming API spotify check logs for more details, transaction: %s", txn.GetTraceMetadata().TraceID)
	}
	for _, playlist := range playLists.Playlists{
		if playlist.Name == "hits of all time" {
			return addSongs(c, txn, auth, songs)
		}
	}
	auth.CreatePlaylistForUser(c, UserID, "hits of all time", "popular songs in history", false, false)	
	return addSongs(c, txn, auth, songs)
}

func addSongs(c *gin.Context, txn *newrelic.Transaction, auth *spotifyauth.Client, songs []string) error {
	for _, song := range songs {
		
	}
}