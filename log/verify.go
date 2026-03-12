package log

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

func VerifyLedger() error {

	file, err := os.ReadFile(LedgerFile)
	if err != nil {
		return err
	}

	lines := split(string(file))

	var prevHash string

	for i, line := range lines {

		var entry Entry

		err := json.Unmarshal([]byte(line), &entry)
		if err != nil {
			return err
		}

		if entry.PrevHash != prevHash {
			return fmt.Errorf("ledger broken at entry %d", i+1)
		}

		expected := hash(fmt.Sprintf("%d%s%s%s",
			entry.Seq,
			entry.EventType,
			entry.PayloadHash,
			entry.PrevHash,
		))

		if entry.EntryHash != expected {
			return fmt.Errorf("hash mismatch at entry %d", i+1)
		}

		prevHash = entry.EntryHash
	}

	fmt.Println("Ledger verification OK")

	return nil
}

func hash(data string) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}
