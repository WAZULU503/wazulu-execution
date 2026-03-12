package log

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Entry struct {
	Seq         int    `json:"seq"`
	Timestamp   int64  `json:"timestamp"`
	EventType   string `json:"event_type"`
	PayloadHash string `json:"payload_hash"`
	PrevHash    string `json:"prev_hash"`
	EntryHash   string `json:"entry_hash"`
}

const LedgerFile = "ledger.jsonl"

func AppendEvent(eventType string, payload string) error {

	payloadHash := hash(payload)

	prevHash := lastHash()

	entry := Entry{
		Seq:         nextSeq(),
		Timestamp:   time.Now().Unix(),
		EventType:   eventType,
		PayloadHash: payloadHash,
		PrevHash:    prevHash,
	}

	entry.EntryHash = hash(fmt.Sprintf("%d%s%s%s",
		entry.Seq,
		entry.EventType,
		entry.PayloadHash,
		entry.PrevHash,
	))

	file, err := os.OpenFile(LedgerFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, _ := json.Marshal(entry)

	_, err = file.Write(append(data, '\n'))

	return err
}

func hash(data string) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

func lastHash() string {

	file, err := os.ReadFile(LedgerFile)
	if err != nil {
		return ""
	}

	lines := string(file)
	if len(lines) == 0 {
		return ""
	}

	var last Entry
	for _, line := range split(lines) {
		json.Unmarshal([]byte(line), &last)
	}

	return last.EntryHash
}

func nextSeq() int {

	file, err := os.ReadFile(LedgerFile)
	if err != nil {
		return 1
	}

	count := len(split(string(file)))

	return count + 1
}

func split(data string) []string {

	var out []string

	current := ""

	for _, c := range data {
		if c == '\n' {
			if current != "" {
				out = append(out, current)
			}
			current = ""
		} else {
			current += string(c)
		}
	}

	if current != "" {
		out = append(out, current)
	}

	return out
}
