package playlist

import (
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Iplaylist interface {
	SaveSongs(c *gin.Context, txn newrelic.Transaction, songs []string) error
}
