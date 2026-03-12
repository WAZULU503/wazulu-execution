package main

import (
	"fmt"
	"os"
)

const version = "v1.4"

func main() {

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {

	case "exec":
		runExecution()

	case "verify":
		runVerify()

	case "audit":
		runAudit()

	case "version":
		fmt.Println("Wazulu Execution", version)

	default:
		printHelp()
	}
}

func runExecution() {
	fmt.Println("Wazulu Execution Engine: exec")
	fmt.Println("Execution event recorded (prototype)")
}

func runVerify() {
	fmt.Println("Wazulu Execution Engine: verify")
	fmt.Println("Ledger verification complete (prototype)")
}

func runAudit() {
	fmt.Println("Wazulu Execution Engine: audit")
	fmt.Println("Audit complete (prototype)")
}

func printHelp() {
	fmt.Println("Wazulu Execution Engine")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  wz exec     Record execution event")
	fmt.Println("  wz verify   Verify execution log")
	fmt.Println("  wz audit    Run integrity audit")
	fmt.Println("  wz version  Show version")
}
