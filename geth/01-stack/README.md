# 01: Stack - Ethereum RPC Connectivity and Chain Information

## What Is This Project About?

This project teaches you how to establish RPC connectivity to an Ethereum node and retrieve fundamental chain information. You'll learn to query the Chain ID (used for replay protection), Network ID (legacy P2P identifier), and block headers (cryptographic commitments to state). This is the foundation of all Ethereum development—before you can do anything on-chain, you need to connect to a node and understand what network you're on.

Think of this as your "hello world" for Ethereum development. Just as you wouldn't write a web API without first learning to make HTTP requests, you can't build Ethereum applications without understanding how to dial an RPC endpoint and retrieve chain metadata.

## Why Is This Important?

Every Ethereum application—whether it's a wallet, DeFi protocol, NFT marketplace, or MEV bot—starts by connecting to an Ethereum node. Understanding this foundational step gives you:

- **Network verification**: Ensure you're on the right chain (mainnet vs testnet) before sending real transactions
- **Block awareness**: Know the current chain height and state root before querying or modifying data
- **Replay protection**: Chain ID prevents transactions from being replayed across different networks
- **Debugging foundation**: When things go wrong, you need to verify basic connectivity first

## Real-World Problems This Solves

- **Multi-chain applications**: Apps that support multiple networks (Ethereum, Polygon, Arbitrum) need to verify which chain they're connected to
- **Transaction safety**: Before sending a high-value transaction, verify you're on mainnet (Chain ID 1) not a testnet
- **Node monitoring**: DevOps teams monitor RPC endpoints by periodically checking chain height and connectivity
- **Development workflows**: Developers switch between local dev chains, testnets, and mainnet—Chain ID verification prevents costly mistakes

## Key Concepts You'll Learn

- **RPC connectivity**: Using `ethclient.DialContext` to connect to Ethereum nodes via HTTP/WebSocket
- **Context propagation**: Go's `context.Context` for timeouts and cancellation
- **Chain ID vs Network ID**: Understanding the difference and when each is used
- **Block headers**: Lightweight structures (~500 bytes) containing cryptographic commitments
- **Error handling**: Go's idiomatic error wrapping with `fmt.Errorf` and `%w`
- **Resource cleanup**: Using `defer` to ensure RPC connections are properly closed

## Prerequisites

- Basic Go knowledge (functions, error handling, interfaces)
- Understanding of Ethereum basics (blocks, transactions, addresses)
- Go 1.21+ installed
- Internet connection to access public RPC endpoints

## Project Structure

```
geth/01-stack/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application with RPC URL argument
│   └── dev/
│       └── main.go          # Debug harness with fixed test inputs
├── internal/
│   └── stack/
│       ├── exercise.go      # Your implementation goes here
│       ├── exercise_test.go # Tests for your implementation
│       ├── solution.reference.go      # Reference solution
│       ├── solution_no_err.reference.go # Simplified reference
│       └── types.go         # Type definitions
└── .vscode/
    └── launch.json          # Debug configurations
```

## How to Run

### Using cmd/app/main.go (CLI Arguments)

```bash
# Query latest block from public RPC
go run ./cmd/app/main.go https://eth.llamarpc.com

# Query specific block number
go run ./cmd/app/main.go https://eth.llamarpc.com 12345

# Query Sepolia testnet
go run ./cmd/app/main.go https://ethereum-sepolia-rpc.publicnode.com

# Query local Geth node
go run ./cmd/app/main.go http://localhost:8545
```

**Arguments:**
- `<RPC_URL>`: Ethereum RPC endpoint URL (required)
- `[block_number]`: Specific block number to query (optional, defaults to latest)

### Using cmd/dev/main.go (Debug Harness)

```bash
# Run with fixed test inputs
go run ./cmd/dev/main.go

# Or use VS Code debugger (F5) with "Debug: cmd/dev (Debug Harness)" configuration
```

**Recommended for learning**: The debug harness has fixed inputs so you can focus on stepping through the code logic without worrying about command-line arguments.

