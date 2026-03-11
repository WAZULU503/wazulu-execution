package main

import (
	"fmt"
	"os"

	"wazulu/execution/execution"
)

func main() {

	fmt.Println("Wazulu Execution Engine v1")

	err := execution.RunExecution()

	if err != nil {
		fmt.Println()
		fmt.Println("ENGINE ERROR:")
		fmt.Println(err)
		os.Exit(1)
	}
}
