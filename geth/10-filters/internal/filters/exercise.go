//go:build !solution && !reference


package filters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const defaultMaxHeads = 5
const defaultPollInterval = time.Second

/*
Problem: Monitor new block headers using subscriptions (WebSocket) or polling (HTTP).

This module teaches you about real-time vs polling approaches for monitoring blockchain
state. You'll implement both WebSocket subscriptions (push) and HTTP polling (pull),
and learn how to detect chain reorganizations.

Computer science principles highlighted:
  - Push vs Pull architecture: Real-time subscriptions vs periodic polling
  - Event-driven programming: Handling asynchronous header updates
  - Reorg detection: Comparing parent hashes to detect chain changes
  - Resource management: Proper cleanup of subscriptions and channels
*/
func Run(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	// TODO: Implement Run
	// See solution.reference.go for reference implementation
	panic("not implemented")
}


func subscribeHeads(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	// TODO: Create channel for receiving headers
	// - make(chan *types.Header) creates an unbuffered channel
	// - Why channel? WebSocket pushes headers asynchronously, channel receives them
	// - Unbuffered channel: Sender blocks until receiver reads (backpressure)

	// TODO: Subscribe to new headers
	// - Call client.SubscribeNewHead(ctx, headCh) to start subscription
	// - Returns Subscription object for managing the subscription
	// - Handle errors (network failures, WebSocket not supported)
	// - Subscription pushes new headers to headCh as they arrive

	// TODO: Ensure subscription cleanup
	// - Use defer sub.Unsubscribe() to clean up when function returns
	// - Why defer? Guarantees cleanup even if we return early or panic
	// - Unsubscribe closes WebSocket connection and stops pushing headers
	// - This is Go's resource management pattern: acquire, defer cleanup, use

	// TODO: Initialize result struct
	// - Create Result with preallocated Heads slice (capacity cfg.MaxHeads)
	// - Set Mode to "subscription" to indicate WebSocket was used
	// - Initialize prevHash to track previous block for reorg detection

	// TODO: Collect headers until we have MaxHeads
	// - Loop while len(result.Heads) < cfg.MaxHeads
	// - Use select statement to handle multiple channels:
	//   * ctx.Done(): Context canceled, return error
	//   * sub.Err(): Subscription error, return error
	//   * headCh: New header received, process it
	// - Why select? Go's way of multiplexing multiple channel operations
	// - This is event-driven programming: react to whichever event happens first

	// TODO: Process each header
	// - Skip nil headers (shouldn't happen, but defensive)
	// - Compute header hash
	// - Detect reorg: if prevHash != 0 and head.ParentHash != prevHash
	// - Create HeadInfo with number, hash, parentHash, reorg flag
	// - Append to result.Heads
	// - Update prevHash for next iteration

	// TODO: Return result
	// - Return Result with all collected headers
	// - Return nil error on success

	return nil, errors.New("not implemented")
}

func pollHeads(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	// TODO: Initialize result and tracking variables
	// - Create Result with preallocated Heads slice
	// - Set Mode to "polling" to indicate HTTP was used
	// - Initialize prevHash and prevNumber for reorg detection and deduplication

	// TODO: Poll for new headers in a loop
	// - Loop while len(result.Heads) < cfg.MaxHeads
	// - Each iteration queries latest header
	// - Continue until we have enough unique headers

	// TODO: Query latest header
	// - Call client.HeaderByNumber(ctx, nil) for latest header
	// - nil means "latest", specific number means historical
	// - Handle errors (network failures, node not synced)
	// - Validate header is not nil

	// TODO: Extract header information
	// - Get block number from header
	// - Compute block hash
	// - These identify the block uniquely

	// TODO: Handle duplicate headers (no new block yet)
	// - If number == prevNumber, we've seen this block already
	// - Wait for cfg.PollInterval before next poll
	// - Use select with ctx.Done() to allow cancellation during wait
	// - Then continue to next iteration
	// - Why wait? Avoid spamming the node with requests when no new blocks

	// TODO: Detect reorgs
	// - Compare head.ParentHash with prevHash
	// - If they don't match (and prevHash != zero), a reorg occurred
	// - Reorg means the chain changed; previous block is no longer canonical

	// TODO: Record header information
	// - Create HeadInfo with number, hash, parentHash, reorg flag
	// - Append to result.Heads
	// - Update prevHash and prevNumber for next iteration

	// TODO: Return result
	// - Return Result with all collected headers
	// - Headers are in chronological order

	return nil, errors.New("not implemented")
}
