//go:build reference

package devnets

/*
Reference Solution

With no typed exercise surface yet, this reference demonstrates a minimal
"devnet profile" selection workflow used in practice.
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
