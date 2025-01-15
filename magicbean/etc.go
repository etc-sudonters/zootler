package magicbean

import (
	"github.com/etc-sudonters/substrate/skelly/graph"
	"sudonters/zootler/internal/skelly/bitset"
)

type nodes struct {
	visited bitset.Bitset32
	// nodes whose edges we haven't completely explored yet
	workset bitset.Bitset32
}

func newtracker(roots []graph.Node) nodes {
	var track nodes
	if len(roots) == 0 {
		panic("no root nodes declared in physical world")
	}

	for _, root := range roots {
		bitset.Set32(&track.workset, root)
		bitset.Set32(&track.visited, root)
	}

	return track
}
