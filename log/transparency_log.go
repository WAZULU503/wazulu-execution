package log

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type LogEntry struct {
	ProtocolVersion int    `json:"protocol_version"`
	DecisionVersion int    `json:"decision_version"`
	Seq             int    `json:"seq"`
	Timestamp       int64  `json:"timestamp"`
	ActorID         string `json:"actor_id"`
	EventType       string `json:"event_type"`
	PayloadHash     string `json:"payload_hash"`
	PrevHash        string `json:"prev_hash"`
	EntryHash       string `json:"entry_hash"`
}

type TransparencyLog struct {
	Path string
}

func NewTransparencyLog(path string) *TransparencyLog {
	return &TransparencyLog{Path: path}
}

func (l *TransparencyLog) loadEntries() ([]LogEntry, error) {

	file, err := os.Open(l.Path)
	if err != nil {

		if os.IsNotExist(err) {
			return []LogEntry{}, nil
		}

		return nil, err
	}

	defer file.Close()

	var entries []LogEntry

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		var e LogEntry

		err := json.Unmarshal(scanner.Bytes(), &e)
		if err != nil {
			return nil, err
		}

		entries = append(entries, e)
	}

	return entries, nil
}

func computeEntryHash(e LogEntry) string {

	data := fmt.Sprintf(
		"WZLOGv1|%d|%d|%d|%s|%s|%s",
		e.Seq,
		e.Timestamp,
		e.ProtocolVersion,
		e.EventType,
		e.PayloadHash,
		e.PrevHash,
	)

	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}

func (l *TransparencyLog) Append(actorID string, eventType string, payloadHash string) (*LogEntry, error) {

	entries, err := l.loadEntries()
	if err != nil {
		return nil, err
	}

	seq := len(entries)

	prevHash := ""
	if seq > 0 {
		prevHash = entries[seq-1].EntryHash
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

	file, err := os.OpenFile(l.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	jsonEntry, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}

	// write entry
	_, err = file.Write(jsonEntry)
	if err != nil {
		return nil, err
	}

	// newline delimiter
	_, err = file.Write([]byte("\n"))
	if err != nil {
		return nil, err
	}

	// force flush to disk (crash-safe write)
	err = file.Sync()
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (l *TransparencyLog) Verify() error {

	entries, err := l.loadEntries()
	if err != nil {
		return err
	}

	for i := 0; i < len(entries); i++ {

		e := entries[i]

		expectedHash := computeEntryHash(e)

		if expectedHash != e.EntryHash {
			return fmt.Errorf("entry hash mismatch at seq %d", e.Seq)
		}

		if i > 0 {

			if e.PrevHash != entries[i-1].EntryHash {
				return fmt.Errorf("prev hash mismatch at seq %d", e.Seq)
			}
		}
	}

	return nil
}
