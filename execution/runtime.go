package execution

import (
	"fmt"

	"wazulu/execution/cas"
	"wazulu/execution/log"
)

func RunExecution() error {

	payload := []byte("hello wazulu")

	hash, err := cas.Store(payload)
	if err != nil {
		return err
	}

	fmt.Println("CAS stored:", hash)

	entry, err := log.AppendEntry("node-01", "execution_intent", hash)
	if err != nil {
		return err
	}

	fmt.Println("Log entry appended:", entry.Seq)

	err = log.VerifyLedger()
	if err != nil {
		return err
	}

	fmt.Println("Ledger verification OK")

	ledger, err := log.LoadLedger()
	if err != nil {
		return err
	}

	var hashes []string
	for _, e := range ledger {
		hashes = append(hashes, e.EntryHash)
	}

	root := log.ComputeMerkleRoot(hashes)

	fmt.Println("Merkle Root:", root)

	sth := log.GenerateSTH(root, len(hashes))

	fmt.Println()
	fmt.Println("Signed Tree Head")
	fmt.Println("Tree Size :", sth.TreeSize)
	fmt.Println("Timestamp :", sth.Timestamp)
	fmt.Println("Signature :", sth.Signature)

	fmt.Println()
	fmt.Println("Witness step skipped (no witness server running)")

	return nil
}
