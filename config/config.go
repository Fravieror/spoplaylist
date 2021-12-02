package config

import (
	"os"

	"github.com/karlseguin/ccache/v2"
	spotifyauth "github.com/zmb3/spotify/v2"
)

type Dependencies struct {
	ClientSpotify *spotifyauth.Client
	CacheAdminPlaylist *ccache.Cache
}

const (
	maxSizeCache = 5000 
)

func BuildDependencies() Dependencies {
	environment := os.Getenv("ENVIRONMENT")
	switch environment{			
		case "PRODUCTION":			
			return Dependencies{}
		default:
			return Dependencies{
				ClientSpotify: &spotifyauth.Client{},
				CacheAdminPlaylist: ccache.New(ccache.Configure().MaxSize(maxSizeCache)),
			}
	}	
}
