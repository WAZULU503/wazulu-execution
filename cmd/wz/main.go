package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"wazulu/execution/cas"
	"wazulu/execution/log"
)

func main() {

	fmt.Println("Wazulu Execution Engine v1")

	// ---- store payload in CAS ----
	payload := []byte("hello wazulu")

	hash, err := cas.Store(payload)
	if err != nil {
		panic(err)
	}

	fmt.Println("CAS stored:", hash)

	// ---- create transparency log ----
	tlog := log.NewTransparencyLog("ledger.jsonl")

	entry, err := tlog.Append(
		"node-01",
		"execution_intent",
		hash,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Log entry appended:", entry.Seq)

	// ---- verify ledger ----
	err = tlog.Verify()
	if err != nil {
		panic(err)
	}

	fmt.Println("Ledger verification OK")

	// ---- load ledger entries for Merkle root ----
	file, err := os.Open("ledger.jsonl")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	var entryHashes []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		var e log.LogEntry

		err := json.Unmarshal(scanner.Bytes(), &e)
		if err != nil {
			panic(err)
		}

		entryHashes = append(entryHashes, e.EntryHash)
	}

	// ---- compute Merkle root ----
	root := log.ComputeMerkleRoot(entryHashes)

	fmt.Println("Merkle Root:", root)
}
