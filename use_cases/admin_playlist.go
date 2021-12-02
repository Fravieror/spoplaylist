package use_cases

import (
	"spoplaylist/Infrastructure/api/playlist"
	"spoplaylist/Infrastructure/api/source_music"
	"time"

	"github.com/gin-gonic/gin"
	ccache "github.com/karlseguin/ccache/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const ttlCacheMinutes = 60

type IAdminPlayList interface {
	PutHot100(c *gin.Context, txn *newrelic.Transaction, date string) (string, error)
}

func NewAdminPlaylist(playlist playlist.Iplaylist, sourceMusic source_music.ISourceMusic, cache *ccache.Cache) IAdminPlayList {
	return &AdminPlaylist{
		Playlist:    playlist,
		SourceMusic: sourceMusic,
		Cache:       cache,
	}
}

type AdminPlaylist struct {
	Playlist    playlist.Iplaylist
	SourceMusic source_music.ISourceMusic
	Cache       *ccache.Cache
}

func (admin *AdminPlaylist) PutHot100(c *gin.Context, txn *newrelic.Transaction, date string) (string, error) {
	var songs []string
	var err error

	itemCache := admin.Cache.Get(date)
	if itemCache == nil || itemCache.Expired() {
		songs, err = admin.SourceMusic.GetHot100Songs(date)
		if err != nil {
			return "", err
		}
		admin.Cache.Set(date, songs, ttlCacheMinutes*time.Minute)
	} else {
		songs = itemCache.Value().([]string)
	}

	snapshotID, err := admin.Playlist.SaveSongsOnPlaylist(c, txn, "hits of all time", songs)
	if err != nil {
		return "", err
	}

	return snapshotID, nil
}
