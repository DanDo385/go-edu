# 06: EIP-1559 Fee Market

## What Is This Project About?

This module teaches you the EIP-1559 transaction fee mechanism introduced in the London hard fork. You'll understand base fees, priority fees (tips), max fees, and how the fee market creates more predictable gas costs while burning a portion of fees.

## Why Is This Important?

EIP-1559 fundamentally changed how transaction fees work:
- More predictable gas pricing
- Better user experience with fee estimation
- ETH burning mechanism affects tokenomics
- Transaction replacement strategies changed

## Real-World Problems This Solves

- **Gas estimation**: Calculate appropriate max and priority fees
- **Fee optimization**: Minimize costs while ensuring timely inclusion
- **Transaction management**: Understand why txs get stuck or replaced
- **Economic analysis**: Track base fee trends and ETH burning

## Key Concepts You'll Learn

- **Base fee**: Protocol-determined fee, burned on inclusion
- **Priority fee (tip)**: Incentive for validators to include tx
- **Max fee**: Maximum total fee you're willing to pay
- **Fee calculation**: Effective gas price = min(base_fee + priority_fee, max_fee)

## Prerequisites

- Completion of `geth/01-stack` through `geth/05-tx-nonces`

## How to Run

```bash
go run ./cmd/app/main.go https://eth.llamarpc.com
```

## Testing

```bash
go test -v ./...
```

## Additional Resources

- [EIP-1559: Fee Market Change](https://eips.ethereum.org/EIPS/eip-1559)
