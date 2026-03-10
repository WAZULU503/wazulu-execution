package cas

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// StoreDir is where CAS artifacts live
const StoreDir = "cas_data"

// Store writes data into the CAS store using SHA256 content addressing.
func Store(data []byte) (string, error) {

	hash := sha256.Sum256(data)
	hashHex := hex.EncodeToString(hash[:])

	path := filepath.Join(StoreDir, hashHex)

	// ensure directory exists
	err := os.MkdirAll(StoreDir, 0755)
	if err != nil {
		return "", err
	}

	// if file already exists we skip writing (deduplication)
	if _, err := os.Stat(path); err == nil {
		return hashHex, nil
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return "", err
	}

	return hashHex, nil
}

// Load retrieves an artifact from CAS by hash.
func Load(hash string) ([]byte, error) {

	path := filepath.Join(StoreDir, hash)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("CAS object not found: %s", hash)
	}

	return data, nil
}
