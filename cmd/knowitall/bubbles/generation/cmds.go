package generation

import (
	"errors"
	"sudonters/libzootr/cmd/knowitall/bubbles/spheres"
	"sudonters/libzootr/components"
	"sudonters/libzootr/magicbean"
	"sudonters/libzootr/magicbean/tracking"
	"sudonters/libzootr/mido/compiler"
	"sudonters/libzootr/mido/vm"
	"sudonters/libzootr/zecs"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/etc-sudonters/substrate/skelly/bitset32"
	"github.com/etc-sudonters/substrate/slipup"
)

type SearchResult spheres.Details

type lowDetailResult struct {
	adult, child nodes
	collected    magicbean.Inventory
	err          error
}

var ErrNoProgress = errors.New("reached 0 new locations")

type nodes struct {
	Visited, Pending, Reached bitset32.Bitset

	Edges []magicbean.EdgeHandle
}

type edge struct {
	crossed bool
	name    components.Name
	id      zecs.Entity
}

func (this nodes) visited(id zecs.Entity) bool {
	return bitset32.IsSet(&this.Visited, id)
}
func (this nodes) pending(id zecs.Entity) bool {
	return bitset32.IsSet(&this.Pending, id)
}
func (this nodes) reached(id zecs.Entity) bool {
	return bitset32.IsSet(&this.Reached, id)
}

type discache map[zecs.Entity]spheres.Disassembly

func disassemble(gen *magicbean.Generation, edge zecs.Entity, cache discache) tea.Cmd {
	return func() tea.Msg {
		dis, exists := cache[edge]
		if exists {
			return dis
		}
		ocm := &gen.Ocm
		values, err := ocm.GetValues(edge,
			zecs.Get[components.Name], zecs.Get[components.RuleCompiled],
			zecs.Get[components.RuleSource], zecs.Get[components.RuleParsed],
			zecs.Get[components.RuleOptimized],
		)
		dis.Id = edge
		dis.Err = err
		dis.Values = values
		dis.Name, _ = values[0].(components.Name)
		dis.Code, _ = values[1].(components.RuleCompiled)
		dis.Src, _ = values[2].(components.RuleSource)
		dis.Ast, _ = values[3].(components.RuleParsed)
		dis.Opt, _ = values[4].(components.RuleOptimized)
		dis.Dis = vm.Disassemble(compiler.Bytecode(dis.Code), &gen.Objects)
		if dis.Err == nil {
			cache[edge] = dis
		}
		return dis
	}
}

var searchExclusion sync.Mutex

func runSearch(gen *magicbean.Generation, searches searches, names tracking.NameTable) tea.Cmd {
	return func() tea.Msg {
		searchExclusion.Lock()
		defer searchExclusion.Unlock()

		var result SearchResult
		initial := doSearch(gen, searches)
		result.Err = initial.err

		moveNodes(&initial, &result, names)
		moveEdges(&initial, &result, names)
		moveItems(&initial, &result, names)

		return result
	}
}

func moveItems(initial *lowDetailResult, result *SearchResult, names tracking.NameTable) {
	result.Tokens = make([]spheres.NamedToken, 0, len(initial.collected))
	for entity, qty := range initial.collected {
		result.Tokens = append(result.Tokens, spheres.NamedToken{
			Name: names[entity],
			Id:   entity,
			Qty:  qty,
		})
	}
}

func moveNodes(initial *lowDetailResult, result *SearchResult, names tracking.NameTable) {
	visited := initial.adult.Visited.Union(initial.child.Visited)
	result.Nodes = make([]spheres.NamedNode, 0, visited.Len())
	for entity := range bitset32.IterT[zecs.Entity](&visited).All {
		named := spheres.NamedNode{
			Name: names[entity],
			Id:   entity,
		}
		ptr := uint32(len(result.Nodes))
		result.Nodes = append(result.Nodes, named)

		if initial.adult.visited(entity) {
			bitset32.Set(&result.Adult.Visited, ptr)
		}
		if initial.adult.pending(entity) {
			bitset32.Set(&result.Adult.Pending, ptr)
		}
		if initial.adult.reached(entity) {
			bitset32.Set(&result.Adult.Reached, ptr)
		}
		if initial.child.visited(entity) {
			bitset32.Set(&result.Child.Visited, ptr)
		}
		if initial.child.pending(entity) {
			bitset32.Set(&result.Child.Pending, ptr)
		}
		if initial.child.reached(entity) {
			bitset32.Set(&result.Child.Reached, ptr)
		}
	}
}

