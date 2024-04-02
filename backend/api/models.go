package api

import "time"

type Response struct {
	Command string                 `json:"command"`
	Status  string                 `json:"status"`
	Data    map[string]interface{} `json:"data"`
}

type ConvoListData struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	LastModified time.Time `json:"lastModified"`
}
