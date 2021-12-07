package config

import (
	"context"
	"log"
	"os"

	"github.com/karlseguin/ccache/v2"
	spotify "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type Dependencies struct {
	ClientSpotify      *spotify.Client
	CacheAdminPlaylist *ccache.Cache
	CachePlaylist	   *ccache.Cache
}

const (
	maxSizeCache = 5000
	UrlToken = "https://accounts.spotify.com/api/token"
)

func BuildDependencies() Dependencies {
	environment := os.Getenv("ENVIRONMENT")

	switch environment {
	case "PRODUCTION":
		return Dependencies{}
	default:
		return Dependencies{
			ClientSpotify:      getClient(),
			CacheAdminPlaylist: ccache.New(ccache.Configure().MaxSize(maxSizeCache)),
			CachePlaylist:  ccache.New(ccache.Configure().MaxSize(maxSizeCache)),
		}
	}
}

func getClient() *spotify.Client {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:   os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		TokenURL:     spotifyauth.TokenURL,		
		Scopes: []string{"playlist-modify-private"},
		AuthStyle: oauth2.AuthStyleInParams,			
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	return spotify.New(httpClient)

}