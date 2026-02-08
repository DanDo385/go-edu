//go:build reference

package devnets

/*
Reference Solution - Devnet Profiles (ChainID, RPC, WS)
========================================================

This file demonstrates devnet profile configuration: map of named profiles
with ChainID, RPC/WS URLs, block time. Used for switching local/dev/testnet.

This connects to the Ethereum ecosystem by showing:
- local/anvil: common dev chain IDs (1337, 31337)
- RPC vs WS: HTTP for one-off calls, WebSocket for subscriptions
- profiles["local"]: map lookup; selected is a copy of the struct

The exercise builds understanding of:
- Map of structs: value is copied on read; selected is independent
- profiles["local"]: returns profile value; modifying selected doesn't affect map
*/
func Run() {
	type profile struct {
		Name      string
		ChainID   uint64
		RPC       string
		WS        string
		BlockTime int
	}

	profiles := map[string]profile{
		"local": {Name: "local", ChainID: 1337, RPC: "http://127.0.0.1:8545", WS: "ws://127.0.0.1:8546", BlockTime: 1},
		"anvil": {Name: "anvil", ChainID: 31337, RPC: "http://127.0.0.1:8545", WS: "ws://127.0.0.1:8545", BlockTime: 1},
	}

	selected := profiles["local"]
	_ = selected
}
