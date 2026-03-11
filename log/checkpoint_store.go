package log

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

type CheckpointRecord struct {
	TreeSize           int    `json:"tree_size"`
	RootHash           string `json:"root_hash"`
	Timestamp          int64  `json:"timestamp"`
	OperatorSignature  string `json:"operator_signature"`
	WitnessSignature   string `json:"witness_signature"`
	PrevCheckpointHash string `json:"prev_checkpoint_hash"`
	CheckpointHash     string `json:"checkpoint_hash"`
}

const checkpointFile = "sth.log"
const genesisFile = "genesis.lock"

func computeCheckpointHash(rec CheckpointRecord) string {

	payload := fmt.Sprintf(
		"%s|%d|%s|%d|%s|%s",
		rec.PrevCheckpointHash,
		rec.TreeSize,
		rec.RootHash,
		rec.Timestamp,
		rec.OperatorSignature,
		rec.WitnessSignature,
	)

	hash := sha256.Sum256([]byte(payload))

	return hex.EncodeToString(hash[:])
}

func loadCheckpointHistory() ([]CheckpointRecord, error) {

	var records []CheckpointRecord

	file, err := os.Open(checkpointFile)
	if err != nil {
		if os.IsNotExist(err) {
			return records, nil
		}
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	for {
		var rec CheckpointRecord

		if err := decoder.Decode(&rec); err != nil {
			break
		}

		records = append(records, rec)
	}

	return records, nil
}

func AppendCheckpoint(
	treeSize int,
	root string,
	timestamp int64,
	operatorSig string,
	witnessSig string,
) error {

	history, err := loadCheckpointHistory()
	if err != nil {
		return err
	}

	var prevHash string

	if len(history) > 0 {
		prevHash = history[len(history)-1].CheckpointHash
	}

	rec := CheckpointRecord{
		TreeSize:           treeSize,
		RootHash:           root,
		Timestamp:          timestamp,
		OperatorSignature:  operatorSig,
		WitnessSignature:   witnessSig,
		PrevCheckpointHash: prevHash,
	}

	rec.CheckpointHash = computeCheckpointHash(rec)

	file, err := os.OpenFile(
		checkpointFile,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(rec); err != nil {
		return err
	}

	// write genesis anchor if first checkpoint
	if len(history) == 0 {
		os.WriteFile(genesisFile, []byte(rec.CheckpointHash), 0644)
	}

	return nil
}
