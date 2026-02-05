//go:build !solution && !reference

package stack

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

/*
Problem: Prove RPC connectivity by reading the network identifiers and latest header.
*/

// Run - TODO: implement this function
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) { // get network id
	networkID, err := client.NetworkID(ctx) // get the network ID
	if err != nil { // if the network ID is not found, return an error
		return nil, fmt.Errorf("failed to get network ID: %w", err) // return the error
	}
	fmt.Printf("Network ID: %s\n", networkID) // print the network ID

	// get latest header
	header, err := client.HeaderByNumber(ctx, cfg.BlockNumber) // get the latest header
	if err != nil { // if the latest header is not found, return an error
		return nil, fmt.Errorf("failed to get latest header: %w", err) // return the error
	} 
	fmt.Printf("Latest Header: %s\n", header) // print the latest header

	return &Result{NetworkID: networkID, Header: header}, nil // return the result
}