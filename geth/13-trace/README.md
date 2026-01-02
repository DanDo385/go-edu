# 13-trace: Transaction Tracing

## Overview

Deep dive into transaction execution using debug_traceTransaction. Understand EVM execution, gas costs, and internal calls.

## Learning Objectives

- Trace transaction execution
- Analyze EVM opcodes and gas costs
- Identify internal contract calls
- Debug failed transactions

## Project Structure

```
13-trace/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/trace/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Trace transaction
go run ./cmd/app/main.go <RPC_URL> <TX_HASH>

# Example
go run ./cmd/app/main.go https://eth.llamarpc.com 0x123...abc

# With specific tracer
go run ./cmd/app/main.go https://eth.llamarpc.com 0x123...abc --tracer callTracer
```

## What the Dev Harness Demonstrates

1. **Call Traces** - Internal contract calls
2. **Opcode Traces** - EVM execution details
3. **Gas Analysis** - Gas consumption breakdown
4. **Failed Transactions** - Why transactions reverted
5. **State Changes** - Storage modifications

## Next Steps

Proceed to **geth/14-explorer** for block explorer functionality.
