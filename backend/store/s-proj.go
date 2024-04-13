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

func CreateProject(
	i *Project,
) (
	*Project,
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

	projID := GenID()
	stamp := int(time.Now().UnixNano())

	i.ProjID = projID
	i.FirstCreated = stamp
	i.LastModified = stamp
	i.Chats = []Chat{}

	projPath := path.Join(dataPath, projID)
	if err := os.Mkdir(projPath, 0755); err != nil {
		return nil, errors.New("failed to mk the proj dir")
	}

	val, err := json.Marshal(i)
	if err != nil {
		return nil, errors.New("failed to marshal proj")
	}

	key := keygen(stamp, projID)

	meta, err := openMeta()
	if err != nil {
		return nil, err
	}
	defer (*meta).Close()

	if err := (*meta).Put([]byte(key), val); err != nil {
		return nil, errors.New("failed to store proj in meta ds")
	}

	return i, nil
}

//

func GetProject(
	projID string,
) (
	*Project,
	error,
) {

	var proj *Project

	meta, err := openMeta()
	if err != nil {
		return nil, err
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

			err = json.Unmarshal(data, &proj)
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

		if proj == nil {
			return nil, errors.New("proj returned nil")
		}

		return proj, nil
	}
}

//

func ListProjects() (
	[]Project,
	error,
) {

	var projArr []Project

	meta, err := openMeta()
	if err != nil {
		return nil, err
	}
	defer (*meta).Close()

	err = (*meta).Scan([]byte(""), func(key bitcask.Key) error {

		strKey := string(key)

		data, err := (*meta).Get([]byte(strKey))
		if err != nil {
			return errors.New("failed to get pair in meta ds")
		}

		var project Project

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
		if projArr == nil {
			projArr = []Project{}
		}
		return projArr, nil
	}
}

//

func UpdateProject(
	projID string,
	i Project,
) (
	*Project,
	error,
) {

	var oldKey string

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

	key := keygen(stamp, project.ProjID)

	val, err := json.Marshal(project)
	if err != nil {
		return nil, errors.New("failed to marshal updated proj")
	}

	meta, err := openMeta()
	if err != nil {
		return nil, err
	}
	defer (*meta).Close()

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

	if project == nil {
		return nil, errors.New("proj returned nil")
	} else {
		return project, nil
	}
}

//

func DeleteProject(
	projID string,
) error {

	var key string

	if projID == "" {
		return errors.New("projID param is required")
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
