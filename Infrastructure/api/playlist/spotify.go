package playlist

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	spotifyauth "github.com/zmb3/spotify/v2"
)

const (
	UserID       = ""
	ClientID     = "83efa!!b411227447*e8a*233*a**1d28882b97"
	ClientSecret = "bc3171a2b!!!95e4b889c52e2*c*efb*b**0ff62"
)

type Spotify struct {
	Client *spotifyauth.Client
}

func NewSpotify(client *spotifyauth.Client) Iplaylist {
	return &Spotify{
		Client: client,
	}
}

func (s *Spotify) SaveSongsOnPlaylist(c *gin.Context, txn *newrelic.Transaction, playlistName string, songs []string) (string, error) {
	playlist, err := s.getPlayList(c, txn, songs, playlistName)
	if err != nil {
		return "", err
	}

	tracks, err := s.getSongs(c, txn, songs)
	if err != nil {
		return "", err
	}
	snapshotID, err := s.Client.AddTracksToPlaylist(c, playlist.ID, tracks...)
	if err != nil {
		fmt.Errorf("error adding tracks to playlist: %w", err)
		return "", fmt.Errorf("error consuming API spotify check logs for more details, transaction: %s", txn.GetTraceMetadata().TraceID)
	}

	fmt.Printf("tracks added to playlist successfully snapshot_id: %s", snapshotID)

	return snapshotID, nil
}

func (s *Spotify) getPlayList(c *gin.Context, txn *newrelic.Transaction, songs []string, playlistName string) (*spotifyauth.SimplePlaylist, error) {
	playLists, err := s.Client.GetPlaylistsForUser(c, UserID)
	if err != nil {
		fmt.Errorf("error getting playlist: %w", err)
		return nil, fmt.Errorf("error consuming API spotify check logs for more details, transaction: %s", txn.GetTraceMetadata().TraceID)
	}
	for _, playlist := range playLists.Playlists {
		if playlist.Name == playlistName {
			return nil, nil
		}
	}
	fullPlayLists, err := s.Client.CreatePlaylistForUser(c, UserID, playlistName, "popular songs in history", false, false)
	if err != nil {
		fmt.Errorf("error creating playlist: %w", err)
		return nil, fmt.Errorf("error creating playlist on API spotify check logs for more details, transaction: %s", txn.GetTraceMetadata().TraceID)
	}
	return &fullPlayLists.SimplePlaylist, nil
}

// getSongs get song concurrently to improve time response using goroutines and avoiding race conditions
func (s *Spotify) getSongs(c *gin.Context, txn *newrelic.Transaction, songs []string) ([]spotifyauth.ID, error) {
	tracks := make([]spotifyauth.ID, 0)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, song := range songs {
		go func(ctx *gin.Context, songParameter string) {
			mu.Lock()
			defer wg.Done()
			defer mu.Unlock()
			result, err := s.Client.Search(ctx, songParameter, spotifyauth.SearchTypeTrack, spotifyauth.RequestOption(spotifyauth.Limit(1)))
			if err != nil {
				fmt.Errorf("error searching for song: %s, error detail:%w", songParameter, err)
			}
			for _, track := range result.Tracks.Tracks {
				tracks = append(tracks, track.ID)
			}
			wg.Wait()
		}(c, song)
	}

	return tracks, nil
}
