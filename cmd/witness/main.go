package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

type SignRequest struct {
	Root string `json:"root"`
	Size int    `json:"size"`
}

type SignResponse struct {
	WitnessID string `json:"witness_id"`
	Signature string `json:"signature"`
}

var privateKey ed25519.PrivateKey
var publicKey ed25519.PublicKey

func init() {
	publicKey, privateKey, _ = ed25519.GenerateKey(rand.Reader)
}

func signHandler(w http.ResponseWriter, r *http.Request) {
	var req SignRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	message := fmt.Sprintf("%s:%d", req.Root, req.Size)
	sig := ed25519.Sign(privateKey, []byte(message))

	resp := SignResponse{
		WitnessID: hex.EncodeToString(publicKey[:8]),
		Signature: hex.EncodeToString(sig),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/sign", signHandler)

	fmt.Println("Witness node running on :9001")

	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		panic(err)
	}
}
