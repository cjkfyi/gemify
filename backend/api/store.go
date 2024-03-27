package api

import (
	"go.mills.io/bitcask/v2"
)

type Store struct {
	db       bitcask.DB
	dataPath string
}

func InitDataStore(dataPath string) (*Store, error) {
	db, err := bitcask.Open(dataPath)
	if err != nil {
		return nil, err
	}

	return &Store{db: db, dataPath: dataPath}, nil
}

func (s *Store) SaveConversation(conversationID string, conversationData []byte) error {
	return s.db.Put([]byte(conversationID), conversationData)
}

func (s *Store) GetConversation(conversationID string) ([]byte, error) {
	return s.db.Get([]byte(conversationID))
}

func (s *Store) Close() error {
	return s.db.Close()
}
