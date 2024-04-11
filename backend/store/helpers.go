package store

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/icza/bitio"
)

// Dumping grounds
const dataPath = ".data"

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
