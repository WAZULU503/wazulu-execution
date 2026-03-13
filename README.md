# Wazulu Execution
Verifiable Execution Engine

Project Status: Research Prototype
Language: Go
License: MIT
Version: v1.4

Wazulu Execution is a command-line tool that records system execution events in an append-only log and produces cryptographic commitments that allow the history of those events to be independently verified.

The system converts runtime activity into cryptographically verifiable execution records.

---

## Overview

Traditional system logs rely on trust and can be modified or deleted.

Wazulu Execution records events in a hash-chained transparency log, commits the log state using a Merkle tree, and produces signed checkpoints that allow independent verification of the system history.

Execution pipeline:

Execution Event
│
▼
Payload Hashing
│
▼
Append-only Log Entry
│
▼
Merkle Tree Commitment
│
▼
Signed Tree Head
│
▼
Verification / Audit

The goal is to make system behavior cryptographically auditable.

---

## Design Goals

Wazulu Execution focuses on a small set of core properties.

Deterministic verification  
Anyone can recompute hashes and verify log integrity independently.

Tamper-evident history  
Log entries are hash-linked, making modification detectable.

Transparent state commitment  
Merkle roots provide cryptographic commitments to log state.

Independent auditing  
External auditors can verify execution history without trusting the runtime.

Minimal trusted components  
The system minimizes trust assumptions beyond the runtime and cryptographic primitives.

---

## Architecture

+----------------------+
|  Execution Runtime   |
+----------------------+
           |
           v
+----------------------+
|   Payload Hashing    |
+----------------------+
           |
           v
+----------------------+
| CAS Evidence Store   |
| (content addressed)  |
+----------------------+
           |
           v
+----------------------+
| Transparency Log     |
| (append-only ledger) |
+----------------------+
           |
           v
+----------------------+
| Merkle Tree          |
| log commitment       |
+----------------------+
           |
           v
+----------------------+
| Signed Tree Head     |
| checkpoint           |
+----------------------+
           |
           v
+----------------------+
| Witness / Auditor    |
| verification         |
+----------------------+

Each layer contributes to the integrity of the execution record.

---

## Example Log Entry

{
  "seq": 120,
  "timestamp": 1710001234,
  "event_type": "execution_intent",
  "payload_hash": "8f4e1c...",
  "prev_hash": "d1a8e1...",
  "entry_hash": "4bc2a9..."
}

Each entry references the previous entry through prev_hash, creating a verifiable chain of execution records.

---

## Security Guarantees

- Tamper-evident history
- Commitment to log state
- Checkpoint integrity
- Independent verification

---

## Threat Model and Limitations

Assumptions

- execution runtime records events correctly
- cryptographic hashes remain secure
- signing keys are managed securely
- host environment is trusted

Limitations

- no distributed consensus
- no Byzantine fault tolerance
- no hardware attestation
- compromised runtimes can omit events

The system focuses on tamper-evident history, not guaranteeing that every event is recorded.

---

## Repository Structure

wazulu-execution
│
├─ cmd/
│   └─ wz/
│       └─ main.go
│
├─ internal/
│   ├─ execution/
│   │   └─ runtime.go
│   │
│   ├─ ledger/
│   │   ├─ ledger.go
│   │   └─ verify.go
│   │
│   ├─ merkle/
│   │   └─ merkle.go
│   │
│   └─ cas/
│       └─ store.go
│
├─ demo.sh
├─ go.mod
├─ README.md
└─ LICENSE

internal/ follows standard Go conventions — packages inside it are private to the module and cannot be imported by external programs.

---

## Quick Start

Clone repository

git clone https://github.com/WAZULU503/wazulu-execution
cd wazulu-execution

Build

go build -o wz ./cmd/wz

Run

./wz exec "example event"

Verify

./wz verify

Audit

./wz audit

---

## Author

Wazulu  
https://github.com/WAZULU503

---

## License

MIT — see LICENSE
