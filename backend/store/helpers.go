package store

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/icza/bitio"
)

// Dumping grounds
const dataPath = ".data"

//

// Used to validate a specific projID
func isProj(projID string) error {
	path := path.Join(dataPath, projID)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return errors.New(
			"projID is invalid",
		)
	} else {
		return nil
	}
}

// Used to validate a pair of proj & chat IDs
func isChat(projID, chatID string) error {
	projPath := path.Join(dataPath, projID)
	chatPath := path.Join(projPath, chatID)
	_, err := os.Stat(projPath)
	if os.IsNotExist(err) {
		return errors.New(
			"projID is invalid",
		)
	} else {
		_, err := os.Stat(chatPath)
		if os.IsNotExist(err) {
			return errors.New(
				"chatID is invalid",
			)
		} else {
			return nil
		}
	}
}

// as long as we have strict creations,
// these helpers should make good sense

//

func GenID() string {
	randomBytes := make([]byte, 3)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(randomBytes)
}

func keygen(stamp int, id string) string {
	return fmt.Sprintf("%d:%s", stamp, id)
}

func extractKey(key string) (string, string, error) {
	parts := strings.Split(key, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf(
			"invalid key fmt: %s",
			key,
		)
	}
	timestamp := parts[0]
	id := parts[1]
	return id, timestamp, nil
}

func pack(data []byte) []byte {
	var packed bytes.Buffer
	writer := bitio.NewWriter(&packed)
	writer.Write(data)
	writer.Close()
	return packed.Bytes()
}

func unpack(data []byte) ([]byte, error) {
	r := bitio.NewReader(bytes.NewReader(data))
	var unpacked bytes.Buffer
	_, err := io.Copy(&unpacked, r)
	if err != nil {
		return nil, err
	}
	return unpacked.Bytes(), nil
}
