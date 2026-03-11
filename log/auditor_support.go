package log

import (
	"encoding/json"
	"fmt"
	"os"
)

type LedgerEntry = Entry

func LoadLedgerFrom(path string) ([]LedgerEntry, error) {

	var entries []LedgerEntry

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return entries, nil
		}
		return nil, fmt.Errorf("open ledger: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for {
		var e LedgerEntry
		if err := decoder.Decode(&e); err != nil {
			break
		}
		entries = append(entries, e)
	}

	return entries, nil
}
