//go:build reference

package reorgs

/*
Reference Solution

This placeholder lesson currently has no config/result types. The reference
captures the core reorg rule: if a new head's parent is not the current tip,
rewind until ancestry reconnects and then apply the new branch.
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