## How to Debug

1. **Set breakpoints** at `// BREAKPOINT:` comments in the code
2. **Press F5** in VS Code and select appropriate configuration:
   - `Debug: cmd/app` - Debug with CLI arguments
   - `Debug: cmd/dev (Debug Harness)` - Debug with fixed inputs (recommended)
   - `Test: Run All Tests` - Debug tests
3. **Step through code**:
   - F10 (Step Over) - Execute current line
   - F11 (Step Into) - Enter function calls
   - Shift+F11 (Step Out) - Return to caller
4. **Watch variables** in the Variables panel
5. **Inspect call stack** to understand function call hierarchy

**Key debugging points:**
- Watch how `ethclient.DialContext` establishes the RPC connection
- Inspect the `chainID` and `networkID` values (1 = mainnet, 11155111 = Sepolia)
- Examine the block header structure to see all the cryptographic commitments

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -v -run TestRun ./internal/stack

# Run with reference implementation
go test -tags=reference -v ./...
```

## Exercises

The main exercise is in `internal/stack/exercise.go`. You need to implement the `Run` function that:

1. Validates inputs (context and client)
2. Retrieves Chain ID using `client.ChainID(ctx)`
3. Retrieves Network ID using `client.NetworkID(ctx)`
4. Retrieves block header using `client.HeaderByNumber(ctx, cfg.BlockNumber)`
5. Returns a `Result` struct with all the information

**Tips:**
- Handle nil context by providing `context.Background()` as default
- Check for nil client and return an error early
- Wrap errors with context: `fmt.Errorf("chain id: %w", err)`
- Make defensive copies of big.Int values to prevent mutation
- Use nil for block number to get latest block

## What You'll Learn

### Technical Skills
- Connecting to Ethereum nodes via RPC
- Using the `go-ethereum/ethclient` package
- Working with `*big.Int` for large numbers
- Understanding Ethereum block headers
- Proper Go error handling and wrapping

### Ethereum Concepts
- **Chain ID**: EIP-155 replay protection (prevents cross-chain transaction replay)
- **Network ID**: Legacy P2P network identifier (predates Chain ID)
- **Block headers**: Contain state root, parent hash, timestamp, gas data, and more
- **RPC methods**: Behind the scenes, this uses `eth_chainId`, `net_version`, and `eth_getBlockByNumber`

### Design Patterns
- **Input validation**: Defensive programming for library functions
- **Context propagation**: Idiomatic Go for cancellation and timeouts
- **Error wrapping**: Maintaining error chains for better debugging
- **Interface-based design**: `RPCClient` interface allows for testing and mocking

## Common Issues

**"Error connecting to RPC endpoint"**
- Check your internet connection
- Try a different public RPC endpoint (RPC endpoints can be rate-limited or down)
- For local nodes, verify Geth is running: `geth attach http://localhost:8545`

**"Context deadline exceeded"**
- The RPC endpoint is too slow or unreachable
- Increase timeout in the code: `context.WithTimeout(context.Background(), 60*time.Second)`

**"Chain ID response was nil"**
- The RPC endpoint might not support `eth_chainId` (very old nodes)
- Verify the endpoint URL is correct

## Additional Resources

- [Go-Ethereum Documentation](https://geth.ethereum.org/docs)
- [Ethereum JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [EIP-155: Simple Replay Attack Protection](https://eips.ethereum.org/EIPS/eip-155)
- [Ethereum Block Structure](https://ethereum.org/en/developers/docs/blocks/)
- [Public RPC Endpoints](https://chainlist.org/)

## Next Steps

After completing this module, you'll be ready for:

- **geth/02-rpc-basics**: Understanding different RPC methods and retry logic
- **geth/03-keys-addresses**: Working with private keys and Ethereum addresses
- **geth/04-accounts-balances**: Querying account balances and nonces
- **geth/05-tx-nonces**: Understanding transaction nonces and sequencing

**Congratulations!** You've taken your first step into Ethereum development with Go. 🚀
