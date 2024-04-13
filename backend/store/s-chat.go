package store

import (
	"encoding/json"
	"errors"
	"time"

	"go.mills.io/bitcask/v2"
)

func CreateChat(
	i *Chat,
) (
	*Chat,
	error,
) {

	if i.Name == "" {
		return nil, errors.New("name param is required")
	} else if len(i.Name) > 160 {
		return nil, errors.New("name cannot exceed 160 chars")
	}

	if i.Desc == "" {
		return nil, errors.New("desc param is required")
	} else if len(i.Desc) > 260 {
		return nil, errors.New("desc cannot exceed 260 chars")
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

	err = updateProject(proj)
	if err != nil {
		return nil, err
	}

	return i, nil
}

//

func GetChat(
	projID,
	chatID string,
) (
	*Chat,
	error,
) {

	project, err := GetProject(projID)
	if err != nil {
		return nil, err
	}

	for _, chat := range project.Chats {
		if chat.ChatID == chatID {
			return &chat, nil
		}
	}

	return nil, errors.New("failed to find chat with chatID")
}

//

func ListChats(
	projID string,
) (
	[]Chat,
	error,
) {

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

func UpdateChat(
	projID,
	chatID string,
	i Chat,
) (
	*Chat,
	error,
) {

	var chatRes *Chat

	if projID == "" {
		return nil, errors.New("projID is required")
	}
	if chatID == "" {
		return nil, errors.New("chatID is required")
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
		return nil, errors.New("chat not found with chatID")
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

	newKey := keygen(stamp, proj.ProjID)

	proj.Chats[chatIndex].LastModified = stamp

	val, err := json.Marshal(proj)
	if err != nil {
		return nil, errors.New("failed to marshal proj")
	}

	meta, err := openMeta()
	if err != nil {
		return nil, err
	}
	defer (*meta).Close()

	var oldKey string

	err = (*meta).Scan([]byte(""), func(key bitcask.Key) error {

		metaKey := string(key)

		projID, _, err := extractKey(metaKey)
		if err == nil && projID == proj.ProjID {
			oldKey = metaKey
		}
		return nil
	})
	if err != nil || oldKey == "" {
		return nil, errors.New("failed to find proj with projID")
	}

	err = (*meta).Delete([]byte(oldKey))
	if err != nil {
		return nil, errors.New("failed to delete old chat entity")
	}

	err = (*meta).Put([]byte(newKey), val)
	if err != nil {
		return nil, errors.New("failed to store new chat entity")
	}

	return chatRes, nil
}

//

func DeleteChat(
	projID,
	chatID string,
) error {

	var project *Project

	if projID == "" {
		return errors.New("projID param is required")
	}
	if chatID == "" {
		return errors.New("chatID param is required")
	}

	meta, err := openMeta()
	if err != nil {
		return err
	}
	defer (*meta).Close()

	err = (*meta).Scan([]byte(""), func(key bitcask.Key) error {

		metaKey := string(key)

		keyID, _, err := extractKey(metaKey)
		if err == nil && keyID == projID {

			data, err := (*meta).Get(key)
			if err != nil {
				return errors.New("failed to find proj with projID")
			}

			err = json.Unmarshal(data, &project)
			if err != nil {
				return errors.New("failed to unmarshal proj")
			}

			for i, chat := range project.Chats {
				if chat.ChatID == chatID {
					project.Chats = append(
						project.Chats[:i],
						project.Chats[i+1:]...,
					)
					break
				} else {
					return errors.New("failed to find chat with chatID")
				}
			}

			val, err := json.Marshal(project)
			if err != nil {
				return errors.New("failed to marshal proj")
			}

			err = (*meta).Put([]byte(key), val)
			if err != nil {
				return errors.New("failed to store new proj")
			}

			return nil
		}
		return nil
	})
	if project == nil {
		return errors.New("proj returned nil")
	}
	if err != nil {
		return err
	} else {
		return nil
	}
}
