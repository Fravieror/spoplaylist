package config

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/karlseguin/ccache/v2"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2"
)

type Dependencies struct {
	ClientSpotify      *spotifyauth.Client
	CacheAdminPlaylist *ccache.Cache
}

const (
	maxSizeCache = 5000
	UrlToken = "https://accounts.spotify.com/api/token"
)

func BuildDependencies() Dependencies {
	environment := os.Getenv("ENVIRONMENT")
	
	client := spotifyauth.New(http.DefaultClient)
	client.CheckRedirect()

	switch environment {
	case "PRODUCTION":
		return Dependencies{}
	default:
		return Dependencies{
			ClientSpotify:      client,
			CacheAdminPlaylist: ccache.New(ccache.Configure().MaxSize(maxSizeCache)),
		}
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// use the same state string here that you used to generate the URL
	token, err := auth.Token(r.Context(), state, r)
	if err != nil {
		  http.Error(w, "Couldn't get token", http.StatusNotFound)
		  return
	}
	// create a client using the specified token
	client := spotify.New(auth.Client(r.Context(), token))

	// the client can now be used to make authenticated requests
}


func getToken() (string, error) {
	cli := http.DefaultClient
	ClientID := os.Getenv("CLIENT_ID")
	ClientSecret := os.Getenv("CLIENT_SECRET")
	payload := strings.NewReader("grant_type=client_credentials")
	req, err := http.NewRequest("POST", UrlToken, payload)
	if err != nil {

	}
	sEnc := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", ClientID, ClientSecret)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", sEnc))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := cli.Do(req)
	if err != nil {

	}
	return res.Status
}