package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

/*
===================================================================
Ethereum Development Networks (Devnets) Explorer
===================================================================

This program demonstrates how to interact with different Ethereum networks
including development networks (devnets), testnets, and mainnet. Understanding
network differences is essential for developing and testing Ethereum applications
before deploying to mainnet.

USAGE:
    go run ./cmd/app/main.go <RPC_URL>
    go run ./cmd/app/main.go <RPC_URL> <demo>

Arguments:
    <RPC_URL>       - Ethereum RPC endpoint (required)
    [demo]          - Demo to run: all, identify, compare, local (default: "all")

Examples:
    # Ethereum Mainnet
    go run ./cmd/app/main.go https://eth.llamarpc.com

    # Sepolia Testnet
    go run ./cmd/app/main.go https://rpc.sepolia.org identify

    # Local development network (if running)
    go run ./cmd/app/main.go http://localhost:8545 local

    # Polygon Mainnet
    go run ./cmd/app/main.go https://polygon-rpc.com compare

Network Types:
    1. Mainnet: Production Ethereum network with real ETH
    2. Testnets: Public test networks (Sepolia, Holesky)
    3. Devnets: Private development networks for testing
    4. L2s: Layer 2 networks (Optimism, Arbitrum, Base, etc.)

Common RPC Endpoints:
    Ethereum Mainnet:
      - https://eth.llamarpc.com
      - https://rpc.ankr.com/eth

    Sepolia Testnet:
      - https://rpc.sepolia.org
      - https://sepolia.gateway.tenderly.co

    Holesky Testnet:
      - https://ethereum-holesky.publicnode.com

    Local Development:
      - http://localhost:8545 (Hardhat, Ganache, Geth devnet)

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: cmd/app")
- Step Into (F11) to enter network detection functions
- Watch Variables panel to see network properties
- Inspect chain IDs and network configurations

KEY CONCEPTS:
- Chain ID: Unique identifier for each Ethereum network (EIP-155)
- Network ID: Legacy identifier (usually same as Chain ID)
- Genesis Block: First block (block 0) - different per network
- Consensus: Mainnet uses PoS, testnets use PoS, devnets often use PoA
- Gas Prices: Mainnet expensive, testnets free, devnets configurable
- Block Time: ~12s on mainnet/testnets, configurable on devnets

COMPUTER SCIENCE CONCEPTS:
- Network Segmentation: Isolating different environments (dev, test, prod)
- Configuration Management: Different settings per environment
- Test-Driven Development: Testing on devnets before mainnet deployment
- CI/CD Pipeline: Automated testing across multiple networks
- Environment Parity: Keeping dev/test/prod environments similar

WHY USE DIFFERENT NETWORKS:
1. Mainnet: Production deployment, real value, real users
2. Testnets: Integration testing, public testing, audits
3. Devnets: Local development, fast iteration, debugging
4. Forks: Test against mainnet state without spending real ETH

NETWORK SELECTION STRATEGY:
- Initial development: Local devnet (Hardhat, Ganache)
- Contract testing: Local fork of mainnet
- Integration testing: Public testnet (Sepolia)
- Pre-production: Testnet with production-like setup
- Production: Mainnet deployment

REAL-WORLD WORKFLOW:
1. Develop smart contracts locally on Hardhat network
2. Write tests and run on local network
3. Fork mainnet to test against real contracts/state
4. Deploy to Sepolia testnet for integration testing
5. Conduct security audit on testnet
6. Deploy to mainnet after successful audit
*/

// NetworkInfo represents information about an Ethereum network
type NetworkInfo struct {
	ChainID      *big.Int
	NetworkID    *big.Int
	NetworkName  string
	NetworkType  string // mainnet, testnet, devnet, L2
	LatestBlock  *big.Int
	GenesisHash  string
	IsTestnet    bool
	GasPrice     *big.Int
	Description  string
}