func uniqueEdges(initial *lowDetailResult, names tracking.NameTable) map[zecs.Entity]spheres.NamedEdge {
	uniqueEdges := make(map[zecs.Entity]spheres.NamedEdge, len(initial.adult.Edges)+len(initial.child.Edges))
	for _, edges := range [][]magicbean.EdgeHandle{initial.adult.Edges, initial.child.Edges} {
		for _, edge := range edges {
			if _, exists := uniqueEdges[edge.Id]; exists {
				continue
			}
			named := spheres.NamedEdge{
				Id:   edge.Id,
				Name: names[edge.Id],
			}
			uniqueEdges[edge.Id] = named
		}
	}

	return uniqueEdges
}

func moveEdgeSet(
	src []magicbean.EdgeHandle,
	dest []spheres.VisitedEdge,
	ptrs map[zecs.Entity]int,
	visited *bitset32.Bitset,
	origins bitset32.Bitset,
) {

	for _, edge := range src {
		if !bitset32.IsSet(&origins, edge.Def.From) {
			continue
		}

		visited := spheres.VisitedEdge{
			Index:   ptrs[edge.Id],
			Crossed: bitset32.IsSet(visited, edge.Def.To),
		}

		dest = append(dest, visited)
	}
}

func moveEdges(initial *lowDetailResult, result *SearchResult, names tracking.NameTable) {
	uniqueEdges := uniqueEdges(initial, names)

	result.Edges = make([]spheres.NamedEdge, 0, len(uniqueEdges))
	edgePtrs := make(map[zecs.Entity]int, len(uniqueEdges))
	for _, edge := range uniqueEdges {
		edgePtrs[edge.Id] = len(result.Edges)
		result.Edges = append(result.Edges, edge)
	}

	result.Adult.Edges = make([]spheres.VisitedEdge, 0, len(initial.adult.Edges))
	result.Child.Edges = make([]spheres.VisitedEdge, 0, len(initial.child.Edges))

	moveEdgeSet(initial.adult.Edges, result.Adult.Edges, edgePtrs, &initial.adult.Visited, result.Adult.Pending.Union(result.Adult.Reached))
	moveEdgeSet(initial.child.Edges, result.Child.Edges, edgePtrs, &initial.child.Visited, result.Child.Pending.Union(result.Child.Reached))
}

func doSearch(gen *magicbean.Generation, searches searches) lowDetailResult {
	var result lowDetailResult
	reached := visitAll(searches, &result)

	if reached.Len() == 0 {
		result.err = ErrNoProgress
		return result
	}
	result.collected, result.err = collect(reached, gen)
	return result
}

func visitAll(searches searches, result *lowDetailResult) bitset32.Bitset {
	visitOne(searches[magicbean.AgeAdult], &result.adult)
	visitOne(searches[magicbean.AgeChild], &result.child)
	return result.adult.Reached.Union(result.child.Reached)
}

func visitOne(search *magicbean.Search, nodes *nodes) {
	result := search.Visit()
	nodes.Reached = result.Reached
	nodes.Pending = bitset32.Copy(search.Pending)
	nodes.Visited = bitset32.Copy(search.Visited)
	nodes.Edges = result.Edges
}

func collect(reached bitset32.Bitset, gen *magicbean.Generation) (magicbean.Inventory, error) {
	precollect := magicbean.CopyInventory(gen.Inventory)
	collectErr := magicbean.CollectTokensFrom(&gen.Ocm, reached, gen.Inventory)
	if collectErr != nil {
		return nil, slipup.Describe(collectErr, "while collecting tokens")
	}
	return magicbean.DiffInventories(precollect, gen.Inventory), nil
}
