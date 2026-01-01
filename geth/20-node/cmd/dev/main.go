package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	// BREAKPOINT
	rpcURL := "https://eth.llamarpc.com"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := rpc.DialContext(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	var v string
	if err := c.CallContext(ctx, &v, "web3_clientVersion"); err != nil {
		panic(err)
	}

	fmt.Println(v)
}
