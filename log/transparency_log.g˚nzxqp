package log

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const LedgerFile = "ledger.jsonl"

type LogEntry struct {
	ProtocolVersion int    `json:"protocol_version"`
	DecisionVersion int    `json:"decision_version"`

	Seq       uint64 `json:"seq"`
	Timestamp int64  `json:"timestamp"`

	ActorID   string `json:"actor_id"`
	EventType string `json:"event_type"`

	PayloadHash string `json:"payload_hash"`

	PrevHash  string `json:"prev_hash"`
	EntryHash string `json:"entry_hash"`
}

func computeEntryHash(e LogEntry) string {

	data := fmt.Sprintf(
		"WZLOGv1|%d|%d|%d|%d|%s|%s|%s|%s",
		e.ProtocolVersion,
		e.DecisionVersion,
		e.Seq,
		e.Timestamp,
		e.ActorID,
		e.EventType,
		e.PayloadHash,
		e.PrevHash,
	)

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func Append(actorID, eventType, payloadHash string) (*LogEntry, error) {

	var prevHash string
	var seq uint64 = 0

	if data, err := os.ReadFile(LedgerFile); err == nil {

		lines := bytesSplitLines(data)

		if len(lines) > 0 {

			var lastEntry LogEntry
			json.Unmarshal(lines[len(lines)-1], &lastEntry)

			prevHash = lastEntry.EntryHash
			seq = lastEntry.Seq + 1
		}
	}

	entry := LogEntry{
		ProtocolVersion: 1,
		DecisionVersion: 1,
		Seq:             seq,
		Timestamp:       time.Now().Unix(),
		ActorID:         actorID,
		EventType:       eventType,
		PayloadHash:     payloadHash,
		PrevHash:        prevHash,
	}

	entry.EntryHash = computeEntryHash(entry)

	file, err := os.OpenFile(LedgerFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	jsonEntry, _ := json.Marshal(entry)

	file.Write(jsonEntry)
	file.Write([]byte("\n"))

	return &entry, nil
}

func bytesSplitLines(data []byte) [][]byte {

	return bytes.Split(bytes.TrimSpace(data), []byte("\n"))
}
func Verify() error {

	data, err := os.ReadFile(LedgerFile)
	if err != nil {
		return err
	}

	lines := bytesSplitLines(data)

	var prevEntry *LogEntry

	for i, line := range lines {

		var entry LogEntry
		err := json.Unmarshal(line, &entry)
		if err != nil {
			return fmt.Errorf("invalid JSON at line %d", i)
		}

		// 1️⃣ sequence continuity
		if uint64(i) != entry.Seq {
			return fmt.Errorf("sequence mismatch at %d", i)
		}

		// 2️⃣ recompute hash
		expectedHash := computeEntryHash(entry)

		if expectedHash != entry.EntryHash {
			return fmt.Errorf("entry hash mismatch at %d", i)
		}

		// 3️⃣ check chain linkage
		if i > 0 {
			if entry.PrevHash != prevEntry.EntryHash {
				return fmt.Errorf("prev_hash mismatch at %d", i)
			}
		}

		prevEntry = &entry
	}

	return nil
}
