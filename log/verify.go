package log

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// VerifyLedger scans the ledger and ensures the hash chain is intact
func VerifyLedger() error {

	file, err := os.Open(LedgerFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var prev Entry
	first := true

	for scanner.Scan() {

		var entry Entry
		err := json.Unmarshal(scanner.Bytes(), &entry)
		if err != nil {
			return err
		}

		if !first {
			if entry.PrevHash != prev.EntryHash {
				return fmt.Errorf("ledger integrity violation at seq %d", entry.Seq)
			}
		}

		prev = entry
		first = false
	}

	fmt.Println("Ledger verification OK")
	return nil
}
