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
Debug harness for geth/05-tx-nonces.

NOTE: This requires a valid private key with testnet ETH.
For safety, this uses NoSend=true to only create the transaction without sending.
*/

func main() {
	fmt.Println("geth/05-tx-nonces: Debug Harness")
	fmt.Println("================================")
	fmt.Println()

	// WARNING: Replace with a test private key that has testnet ETH
	// This is for debugging only - never use mainnet keys!
	privateKeyHex := os.Getenv("TEST_PRIVATE_KEY")
	if privateKeyHex == "" {
		fmt.Println("WARNING: TEST_PRIVATE_KEY not set. Using NoSend=true mode.")
		fmt.Println("Set TEST_PRIVATE_KEY environment variable to test actual sending.")
		fmt.Println()
	}

	rpcURL := "https://eth.llamarpc.com"
	toAddress := common.HexToAddress("0x0000000000000000000000000000000000000001")
	amountWei := big.NewInt(1000000000000000) // 0.001 ETH

	var privateKey *ecdsa.PrivateKey
	var err error
	if privateKeyHex != "" {
		privateKey, err = crypto.HexToECDSA(privateKeyHex[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid private key: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Generate a random key for testing (won't have funds)
		privateKey, _ = crypto.GenerateKey()
		fmt.Println("Using randomly generated key (NoSend mode)")
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
		NoSend:     privateKeyHex == "", // Don't send if no key provided
	}

	result, err := txnonces.Run(ctx, client, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Result ===")
	fmt.Printf("From:    %s\n", result.FromAddress.Hex())
	fmt.Printf("Nonce:   %d\n", result.Nonce)
	fmt.Printf("Tx Hash: %s\n", result.Tx.Hash().Hex())
	if cfg.NoSend {
		fmt.Println("\n(Transaction created but not sent - NoSend=true)")
	}
}
