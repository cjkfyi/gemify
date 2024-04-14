package store

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strings"
	"time"

	"go.mills.io/bitcask/v2"
)

func CreateProject(i *Project) (*Project, error) {

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

	projID := GenID()
	stamp := int(time.Now().UnixNano())

	i.ProjID = projID
	i.FirstCreated = stamp
	i.LastModified = stamp
	i.Chats = []Chat{}

	projPath := path.Join(dataPath, projID)
	if err := os.Mkdir(projPath, 0755); err != nil {
		return nil, errors.New("failed ds op")
	}

	val, err := json.Marshal(i)
	if err != nil {
		return nil, errors.New("failed ds op")
	}

	key := keygen(stamp, projID)

	meta, err := openMeta()
	if err != nil {
		return nil, err
	}
	defer (*meta).Close()

	if err := (*meta).Put([]byte(key), val); err != nil {
		return nil, errors.New("failed ds op")
	}

	return i, nil
}

//

func GetProject(projID string) (*Project, error) {

	var proj *Project

	meta, err := OpenMeta(projID)
	if err != nil {
		return nil, err
	}
	defer (*meta).Close()

	err = (*meta).Scan([]byte(""), func(k bitcask.Key) error {
		metaKey := string(k)
		parts := strings.Split(metaKey, ":")
		if len(parts) == 2 && parts[1] == projID {
			data, err := (*meta).Get(k)
			if err != nil {
				return errors.New("failed ds op")
			}
			err = json.Unmarshal(data, &proj)
			if err != nil {
				return errors.New("failed ds op")
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	} else {
		return proj, nil
	}
}

//

func ListProjects() ([]Project, error) {

	var projArr []Project

	meta, err := openMeta()
	if err != nil {
		return nil, err
	}
	defer (*meta).Close()

	err = (*meta).Scan([]byte(""), func(k bitcask.Key) error {
		metaKey := string(k)
		data, err := (*meta).Get([]byte(metaKey))
		if err != nil {
			return errors.New("failed ds op")
		}
		var project Project
		err = json.Unmarshal(data, &project)
		if err != nil {
			return errors.New("failed ds op")
		}
		projArr = append(projArr, project)
		return nil
	})

	if projArr == nil {
		projArr = []Project{}
	}

	if err != nil {
		return nil, err
	} else {
		return projArr, nil
	}
}

//

func UpdateProject(projID string, i Project) (*Project, error) {

	var oldKey string

	proj, err := GetProject(projID)
	if err != nil {
		return nil, err
	}

	stamp := int(time.Now().UnixNano())
	key := keygen(stamp, proj.ProjID)
	proj.LastModified = stamp

	if i.Name != "" {
		if len(i.Name) > 160 {
			return nil, errors.New("name cannot exceed 160 chars")
		} else {
			proj.Name = i.Name
		}
	}
	if i.Desc != "" {
		if len(i.Desc) > 260 {
			return nil, errors.New("desc cannot exceed 260 chars")
		} else {
			proj.Desc = i.Desc
		}
	}

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

	err = (*meta).Put([]byte(key), val)
	if err != nil {
		return nil, errors.New("failed ds op")
	}

	return proj, nil
}

//

func DeleteProject(projID string) error {

	var key string

	meta, err := OpenMeta(projID)
	if err != nil {
		return err
	}
	defer (*meta).Close()

	err = (*meta).Scan([]byte(""), func(k bitcask.Key) error {
		metaKey := string(k)
		keyID, _, err := extractKey(metaKey)
		if err == nil && keyID == projID {
			key = metaKey
			return nil
		}
		return nil
	})
	if key == "" {
		return errors.New("invalid projID parameter")
	} else if err != nil {
		return errors.New("failed ds op")
	}

	err = (*meta).Delete([]byte(key))
	if err != nil {
		return errors.New("failed ds op")
	}

	projDir := path.Join(dataPath, projID)
	err = os.RemoveAll(projDir)
	if err != nil {
		return errors.New("failed ds op")
	} else {
		return nil
	}
}
