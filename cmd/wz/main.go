package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"wazulu/execution/cas"
	"wazulu/execution/log"
)

type WitnessRequest struct {
	Root string `json:"root"`
	Size int    `json:"size"`
}

type WitnessResponse struct {
	WitnessID string `json:"witness_id"`
	Signature string `json:"signature"`
}

func main() {

	fmt.Println("Wazulu Execution Engine v1")

	payload := []byte("hello wazulu")

	hash, err := cas.Store(payload)
	if err != nil {
		panic(err)
	}

	fmt.Println("CAS stored:", hash)

	entry, err := log.AppendEntry(
		"node-01",
		"execution_intent",
		hash,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Log entry appended:", entry.Seq)

	err = log.VerifyLedger()
	if err != nil {
		panic(err)
	}

	fmt.Println("Ledger verification OK")

	ledger, err := log.LoadLedger()
	if err != nil {
		panic(err)
	}

	var hashes []string
	for _, e := range ledger {
		hashes = append(hashes, e.EntryHash)
	}

	root := log.ComputeMerkleRoot(hashes)

	fmt.Println("Merkle Root:", root)

	// USE CORRECT FUNCTION NAME
	sth := log.GenerateSTH(root, len(hashes))

	fmt.Println()
	fmt.Println("Signed Tree Head")
	fmt.Println("Tree Size :", sth.TreeSize)
	fmt.Println("Timestamp :", sth.Timestamp)
	fmt.Println("Signature :", sth.Signature)

	req := WitnessRequest{
		Root: root,
		Size: len(hashes),
	}

	body, _ := json.Marshal(req)

	resp, err := http.Post(
		"http://localhost:9001/sign",
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		fmt.Println()
		fmt.Println("Witness not reachable")
		os.Exit(0)
	}

	defer resp.Body.Close()

	var witnessResp WitnessResponse
	json.NewDecoder(resp.Body).Decode(&witnessResp)

	fmt.Println()
	fmt.Println("Witness Cosignature")
	fmt.Println("Witness ID :", witnessResp.WitnessID)
	fmt.Println("Signature  :", witnessResp.Signature)
}
