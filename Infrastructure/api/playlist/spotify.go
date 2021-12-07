package playlist

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/karlseguin/ccache/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
	spotifyauth "github.com/zmb3/spotify/v2"
)

const ttlCacheMinutes = 60
const UserID = ""

type Spotify struct {
	Client *spotifyauth.Client
	HttpClient *http.Client
	Cache *ccache.Cache
}

func NewSpotify(client *spotifyauth.Client, httpClient *http.Client, cache *ccache.Cache) Iplaylist {
	return &Spotify{
		Client: client,
		HttpClient: httpClient,
		Cache: cache,
	}
}

func (s *Spotify) SaveSongsOnPlaylist(c *gin.Context, txn *newrelic.Transaction, 
				playlistName string, songs []string) (string, error) {
	playlist, err := s.getPlayList(c, txn, songs, playlistName)
	if err != nil {
		return "", err
	}

	tracks, err := s.getSongs(c, txn, songs)
	if err != nil {
		return "", err
	}
	snapshotID, err := s.Client.AddTracksToPlaylist(c, spotifyauth.ID(playlist.String()), tracks...)
	if err != nil {
		fmt.Println(fmt.Errorf("error adding tracks to playlist: %w", err))
		return "", fmt.Errorf("error consuming API spotify check logs for more details, transaction: %s", txn.GetTraceMetadata().TraceID)
	}

	fmt.Printf("tracks added to playlist successfully snapshot_id: %s", snapshotID)

	return snapshotID, nil
}

func (s *Spotify) getPlayList(c *gin.Context, txn *newrelic.Transaction, songs []string, playlistName string) (*spotifyauth.ID, error) {
	UserID := os.Getenv("USER_ID")
	playLists, err := s.Client.GetPlaylistsForUser(c, UserID)
	if err != nil {
		fmt.Println(fmt.Errorf("error getting playlist: %w", err))
		return nil, fmt.Errorf("error consuming API spotify check logs for more details, transaction: %s", txn.GetTraceMetadata().TraceID)
	}
	for _, playlist := range playLists.Playlists {
		if playlist.Name == playlistName {
			return &playlist.ID, nil
		}
	}
	fullPlayLists, err := s.Client.CreatePlaylistForUser(c, UserID, playlistName, "popular songs in history", false, false)
	if err != nil {
		fmt.Println(fmt.Errorf("error creating playlist: %w", err))
		return nil, fmt.Errorf("error creating playlist on API spotify check logs for more details, transaction: %s", txn.GetTraceMetadata().TraceID)
	}
	return &fullPlayLists.SimplePlaylist.ID, nil
}

// getSongs get song concurrently to improve time response using goroutines and avoiding race conditions
func (s *Spotify) getSongs(c *gin.Context, txn *newrelic.Transaction, songs []string) ([]spotifyauth.ID, error) {
	tracks := make([]spotifyauth.ID, 0)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, song := range songs {
		wg.Add(1)
		go func(ctx *gin.Context, songParameter string) {
			mu.Lock()
			defer wg.Done()
			defer mu.Unlock()
			trackID := s.getSongFromCache(songParameter)
			if trackID == "" {
				result, err := s.Client.Search(ctx, songParameter, spotifyauth.SearchTypeTrack, spotifyauth.RequestOption(spotifyauth.Limit(1)))
				if err != nil {
					fmt.Println(fmt.Errorf("error searching for song: %s,  error detail:%w", songParameter, err))
				}
				for _, track := range result.Tracks.Tracks {
					tracks = append(tracks, track.SimpleTrack.ID)
					s.Cache.Set(songParameter, track.SimpleTrack.ID, ttlCacheMinutes * time.Minute)
				}					
			} else{
				tracks = append(tracks, trackID)				
			}				
		}(c, song)
	}

	wg.Wait()

	return tracks, nil
}

func (s *Spotify) getSongFromCache(song string) spotifyauth.ID {
	itemCache := s.Cache.Get(song)
	if itemCache == nil || itemCache.Expired() {
		return ""
	} else {
		return itemCache.Value().(spotifyauth.ID)
	}
}