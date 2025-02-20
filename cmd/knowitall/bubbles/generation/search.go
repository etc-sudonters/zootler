package generation

import (
	"errors"
	"sudonters/libzootr/cmd/knowitall/bubbles/explore"
	"sudonters/libzootr/magicbean/tracking"
	"sudonters/libzootr/playthrough"
	"sudonters/libzootr/zecs"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/etc-sudonters/substrate/skelly/bitset32"
)

func runSearch(msg explore.ExploreSphere, searches Searches, names tracking.NameTable) tea.Cmd {
	return func() tea.Msg {
		spheres := []playthrough.SearchSphere{
			searches.Adult.Explore(),
			searches.Child.Explore(),
		}
		adult := spheres[0]
		child := spheres[1]

		edges := adult.Edges.All().Union(child.Edges.All())
		nodes := adult.Nodes.All().Union(child.Nodes.All())
		named := explore.NamedSphere{
			Edges: make([]explore.NamedEdge, 0, edges.Len()),
			Nodes: make([]explore.NamedNode, 0, nodes.Len()),
		}
		for edge := range bitset32.IterT[zecs.Entity](&edges).UntilEmpty {
			name := names[edge]
			index := uint32(len(named.Edges))
			named.Edges = append(named.Edges, explore.NamedEdge{edge, name})

			if bitset32.IsSet(&adult.Edges.Crossed, edge) {
				bitset32.Set(&named.Adult.Edges.Crossed, index)
			}
			if bitset32.IsSet(&adult.Edges.Pended, edge) {
				bitset32.Set(&named.Adult.Edges.Pended, index)
			}
			if bitset32.IsSet(&child.Edges.Crossed, edge) {
				bitset32.Set(&named.Child.Edges.Crossed, index)
			}
			if bitset32.IsSet(&child.Edges.Pended, edge) {
				bitset32.Set(&named.Child.Edges.Pended, index)
			}
		}

		for node := range bitset32.IterT[zecs.Entity](&nodes).UntilEmpty {
			name := names[node]
			index := uint32(len(named.Nodes))
			named.Nodes = append(named.Nodes, explore.NamedNode{node, name})

			if bitset32.IsSet(&adult.Nodes.Reached, node) {
				bitset32.Set(&named.Adult.Nodes.Reached, index)
			}
			if bitset32.IsSet(&adult.Nodes.Pended, node) {
				bitset32.Set(&named.Adult.Nodes.Pended, index)
			}
			if bitset32.IsSet(&child.Nodes.Reached, node) {
				bitset32.Set(&named.Child.Nodes.Reached, index)
			}
			if bitset32.IsSet(&child.Nodes.Pended, node) {
				bitset32.Set(&named.Child.Nodes.Pended, index)
			}
		}

		msg := explore.SphereExplored{
			Sphere: named,
		}

		reached := child.Nodes.Reached.Len() + adult.Nodes.Reached.Len()
		if reached == 0 {
			msg.Err = errors.New("no progress made")
		}

		return msg
	}
}
