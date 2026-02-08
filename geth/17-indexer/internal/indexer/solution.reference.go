//go:build reference

package indexer

/*
Reference Solution - Chain Indexing (Hash->Height, Activity)
============================================================

This file demonstrates in-memory chain indexing: build hash→height map and
address→blockHeights activity map. Foundation for block explorers and analytics.

This connects to the Ethereum ecosystem by showing:
- indexByHash: block hash → height for O(1) lookup
- activity: address → []blockHeights for "blocks this address participated in"
- Per block: index hash; append height to both TxFrom and TxTo activity lists

The exercise builds understanding of:
- Map growth: activity[addr] = append(activity[addr], h) — slice grows dynamically
- make(map[string][]uint64): map to slices; nil slice append works (allocates)

Teaching notes (per .cursorrules):
- activity[b.TxFrom] = append(activity[b.TxFrom], b.Height): when key absent,
  activity[b.TxFrom] is nil; append(nil, x) returns new slice. Map stores it.
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
