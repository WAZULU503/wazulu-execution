package log

import (
	"crypto/sha256"
	"encoding/hex"
)

// ComputeMerkleRoot builds a Merkle root from a list of entry hashes.
// Leaves are domain separated to avoid hash reuse across contexts.
func ComputeMerkleRoot(entryHashes []string) string {

	if len(entryHashes) == 0 {
		return ""
	}

	// convert hex hashes to raw leaf hashes
	var level [][]byte

	for _, h := range entryHashes {

		leaf := sha256.Sum256([]byte("WZMERKLEv1|" + h))
		level = append(level, leaf[:])
	}

	// build tree upward
	for len(level) > 1 {

		var next [][]byte

		for i := 0; i < len(level); i += 2 {

			if i+1 == len(level) {
				next = append(next, level[i])
				continue
			}

			combined := append(level[i], level[i+1]...)
			hash := sha256.Sum256(combined)

			next = append(next, hash[:])
		}

		level = next
	}

	return hex.EncodeToString(level[0])
}
