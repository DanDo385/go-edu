package main

import "fmt"

func main() {
	fmt.Println("Run a local devnet: geth --dev --http --http.api eth,net,web3,personal")
	fmt.Println("Then run cmd/app with http://127.0.0.1:8545")
	fmt.Println("BREAKPOINT: read geth/19-devnets README")
}
