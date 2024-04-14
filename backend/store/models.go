package store

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
	MsgID        string `json:"msgID"`
	ChatID       string `json:"chatID"`
	ProjID       string `json:"projID"`
	IsUser       bool   `json:"isUser"`
	Message      string `json:"message"`
	LastModified int    `json:"lastModified"`
	FirstCreated int    `json:"firstCreated"`
}

const (
	ERR_Internal     = "INTERNAL"
	ERR_MissingInput = "MISSING_INPUT"
	ERR_InvalidInput = "INVALID_INPUT"
	ERR_InvalidParam = "INVALID_PARAMETER"
	ERR_Decode       = "DECODING"
)
