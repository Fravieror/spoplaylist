package entities

type Credential struct {
	ClientID         string `json:"client_id" binding:"required"`
	ClientCredential string `json:"client_credential" binding:"required"`
	UserID           string `json:"client_credential"`
	Playlist         string `json:"playlist_name"`
}