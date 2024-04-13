package models

type Response struct {
	Command string                 `json:"command"`
	Data    map[string]interface{} `json:"data"`
	Status  string                 `json:"status"`
}

type Project struct {
	ProjID       string `json:"projID"`
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	Workspace    string `json:"workspace"`
	LastModified int    `json:"lastModified"`
	FirstCreated int    `json:"firstCreated"`
	Chats        []Chat `json:"chats"`
}

type Chat struct {
	ChatID       string `json:"chatID"`
	ProjID       string `json:"projID"`
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	LastModified int    `json:"lastModified"`
	FirstCreated int    `json:"firstCreated"`
}

type Message struct {
	// rename to MsgID
	ID           string `json:"id"`
	ChatID       string `json:"chatID"`
	ProjID       string `json:"projID"`
	IsUser       bool   `json:"isUser"`
	IsDeleted    bool   `json:"isDeleted"`
	Message      string `json:"message"`
	LastModified int    `json:"lastModified"`
	FirstCreated int    `json:"firstCreated"`
}

const (
	ERR_Internal      = "INTERNAL"
	ERR_MissingInput  = "MISSING_INPUT"
	ERR_InvalidInput  = "INVALID_INPUT"
	ERR_InvalidProjID = "INVALID_PROJ_ID"
	ERR_InvalidChatID = "INVALID_CHAT_ID"
	ERR_Decode        = "DECODING"
)
