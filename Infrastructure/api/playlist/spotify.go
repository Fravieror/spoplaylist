package playlist

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	spotifyauth "github.com/zmb3/spotify/v2"
)

const (
	UserID       = ""
	ClientID     = "83efab411227447e8a233a1d28882b97"
	ClientSecret = "bc3171a2b95e4b889c52e2cefbb0ff62"
)

type Spotify struct {
	Client *spotifyauth.Client
}

func NewSpotify(client *spotifyauth.Client) Iplaylist {
	return &Spotify{
		Client: client,
	}
}

func (s *Spotify) SaveSongs(c *gin.Context, txn newrelic.Transaction, songs []string) error {
	return nil
}

func (s *Spotify) validatePlayList(c *gin.Context, txn newrelic.Transaction, songs []string) (bool, error) {
	playLists, err := s.Client.GetPlaylistsForUser(c, UserID)
	if err != nil {
		fmt.Errorf("error getting playlist: %w", err)
		return false, fmt.Errorf("error consuming API spotify check logs for more details, transaction: %s", txn.GetTraceMetadata().TraceID)
	}
	for _, playlist := range playLists.Playlists {
		if playlist.Name == "hits of all time" {
			return true, nil
		}
	}
	_, err = s.Client.CreatePlaylistForUser(c, UserID, "hits of all time", "popular songs in history", false, false)
	if err != nil {
		fmt.Errorf("error creating playlist: %w", err)
		return false, fmt.Errorf("error creating playlist on API spotify check logs for more details, transaction: %s", txn.GetTraceMetadata().TraceID)
	}
	return true, nil
}