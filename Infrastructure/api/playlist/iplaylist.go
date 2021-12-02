package playlist

import (
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Iplaylist interface {
	SaveSongsOnPlaylist(c *gin.Context, txn *newrelic.Transaction, playlistName string, songs []string) (string, error)
}
