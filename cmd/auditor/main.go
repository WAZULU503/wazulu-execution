package main

import (
	"fmt"
	"os"

	"wazulu/execution/log"
)

func main() {

	fmt.Println("=== Wazulu Auditor Node ===")
	fmt.Println()

	ledger, err := log.LoadLedger()
	if err != nil {
		fmt.Println("AUDIT FAILURE: cannot load ledger")
		fmt.Println(err)
		os.Exit(1)
	}

	if len(ledger) == 0 {
		fmt.Println("AUDIT INFO: empty ledger")
		os.Exit(0)
	}

	var hashes []string

	for i := range ledger {

		hashes = append(hashes, ledger[i].EntryHash)

		if i > 0 && ledger[i].PrevHash != ledger[i-1].EntryHash {
			fmt.Println("AUDIT FAILURE: chain break at entry", i)
			os.Exit(1)
		}
	}

	root := log.ComputeMerkleRoot(hashes)

	fmt.Println("Recomputed Merkle Root:", root)

	fmt.Println()
	fmt.Println("=== AUDIT SUCCESS ===")
}
