package store

import (
	"encoding/json"
	"errors"
	"path"
	"time"

	"go.mills.io/bitcask/v2"
)

func openMeta() (
	*bitcask.DB,
	error,
) {
	meta, err := bitcask.Open(
		path.Join(dataPath, "meta"),
	)
	if err != nil {
		return nil, errors.New(
			"failed to open meta ds",
		)
	}
	return &meta, nil
}

func openChat(
	projID,
	chatID string,
) (
	*bitcask.DB,
	error,
) {
	err := isChat(projID, chatID)
	if err != nil {
		return nil, err
	}

	projPath := path.Join(dataPath, projID)
	chatPath := path.Join(projPath, chatID)

	chat, err := bitcask.Open(chatPath)
	if err != nil {
		return nil, errors.New(
			"failed to open chat ds",
		)
	}
	return &chat, nil
}

func updateProject(
	i *Project,
) error {

	var oldKey string

	stamp := int(time.Now().UnixNano())
	newKey := keygen(stamp, i.ProjID)

	val, err := json.Marshal(i)
	if err != nil {
		return errors.New("failed to marshal")
	}

	meta, err := openMeta()
	if err != nil {
		return err
	}
	defer (*meta).Close()

	err = (*meta).Scan([]byte(""), func(key bitcask.Key) error {

		metaKey := string(key)

		projID, _, err := extractKey(metaKey)
		if err == nil && projID == i.ProjID {
			oldKey = metaKey
		}

		return nil
	})
	if err != nil || oldKey == "" {
		return errors.New("failed ds op")
	}

	err = (*meta).Delete([]byte(oldKey))
	if err != nil {
		return errors.New("failed ds op")
	}

	err = (*meta).Put([]byte(newKey), val)
	if err != nil {
		return errors.New("failed ds op")
	}

	return nil
}
