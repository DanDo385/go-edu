# 24-monitor: Chain Monitoring

## Overview

Build a comprehensive chain monitoring system. Track new blocks, transactions, events, and anomalies in real-time.

## Learning Objectives

- Implement real-time block monitoring
- Track transaction patterns
- Detect anomalies
- Set up alerting
- Build monitoring dashboards

## Project Structure

```
24-monitor/
├── cmd/
│   ├── app/main.go
│   └── dev/main.go
├── internal/monitor/
│   ├── exercise.go
│   ├── exercise_test.go
│   ├── solution.reference.go
│   └── solution_no_err.reference.go
└── README.md
```

## CLI Arguments

```bash
# Start monitoring
go run ./cmd/app/main.go <RPC_URL>

# Monitor specific contracts
go run ./cmd/app/main.go <RPC_URL> --contracts addresses.txt

# With alerting
go run ./cmd/app/main.go <RPC_URL> --alert-on-large-transfer 1000
```

## What the Dev Harness Demonstrates

1. **Block Streaming** - Real-time block monitoring
2. **Event Tracking** - Monitor contract events
3. **Transaction Analysis** - Pattern detection
4. **Anomaly Detection** - Unusual activity alerts
5. **Performance Metrics** - Gas usage trends

## Next Steps

Proceed to **geth/25-toolbox** for the utility collection.
