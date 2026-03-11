package log

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type SignedTreeHead struct {
	TreeSize  int    `json:"tree_size"`
	RootHash  string `json:"root_hash"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

var operatorPrivateKey ed25519.PrivateKey
var operatorPublicKey ed25519.PublicKey

func init() {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	operatorPrivateKey = priv
	operatorPublicKey = pub
}

func GenerateSTH(rootHash string, treeSize int) SignedTreeHead {
	timestamp := time.Now().Unix()

	payload := fmt.Sprintf("WZ-STH-V1|%d|%s|%d", treeSize, rootHash, timestamp)

	hash := sha256.Sum256([]byte(payload))

	sig := ed25519.Sign(operatorPrivateKey, hash[:])

	return SignedTreeHead{
		TreeSize:  treeSize,
		RootHash:  rootHash,
		Timestamp: timestamp,
		Signature: hex.EncodeToString(sig),
	}
}
