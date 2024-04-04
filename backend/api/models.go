package api

import "time"

type Response struct {
	Command string                 `json:"command"`
	Status  string                 `json:"status"`
	Data    map[string]interface{} `json:"data"`
}

type Convo struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	LastModified time.Time `json:"lastModified"`
	FirstCreated time.Time `json:"firstCreated"`
}

type ConvoHistory struct {
	ID           string    `json:"id"`
	IsUser       bool      `json:"isUser"`
	Message      string    `json:"message"`
	LastModified time.Time `json:"lastModified"`
	FirstCreated time.Time `json:"firstCreated"`
}
