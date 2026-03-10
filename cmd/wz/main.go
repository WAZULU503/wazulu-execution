package main

import (
	"fmt"

	"wazulu/execution/cas"
)

func main() {

	fmt.Println("Wazulu Execution Engine v1")

	data := []byte("hello wazulu")

	hash, err := cas.Store(data)
	if err != nil {
		panic(err)
	}

	fmt.Println("Stored CAS object:", hash)

	obj, err := cas.Load(hash)
	if err != nil {
		panic(err)
	}

	fmt.Println("Loaded CAS object:", string(obj))
}
