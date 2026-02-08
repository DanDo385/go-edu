//go:build reference

package reorgs

/*
Reference Solution - Chain Reorg Handling
=========================================

This file demonstrates reorg logic: when a longer chain arrives, find common
ancestor, prune canonical back to it, append incoming. Simulates blockchain
consensus reorg handling.

This connects to the Ethereum ecosystem by showing:
- Reorg: canonical chain replaced by chain with higher total difficulty
- Common ancestor: incoming[0].Parent must exist in canonical
- canonical = canonical[:len(canonical)-1]: prune from tip until we hit ancestor

The exercise builds understanding of:
- Slice truncation: canonical[:len(canonical)-1] removes last element
- append(canonical, incoming...): append all incoming blocks
- Loop invariant: we stop when canonical tip hash matches incoming's parent

Teaching notes (per .cursorrules):
- for len(canonical) > 0 && canonical[len(canonical)-1].Hash != incoming[0].Parent:
  prune while tip isn't the ancestor. len(canonical) > 0 prevents index out of range.
*/
func Run() {
	type head struct {
		Hash   string
		Parent string
	}

	canonical := []head{{Hash: "A", Parent: ""}, {Hash: "B", Parent: "A"}, {Hash: "C", Parent: "B"}}
	incoming := []head{{Hash: "X", Parent: "B"}, {Hash: "Y", Parent: "X"}}

	tip := canonical[len(canonical)-1]
	if incoming[0].Parent != tip.Hash {
		for len(canonical) > 0 && canonical[len(canonical)-1].Hash != incoming[0].Parent {
			canonical = canonical[:len(canonical)-1]
		}
	}
	canonical = append(canonical, incoming...)
	_ = canonical
}
