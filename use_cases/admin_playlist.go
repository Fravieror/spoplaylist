package use_cases

import (
	"spoplaylist/Infrastructure/api/playlist"
	"spoplaylist/Infrastructure/api/source_music"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type IAdminPlayList interface {
	PutHot100(c *gin.Context, txn newrelic.Transaction, date string) error
}

func NewAdminPlaylist(playlist playlist.Iplaylist, sourceMusic source_music.ISourceMusic) IAdminPlayList {
	return &AdminPlaylist{
		Playlist: playlist,
		SourceMusic: sourceMusic,
	}
}

type AdminPlaylist struct {
	Playlist playlist.Iplaylist
	SourceMusic source_music.ISourceMusic
}

func (admin *AdminPlaylist) PutHot100(c *gin.Context, txn newrelic.Transaction, date string) error {		
	songs, err := admin.SourceMusic.GetHot100Songs(date)
	if err != nil {		
		return err
	}

	err = admin.Playlist.SaveSongs(songs)
	if err != nil {
		return err
	}

	return nil
}