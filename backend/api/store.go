package api

import (
	"bytes"
	"fmt"
	"io"
	"sort"

	"github.com/goccy/go-json"
	"github.com/icza/bitio"
	"go.mills.io/bitcask/v2"
)

type Chat struct {
	history bitcask.DB
	list    bitcask.DB
}

func InitDataStores() (*Chat, error) {
	historyDB, err := bitcask.Open("data/history")
	if err != nil {
		return nil, err
	}
	listDB, err := bitcask.Open("data/list")
	if err != nil {
		return nil, err
	}
	return &Chat{
		history: historyDB,
		list:    listDB,
	}, nil
}

func (c *Chat) GracefulClosure() error {

	// TODO: Fix hacky solution
	err1 := c.history.Close()
	err2 := c.list.Close()

	// Check if either closure resulted in an error
	if err1 != nil || err2 != nil {
		// Combine errors if necessary, or return one of them
		return fmt.Errorf("errors during closure: %v, %v", err1, err2)
	}
	return nil
}

//

//

// GetFullConvoList()
// GetShortConvoList()
// AddNewConvo()
// AdjustConvo()

// GetConvoHistory()
// AddConvoHistory()
// AdjustConvoHistory()

//

func (c *Chat) GetConvoHistory(convoID string) ([]byte, error) {
	return c.history.Get([]byte(convoID))
}

func (c *Chat) UpdateConvoHistory(convoID string, convoData []byte) error {
	return c.list.Put([]byte(convoID), convoData)
}

//

func (c *Chat) SaveNewConvo(convoID string, data Convo) error {
	metadataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error serializing conversation data: %v", err)
	}

	var compressedBuffer bytes.Buffer
	compressedWriter := bitio.NewWriter(&compressedBuffer)
	_, err = compressedWriter.Write(metadataBytes)
	if err != nil {
		return fmt.Errorf("error compressing data: %v", err)
	}
	err = compressedWriter.Close()
	if err != nil {
		return fmt.Errorf("error closing compressed writer: %v", err)
	}

	return c.list.Put([]byte(convoID), compressedBuffer.Bytes())
}

func (c *Chat) GetConvoList() ([]Convo, error) {
	var listArr []Convo // For each key we will...
	err := c.list.ForEach(func(key bitcask.Key) error {
		compressedBytes, err := c.list.Get(key)
		if err != nil {
			return fmt.Errorf("err compressing data: %v", err)
		}
		reader := bitio.NewReader(bytes.NewReader(compressedBytes))
		decompressedBytes, err := io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("err decompressing data: %v", err)
		}
		// Process metadata
		var listItem Convo
		err = json.Unmarshal(decompressedBytes, &listItem)
		if err != nil {
			return fmt.Errorf("err parsing metadata for key %s: %v", key, err)
		}
		// Add to the list of conversations
		listArr = append(listArr, listItem)
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(listArr, func(i, j int) bool {
		return listArr[i].LastModified.After(listArr[j].LastModified)
	})

	return listArr, nil
}

// Current datastore design:

// 2 buckets
//    - history
//  	 - convoID (k)
//  	 	- Message
// 			- Sender
// 			- IsUser
//			- Stamp

//    - list
//		 - convoID (k)
//			- ID
//			- Title
//			- LastModified
