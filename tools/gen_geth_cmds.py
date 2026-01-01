#!/usr/bin/env python3
"""
Generates geth/*/cmd/{app,dev}/main.go implementations.

This repo is intentionally educational. The internal packages contain the actual
exercise logic; these cmd entrypoints are thin CLIs / debug harnesses that call
into that logic with sensible defaults.
"""

from __future__ import annotations

from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
GETH = ROOT / "geth"

DEFAULT_RPC = "https://eth.llamarpc.com"
USDC = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
VITALIK = "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"


def write(path: Path, content: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content.rstrip() + "\n", encoding="utf-8")


def gen() -> None:
    # NOTE: These files are intentionally verbose, with BREAKPOINT markers.
    files: dict[str, tuple[str, str]] = {}

    # 06-smart-contracts is console-first.
    files["06-smart-contracts"] = (
        """package main

import (
\t\"fmt\"
\t\"os\"
)

func main() {
\t// geth/06-smart-contracts (console-first)
\t//
\t// This module is intentionally taught via the Geth JavaScript console.
\t// Open geth/06-smart-contracts/README.md and follow the tutorial.
\t//
\t// Usage:
\t//   go run ./geth/06-smart-contracts/cmd/app
\t//
\t// BREAKPOINT: open README.md
\tfmt.Fprintln(os.Stdout, \"See geth/06-smart-contracts/README.md for the console tutorial.\")
}
""",
        """package main

import \"fmt\"

func main() {
\tfmt.Println(\"BREAKPOINT: open geth/06-smart-contracts/README.md and run the console steps\")
}
""",
    )

    # Helper: positional RPC_URL and environment fallback.
    helper = f"""func rpcURLFromArgs(defaultURL string) string {{
\tif len(os.Args) >= 2 && os.Args[1] != \"\" {{
\t\treturn os.Args[1]
\t}}
\tif v := os.Getenv(\"RPC_URL\"); v != \"\" {{
\t\treturn v
\t}}
\treturn defaultURL
}}
"""

    # 01-stack
    files["01-stack"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"math/big\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/01-stack/internal/stack\"
)

{helper}

