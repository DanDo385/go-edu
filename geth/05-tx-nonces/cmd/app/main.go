package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"geth/05-tx-nonces/internal/txnonces"
)

/*
geth/05-tx-nonces: Send transactions with proper nonce handling.

This module teaches:
- Transaction nonce management
- Signing and sending transactions
- Understanding transaction lifecycle

Usage:
  go run ./cmd/app/main.go <RPC_URL> <private_key> <to_address> <amount_wei>

Examples:
  go run ./cmd/app/main.go https://eth.llamarpc.com 0x... 0x... 1000000000000000000
*/

func main() {
	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "Usage: %s <RPC_URL> <private_key> <to_address> <amount_wei>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s https://eth.llamarpc.com 0x... 0x... 1000000000000000000\n", os.Args[0])
		os.Exit(1)
	}

	rpcURL := os.Args[1]
	privateKeyHex := os.Args[2]
	toAddressHex := os.Args[3]
	amountWeiStr := os.Args[4]

	// Parse private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex[2:]) // Remove 0x prefix
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid private key: %v\n", err)
		os.Exit(1)
	}

	// Parse to address
	toAddress := common.HexToAddress(toAddressHex)

	// Parse amount
	amountWei, ok := new(big.Int).SetString(amountWeiStr, 10)
	if !ok {
		fmt.Fprintf(os.Stderr, "Invalid amount: %s\n", amountWeiStr)
		os.Exit(1)
	}

	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Printf("To:      %s\n", toAddress.Hex())
	fmt.Printf("Amount:  %s wei\n", amountWei.String())
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rpcClient, err := rpc.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer rpcClient.Close()

	client := ethclient.NewClient(rpcClient)

	cfg := txnonces.Config{
		PrivateKey: privateKey,
		To:         toAddress,
		AmountWei:  amountWei,
	}

	result, err := txnonces.Run(ctx, client, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Transaction Result ===")
	fmt.Printf("From:    %s\n", result.FromAddress.Hex())
	fmt.Printf("Nonce:   %d\n", result.Nonce)
	fmt.Printf("Tx Hash: %s\n", result.Tx.Hash().Hex())
}