// identifyNetwork determines which network we're connected to based on Chain ID
// BREAKPOINT: Set breakpoint here to see network identification logic
func identifyNetwork(chainID *big.Int) NetworkInfo {
	info := NetworkInfo{
		ChainID:   new(big.Int).Set(chainID),
		IsTestnet: false,
	}

	// BREAKPOINT: Set breakpoint here
	// DEBUG: Watch how different chain IDs map to different networks

	switch chainID.Int64() {
	case 1:
		info.NetworkName = "Ethereum Mainnet"
		info.NetworkType = "mainnet"
		info.Description = "Production Ethereum network with real ETH value"

	case 11155111:
		info.NetworkName = "Sepolia Testnet"
		info.NetworkType = "testnet"
		info.IsTestnet = true
		info.Description = "Ethereum testnet for application testing (free ETH from faucets)"

	case 17000:
		info.NetworkName = "Holesky Testnet"
		info.NetworkType = "testnet"
		info.IsTestnet = true
		info.Description = "Ethereum testnet for staking and infrastructure testing"

	case 5:
		info.NetworkName = "Goerli Testnet (deprecated)"
		info.NetworkType = "testnet"
		info.IsTestnet = true
		info.Description = "Deprecated Ethereum testnet (use Sepolia instead)"

	case 137:
		info.NetworkName = "Polygon Mainnet"
		info.NetworkType = "L2"
		info.Description = "Polygon PoS sidechain for scaling Ethereum"

	case 80001:
		info.NetworkName = "Polygon Mumbai Testnet (deprecated)"
		info.NetworkType = "testnet"
		info.IsTestnet = true
		info.Description = "Deprecated Polygon testnet"

	case 10:
		info.NetworkName = "Optimism Mainnet"
		info.NetworkType = "L2"
		info.Description = "Optimism Layer 2 rollup for Ethereum scaling"

	case 42161:
		info.NetworkName = "Arbitrum One"
		info.NetworkType = "L2"
		info.Description = "Arbitrum Layer 2 rollup for Ethereum scaling"

	case 8453:
		info.NetworkName = "Base Mainnet"
		info.NetworkType = "L2"
		info.Description = "Base Layer 2 rollup by Coinbase"

	case 1337:
		info.NetworkName = "Local Devnet (Ganache/Hardhat default)"
		info.NetworkType = "devnet"
		info.IsTestnet = true
		info.Description = "Local development network for testing"

	case 31337:
		info.NetworkName = "Local Devnet (Hardhat)"
		info.NetworkType = "devnet"
		info.IsTestnet = true
		info.Description = "Hardhat Network for local development"

	default:
		info.NetworkName = "Unknown Network"
		info.NetworkType = "unknown"
		info.Description = "Unknown or custom network"
	}

	return info
}

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect command-line arguments
	// DEBUG: os.Args[0] is the program name
	// DEBUG: os.Args[1] is the RPC URL (required)
	// DEBUG: os.Args[2] is the demo type (optional)

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <RPC_URL> [demo]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s https://rpc.sepolia.org identify\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s http://localhost:8545 local\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDemo options: all, identify, compare, local\n")
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	demo := "all"
	if len(os.Args) > 2 {
		demo = os.Args[2]
	}

	// BREAKPOINT: Set breakpoint here to verify parsed arguments
	// DEBUG: Check 'rpcURL' and 'demo' variables in Variables panel
	fmt.Printf("=== Ethereum Development Networks Explorer ===\n")
	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Printf("Demo: %s\n\n", demo)

	// ============================================================================
	// STEP 2: Connect to Ethereum RPC Endpoint
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before connection

	ctx := context.Background()
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// BREAKPOINT: Set breakpoint here
	client, err := ethclient.Dial(dialCtx, rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to RPC endpoint: %v\n", err)
	}
	defer client.Close()

	fmt.Println("✓ Connected to RPC endpoint")
	fmt.Println()

	// ============================================================================
	// STEP 3: Example 1 - Identify Network
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through network identification
	// DEBUG: This demonstrates how to detect which network you're connected to

	if demo == "all" || demo == "identify" {
		fmt.Println("=== Example 1: Network Identification ===")
		fmt.Println("Detecting network based on Chain ID...")
		fmt.Println()

		// Fetch Chain ID
		// BREAKPOINT: Set breakpoint here
		chainCtx, chainCancel := context.WithTimeout(ctx, 5*time.Second)
		chainID, err := client.ChainID(chainCtx)
		chainCancel()
		if err != nil {
			log.Fatalf("Failed to get Chain ID: %v\n", err)
		}

		// Fetch Network ID
		networkCtx, networkCancel := context.WithTimeout(ctx, 5*time.Second)
		networkID, err := client.NetworkID(networkCtx)
		networkCancel()
		if err != nil {
			log.Fatalf("Failed to get Network ID: %v\n", err)
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Inspect chainID and networkID values
		// DEBUG: Step Into (F11) identifyNetwork to see mapping logic
		info := identifyNetwork(chainID)
		info.NetworkID = networkID

		// Fetch additional network information
		// BREAKPOINT: Set breakpoint here
		headerCtx, headerCancel := context.WithTimeout(ctx, 5*time.Second)
		latestHeader, err := client.HeaderByNumber(headerCtx, nil)
		headerCancel()
		if err != nil {
			log.Fatalf("Failed to get latest block: %v\n", err)
		}

		info.LatestBlock = latestHeader.Number

		// Fetch genesis block
		genesisCtx, genesisCancel := context.WithTimeout(ctx, 5*time.Second)
		genesisHeader, err := client.HeaderByNumber(genesisCtx, big.NewInt(0))
		genesisCancel()
		if err != nil {
			log.Printf("Warning: Could not fetch genesis block: %v\n", err)
		} else {
			info.GenesisHash = genesisHeader.Hash().Hex()
		}

		// Fetch gas price
		gasPriceCtx, gasPriceCancel := context.WithTimeout(ctx, 5*time.Second)
		gasPrice, err := client.SuggestGasPrice(gasPriceCtx)
		gasPriceCancel()
		if err != nil {
			log.Printf("Warning: Could not fetch gas price: %v\n", err)
			info.GasPrice = big.NewInt(0)
		} else {
			info.GasPrice = gasPrice
		}

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Inspect complete NetworkInfo structure

		// Display network information
		fmt.Println("Network Information:")
		fmt.Printf("  Name:        %s\n", info.NetworkName)
		fmt.Printf("  Type:        %s\n", info.NetworkType)
		fmt.Printf("  Chain ID:    %s\n", info.ChainID.String())
		fmt.Printf("  Network ID:  %s\n", info.NetworkID.String())
		fmt.Printf("  Latest Block: #%s\n", info.LatestBlock.String())
		if info.GenesisHash != "" {
			fmt.Printf("  Genesis Hash: %s\n", info.GenesisHash)
		}

		// Convert gas price to gwei for readability
		gweiPrice := new(big.Int).Div(info.GasPrice, big.NewInt(1e9))
		fmt.Printf("  Gas Price:   %s gwei\n", gweiPrice.String())
		fmt.Printf("  Is Testnet:  %v\n", info.IsTestnet)
		fmt.Println()
		fmt.Printf("Description: %s\n", info.Description)
		fmt.Println()

		// Educational explanation
		fmt.Println("Network Identification Insights:")
		fmt.Println("  1. Chain ID uniquely identifies each network (EIP-155)")
		fmt.Println("  2. Network ID is legacy identifier, usually same as Chain ID")
		fmt.Println("  3. Genesis hash verifies you're on correct network fork")
		fmt.Println("  4. Always check Chain ID before signing transactions")
		fmt.Println("  5. Wrong network = wrong Chain ID = invalid transactions")
		fmt.Println()
	}

	// ============================================================================
	// STEP 4: Example 2 - Compare Network Characteristics
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through network comparison
	// DEBUG: This demonstrates differences between network types

	if demo == "all" || demo == "compare" {
		fmt.Println("=== Example 2: Network Characteristics Comparison ===")
		fmt.Println("Comparing characteristics of different network types...")
		fmt.Println()

		// Fetch current network info
		chainCtx, chainCancel := context.WithTimeout(ctx, 5*time.Second)
		chainID, err := client.ChainID(chainCtx)
		chainCancel()
		if err != nil {
			log.Fatalf("Failed to get Chain ID: %v\n", err)
		}

		info := identifyNetwork(chainID)

		// BREAKPOINT: Set breakpoint here
		// DEBUG: Watch how different networks have different characteristics

		fmt.Printf("Connected Network: %s (Chain ID: %s)\n\n", info.NetworkName, info.ChainID.String())

		fmt.Println("Network Type Characteristics:")
		fmt.Println()

		fmt.Println("1. MAINNET (Chain ID: 1)")
		fmt.Println("   Purpose: Production deployment")
		fmt.Println("   ETH Value: Real monetary value")
		fmt.Println("   Gas Costs: Expensive (varies with demand)")
		fmt.Println("   Block Time: ~12 seconds")
		fmt.Println("   Consensus: Proof of Stake (post-Merge)")
		fmt.Println("   Finality: ~12-15 minutes (2 epochs)")
		fmt.Println("   Use Cases: Production dapps, real transactions")
		fmt.Println("   Risks: Bugs = real money lost")
		fmt.Println()

		fmt.Println("2. TESTNETS (Sepolia: 11155111, Holesky: 17000)")
		fmt.Println("   Purpose: Public testing and integration")
		fmt.Println("   ETH Value: No monetary value (free from faucets)")
		fmt.Println("   Gas Costs: Free")
		fmt.Println("   Block Time: ~12 seconds (mimics mainnet)")
		fmt.Println("   Consensus: Proof of Stake")
		fmt.Println("   Finality: Similar to mainnet")
		fmt.Println("   Use Cases: Testing, audits, user acceptance testing")
		fmt.Println("   Benefits: Safe testing environment, public visibility")
		fmt.Println()

		fmt.Println("3. DEVNETS (Chain ID: 1337, 31337, custom)")
		fmt.Println("   Purpose: Local development and debugging")
		fmt.Println("   ETH Value: No value (unlimited supply)")
		fmt.Println("   Gas Costs: Usually disabled or very cheap")
		fmt.Println("   Block Time: Configurable (instant to 12s)")
		fmt.Println("   Consensus: Usually Proof of Authority or instant")
		fmt.Println("   Finality: Immediate (no risk of reorgs)")
		fmt.Println("   Use Cases: Development, unit testing, debugging")
		fmt.Println("   Benefits: Fast iteration, full control, no network dependency")
		fmt.Println()

		fmt.Println("4. LAYER 2 NETWORKS (Optimism: 10, Arbitrum: 42161, Base: 8453)")
		fmt.Println("   Purpose: Scaling Ethereum with lower fees")
		fmt.Println("   ETH Value: Real value (bridged from mainnet)")
		fmt.Println("   Gas Costs: Much cheaper than mainnet")
		fmt.Println("   Block Time: Faster than mainnet (~1-2 seconds)")
		fmt.Println("   Consensus: Sequencer + fraud proofs or validity proofs")
		fmt.Println("   Finality: Delayed due to challenge period (7 days for Optimism)")
		fmt.Println("   Use Cases: High-throughput dapps, gaming, social")
		fmt.Println("   Benefits: Lower fees, faster transactions, Ethereum security")
		fmt.Println()

		// Show current network type
		var networkTypeDescription string
		switch info.NetworkType {
		case "mainnet":
			networkTypeDescription = "You're on MAINNET - use caution, real money at stake!"
		case "testnet":
			networkTypeDescription = "You're on a TESTNET - safe for testing, free ETH available"
		case "devnet":
			networkTypeDescription = "You're on a DEVNET - local development environment"
		case "L2":
			networkTypeDescription = "You're on a LAYER 2 - lower fees than mainnet"
		default:
			networkTypeDescription = "Unknown network type - verify before proceeding"
		}

		fmt.Printf("Current Network Type: %s\n", info.NetworkType)
		fmt.Printf("⚠ %s\n", networkTypeDescription)
		fmt.Println()
	}

	// ============================================================================
	// STEP 5: Example 3 - Local Development Network Features
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to step through devnet features
	// DEBUG: This demonstrates special features of development networks

	if demo == "all" || demo == "local" {
		fmt.Println("=== Example 3: Local Development Network Features ===")
		fmt.Println("Exploring development network capabilities...")
		fmt.Println()

		// Fetch network info
		chainCtx, chainCancel := context.WithTimeout(ctx, 5*time.Second)
		chainID, err := client.ChainID(chainCtx)
		chainCancel()
		if err != nil {
			log.Fatalf("Failed to get Chain ID: %v\n", err)
		}

		info := identifyNetwork(chainID)

		if info.NetworkType == "devnet" {
			fmt.Printf("✓ Connected to local development network: %s\n\n", info.NetworkName)

			// Check balance of default accounts
			// Hardhat/Ganache provide pre-funded accounts
			// BREAKPOINT: Set breakpoint here
			fmt.Println("Checking pre-funded accounts:")

			// Common default account (Hardhat's first default account)
			defaultAccount := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")

			balanceCtx, balanceCancel := context.WithTimeout(ctx, 5*time.Second)
			balance, err := client.BalanceAt(balanceCtx, defaultAccount, nil)
			balanceCancel()

			if err != nil {
				fmt.Printf("  Could not fetch balance: %v\n", err)
			} else {
				// Convert to ETH
				ethBalance := new(big.Float).Quo(
					new(big.Float).SetInt(balance),
					big.NewFloat(1e18),
				)
				fmt.Printf("  Account: %s\n", defaultAccount.Hex())
				fmt.Printf("  Balance: %s ETH\n", ethBalance.Text('f', 4))
			}
			fmt.Println()

		} else {
			fmt.Printf("Not connected to a local devnet (current: %s)\n\n", info.NetworkName)
		}

		// Educational explanation about devnet features
		fmt.Println("Local Development Network Features:")
		fmt.Println()

		fmt.Println("1. Pre-funded Accounts:")
		fmt.Println("   - Hardhat: 20 accounts with 10,000 ETH each")
		fmt.Println("   - Ganache: 10 accounts with 100 ETH each")
		fmt.Println("   - Known private keys for easy testing")
		fmt.Println()

		fmt.Println("2. Instant Mining:")
		fmt.Println("   - Transactions mined immediately (no waiting)")
		fmt.Println("   - Configurable: instant, interval, or on-demand")
		fmt.Println("   - Perfect for rapid development iteration")
		fmt.Println()

		fmt.Println("3. Forking Capability:")
		fmt.Println("   - Fork mainnet state at any block number")
		fmt.Println("   - Test against real contracts and data")
		fmt.Println("   - Simulate mainnet without spending real ETH")
		fmt.Println("   - Example: hardhat fork from mainnet")
		fmt.Println()

		fmt.Println("4. Advanced Debugging:")
		fmt.Println("   - console.log() in Solidity (Hardhat)")
		fmt.Println("   - Transaction traces and call stacks")
		fmt.Println("   - Revert reasons and error messages")
		fmt.Println("   - Stack traces for failed transactions")
		fmt.Println()

		fmt.Println("5. Time Manipulation:")
		fmt.Println("   - Fast-forward blockchain time (evm_increaseTime)")
		fmt.Println("   - Mine blocks on demand (evm_mine)")
		fmt.Println("   - Test time-dependent contracts")
		fmt.Println("   - Useful for vesting, voting periods, etc.")
		fmt.Println()

		fmt.Println("6. Snapshot and Revert:")
		fmt.Println("   - Take state snapshots (evm_snapshot)")
		fmt.Println("   - Revert to snapshots (evm_revert)")
		fmt.Println("   - Reset between tests for isolation")
		fmt.Println("   - Faster than restarting network")
		fmt.Println()

		fmt.Println("7. Account Impersonation:")
		fmt.Println("   - Impersonate any address (hardhat_impersonateAccount)")
		fmt.Println("   - Test multi-signature wallets")
		fmt.Println("   - Simulate interactions from any account")
		fmt.Println("   - No need for private keys")
		fmt.Println()

		fmt.Println("Popular Development Networks:")
		fmt.Println("  - Hardhat Network: Built-in, console.log, stack traces")
		fmt.Println("  - Ganache: Standalone, GUI available, easy setup")
		fmt.Println("  - Anvil (Foundry): Rust-based, fast, powerful")
		fmt.Println("  - Geth --dev: Full Geth node, production-like")
		fmt.Println()
	}

	// ============================================================================
	// STEP 6: Development Workflow Best Practices
	// ============================================================================
	// BREAKPOINT: Set breakpoint here

	if demo == "all" {
		fmt.Println("=== Development Workflow Best Practices ===")
		fmt.Println()

		fmt.Println("Phase 1: Local Development")
		fmt.Println("  Network: Local devnet (Hardhat, Ganache)")
		fmt.Println("  Activities:")
		fmt.Println("    ✓ Write smart contracts")
		fmt.Println("    ✓ Write unit tests")
		fmt.Println("    ✓ Debug contract logic")
		fmt.Println("    ✓ Rapid iteration (instant mining)")
		fmt.Println("  Benefits: Fast, free, full control, debugging tools")
		fmt.Println()

		fmt.Println("Phase 2: Mainnet Fork Testing")
		fmt.Println("  Network: Local fork of mainnet")
		fmt.Println("  Activities:")
		fmt.Println("    ✓ Test against real contracts (Uniswap, USDC, etc.)")
		fmt.Println("    ✓ Verify integration with existing protocols")
		fmt.Println("    ✓ Test with real token balances and liquidity")
		fmt.Println("    ✓ Simulate mainnet conditions without risk")
		fmt.Println("  Benefits: Realistic testing, no real money at risk")
		fmt.Println()

		fmt.Println("Phase 3: Testnet Deployment")
		fmt.Println("  Network: Public testnet (Sepolia)")
		fmt.Println("  Activities:")
		fmt.Println("    ✓ Deploy contracts publicly")
		fmt.Println("    ✓ Integration testing with other testnet contracts")
		fmt.Println("    ✓ User acceptance testing (UAT)")
		fmt.Println("    ✓ Frontend integration testing")
		fmt.Println("    ✓ Monitor gas usage and optimization")
		fmt.Println("  Benefits: Public visibility, realistic network conditions")
		fmt.Println()

		fmt.Println("Phase 4: Security Audit")
		fmt.Println("  Network: Testnet (for auditors to review)")
		fmt.Println("  Activities:")
		fmt.Println("    ✓ Professional security audit")
		fmt.Println("    ✓ Formal verification (if applicable)")
		fmt.Println("    ✓ Bug bounty program on testnet")
		fmt.Println("    ✓ Fix identified issues")
		fmt.Println("  Benefits: Catch vulnerabilities before mainnet")
		fmt.Println()

		fmt.Println("Phase 5: Mainnet Deployment")
		fmt.Println("  Network: Ethereum Mainnet")
		fmt.Println("  Activities:")
		fmt.Println("    ✓ Deploy to mainnet")
		fmt.Println("    ✓ Gradual rollout (if applicable)")
		fmt.Println("    ✓ Monitor closely for issues")
		fmt.Println("    ✓ Emergency pause mechanisms ready")
		fmt.Println("  Benefits: Production ready, real users, real value")
		fmt.Println()

		fmt.Println("Network Selection Checklist:")
		fmt.Println("  ✓ Local devnet: Use for all initial development")
		fmt.Println("  ✓ Mainnet fork: Use for integration testing")
		fmt.Println("  ✓ Testnet: Use for public testing and audits")
		fmt.Println("  ✓ Mainnet: Only after thorough testing and audits")
		fmt.Println("  ✓ Always verify Chain ID before transactions")
		fmt.Println("  ✓ Use different addresses for different networks")
		fmt.Println("  ✓ Never reuse private keys across networks")
		fmt.Println()

		fmt.Println("Common Mistakes to Avoid:")
		fmt.Println("  ⚠ Deploying to mainnet without testnet testing")
		fmt.Println("  ⚠ Using same addresses across all networks")
		fmt.Println("  ⚠ Not checking Chain ID before signing transactions")
		fmt.Println("  ⚠ Skipping security audits for valuable contracts")
		fmt.Println("  ⚠ Inadequate testing of edge cases")
		fmt.Println("  ⚠ Not having emergency pause mechanisms")
		fmt.Println()
	}

	// ============================================================================
	// DEBUGGING TIPS
	// ============================================================================
	//
	// KEY DEBUGGING TECHNIQUES:
	//
	// 1. Set breakpoints at "// BREAKPOINT:" comments
	// 2. Press F5 to start debugging
	// 3. Use F10 (Step Over) to execute line by line
	// 4. Use F11 (Step Into) to enter network identification functions
	// 5. Watch Variables panel to see network properties
	// 6. Inspect Chain ID and network characteristics
	//
	// NETWORK VERIFICATION:
	// - Always check Chain ID before any transaction
	// - Verify you're on intended network (dev/test/prod)
	// - Check genesis hash to confirm network fork
	// - Monitor gas prices to detect network congestion
	//
	// DEVELOPMENT WORKFLOW:
	// - Start with local devnet for rapid iteration
	// - Use mainnet forks for realistic testing
	// - Deploy to testnet for integration testing
	// - Conduct security audit before mainnet
	// - Deploy to mainnet only after thorough testing
	//
	// ETHEREUM CONCEPTS TO EXPLORE:
	// - Why different networks? Separation of environments
	// - What's Chain ID? Replay protection mechanism (EIP-155)
	// - How to get testnet ETH? Use faucets (search "Sepolia faucet")
	// - What's a fork? Copy of network state at specific block
	// - Why testnets reset? To prevent state bloat and maintain utility
}
