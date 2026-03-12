package main

import (
	"fmt"
	"os"

	wlog "github.com/WAZULU503/wazulu-execution/log"
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
		fmt.Println("Wazulu Execution Engine", version)

	default:
		printHelp()
	}
}

func runExecution() {

	err := wlog.AppendEvent("execution", "prototype payload")

	if err != nil {
		fmt.Println("Execution failed:", err)
		return
	}

	fmt.Println("Execution event recorded")
}

func runVerify() {

	err := wlog.VerifyLedger()

	if err != nil {
		fmt.Println("Verification failed:", err)
		return
	}

	fmt.Println("Ledger verification OK")
}

func runAudit() {

	err := wlog.VerifyLedger()

	if err != nil {
		fmt.Println("Audit failed:", err)
		return
	}

	fmt.Println("Ledger audit complete")
}

func printHelp() {
	fmt.Println("Wazulu Execution Engine")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  wz exec     Record execution event")
	fmt.Println("  wz verify   Verify ledger integrity")
	fmt.Println("  wz audit    Run integrity audit")
	fmt.Println("  wz version  Show version")
}
