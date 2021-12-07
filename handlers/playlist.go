package handlers

import (
	"fmt"
	"net/http"
	"os"
	"spoplaylist/entities"
	"spoplaylist/use_cases"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
)


var ( 
    words = make(map[string]int)
)

type Handler struct {
	NewRelicApp   *newrelic.Application
	AdminPlaylist use_cases.IAdminPlayList
}

func (h *Handler) PutHot100(c *gin.Context) {
	txn := h.NewRelicApp.StartTransaction("put_hot_100")
	date := c.Param("date")
	snapshotID, err := h.AdminPlaylist.PutHot100(c, txn, date)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("songs added successfully to playlist #snapshot: %s", snapshotID))
}

func (h *Handler) GetHot100(c *gin.Context) {
	txn := h.NewRelicApp.StartTransaction("get_hot_100")
	date := c.Param("date")
	songs, err := h.AdminPlaylist.GetHot100(c, txn, date)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, songs)
}

func (h *Handler) SetCredential(c *gin.Context) {
	var credential entities.Credential
	if err := c.BindJSON(&credential); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ResponseError{
			Error: "client_id or client_credential empty",
			ErrorDescription: "client_id and client_credential required",
		})
		return
	}
	err := os.Setenv("CLIENT_ID", credential.ClientID)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("error setting client_id, error: %w", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ResponseError{
			Error: "error setting credential",
			ErrorDescription: "error setting client_id, try again",
		})
	}
	err = os.Setenv("CLIENT_SECRET", credential.ClientCredential)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("error setting client_secret, error: %w", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ResponseError{
			Error: "error setting credential",
			ErrorDescription: "error setting client_secret, try again",
		})
	}
	err = os.Setenv("USER_ID", credential.UserID)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("error setting user_id, error: %w", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ResponseError{
			Error: "error setting credential",
			ErrorDescription: "error setting user_id, try again",
		})
	}
	err = os.Setenv("PLAYLIST", credential.Playlist)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("error setting playlist_name, error: %w", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ResponseError{
			Error: "error setting credential",
			ErrorDescription: "error setting playlist_name, try again",
		})
	}
	c.JSON(http.StatusOK, "Credential updated")
}