func main() {{
\t// geth/01-stack
\t//
\t// Usage:
\t//   go run ./geth/01-stack/cmd/app <RPC_URL> [block_number]
\t//
\t// Examples:
\t//   go run ./geth/01-stack/cmd/app {DEFAULT_RPC}
\t//   go run ./geth/01-stack/cmd/app {DEFAULT_RPC} 19000000
\t//
\t// BREAKPOINT: parse args
\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")

\tvar blockNumber *big.Int
\tif len(os.Args) >= 3 {{
\t\t// BREAKPOINT: parse optional block number
\t\tvar n big.Int
\t\tif _, ok := n.SetString(os.Args[2], 10); ok {{
\t\t\tblockNumber = &n
\t\t}} else {{
\t\t\tfmt.Fprintln(os.Stderr, \"invalid block_number (expected base-10 uint)\")
\t\t\tos.Exit(2)
\t\t}}
\t}}

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\t// BREAKPOINT: dial RPC
\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer client.Close()

\t// BREAKPOINT: call internal module logic
\tout, err := stack.Run(ctx, client, stack.Config{{BlockNumber: blockNumber}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}}

\tfmt.Println(\"ChainID:\", out.ChainID)
\tfmt.Println(\"NetworkID:\", out.NetworkID)
\tfmt.Println(\"Header.Number:\", out.Header.Number)
\tfmt.Println(\"Header.Hash:\", out.Header.Hash())
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"math/big\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/01-stack/internal/stack\"
)

func main() {{
\t// Deterministic debug harness:
\t// - fixed RPC URL
\t// - fixed optional block number (nil => latest)
\t//
\t// BREAKPOINT: change inputs here
\trpcURL := \"{DEFAULT_RPC}\"
\tvar blockNumber *big.Int = nil

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\t// BREAKPOINT: dial RPC
\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer client.Close()

\t// BREAKPOINT: run
\tout, err := stack.Run(ctx, client, stack.Config{{BlockNumber: blockNumber}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(\"ChainID:\", out.ChainID)
\tfmt.Println(\"Header.Number:\", out.Header.Number)
}}
""",
    )

    # 02-rpc-basics
    files["02-rpc-basics"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"os\"
\t\"strconv\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/02-rpc-basics/internal/rpcbasics\"
)

{helper}

func main() {{
\t// geth/02-rpc-basics
\t//
\t// Usage:
\t//   go run ./geth/02-rpc-basics/cmd/app <RPC_URL> [retries]
\t//
\t// BREAKPOINT: parse args
\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")

\tretries := 3
\tif len(os.Args) >= 3 {{
\t\tif v, err := strconv.Atoi(os.Args[2]); err == nil {{
\t\t\tretries = v
\t\t}}
\t}}

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\t// BREAKPOINT: dial
\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer client.Close()

\t// BREAKPOINT: run
\tout, err := rpcbasics.Run(ctx, client, rpcbasics.Config{{Retries: retries}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}}

\tfmt.Println(\"NetworkID:\", out.NetworkID)
\tfmt.Println(\"Latest BlockNumber:\", out.BlockNumber)
\tif out.Block != nil {{
\t\tfmt.Println(\"Block.Hash:\", out.Block.Hash())
\t\tfmt.Println(\"Block.TxCount:\", len(out.Block.Transactions()))
\t\tfmt.Println(\"Block.Number:\", out.Block.Number())
\t}}
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/02-rpc-basics/internal/rpcbasics\"
)

func main() {{
\t// BREAKPOINT: deterministic inputs
\trpcURL := \"{DEFAULT_RPC}\"
\tretries := 3

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer client.Close()

\tout, err := rpcbasics.Run(ctx, client, rpcbasics.Config{{Retries: retries}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(\"NetworkID:\", out.NetworkID)
\tfmt.Println(\"BlockNumber:\", out.BlockNumber)
}}
""",
    )

    # 03-keys-addresses
    files["03-keys-addresses"] = (
        """package main

import (
\t\"flag\"
\t\"fmt\"
\t\"os\"

\t\"github.com/example/go-10x-minis/geth/03-keys-addresses/internal/keysaddresses\"
)

func main() {
\t// geth/03-keys-addresses
\t//
\t// This module does NOT need an RPC endpoint. It demonstrates:
\t// - generating a keypair
\t// - deriving an Ethereum address
\t// - writing an encrypted keystore file
\t//
\t// Usage:
\t//   go run ./geth/03-keys-addresses/cmd/app --out ./keystore-demo --pass changeit
\t//
\t// BREAKPOINT: parse flags
\toutDir := flag.String(\"out\", \"./keystore-demo\", \"Directory to write keystore file\")
\tpass := flag.String(\"pass\", \"changeit\", \"Keystore passphrase (demo default is intentionally weak)\")
\tflag.Parse()

\tres, err := keysaddresses.Run(keysaddresses.Config{OutputDir: *outDir, Passphrase: *pass})
\tif err != nil {
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}

\tfmt.Println(\"Address:\", res.Address.Hex())
\tfmt.Println(\"PrivateKey (hex):\", res.PrivateKeyHex)
\tfmt.Println(\"KeystorePath:\", res.KeystorePath)
}
""",
        """package main

import (
\t\"fmt\"

\t\"github.com/example/go-10x-minis/geth/03-keys-addresses/internal/keysaddresses\"
)

func main() {
\t// Deterministic-ish harness: fixed output dir and passphrase.
\t// (Key generation is random by design.)
\t//
\t// BREAKPOINT: change inputs
\tres, err := keysaddresses.Run(keysaddresses.Config{OutputDir: \"./keystore-demo\", Passphrase: \"changeit\"})
\tif err != nil {
\t\tpanic(err)
\t}

\tfmt.Println(\"Address:\", res.Address.Hex())
\tfmt.Println(\"KeystorePath:\", res.KeystorePath)
}
""",
    )

    # 04-accounts-balances
    files["04-accounts-balances"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/04-accounts-balances/internal/accountsbalances\"
)

{helper}

func main() {{
\t// geth/04-accounts-balances
\t//
\t// Usage:
\t//   go run ./geth/04-accounts-balances/cmd/app <RPC_URL> [addr1] [addr2] ...
\t//
\t// If no addresses are provided, we query a known EOA + a known contract.
\t//
\t// BREAKPOINT: parse args
\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")

\taddrs := []common.Address{{common.HexToAddress(\"{VITALIK}\"), common.HexToAddress(\"{USDC}\")}}
\tif len(os.Args) >= 3 {{
\t\taddrs = nil
\t\tfor _, s := range os.Args[2:] {{
\t\t\taddrs = append(addrs, common.HexToAddress(s))
\t\t}}
\t}}

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer client.Close()

\t// BREAKPOINT: run
\tout, err := accountsbalances.Run(ctx, client, accountsbalances.Config{{Addresses: addrs}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}}

\tfor _, a := range out.Accounts {{
\t\tfmt.Println(\"Address:\", a.Address.Hex())
\t\tfmt.Println(\"Type:\", a.Type)
\t\tfmt.Println(\"BalanceWei:\", a.Balance)
\t\tfmt.Println(\"CodeBytes:\", len(a.Code))
\t\tfmt.Println()
\t}}
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/04-accounts-balances/internal/accountsbalances\"
)

func main() {{
\t// BREAKPOINT: deterministic inputs
\trpcURL := \"{DEFAULT_RPC}\"
\taddrs := []common.Address{{
\t\tcommon.HexToAddress(\"{VITALIK}\"),
\t\tcommon.HexToAddress(\"{USDC}\"),
\t}}

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer client.Close()

\tout, err := accountsbalances.Run(ctx, client, accountsbalances.Config{{Addresses: addrs}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(\"Accounts:\", len(out.Accounts))
\tfor _, a := range out.Accounts {{
\t\tfmt.Println(a.Address.Hex(), a.Type, a.Balance)
\t}}
}}
""",
    )

    # 05-tx-nonces
    files["05-tx-nonces"] = (
        """package main

import (
\t\"context\"
\t\"fmt\"
\t\"math/big\"
\t\"os\"
\t\"strings\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/crypto\"
\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/05-tx-nonces/internal/txnonces\"
)

func main() {
\t// geth/05-tx-nonces
\t//
\t// Usage:
\t//   go run ./geth/05-tx-nonces/cmd/app <RPC_URL> <private_key_hex> <to_address> <amount_wei> [--send]
\t//
\t// By default this command builds + signs the tx but does NOT broadcast it.
\t// Pass --send to broadcast.
\t//
\t// BREAKPOINT: parse args
\tif len(os.Args) < 5 {
\t\tfmt.Fprintln(os.Stderr, \"usage: <RPC_URL> <private_key_hex> <to_address> <amount_wei> [--send]\")
\t\tos.Exit(2)
\t}

\trpcURL := os.Args[1]
\tpkHex := strings.TrimPrefix(os.Args[2], \"0x\")
\tto := common.HexToAddress(os.Args[3])
\tamount, ok := new(big.Int).SetString(os.Args[4], 10)
\tif !ok {
\t\tfmt.Fprintln(os.Stderr, \"invalid amount_wei\")
\t\tos.Exit(2)
\t}

\tsend := false
\tif len(os.Args) >= 6 && os.Args[5] == \"--send\" {
\t\tsend = true
\t}

\tpk, err := crypto.HexToECDSA(pkHex)
\tif err != nil {
\t\tfmt.Fprintln(os.Stderr, \"invalid private key:\", err)
\t\tos.Exit(2)
\t}

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}
\tdefer client.Close()

\t// BREAKPOINT: build + sign tx (optionally broadcast)
\tout, err := txnonces.Run(ctx, client, txnonces.Config{
\t\tPrivateKey: pk,
\t\tTo:         to,
\t\tAmountWei:  amount,
\t\tGasLimit:   21000,
\t\tNoSend:     !send,
\t})
\tif err != nil {
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}

\tfmt.Println(\"From:\", out.FromAddress.Hex())
\tfmt.Println(\"Nonce:\", out.Nonce)
\tfmt.Println(\"TxHash:\", out.Tx.Hash().Hex())
\tfmt.Println(\"Sent:\", send)
}
""",
        f"""package main

import (
\t\"context\"
\t\"crypto/ecdsa\"
\t\"fmt\"
\t\"math/big\"
\t\"strings\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/crypto\"
\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/05-tx-nonces/internal/txnonces\"
)

func mustKey(hex string) *ecdsa.PrivateKey {{
\thex = strings.TrimPrefix(hex, \"0x\")
\tk, err := crypto.HexToECDSA(hex)
\tif err != nil {{
\t\tpanic(err)
\t}}
\treturn k
}}

func main() {{
\t// Deterministic harness: NoSend=true so this is safe.
\t//
\t// BREAKPOINT: set a throwaway private key for local devnets.
\trpcURL := \"{DEFAULT_RPC}\"
\tprivateKeyHex := \"0x\" // set me
\tto := common.HexToAddress(\"{VITALIK}\")
\tamountWei := big.NewInt(1)

\tif privateKeyHex == \"0x\" {{
\t\tfmt.Println(\"Set privateKeyHex in cmd/dev before running (use a devnet key).\")\n\t\treturn
\t}}

\tpk := mustKey(privateKeyHex)

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer client.Close()

\tout, err := txnonces.Run(ctx, client, txnonces.Config{{
\t\tPrivateKey: pk,
\t\tTo:         to,
\t\tAmountWei:  amountWei,
\t\tGasLimit:   21000,
\t\tNoSend:     true,
\t}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(\"From:\", out.FromAddress.Hex())
\tfmt.Println(\"Nonce:\", out.Nonce)
\tfmt.Println(\"TxHash:\", out.Tx.Hash().Hex())
}}
""",
    )

    # 06-eip1559
    files["06-eip1559"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"math/big\"
\t\"os\"
\t\"strings\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/crypto\"
\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/06-eip1559/internal/eip1559\"
)

func main() {{
\t// geth/06-eip1559
\t//
\t// Usage:
\t//   go run ./geth/06-eip1559/cmd/app <RPC_URL> <private_key_hex> <to_address> <amount_wei> [--send]
\t//
\t// By default we do NOT broadcast. Pass --send to broadcast.
\t//
\t// BREAKPOINT: parse args
\tif len(os.Args) < 5 {{
\t\tfmt.Fprintln(os.Stderr, \"usage: <RPC_URL> <private_key_hex> <to_address> <amount_wei> [--send]\")
\t\tos.Exit(2)
\t}}

\trpcURL := os.Args[1]
\tpkHex := strings.TrimPrefix(os.Args[2], \"0x\")
\tto := common.HexToAddress(os.Args[3])
\tamount, ok := new(big.Int).SetString(os.Args[4], 10)
\tif !ok {{
\t\tfmt.Fprintln(os.Stderr, \"invalid amount_wei\")
\t\tos.Exit(2)
\t}}

\tsend := len(os.Args) >= 6 && os.Args[5] == \"--send\"

\tpk, err := crypto.HexToECDSA(pkHex)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"invalid private key:\", err)
\t\tos.Exit(2)
\t}}

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer client.Close()

\t// BREAKPOINT: run
\tout, err := eip1559.Run(ctx, client, eip1559.Config{{
\t\tPrivateKey: pk,
\t\tTo:         to,
\t\tAmountWei:  amount,
\t\tGasLimit:   21000,
\t\tNoSend:     !send,
\t}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}}

\tfmt.Println(\"From:\", out.FromAddress.Hex())
\tfmt.Println(\"Nonce:\", out.Nonce)
\tfmt.Println(\"BaseFee:\", out.BaseFee)
\tfmt.Println(\"TxHash:\", out.Tx.Hash().Hex())
\tfmt.Println(\"Sent:\", send)
}}
""",
        f"""package main

import (
\t\"context\"
\t\"crypto/ecdsa\"
\t\"fmt\"
\t\"math/big\"
\t\"strings\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/crypto\"
\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/06-eip1559/internal/eip1559\"
)

func mustKey(hex string) *ecdsa.PrivateKey {{
\thex = strings.TrimPrefix(hex, \"0x\")
\tk, err := crypto.HexToECDSA(hex)
\tif err != nil {{
\t\tpanic(err)
\t}}
\treturn k
}}

func main() {{
\t// BREAKPOINT: set a throwaway private key for devnets.
\trpcURL := \"{DEFAULT_RPC}\"
\tprivateKeyHex := \"0x\" // set me
\tto := common.HexToAddress(\"{VITALIK}\")
\tamountWei := big.NewInt(1)

\tif privateKeyHex == \"0x\" {{
\t\tfmt.Println(\"Set privateKeyHex in cmd/dev before running (use a devnet key).\")\n\t\treturn
\t}}

\tpk := mustKey(privateKeyHex)

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer client.Close()

\tout, err := eip1559.Run(ctx, client, eip1559.Config{{
\t\tPrivateKey: pk,
\t\tTo:         to,
\t\tAmountWei:  amountWei,
\t\tGasLimit:   21000,
\t\tNoSend:     true,
\t}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(\"BaseFee:\", out.BaseFee)
\tfmt.Println(\"TxHash:\", out.Tx.Hash().Hex())
}}
""",
    )

    # 07-eth-call
    files["07-eth-call"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/07-eth-call/internal/ethcall\"
)

{helper}

func main() {{
\t// geth/07-eth-call
\t//
\t// Prerequisite: geth/06-smart-contracts (console intuition for calls vs txs).
\t//
\t// Usage:
\t//   go run ./geth/07-eth-call/cmd/app <RPC_URL> [contract_address]
\t//
\t// Example:
\t//   go run ./geth/07-eth-call/cmd/app {DEFAULT_RPC} {USDC}
\t//
\t// BREAKPOINT: parse args
\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")

\tcontract := common.HexToAddress(\"{USDC}\")
\tif len(os.Args) >= 3 {{
\t\tcontract = common.HexToAddress(os.Args[2])
\t}}

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer client.Close()

\t// BREAKPOINT: run
\tout, err := ethcall.Run(ctx, client, ethcall.Config{{Contract: contract}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}}

\tfmt.Println(\"Name:\", out.Name)
\tfmt.Println(\"Symbol:\", out.Symbol)
\tfmt.Println(\"Decimals:\", out.Decimals)
\tfmt.Println(\"TotalSupply:\", out.TotalSupply)
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/07-eth-call/internal/ethcall\"
)

func main() {{
\t// BREAKPOINT: deterministic inputs
\trpcURL := \"{DEFAULT_RPC}\"
\tcontract := common.HexToAddress(\"{USDC}\")

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer client.Close()

\tout, err := ethcall.Run(ctx, client, ethcall.Config{{Contract: contract}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(out.Name, out.Symbol, out.Decimals)
}}
""",
    )

    # 08-abigen
    files["08-abigen"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/08-abigen/internal/abigen\"
)

{helper}

func main() {{
\t// geth/08-abigen
\t//
\t// Prerequisites:
\t// - geth/06-smart-contracts (console calls/txs)
\t// - geth/07-eth-call (manual ABI intuition)
\t//
\t// Usage:
\t//   go run ./geth/08-abigen/cmd/app <RPC_URL> [contract_address] [holder_address]
\t//
\t// BREAKPOINT: parse args
\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")

\tcontract := common.HexToAddress(\"{USDC}\")
\tif len(os.Args) >= 3 {{
\t\tcontract = common.HexToAddress(os.Args[2])
\t}}

\tvar holder *common.Address
\tif len(os.Args) >= 4 {{
\t\th := common.HexToAddress(os.Args[3])
\t\tholder = &h
\t}}

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer client.Close()

\t// BREAKPOINT: run
\tout, err := abigen.Run(ctx, client, abigen.Config{{Contract: contract, Holder: holder}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}}

\tfmt.Println(\"Name:\", out.Name)
\tfmt.Println(\"Symbol:\", out.Symbol)
\tfmt.Println(\"Decimals:\", out.Decimals)
\tfmt.Println(\"TotalSupply:\", out.TotalSupply)
\tif out.Balance != nil {{
\t\tfmt.Println(\"Balance:\", out.Balance)
\t}}
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/08-abigen/internal/abigen\"
)

func main() {{
\t// BREAKPOINT: deterministic inputs
\trpcURL := \"{DEFAULT_RPC}\"
\tcontract := common.HexToAddress(\"{USDC}\")

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer client.Close()

\tout, err := abigen.Run(ctx, client, abigen.Config{{Contract: contract}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(out.Name, out.Symbol, out.Decimals)
}}
""",
    )

    # 09-events
    files["09-events"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"math/big\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/09-events/internal/events\"
)

{helper}

func main() {{
\t// geth/09-events
\t//
\t// Prerequisites:
\t// - geth/06-smart-contracts (receipts/logs in console)
\t// - geth/07-eth-call
\t// - geth/08-abigen
\t//
\t// Usage:
\t//   go run ./geth/09-events/cmd/app <RPC_URL> [token_address] [from_block] [to_block]
\t//
\t// BREAKPOINT: parse args
\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")

\ttoken := common.HexToAddress(\"{USDC}\")
\tif len(os.Args) >= 3 {{
\t\ttoken = common.HexToAddress(os.Args[2])
\t}}

\tctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer client.Close()

\t// Pick a recent window if not provided.
\theader, err := client.HeaderByNumber(ctx, nil)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"header:\", err)
\t\tos.Exit(1)
\t}}
\ttoBlock := new(big.Int).Set(header.Number)
\tfromBlock := new(big.Int).Sub(toBlock, big.NewInt(500))

\tif len(os.Args) >= 5 {{
\t\tif v, ok := new(big.Int).SetString(os.Args[3], 10); ok {{
\t\t\tfromBlock = v
\t\t}}
\t\tif v, ok := new(big.Int).SetString(os.Args[4], 10); ok {{
\t\t\ttoBlock = v
\t\t}}
\t}}

\t// BREAKPOINT: run
\tout, err := events.Run(ctx, client, events.Config{{Token: token, FromBlock: fromBlock, ToBlock: toBlock}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}}

\tfmt.Println(\"Events:\", len(out.Events))
\tfor i, ev := range out.Events {{
\t\tif i >= 5 {{
\t\t\tbreak
\t\t}}
\t\tfmt.Println(ev.BlockNumber, ev.TxHash.Hex(), ev.From.Hex(), \"->\", ev.To.Hex(), ev.Value)
\t}}
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"math/big\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/09-events/internal/events\"
)

func main() {{
\t// BREAKPOINT: deterministic inputs
\trpcURL := \"{DEFAULT_RPC}\"
\ttoken := common.HexToAddress(\"{USDC}\")
\tfromBlock := big.NewInt(0) // set me if you want

\tctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer client.Close()

\theader, err := client.HeaderByNumber(ctx, nil)
\tif err != nil {{
\t\tpanic(err)
\t}}
\ttoBlock := header.Number
\tif fromBlock.Sign() == 0 {{
\t\tfromBlock = new(big.Int).Sub(toBlock, big.NewInt(500))
\t}}

\tout, err := events.Run(ctx, client, events.Config{{Token: token, FromBlock: fromBlock, ToBlock: toBlock}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(\"Events:\", len(out.Events))
}}
""",
    )

    # 10-filters (polling mode default to work with HTTP RPC endpoints).
    files["10-filters"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/10-filters/internal/filters\"
)

{helper}

func main() {{
\t// geth/10-filters
\t//
\t// This module can use either:
\t// - WebSocket subscriptions (real-time)
\t// - HTTP polling (works everywhere; slower)
\t//
\t// Usage:
\t//   go run ./geth/10-filters/cmd/app <RPC_URL>
\t//
\t// BREAKPOINT: parse args
\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")

\tctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer client.Close()

\t// Default to polling so this works with HTTP endpoints.
\tout, err := filters.Run(ctx, client, filters.Config{{MaxHeads: 5, PollInterval: 2 * time.Second, PollMode: true}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}}

\tfmt.Println(\"Mode:\", out.Mode)
\tfor _, h := range out.Heads {{
\t\tfmt.Println(h.Number, h.Hash.Hex(), \"reorg=\", h.Reorg)
\t}}
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/10-filters/internal/filters\"
)

func main() {{
\t// BREAKPOINT: deterministic inputs
\trpcURL := \"{DEFAULT_RPC}\"

\tctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer client.Close()

\tout, err := filters.Run(ctx, client, filters.Config{{MaxHeads: 5, PollInterval: 2 * time.Second, PollMode: true}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(\"Heads:\", len(out.Heads), \"mode:\", out.Mode)
}}
""",
    )

    # 11-storage
    files["11-storage"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"math/big\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/11-storage/internal/storage\"
)

{helper}

func main() {{
\t// geth/11-storage
\t//
\t// Usage:
\t//   go run ./geth/11-storage/cmd/app <RPC_URL> [contract_address] [slot]
\t//
\t// BREAKPOINT: parse args
\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")

\tcontract := common.HexToAddress(\"{USDC}\")
\tif len(os.Args) >= 3 {{
\t\tcontract = common.HexToAddress(os.Args[2])
\t}}

\tslot := big.NewInt(0)
\tif len(os.Args) >= 4 {{
\t\tif v, ok := new(big.Int).SetString(os.Args[3], 10); ok {{
\t\t\tslot = v
\t\t}}
\t}}

\tctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer client.Close()

\t// BREAKPOINT: run
\tout, err := storage.Run(ctx, client, storage.Config{{Contract: contract, Slot: slot}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}}

\tfmt.Println(\"ResolvedSlot:\", out.ResolvedSlot.Hex())
\tfmt.Printf(\"Value (hex): 0x%x\\n\", out.Value)
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"math/big\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/11-storage/internal/storage\"
)

func main() {{
\t// BREAKPOINT: deterministic inputs
\trpcURL := \"{DEFAULT_RPC}\"
\tcontract := common.HexToAddress(\"{USDC}\")
\tslot := big.NewInt(0)

\tctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer client.Close()

\tout, err := storage.Run(ctx, client, storage.Config{{Contract: contract, Slot: slot}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(out.ResolvedSlot.Hex())
}}
""",
    )

    # 12-proofs (eth_getProof via gethclient)
    files["12-proofs"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/ethclient/gethclient\"
\t\"github.com/ethereum/go-ethereum/rpc\"

\t\"github.com/example/go-10x-minis/geth/12-proofs/internal/proofs\"
)

{helper}

func main() {{
\t// geth/12-proofs
\t//
\t// Usage:
\t//   go run ./geth/12-proofs/cmd/app <RPC_URL> [account_address]
\t//
\t// BREAKPOINT: parse args
\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")

\taccount := common.HexToAddress(\"{VITALIK}\")
\tif len(os.Args) >= 3 {{
\t\taccount = common.HexToAddress(os.Args[2])
\t}}

\tctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
\tdefer cancel()

\t// BREAKPOINT: dial raw RPC (needed for gethclient)
\trpcClient, err := rpc.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial rpc:\", err)
\t\tos.Exit(1)
\t}}
\tdefer rpcClient.Close()

\tgc := gethclient.New(rpcClient)

\t// BREAKPOINT: run
\tout, err := proofs.Run(ctx, gc, proofs.Config{{Account: account}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}}

\tfmt.Println(\"Balance:\", out.Account.Balance)
\tfmt.Println(\"Nonce:\", out.Account.Nonce)
\tfmt.Println(\"CodeHash:\", out.Account.CodeHash.Hex())
\tfmt.Println(\"StorageHash:\", out.Account.StorageHash.Hex())
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/ethclient/gethclient\"
\t\"github.com/ethereum/go-ethereum/rpc\"

\t\"github.com/example/go-10x-minis/geth/12-proofs/internal/proofs\"
)

func main() {{
\t// BREAKPOINT: deterministic inputs
\trpcURL := \"{DEFAULT_RPC}\"
\taccount := common.HexToAddress(\"{VITALIK}\")

\tctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
\tdefer cancel()

\trpcClient, err := rpc.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer rpcClient.Close()

\tgc := gethclient.New(rpcClient)

\tout, err := proofs.Run(ctx, gc, proofs.Config{{Account: account}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(\"Nonce:\", out.Account.Nonce)
}}
""",
    )

    # 13-trace (best-effort debug_traceTransaction)
    files["13-trace"] = (
        f"""package main

import (
\t\"context\"
\t\"encoding/json\"
\t\"fmt\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/ethclient\"
\t\"github.com/ethereum/go-ethereum/rpc\"

\t\"github.com/example/go-10x-minis/geth/13-trace/internal/trace\"
)

{helper}

type rpcTraceClient struct {{
\trpc *rpc.Client
}}

func (c rpcTraceClient) TraceTransaction(ctx context.Context, txHash common.Hash) (json.RawMessage, error) {{
\tvar out json.RawMessage
\terr := c.rpc.CallContext(ctx, &out, \"debug_traceTransaction\", txHash)
\treturn out, err
}}

func main() {{
\t// geth/13-trace
\t//
\t// NOTE: Many public RPC providers disable debug tracing.
\t// This program is best-effort: it will explain failures clearly.
\t//
\t// Usage:
\t//   go run ./geth/13-trace/cmd/app <RPC_URL> [tx_hash]
\t//
\t// BREAKPOINT: parse args
\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")

\tctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
\tdefer cancel()

\tvar txHash common.Hash
\tif len(os.Args) >= 3 {{
\t\ttxHash = common.HexToHash(os.Args[2])
\t}} else {{
\t\t// Pick a recent tx hash (first tx in latest block), if available.
\t\tc, err := ethclient.DialContext(ctx, rpcURL)
\t\tif err != nil {{
\t\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\t\tos.Exit(1)
\t\t}}
\t\tdefer c.Close()

\t\tb, err := c.BlockByNumber(ctx, nil)
\t\tif err == nil && b != nil && len(b.Transactions()) > 0 {{
\t\t\ttxHash = b.Transactions()[0].Hash()
\t\t}} else {{
\t\t\tfmt.Fprintln(os.Stderr, \"provide a tx_hash: could not auto-select from latest block\")
\t\t\tos.Exit(2)
\t\t}}
\t}}

\t// BREAKPOINT: dial raw RPC for debug_* methods
\trpcClient, err := rpc.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial rpc:\", err)
\t\tos.Exit(1)
\t}}
\tdefer rpcClient.Close()

\tout, err := trace.Run(ctx, rpcTraceClient{{rpc: rpcClient}}, trace.Config{{TxHash: txHash}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"trace failed:\", err)
\t\tfmt.Fprintln(os.Stderr, \"Hint: use a local geth node with --http.api debug,eth,net,web3 (and often an archive node for older txs).\")\n\t\tos.Exit(1)
\t}}

\tfmt.Println(\"TxHash:\", out.TxHash.Hex())
\tfmt.Printf(\"TraceBytes: %d\\n\", len(out.Trace))
}}
""",
        f"""package main

import (
\t\"context\"
\t\"encoding/json\"
\t\"fmt\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/rpc\"

\t\"github.com/example/go-10x-minis/geth/13-trace/internal/trace\"
)

type rpcTraceClient struct {{
\trpc *rpc.Client
}}

func (c rpcTraceClient) TraceTransaction(ctx context.Context, txHash common.Hash) (json.RawMessage, error) {{
\tvar out json.RawMessage
\terr := c.rpc.CallContext(ctx, &out, \"debug_traceTransaction\", txHash)
\treturn out, err
}}

func main() {{
\t// BREAKPOINT: deterministic inputs (requires debug-enabled endpoint)
\trpcURL := \"{DEFAULT_RPC}\"
\ttxHash := common.Hash{{}} // set me

\tif txHash == (common.Hash{{}}) {{
\t\tfmt.Println(\"Set txHash in cmd/dev to a real tx hash (and use a debug-enabled node).\")\n\t\treturn
\t}}

\tctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
\tdefer cancel()

\trpcClient, err := rpc.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer rpcClient.Close()

\tout, err := trace.Run(ctx, rpcTraceClient{{rpc: rpcClient}}, trace.Config{{TxHash: txHash}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(\"TraceBytes:\", len(out.Trace))
}}
""",
    )

    # 14-explorer
    files["14-explorer"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"math/big\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/14-explorer/internal/explorer\"
)

{helper}

func main() {{
\t// geth/14-explorer
\t//
\t// Usage:
\t//   go run ./geth/14-explorer/cmd/app <RPC_URL> [block_number] [--txs]
\t//
\t// BREAKPOINT: parse args
\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")
\tincludeTxs := false
\tvar block *big.Int
\tfor _, a := range os.Args[2:] {{
\t\tif a == \"--txs\" {{
\t\t\tincludeTxs = true
\t\t\tcontinue
\t\t}}
\t\tif v, ok := new(big.Int).SetString(a, 10); ok {{
\t\t\tblock = v
\t\t}}
\t}}

\tctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer client.Close()

\tout, err := explorer.Run(ctx, client, explorer.Config{{Number: block, IncludeTxs: includeTxs}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}}

\tfmt.Println(\"Block:\", out.Number, out.Hash.Hex())
\tfmt.Println(\"TxCount:\", out.TxCount)
\tif includeTxs {{
\t\tfor i, tx := range out.Txs {{
\t\t\tif i >= 5 {{
\t\t\t\tbreak
\t\t\t}}
\t\t\tfmt.Println(tx.Hash.Hex(), tx.Gas)
\t\t}}
\t}}
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/14-explorer/internal/explorer\"
)

func main() {{
\t// BREAKPOINT
\trpcURL := \"{DEFAULT_RPC}\"

\tctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer client.Close()

\tout, err := explorer.Run(ctx, client, explorer.Config{{IncludeTxs: true}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(\"Block:\", out.Number, \"txs:\", out.TxCount)
}}
""",
    )

    # 15-receipts
    files["15-receipts"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/15-receipts/internal/receipts\"
)

{helper}

func main() {{
\t// geth/15-receipts
\t//
\t// Usage:
\t//   go run ./geth/15-receipts/cmd/app <RPC_URL> [tx_hash]
\t//
\t// BREAKPOINT: parse args
\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")

\tctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer client.Close()

\tvar txHash common.Hash
\tif len(os.Args) >= 3 {{
\t\ttxHash = common.HexToHash(os.Args[2])
\t}} else {{
\t\tb, err := client.BlockByNumber(ctx, nil)
\t\tif err != nil || b == nil || len(b.Transactions()) == 0 {{
\t\t\tfmt.Fprintln(os.Stderr, \"provide tx_hash: could not auto-select from latest block\")
\t\t\tos.Exit(2)
\t\t}}
\t\ttxHash = b.Transactions()[0].Hash()
\t}}

\t// BREAKPOINT: run
\tout, err := receipts.Run(ctx, client, receipts.Config{{TxHash: txHash}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}}

\tfmt.Println(\"TxHash:\", out.TxHash.Hex())
\tfmt.Println(\"StatusOK:\", out.StatusOK)
\tfmt.Println(\"GasUsed:\", out.GasUsed)
\tfmt.Println(\"Logs:\", len(out.Logs))
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/common\"
\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/15-receipts/internal/receipts\"
)

func main() {{
\t// BREAKPOINT: deterministic inputs
\trpcURL := \"{DEFAULT_RPC}\"
\ttxHash := common.Hash{{}} // optionally set

\tctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer client.Close()

\tif txHash == (common.Hash{{}}) {{
\t\tb, err := client.BlockByNumber(ctx, nil)
\t\tif err != nil || b == nil || len(b.Transactions()) == 0 {{
\t\t\tpanic(\"no txs in latest block\")
\t\t}}
\t\ttxHash = b.Transactions()[0].Hash()
\t}}

\tout, err := receipts.Run(ctx, client, receipts.Config{{TxHash: txHash}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(\"Logs:\", len(out.Logs))
}}
""",
    )

    # 16-concurrency
    files["16-concurrency"] = (
        """package main

import (
\t\"bytes\"
\t\"context\"
\t\"encoding/json\"
\t\"flag\"
\t\"fmt\"
\t\"net/http\"
\t\"os\"
\t\"strings\"
\t\"time\"

\t\"github.com/example/go-10x-minis/geth/16-concurrency/internal/concurrency\"
)

type jsonRPCProber struct {
\thc *http.Client
}

func (p jsonRPCProber) Probe(ctx context.Context, endpoint string) error {
\tpayload := []byte(`{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"eth_chainId\",\"params\":[]}`)
\treq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
\tif err != nil {
\t\treturn err
\t}
\treq.Header.Set(\"Content-Type\", \"application/json\")

\tresp, err := p.hc.Do(req)
\tif err != nil {
\t\treturn err
\t}
\tdefer resp.Body.Close()

\tif resp.StatusCode < 200 || resp.StatusCode >= 300 {
\t\treturn fmt.Errorf(\"non-2xx status: %s\", resp.Status)
\t}

\tvar out struct {
\t\tResult string          `json:\"result\"`
\t\tError  json.RawMessage `json:\"error\"`
\t}
\tif err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
\t\treturn err
\t}
\tif len(out.Error) != 0 {
\t\treturn fmt.Errorf(\"rpc error: %s\", string(out.Error))
\t}
\tif out.Result == \"\" {
\t\treturn fmt.Errorf(\"missing result\")
\t}
\treturn nil
}

func main() {
\t// geth/16-concurrency
\t//
\t// Usage:
\t//   go run ./geth/16-concurrency/cmd/app --endpoints <url1,url2,...> [--workers N]
\t//
\t// BREAKPOINT: parse flags
\tendpointsCSV := flag.String(\"endpoints\", \"https://eth.llamarpc.com,https://rpc.ankr.com/eth\", \"comma-separated RPC URLs\")
\tworkers := flag.Int(\"workers\", 4, \"worker count\")
\ttimeout := flag.Duration(\"timeout\", 5*time.Second, \"overall timeout\")
\tflag.Parse()

\tendpoints := []string{}
\tfor _, s := range strings.Split(*endpointsCSV, \",\") {
\t\tif t := strings.TrimSpace(s); t != \"\" {
\t\t\tendpoints = append(endpoints, t)
\t\t}
\t}
\tif len(endpoints) == 0 {
\t\tfmt.Fprintln(os.Stderr, \"no endpoints\")
\t\tos.Exit(2)
\t}

\tctx, cancel := context.WithTimeout(context.Background(), *timeout)
\tdefer cancel()

\tp := jsonRPCProber{hc: &http.Client{Timeout: *timeout / 2}}

\t// BREAKPOINT: run
\tout, err := concurrency.Run(ctx, p, concurrency.Config{Endpoints: endpoints, Workers: *workers, Timeout: *timeout})
\tif err != nil {
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t}

\tfmt.Println(\"Successes:\", len(out.Successes))
\tfmt.Println(\"Failures:\", len(out.Failures))
\tfor ep, d := range out.Successes {
\t\tfmt.Println(\"OK\", ep, d)
\t}
\tfor ep, e := range out.Failures {
\t\tfmt.Println(\"ERR\", ep, e)
\t}
}
""",
        f"""package main

import (
\t\"bytes\"
\t\"context\"
\t\"encoding/json\"
\t\"fmt\"
\t\"net/http\"
\t\"strings\"
\t\"time\"

\t\"github.com/example/go-10x-minis/geth/16-concurrency/internal/concurrency\"
)

type jsonRPCProber struct {{
\thc *http.Client
}}

func (p jsonRPCProber) Probe(ctx context.Context, endpoint string) error {{
\tpayload := []byte(`{{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"eth_chainId\",\"params\":[]}}`)
\treq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
\tif err != nil {{
\t\treturn err
\t}}
\treq.Header.Set(\"Content-Type\", \"application/json\")

\tresp, err := p.hc.Do(req)
\tif err != nil {{
\t\treturn err
\t}}
\tdefer resp.Body.Close()

\tif resp.StatusCode < 200 || resp.StatusCode >= 300 {{
\t\treturn fmt.Errorf(\"non-2xx status: %s\", resp.Status)
\t}}

\tvar out struct {{
\t\tResult string          `json:\"result\"`
\t\tError  json.RawMessage `json:\"error\"`
\t}}
\tif err := json.NewDecoder(resp.Body).Decode(&out); err != nil {{
\t\treturn err
\t}}
\tif len(out.Error) != 0 {{
\t\treturn fmt.Errorf(\"rpc error: %s\", string(out.Error))
\t}}
\tif out.Result == \"\" {{
\t\treturn fmt.Errorf(\"missing result\")
\t}}
\treturn nil
}}

func main() {{
\t// BREAKPOINT: deterministic inputs
\tendpoints := []string{{\"{DEFAULT_RPC}\", \"https://rpc.ankr.com/eth\"}}
\tworkers := 4
\ttimeout := 5 * time.Second

\tctx, cancel := context.WithTimeout(context.Background(), timeout)
\tdefer cancel()

\tp := jsonRPCProber{{hc: &http.Client{{Timeout: timeout / 2}}}}

\tout, err := concurrency.Run(ctx, p, concurrency.Config{{Endpoints: endpoints, Workers: workers, Timeout: timeout}})
\tif err != nil {{
\t\tfmt.Println(\"run error:\", err)
\t}}

\tfmt.Println(\"Successes:\", len(out.Successes), \"Failures:\", len(out.Failures))
\tfor ep, d := range out.Successes {{
\t\tfmt.Println(\"OK\", strings.TrimSpace(ep), d)
\t}}
}}
""",
    )

    # 17-indexer / 18-reorgs / 19-devnets / 20-node are placeholder internal packages in this snapshot.
    # cmd/* provides a small, runnable demo anyway.
    files["17-indexer"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"math/big\"
\t\"os\"
\t\"strconv\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"
)

{helper}

func main() {{
\t// geth/17-indexer (lightweight demo)
\t//
\t// This module's internal exercise is a placeholder in this repo snapshot.
\t// The cmd/app demonstrates a tiny \"index\" pass: fetch N recent blocks and count txs.
\t//
\t// Usage:
\t//   go run ./geth/17-indexer/cmd/app <RPC_URL> [n]
\t//
\t// BREAKPOINT: parse args
\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")
\tn := 5
\tif len(os.Args) >= 3 {{
\t\tif v, err := strconv.Atoi(os.Args[2]); err == nil {{
\t\t\tn = v
\t\t}}
\t}}

\tctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
\tdefer cancel()

\tc, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer c.Close()

\thead, err := c.HeaderByNumber(ctx, nil)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"header:\", err)
\t\tos.Exit(1)
\t}}

\tstart := new(big.Int).Sub(head.Number, big.NewInt(int64(n-1)))
\tvar total int
\tfor i := new(big.Int).Set(start); i.Cmp(head.Number) <= 0; i.Add(i, big.NewInt(1)) {{
\t\tb, err := c.BlockByNumber(ctx, i)
\t\tif err != nil {{
\t\t\tfmt.Fprintln(os.Stderr, \"block:\", i, err)
\t\t\tcontinue
\t\t}}
\t\ttotal += len(b.Transactions())
\t\tfmt.Println(\"block\", b.NumberU64(), \"txs\", len(b.Transactions()))
\t}}
\tfmt.Println(\"total txs:\", total)
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"
)

func main() {{
\t// BREAKPOINT
\trpcURL := \"{DEFAULT_RPC}\"

\tctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
\tdefer cancel()

\tc, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer c.Close()

\th, err := c.HeaderByNumber(ctx, nil)
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(\"latest:\", h.Number.Uint64())
}}
""",
    )

    files["18-reorgs"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"math/big\"
\t\"os\"
\t\"strconv\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"
)

{helper}

func main() {{
\t// geth/18-reorgs (continuity check demo)
\t//
\t// This demo walks back N headers and verifies parent-hash links.
\t// Real reorg handling requires persistence + backfill.
\t//
\t// Usage:
\t//   go run ./geth/18-reorgs/cmd/app <RPC_URL> [n]
\t//
\t// BREAKPOINT: parse args
\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")
\tn := 20
\tif len(os.Args) >= 3 {{
\t\tif v, err := strconv.Atoi(os.Args[2]); err == nil {{
\t\t\tn = v
\t\t}}
\t}}

\tctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
\tdefer cancel()

\tc, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer c.Close()

\thead, err := c.HeaderByNumber(ctx, nil)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"header:\", err)
\t\tos.Exit(1)
\t}}

\tprev := head
\tok := true
\tfor i := 0; i < n; i++ {{
\t\tnum := new(big.Int).Sub(prev.Number, big.NewInt(1))
\t\th, err := c.HeaderByNumber(ctx, num)
\t\tif err != nil {{
\t\t\tfmt.Fprintln(os.Stderr, \"header:\", num, err)
\t\t\tbreak
\t\t}}
\t\tif prev.ParentHash != h.Hash() {{
\t\t\tok = false
\t\t\tfmt.Println(\"MISMATCH at\", prev.Number.Uint64(), \"parent\", prev.ParentHash.Hex(), \"expected\", h.Hash().Hex())
\t\t\tbreak
\t\t}}
\t\tprev = h
\t}}

\tfmt.Println(\"continuity_ok:\", ok)
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"
)

func main() {{
\t// BREAKPOINT
\trpcURL := \"{DEFAULT_RPC}\"

\tctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
\tdefer cancel()

\tc, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer c.Close()

\th, err := c.HeaderByNumber(ctx, nil)
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(\"latest:\", h.Number.Uint64())
}}
""",
    )

    files["19-devnets"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"
)

{helper}

func main() {{
\t// geth/19-devnets
\t//
\t// This module is typically used with a local devnet (e.g. `geth --dev`).
\t// This cmd/app confirms connectivity + prints chain metadata.
\t//
\t// Usage:
\t//   go run ./geth/19-devnets/cmd/app <RPC_URL>
\t//
\t// Example (local dev chain):
\t//   geth --dev --http --http.api eth,net,web3,personal
\t//   go run ./geth/19-devnets/cmd/app http://127.0.0.1:8545
\t//
\t// BREAKPOINT
\trpcURL := rpcURLFromArgs(\"http://127.0.0.1:8545\")

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tc, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer c.Close()

\tchainID, _ := c.ChainID(ctx)
\tnetID, _ := c.NetworkID(ctx)
\thead, _ := c.HeaderByNumber(ctx, nil)

\tfmt.Println(\"ChainID:\", chainID)
\tfmt.Println(\"NetworkID:\", netID)
\tif head != nil {{
\t\tfmt.Println(\"Head:\", head.Number.Uint64())
\t}}
}}
""",
        """package main

import \"fmt\"

func main() {
\tfmt.Println(\"Run a local devnet: geth --dev --http --http.api eth,net,web3,personal\")
\tfmt.Println(\"Then run cmd/app with http://127.0.0.1:8545\")
\tfmt.Println(\"BREAKPOINT: read geth/19-devnets README\")
}
""",
    )

    files["20-node"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/rpc\"
)

{helper}

func main() {{
\t// geth/20-node
\t//
\t// This demo calls `web3_clientVersion` via raw JSON-RPC.
\t//
\t// Usage:
\t//   go run ./geth/20-node/cmd/app <RPC_URL>
\t//
\t// BREAKPOINT
\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tc, err := rpc.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer c.Close()

\tvar v string
\tif err := c.CallContext(ctx, &v, \"web3_clientVersion\"); err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"web3_clientVersion:\", err)
\t\tos.Exit(1)
\t}}
\tfmt.Println(\"clientVersion:\", v)
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/rpc\"
)

func main() {{
\t// BREAKPOINT
\trpcURL := \"{DEFAULT_RPC}\"

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tc, err := rpc.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer c.Close()

\tvar v string
\tif err := c.CallContext(ctx, &v, \"web3_clientVersion\"); err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(v)
}}
""",
    )

    # 21-sync
    files["21-sync"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/21-sync/internal/sync\"
)

{helper}

func main() {{
\t// geth/21-sync
\t//\n\t// Usage:\n\t//   go run ./geth/21-sync/cmd/app <RPC_URL>\n\t//\n\t// BREAKPOINT\n+\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer client.Close()

\tout, err := sync.Run(ctx, client, sync.Config{{}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}}

\tfmt.Println(\"IsSyncing:\", out.IsSyncing)
\tif out.Progress != nil {{
\t\tfmt.Printf(\"Progress: %+v\\n\", *out.Progress)
\t}}
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/21-sync/internal/sync\"
)

func main() {{
\t// BREAKPOINT
\trpcURL := \"{DEFAULT_RPC}\"

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer client.Close()

\tout, err := sync.Run(ctx, client, sync.Config{{}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(out.IsSyncing)
}}
""",
    )

    # 22-peers
    files["22-peers"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/22-peers/internal/peers\"
)

{helper}

func main() {{
\t// geth/22-peers\n\t//\n\t// Usage:\n\t//   go run ./geth/22-peers/cmd/app <RPC_URL>\n\t//\n\t// BREAKPOINT\n+\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer client.Close()

\tout, err := peers.Run(ctx, client, peers.Config{{}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}}

\tfmt.Println(\"PeerCount:\", out.PeerCount)
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/22-peers/internal/peers\"
)

func main() {{
\t// BREAKPOINT
\trpcURL := \"{DEFAULT_RPC}\"

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer client.Close()

\tout, err := peers.Run(ctx, client, peers.Config{{}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(out.PeerCount)
}}
""",
    )

    # 23-mempool
    files["23-mempool"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/23-mempool/internal/mempool\"
)

{helper}

func main() {{
\t// geth/23-mempool\n\t//\n\t// Usage:\n\t//   go run ./geth/23-mempool/cmd/app <RPC_URL>\n\t//\n\t// BREAKPOINT\n+\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer client.Close()

\tout, err := mempool.Run(ctx, client, mempool.Config{{Limit: 0}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}}

\tfmt.Println(\"PendingCount:\", out.PendingCount)
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/23-mempool/internal/mempool\"
)

func main() {{
\t// BREAKPOINT
\trpcURL := \"{DEFAULT_RPC}\"

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer client.Close()

\tout, err := mempool.Run(ctx, client, mempool.Config{{Limit: 0}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(out.PendingCount)
}}
""",
    )

    # 24-monitor
    files["24-monitor"] = (
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/24-monitor/internal/monitor\"
)

{helper}

func main() {{
\t// geth/24-monitor\n\t//\n\t// Usage:\n\t//   go run ./geth/24-monitor/cmd/app <RPC_URL>\n\t//\n\t// BREAKPOINT\n+\trpcURL := rpcURLFromArgs(\"{DEFAULT_RPC}\")

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer client.Close()

\tout, err := monitor.Run(ctx, client, monitor.Config{{MaxLagSeconds: 120}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}}

\tfmt.Println(\"Status:\", out.Status)
\tfmt.Println(\"BlockNumber:\", out.BlockNumber)
\tfmt.Println(\"LagSeconds:\", out.LagSeconds)
}}
""",
        f"""package main

import (
\t\"context\"
\t\"fmt\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/24-monitor/internal/monitor\"
)

func main() {{
\t// BREAKPOINT
\trpcURL := \"{DEFAULT_RPC}\"

\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer client.Close()

\tout, err := monitor.Run(ctx, client, monitor.Config{{MaxLagSeconds: 120}})
\tif err != nil {{
\t\tpanic(err)
\t}}

\tfmt.Println(out.Status, out.LagSeconds)
}}
""",
    )

    # 25-toolbox
    files["25-toolbox"] = (
        f"""package main

import (
\t\"context\"
\t\"encoding/json\"
\t\"fmt\"
\t\"os\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/25-toolbox/internal/toolbox\"
)

func main() {{
\t// geth/25-toolbox
\t//
\t// Usage:
\t//   go run ./geth/25-toolbox/cmd/app <RPC_URL> <command> [args...]
\t//
\t// Commands:
\t//   status
\t//   block <block_number>
\t//   tx <tx_hash>
\t//
\t// BREAKPOINT: parse args
\tif len(os.Args) < 3 {{
\t\tfmt.Fprintln(os.Stderr, \"usage: <RPC_URL> <command> [args...]\")
\t\tos.Exit(2)
\t}}

\trpcURL := os.Args[1]
\tcmd := os.Args[2]
\targs := []string{{}}
\tif len(os.Args) > 3 {{
\t\targs = os.Args[3:]
\t}}

\tctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"dial:\", err)
\t\tos.Exit(1)
\t}}
\tdefer client.Close()

\t// BREAKPOINT: run
\tout, err := toolbox.Run(ctx, client, toolbox.Config{{Command: cmd, Args: args}})
\tif err != nil {{
\t\tfmt.Fprintln(os.Stderr, \"run:\", err)
\t\tos.Exit(1)
\t}}

\tb, _ := json.MarshalIndent(out.Output, \"\", \"  \")
\tfmt.Println(string(b))
}}
""",
        f"""package main

import (
\t\"context\"
\t\"encoding/json\"
\t\"fmt\"
\t\"time\"

\t\"github.com/ethereum/go-ethereum/ethclient\"

\t\"github.com/example/go-10x-minis/geth/25-toolbox/internal/toolbox\"
)

func main() {{
\t// BREAKPOINT
\trpcURL := \"{DEFAULT_RPC}\"
\tcfg := toolbox.Config{{Command: \"status\"}}

\tctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
\tdefer cancel()

\tclient, err := ethclient.DialContext(ctx, rpcURL)
\tif err != nil {{
\t\tpanic(err)
\t}}
\tdefer client.Close()

\tout, err := toolbox.Run(ctx, client, cfg)
\tif err != nil {{
\t\tpanic(err)
\t}}

\tb, _ := json.MarshalIndent(out.Output, \"\", \"  \")
\tfmt.Println(string(b))
}}
""",
    )

    # Write everything we defined above.
    for mod, (app, dev) in files.items():
        mod_dir = next((p for p in GETH.iterdir() if p.is_dir() and p.name.endswith(mod)), None)
        if mod_dir is None:
            raise SystemExit(f"could not find module dir for {mod}")
        write(mod_dir / "cmd/app/main.go", app)
        write(mod_dir / "cmd/dev/main.go", dev)


if __name__ == "__main__":
    gen()

