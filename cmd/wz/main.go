package main

import (
	"fmt"

	"wazulu/execution/cas"
	"wazulu/execution/log"
)

func main() {

	fmt.Println("Wazulu Execution Engine v1")

	data := []byte("hello wazulu")

	hash, err := cas.Store(data)
	if err != nil {
		panic(err)
	}

	fmt.Println("CAS stored:", hash)

	entry, err := log.Append(
		"node-01",
		"execution_intent",
		hash,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Log entry appended:", entry.Seq)

	err = log.Verify()
	if err != nil {
		panic(err)
	}

	fmt.Println("Ledger verification OK")
}
