package log

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

const ledgerFile = "ledger.jsonl"

type Entry struct {
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

func computeHash(data string) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

func LoadLedger() ([]Entry, error) {

	file, err := os.OpenFile(ledgerFile, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []Entry

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		var e Entry

		err := json.Unmarshal(scanner.Bytes(), &e)
		if err != nil {
			return nil, err
		}

		entries = append(entries, e)
	}

	return entries, nil
}

func AppendEntry(actor string, event string, payload string) (Entry, error) {

	entries, err := LoadLedger()
	if err != nil {
		return Entry{}, err
	}

	seq := len(entries)

	prev := ""
	if seq > 0 {
		prev = entries[seq-1].EntryHash
	}

	ts := time.Now().Unix()

	data := strconv.Itoa(seq) + actor + event + payload + prev + strconv.FormatInt(ts, 10)

	hash := computeHash(data)

	entry := Entry{
		ProtocolVersion: 1,
		DecisionVersion: 1,
		Seq:             seq,
		Timestamp:       ts,
		ActorID:         actor,
		EventType:       event,
		PayloadHash:     payload,
		PrevHash:        prev,
		EntryHash:       hash,
	}

	file, err := os.OpenFile(ledgerFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return Entry{}, err
	}
	defer file.Close()

	b, _ := json.Marshal(entry)

	file.Write(b)
	file.Write([]byte("\n"))

	return entry, nil
}

func VerifyLedger() error {

	entries, err := LoadLedger()
	if err != nil {
		return err
	}

	for i := range entries {

		if i == 0 {
			continue
		}

		if entries[i].PrevHash != entries[i-1].EntryHash {
			return fmt.Errorf("ledger chain broken at seq %d", i)
		}
	}

	return nil
}
