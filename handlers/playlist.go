package handlers

import (
	"net/http"
	"spoplaylist/use_cases"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Handler struct {
	NewRelicApp *newrelic.Application
	AdminPlaylist use_cases.IAdminPlayList
}

func (h *Handler) PutHot100(c *gin.Context){		
	txn := h.NewRelicApp.StartTransaction("put_hot_100")
	date := c.Param("date")
	err := h.AdminPlaylist.PutHot100(c, *txn, date)
	if err != nil {		
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
		
	c.JSON(http.StatusOK, songs)
}