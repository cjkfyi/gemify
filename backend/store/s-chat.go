package store

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"time"

	"go.mills.io/bitcask/v2"
)

func CreateChat(i *Chat) (*Chat, error) {

	if i.Name == "" {
		return nil, errors.New("`name` param is required")
	} else if len(i.Name) > 160 {
		return nil, errors.New("`name` cannot exceed 160 chars")
	}
	if i.Desc == "" {
		return nil, errors.New("`desc` param is required")
	} else if len(i.Desc) > 260 {
		return nil, errors.New("`desc` cannot exceed 260 chars")
	}

	chatID := GenID()
	stamp := int(time.Now().UnixNano())

	i.ChatID = chatID
	i.FirstCreated = stamp
	i.LastModified = stamp

	proj, err := GetProject(i.ProjID)
	if err != nil {
		return nil, err
	}

	proj.Chats = append(proj.Chats, *i)
	proj.LastModified = stamp

	projDir := path.Join(dataPath, proj.ProjID)
	chat := path.Join(projDir, chatID)
	meta, err := bitcask.Open(chat)
	if err != nil {
		return nil, errors.New("failed ds op")
	}
	defer meta.Close()

	err = updateProject(proj)
	if err != nil {
		return nil, err
	}

	return i, nil
}

//

func GetChat(projID, chatID string) (*Chat, error) {

	err := isChat(projID, chatID)
	if err != nil {
		return nil, err
	}

	project, err := GetProject(projID)
	if err != nil {
		return nil, err
	}

	for _, chat := range project.Chats {
		if chat.ChatID == chatID {
			return &chat, nil
		}
	}

	return nil, errors.New("failed ds op")
}

//

func ListChats(projID string) ([]Chat, error) {

	var chatArr []Chat

	project, err := GetProject(projID)
	if err != nil {
		return nil, err
	}

	chatArr = append(
		chatArr,
		project.Chats...,
	)

	if chatArr == nil {
		chatArr = []Chat{}
	}

	return chatArr, nil
}

//

func UpdateChat(projID, chatID string, i Chat) (*Chat, error) {

	var chatRes *Chat
	var oldKey string

	if projID == "" {
		return nil, errors.New("`projID` param is required")
	}
	if chatID == "" {
		return nil, errors.New("`chatID` param is required")
	}

	err := isChat(projID, chatID)
	if err != nil {
		return nil, err
	}
	proj, err := GetProject(projID)
	if err != nil {
		return nil, err
	}

	chatIndex := -1

	for i, chat := range proj.Chats {
		if chat.ChatID == chatID {
			chatRes = &proj.Chats[i]
			chatIndex = i
			break
		}
	}
	if chatIndex == -1 {
		return nil, errors.New("failed ds op")
	}

	if i.Name != "" {
		if len(i.Name) > 160 {
			return nil, errors.New("name cannot exceed 160 chars")
		} else {
			proj.Chats[chatIndex].Name = i.Name
		}
	}
	if i.Desc != "" {
		if len(i.Desc) > 260 {
			return nil, errors.New("desc cannot exceed 260 chars")
		} else {
			proj.Chats[chatIndex].Desc = i.Desc
		}
	}

	stamp := int(time.Now().UnixNano())
	proj.Chats[chatIndex].LastModified = stamp
	newKey := keygen(stamp, proj.ProjID)

	val, err := json.Marshal(proj)
	if err != nil {
		return nil, errors.New("failed ds op")
	}

	meta, err := openMeta()
	if err != nil {
		return nil, err
	}
	defer (*meta).Close()

	err = (*meta).Scan([]byte(""), func(k bitcask.Key) error {
		metaKey := string(k)
		projID, _, err := extractKey(metaKey)
		if err == nil && projID == proj.ProjID {
			oldKey = metaKey
		}
		return nil
	})
	if err != nil || oldKey == "" {
		return nil, errors.New("failed ds op")
	}

	err = (*meta).Delete([]byte(oldKey))
	if err != nil {
		return nil, errors.New("failed ds op")
	}

	err = (*meta).Put([]byte(newKey), val)
	if err != nil {
		return nil, errors.New("failed ds op")
	}

	return chatRes, nil
}

//

func DeleteChat(projID, chatID string) error {

	var project *Project

	if projID == "" {
		return errors.New("`projID` param is required")
	}
	if chatID == "" {
		return errors.New("`chatID` param is required")
	}

	err := isChat(projID, chatID)
	if err != nil {
		return err
	}

	meta, err := openMeta()
	if err != nil {
		return err
	}
	defer (*meta).Close()

	err = (*meta).Scan([]byte(""), func(k bitcask.Key) error {
		metaKey := string(k)
		keyID, _, err := extractKey(metaKey)
		if err == nil && keyID == projID {

			data, err := (*meta).Get(k)
			if err != nil {
				return errors.New("failed ds op")
			}

			err = json.Unmarshal(data, &project)
			if err != nil {
				return errors.New("failed ds op")
			}

			for i, chat := range project.Chats {
				if chat.ChatID == chatID {
					project.Chats = append(
						project.Chats[:i],
						project.Chats[i+1:]...,
					)
					break
				}
			}

			val, err := json.Marshal(project)
			if err != nil {
				return errors.New("failed ds op")
			}

			err = (*meta).Put([]byte(k), val)
			if err != nil {
				return errors.New("failed ds op")
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	} else {
		projDir := path.Join(dataPath, projID)
		chatDir := path.Join(projDir, chatID)
		err = os.RemoveAll(chatDir)
		if err != nil {
			return errors.New("failed ds op")
		} else {
			return nil
		}
	}
}
