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
		return nil, errors.New("message input is required")
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
		return nil, errors.New("failed to marshal msg")
	}

	err = (*chat).Put([]byte(key), val)
	if err != nil {
		return nil, errors.New("failed to store msg")
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
		return nil, errors.New("failed to open chat ds")
	}
	defer (*chat).Close()

	err = (*chat).Scan([]byte(""), func(key bitcask.Key) error {

		strKey := string(key)

		if strings.HasSuffix(strKey, msgID) {
			data, err := (*chat).Get(key)
			if err != nil {
				return errors.New("failed to pull msg with key")
			}

			err = json.Unmarshal(data, &msg)
			if err != nil {
				return errors.New("failed to unmarshal msg")
			}

			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	} else {
		return msg, nil
	}
}

//

func ListMessages(
	chatID,
	projID string,
) (
	[]Message,
	error,
) {

	var msgArr []Message
	var msg Message

	chat, err := openChat(projID, chatID)
	if err != nil {
		return nil, errors.New("failed to open chat ds")
	}
	defer (*chat).Close()

	err = (*chat).Scan([]byte(""), func(key bitcask.Key) error {

		data, err := (*chat).Get(key)
		if err != nil {
			return errors.New("failed to pull msg with key")
		}

		err = json.Unmarshal(data, &msg)
		if err != nil {
			return errors.New("failed to unmarshal msg")
		}

		msgArr = append(msgArr, msg)
		return nil
	})
	if err != nil {
		return nil, err
	} else {

		sort.Slice(msgArr, func(i, j int) bool {
			return msgArr[i].LastModified < msgArr[j].LastModified
		})

		return msgArr, nil
	}
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
		return nil, errors.New("failed to open chat ds")
	}
	defer (*chat).Close()

	err = (*chat).Scan([]byte(""), func(k bitcask.Key) error {

		strKey := string(k)

		if strings.HasSuffix(strKey, msgID) {

			key = strKey

			data, err := (*chat).Get([]byte(key))
			if err != nil {
				return errors.New("failed to pull msg with keys")
			}

			err = json.Unmarshal(data, &msg)
			if err != nil {
				return errors.New("failed to unmarshal msg")
			}

			msg = i
			msg.LastModified = int(time.Now().UnixNano())

			val, err := json.Marshal(msg)
			if err != nil {
				return errors.New("failed to marshal msg")
			}

			err = (*chat).Put([]byte(key), val)
			if err != nil {
				return errors.New("failed to store msg")
			}

			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	} else if key == "" {
		return nil, errors.New("failed to find msg with keys")
	} else {
		return &msg, nil
	}
}

//

func DeleteMessage(
	projID,
	chatID,
	msgID string,
) error {

	var msg Message
	var key string

	chat, err := openChat(projID, chatID)
	if err != nil {
		return err
	}
	defer (*chat).Close()

	err = (*chat).Scan([]byte(""), func(k bitcask.Key) error {

		strKey := string(k)

		if strings.HasSuffix(strKey, msgID) {

			key = strKey

			data, err := (*chat).Get([]byte(key))
			if err != nil {
				return errors.New("failed to pull msg with key")
			}

			err = json.Unmarshal(data, &msg)
			if err != nil {
				return errors.New("failed to unmarshal msg")
			}

			msg.IsDeleted = true

			val, err := json.Marshal(msg)
			if err != nil {
				return errors.New("failed to marshal msg")
			}

			err = (*chat).Put([]byte(key), val)
			if err != nil {
				return errors.New("failed to store updated msg")
			}

			return nil
		}
		return nil
	})
	if err != nil {
		return err
	} else {
		return nil
	}
}
