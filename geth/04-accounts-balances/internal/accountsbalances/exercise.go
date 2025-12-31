//go:build !solution
// +build !solution

package accountsbalances

import (
	"context"
	"errors"
)

/*
Problem: Query the state of multiple Ethereum accounts and classify them as
either Externally Owned Accounts (EOAs) or Smart Contracts.

This module dives into the fundamental difference between the two account types
on Ethereum. You'll learn how to query the balance and code of an account, and
use that information to determine its type. This is a core skill for building
any application that interacts with Ethereum accounts.

Computer science principles highlighted:
  - Type systems (classifying accounts based on their properties)
  - State machines (querying the state of an account at a specific block)
  - Defensive programming (handling nil inputs, copying reference types)
*/
func Run(ctx context.Context, client AccountClient, cfg Config) (*Result, error) {
	// TODO: Implement Run
	// See solution.reference.go for reference implementation
	panic("not implemented")
}

