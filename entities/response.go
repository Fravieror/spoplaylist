package entities

type ResponseError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"errorDescription"`
}