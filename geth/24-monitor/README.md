# geth-24-monitor: Node Health Monitoring

**Goal:** implement basic node health checks (block freshness/lag) and discuss alerting patterns.

## Big Picture

Monitoring a node involves tracking head freshness, RPC latency, and error rates. Simple checks catch stale nodes early; production systems export metrics to Prometheus/Grafana and alert on thresholds.

**Computer Science principle:** Threshold-based health checks convert continuous metrics (lag time) into discrete states (OK/STALE) for actionable alerting.

## Learning Objectives

1. **Fetch latest header** and extract timestamp
2. **Calculate lag** between block time and current time
3. **Classify status** using configurable thresholds
4. **Understand monitoring patterns** for production systems

## Prerequisites
- **Module 01:** HeaderByNumber pattern
- **Module 21:** Health check concepts (sync status)
- **Comfort with:** Time calculations, alerting principles

## Files

- **Exercise:** `exercise/exercise.go` - TODOs guide implementation
- **Solution:** `exercise/solution.go` - Full implementation with educational comments
- **Types:** `exercise/types.go` - Interface and struct definitions
- **Tests:** `exercise/exercise_test.go` - Test suite

## How to Run Tests

```bash
cd geth/24-monitor
go test ./exercise/
go test -v ./exercise/
go test -tags solution -v ./exercise/
```

## Key Concepts

### Block Lag Calculation

```go
currentTime := time.Now()
blockTime := time.Unix(header.Time, 0)
lagSeconds := currentTime.Sub(blockTime).Seconds()
```

### Status Classification

- **Lag < 60s:** OK (node is synced, receiving new blocks)
- **Lag >= 60s:** STALE (node behind, may have sync issues)
- **Negative lag:** OK (slight clock skew acceptable)

### Production Monitoring Patterns

1. **Prometheus Metrics:**
   - Gauge: `eth_block_lag_seconds`
   - Counter: `eth_stale_node_count`
   - Histogram: `eth_rpc_latency_seconds`

2. **Alerting Rules:**
   - Fire if lag > 60s for 5 consecutive minutes
   - Reduce false positives from transient issues

3. **Multi-Node Monitoring:**
   - Check multiple RPC endpoints
   - Alert only if all nodes are stale

## Common Pitfalls

### Pitfall 1: Single Point-in-Time Check
- One check can give false positives (transient network blip)
- Solution: Require N consecutive failures before alerting

### Pitfall 2: Not Tuning Thresholds by Network
- Mainnet: 60s = ~5 blocks behind
- Fast L2s: 60s = ~30 blocks behind
- Solution: Adjust MaxLagSeconds per network

### Pitfall 3: Ignoring Clock Skew
- Block timestamps can be slightly in future
- Solution: Accept negative lag (don't mark as STALE)

## How Concepts Build

1. **From Module 01:** HeaderByNumber pattern → Used for health monitoring
2. **From Module 21:** Health status → Extended to time-based checks
3. **New:** Time calculations, threshold-based classification, observability

## Fun Facts

- **Block time variance:** Ethereum targets 12s but can vary 1-30s
- **Timestamp rules:** Blocks can't be > 15s in future (rejected by network)
- **Validator timing:** Post-Merge, validators have strict 12s slot timing

## Next Steps

After completing this module, you'll move to **25-toolbox** where you'll:
- Build a Swiss Army CLI combining multiple operations
- Implement subcommands: status, block, tx, events
- Reuse patterns from all previous modules
