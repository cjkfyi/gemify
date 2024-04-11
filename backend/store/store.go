package store

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"go.mills.io/bitcask/v2"

	"gemify/models"
)

//
// Project

func openMeta() (
	*bitcask.DB,
	error,
) {
	meta, err := bitcask.Open(path.Join(dataPath, "meta"))
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

func CreateProject(
	i *models.Project,
) (
	*models.Project,
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

	meta, err := openMeta()
	if err != nil {
		return nil, errors.New("failed to open meta ds")
	}
	defer (*meta).Close()

	projID := GenID()
	stamp := int(time.Now().UnixNano())

	i.ProjID = projID
	i.FirstCreated = stamp
	i.LastModified = stamp
	i.Chats = []models.Chat{}

	projPath := path.Join(dataPath, projID)
	if err := os.Mkdir(projPath, 0755); err != nil {
		return nil, errors.New("failed to mk the proj dir")
	}

	val, err := json.Marshal(i)
	if err != nil {
		return nil, errors.New("failed to marshal proj")
	}

	key := keygen(stamp, projID)

	if err := (*meta).Put([]byte(key), val); err != nil {
		return nil, errors.New("failed to store proj in meta ds")
	}

	return i, nil
}

func GetProject(
	projID string,
) (
	*models.Project,
	error,
) {

	var project models.Project

	if projID == "" {
		return nil, errors.New("projID param is required")
	}

	meta, err := openMeta()
	if err != nil {
		return nil, errors.New("failed to open meta ds")
	}
	defer (*meta).Close()

	err = (*meta).Scan([]byte(""), func(key bitcask.Key) error {

		strKey := string(key)

		parts := strings.Split(strKey, ":")
		if len(parts) == 2 && parts[1] == projID {

			data, err := (*meta).Get(key)
			if err != nil {
				return errors.New("failed to find proj with projID")
			}

			err = json.Unmarshal(data, &project)
			if err != nil {
				return errors.New("failed to unmarshal proj")
			}

			return nil
		}
		return nil
	})

	if err != nil {
		return nil, err
	} else {
		return &project, nil
	}
}

func ListProjects() (
	[]models.Project,
	error,
) {

	var projArr []models.Project

	meta, err := openMeta()
	if err != nil {
		return nil, errors.New("failed to open meta ds")
	}
	defer (*meta).Close()

	err = (*meta).Scan([]byte(""), func(key bitcask.Key) error {

		strKey := string(key)

		data, err := (*meta).Get([]byte(strKey))
		if err != nil {
			return errors.New("failed to get pair in meta ds")
		}

		var project models.Project

		err = json.Unmarshal(data, &project)
		if err != nil {
			return errors.New("failed to unmarshal proj")
		}

		projArr = append(projArr, project)

		return nil
	})

	if err != nil {
		return nil, err
	} else {
		return projArr, nil
	}
}

func UpdateProject(
	projID string,
	i models.Project,
) (
	*models.Project,
	error,
) {

	if projID == "" {
		return nil, errors.New("projID param is required")
	}

	project, err := GetProject(projID)
	if err != nil {
		return nil, err
	}

	stamp := int(time.Now().UnixNano())

	project.LastModified = stamp

	if i.Name != "" {
		if len(i.Name) > 160 {
			return nil, errors.New("name cannot exceed 160 chars")
		} else {
			project.Name = i.Name
		}
	}
	if i.Desc != "" {
		if len(i.Desc) > 260 {
			return nil, errors.New("desc cannot exceed 260 chars")
		} else {
			project.Desc = i.Desc
		}
	}

	meta, err := openMeta()
	if err != nil {
		return nil, errors.New("failed to open meta ds")
	}
	defer (*meta).Close()

	var oldKey string

	key := keygen(stamp, project.ProjID)

	val, err := json.Marshal(project)
	if err != nil {
		return nil, errors.New("failed to marshal updated proj")
	}

	err = (*meta).Scan([]byte(""), func(k bitcask.Key) error {

		metaKey := string(k)

		projID, _, err := extractKey(metaKey)
		if err == nil && projID == project.ProjID {
			oldKey = metaKey
		}
		return nil
	})
	if err != nil || oldKey == "" {
		return nil, errors.New("failed to scan meta ds for key")
	}

	err = (*meta).Delete([]byte(oldKey))
	if err != nil {
		return nil, errors.New("failed to delete old proj entry")
	}

	err = (*meta).Put([]byte(key), val)
	if err != nil {
		return nil, errors.New("failed to store new proj entry")
	}

	return project, nil
}

func DeleteProject(
	projID string,
) error {

	if projID == "" {
		return errors.New("projID param is required")
	}

	meta, err := openMeta()
	if err != nil {
		return errors.New("failed to open meta ds")
	}
	defer (*meta).Close()

	var key string

	err = (*meta).Scan([]byte(""), func(k bitcask.Key) error {

		metaKey := string(k)

		keyID, _, err := extractKey(metaKey)
		if err == nil && keyID == projID {

			key = metaKey

			return nil
		}
		return nil
	})
	if err != nil || key == "" {
		return errors.New("failed to find proj with projID")
	}

	err = (*meta).Delete([]byte(key))
	if err != nil {
		return errors.New("failed to del proj entry")
	} else {
		return nil
	}
}

//
// Chat

func openChat(
	projID, chatID string,
) (
	*bitcask.DB,
	error,
) {

	projPath := path.Join(dataPath, projID)

	if _, err := os.Stat(projPath); os.IsNotExist(err) {
		return nil, errors.New("failed to find proj with projID")
	}

	chatPath := path.Join(projPath, chatID)
	chat, err := bitcask.Open(chatPath)
	if err != nil {
		return nil, errors.New("failed to open chat ds")
	}

	return &chat, nil
}

func addChat(
	project *models.Project,
) error {

	val, err := json.Marshal(project)
	if err != nil {
		return errors.New("failed to marshal proj")
	}

	metaDB, err := openMeta()
	if err != nil {
		return errors.New("failed to open meta ds")
	}
	defer (*metaDB).Close()

	var oldKey string

	err = (*metaDB).Scan([]byte(""), func(key bitcask.Key) error {

		metaKey := string(key)

		projID, _, err := extractKey(metaKey)
		if err == nil && projID == project.ProjID {
			oldKey = metaKey
		}

		return nil
	})

	if err != nil || oldKey == "" {
		return errors.New("failed to find proj with projID")
	}

	stamp := int(time.Now().UnixNano())
	newKey := keygen(stamp, project.ProjID)

	err = (*metaDB).Put([]byte(newKey), val)
	if err != nil {
		return errors.New("failed to store new chat entity")
	}

	err = (*metaDB).Delete([]byte(oldKey))
	if err != nil {
		return errors.New("failed to delete old chat entity")
	}

	return nil
}

func CreateChat(
	i *models.Chat,
) (
	*models.Chat,
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

	err = addChat(proj)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func GetChat(
	projID, chatID string,
) (
	*models.Chat,
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

func ListChats(
	projID string,
) (
	[]models.Chat,
	error,
) {

	var chatArr []models.Chat

	project, err := GetProject(projID)
	if err != nil {
		return nil, err
	}

	chatArr = append(
		chatArr,
		project.Chats...,
	)

	return chatArr, nil
}

func UpdateChat(
	projID, chatID string,
	i models.Chat,
) (
	*models.Chat,
	error,
) {

	var chatRes *models.Chat

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
		return nil, errors.New("failed to open meta ds")
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

func DeleteChat(
	projID, chatID string,
) error {

	if projID == "" {
		return errors.New("projID param is required")
	}
	if chatID == "" {
		return errors.New("chatID param is required")
	}

	meta, err := openMeta()
	if err != nil {
		return errors.New("failed to open meta ds")
	}
	defer (*meta).Close()

	err = (*meta).Scan([]byte(""), func(key bitcask.Key) error {

		metaKey := string(key)

		keyID, _, err := extractKey(metaKey)
		if err == nil && keyID == projID {

			var project models.Project

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
	if err != nil {
		return err
	} else {
		return nil
	}
}

//
// Message

func AddMessage(
	chatID, projID string,
	i *models.Message,
) (
	*models.Message,
	error,
) {
	chat, err := openChat(projID, chatID)
	if err != nil {
		return nil, errors.New("failed to open chat ds")
	}
	defer (*chat).Close()

	//

	//
	//
	return nil, nil
}

func CreateMessage(
	chatID, projID string,
	i *models.Message,
) (
	*models.Message,
	error,
) {

	chat, err := openChat(projID, chatID)
	if err != nil {
		return nil, errors.New("failed to open chat ds")
	}
	defer (*chat).Close()

	msgID := GenID()

	messaged := i

	msgJSON, err := json.Marshal(messaged)
	if err != nil {
		return nil, errors.New("failed to marshal msg")
	}

	// messageKey := fmt.Sprintf("%s:%d:%s", chatID, stamp, msgID)
	err = (*chat).Put([]byte(msgID), msgJSON)
	if err != nil {
		return nil, errors.New("failed to store msg")
	}

	return messaged, nil
}

func GetMessage(
	projID, chatID, msgID string,
) (
	*models.Message,
	error,
) {

	var msg *models.Message

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

func ListMessages(
	chatID, projID string,
) (
	[]models.Message,
	error,
) {

	var msgArr []models.Message
	var msg models.Message

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

func UpdateMessage(
	projID, chatID, msgID string,
	i models.Message,
) (
	*models.Message,
	error,
) {

	var msg models.Message
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

func DeleteMessage(
	projID, chatID, msgID string,
) error {

	var msg models.Message
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
