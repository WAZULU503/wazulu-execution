package main

import (
	"fmt"
	"os"

	"wazulu/execution/execution"
	"wazulu/execution/log"
)

func main() {

	if len(os.Args) < 2 {
		runExecution()
		return
	}

	switch os.Args[1] {

	case "exec":
		runExecution()

	case "verify":
		runVerify()

	default:
		printHelp()
	}
}

func runExecution() {

	fmt.Println("Wazulu Execution Engine v1")

	err := execution.RunExecution()

	if err != nil {
		fmt.Println()
		fmt.Println("ENGINE ERROR:")
		fmt.Println(err)
		os.Exit(1)
	}
}

func runVerify() {

	fmt.Println("Verifying ledger integrity...")

	err := log.VerifyLedger()

	if err != nil {
		fmt.Println("Verification failed:")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Ledger verification OK")
}

func printHelp() {
	fmt.Println("Wazulu Execution Engine")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  wz exec     run execution pipeline")
	fmt.Println("  wz verify   verify ledger integrity")
}
