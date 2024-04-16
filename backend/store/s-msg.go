package store

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"go.mills.io/bitcask/v2"
)

func CreateMessage(
	chatID,
	projID,
	msg string,
	isUser bool,
) (
	*Message,
	error,
) {
	if msg == "" {
		return nil, errors.New("`message` field is required")
	}

	msgID := GenID()
	stamp := int(time.Now().UnixNano())
	key := keygen(stamp, msgID)

	new := &Message{
		MsgID:        msgID,
		ChatID:       chatID,
		ProjID:       projID,
		IsUser:       isUser,
		Message:      msg,
		LastModified: stamp,
		FirstCreated: stamp,
	}

	chat, err := openChat(projID, chatID)
	if err != nil {
		return nil, err
	}
	defer (*chat).Close()

	val, err := json.Marshal(new)
	if err != nil {
		return nil, errors.New("failed ds op")
	}

	err = (*chat).Put([]byte(key), val)
	if err != nil {
		return nil, errors.New("failed ds op")
	}

	return new, nil
}

//

func GetMessage(
	projID,
	chatID,
	msgID string,
) (
	*Message,
	error,
) {
	var msg *Message

	chat, err := openChat(projID, chatID)
	if err != nil {
		return nil, err
	}
	defer (*chat).Close()

	err = (*chat).Scan([]byte(""), func(key bitcask.Key) error {
		chatKey := string(key)
		if strings.HasSuffix(chatKey, msgID) {
			data, err := (*chat).Get(key)
			if err != nil {
				return errors.New("failed ds op")
			}

			err = json.Unmarshal(data, &msg)
			if err != nil {
				return errors.New("failed ds op")
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	} else if msg == nil {
		return nil, errors.New("invalid `msgID` parameter")
	}
	return msg, nil
}

//

func ListMessages(
	projID,
	chatID string,
) (
	[]Message,
	error,
) {
	var msgArr []Message

	chat, err := openChat(projID, chatID)
	if err != nil {
		return nil, err
	}
	defer (*chat).Close()

	err = (*chat).Scan([]byte(""), func(key bitcask.Key) error {

		data, err := (*chat).Get(key)
		if err != nil {
			return errors.New("failed ds op")
		}

		var msg Message

		err = json.Unmarshal(data, &msg)
		if err != nil {
			return errors.New("failed ds op")
		}

		msgArr = append(msgArr, msg)
		return nil
	})
	if err != nil {
		return nil, err
	}

	if msgArr == nil {
		msgArr = []Message{}
		return msgArr, nil
	}

	sort.Slice(msgArr, func(i, j int) bool {
		return msgArr[i].LastModified < msgArr[j].LastModified
	})
	return msgArr, nil
}

//

func UpdateMessage(
	projID,
	chatID,
	msgID string,
	i Message,
) (
	*Message,
	error,
) {
	var msg Message
	var key string

	chat, err := openChat(projID, chatID)
	if err != nil {
		return nil, err
	}
	defer (*chat).Close()

	err = (*chat).Scan([]byte(""), func(k bitcask.Key) error {

		chatKey := string(k)

		if strings.HasSuffix(chatKey, msgID) {

			key = chatKey

			data, err := (*chat).Get([]byte(key))
			if err != nil {
				return errors.New("failed ds op")
			}

			err = json.Unmarshal(data, &msg)
			if err != nil {
				return errors.New("failed ds op")
			}

			if i.Message != "" {
				msg.Message = i.Message
			}
			msg.IsUser = i.IsUser

			msg.LastModified = int(time.Now().UnixNano())

			val, err := json.Marshal(msg)
			if err != nil {
				return errors.New("failed ds op")
			}

			err = (*chat).Put([]byte(key), val)
			if err != nil {
				return errors.New("failed ds op")
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	} else if msg.MsgID == "" {
		return nil, errors.New("invalid `msgID` parameter")
	}

	return &msg, nil
}

//

func DeleteMessage(
	projID,
	chatID,
	msgID string,
) error {

	chat, err := openChat(projID, chatID)
	if err != nil {
		return err
	}
	defer (*chat).Close()

	found := false

	err = (*chat).Scan([]byte(""), func(k bitcask.Key) error {
		chatKey := string(k)
		if strings.HasSuffix(chatKey, msgID) {
			found = true
			err := (*chat).Delete([]byte(chatKey))
			if err != nil {
				return errors.New("failed ds op")
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	} else if !found {
		return errors.New("invalid `msgID` parameter")
	} else {
		return nil
	}
}
