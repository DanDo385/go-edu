//go:build reference

package indexer

/*
Reference Solution

This lesson currently exposes only a `Run()` entrypoint with no typed inputs or
outputs. The reference still demonstrates the expected indexing shape:
- scan canonical blocks in order,
- maintain hash->height and address activity indexes,
- keep writes deterministic.
*/
func Run() {
	type block struct {
		Height uint64
		Hash   string
		TxFrom string
		TxTo   string
	}

	chain := []block{
		{Height: 100, Hash: "0xaaa", TxFrom: "0x1", TxTo: "0x2"},
		{Height: 101, Hash: "0xbbb", TxFrom: "0x2", TxTo: "0x3"},
	}

	indexByHash := make(map[string]uint64, len(chain))
	activity := make(map[string][]uint64)

	for _, b := range chain {
		indexByHash[b.Hash] = b.Height
		activity[b.TxFrom] = append(activity[b.TxFrom], b.Height)
		activity[b.TxTo] = append(activity[b.TxTo], b.Height)
	}

	_ = indexByHash
	_ = activity
}
