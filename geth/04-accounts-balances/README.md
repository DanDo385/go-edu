# 04: Accounts and Balances

## What Is This Project About?

This module teaches you how to query Ethereum account state—specifically balances and nonces. You'll understand the difference between externally owned accounts (EOAs) and contract accounts, how Ether balances are stored and represented, and what transaction nonces are used for.

## Why Is This Important?

Querying account state is fundamental to almost every Ethereum application:
- Wallets display balances
- DApps check if users can afford transactions
- Transaction builders need nonces for sequencing
- Analytics tools track account activity

## Real-World Problems This Solves

- **Wallet balance display**: Show users their ETH holdings
- **Transaction validation**: Ensure sufficient funds before sending
- **Nonce management**: Get the correct nonce for new transactions
- **Account monitoring**: Track account activity over time

## Key Concepts You'll Learn

- **Account types**: EOA vs Contract accounts
- **Balance representation**: Wei (10^-18 ETH) as the base unit
- **Nonces**: Sequential transaction counter for replay protection
- **State queries**: eth_getBalance and eth_getTransactionCount

## Prerequisites

- Completion of `geth/01-stack` through `geth/03-keys-addresses`

## How to Run

```bash
# Query a specific address
go run ./cmd/app/main.go https://eth.llamarpc.com 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045
```

## Testing

```bash
go test -v ./...
```
