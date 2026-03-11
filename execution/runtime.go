package execution

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"wazulu/execution/cas"
	"wazulu/execution/log"
)

type WitnessRequest struct {
	Root              string `json:"root"`
	Size              int    `json:"size"`
	Timestamp         int64  `json:"timestamp"`
	OperatorSignature string `json:"operator_signature"`
}

type WitnessResponse struct {
	WitnessID string `json:"witness_id"`
	Signature string `json:"signature"`
}

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

	req := WitnessRequest{
		Root:              root,
		Size:              len(hashes),
		Timestamp:         sth.Timestamp,
		OperatorSignature: sth.Signature,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		"http://localhost:9001/sign",
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		return fmt.Errorf("witness unavailable: %w", err)
	}

	defer resp.Body.Close()

	var witnessResp WitnessResponse

	err = json.NewDecoder(resp.Body).Decode(&witnessResp)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Witness Cosignature")
	fmt.Println("Witness ID :", witnessResp.WitnessID)
	fmt.Println("Signature  :", witnessResp.Signature)

	// Persist checkpoint
	err = log.AppendCheckpoint(
		sth.TreeSize,
		sth.RootHash,
		sth.Timestamp,
		sth.Signature,
		witnessResp.Signature,
	)

	if err != nil {
		fmt.Println("Checkpoint write failed:", err)
		return err
	}

	fmt.Println()
	fmt.Println("Checkpoint persisted to sth.log")

	return nil
}
