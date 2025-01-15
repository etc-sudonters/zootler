package magicbean

import (
	"errors"
	"fmt"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/skelly/bitset"
	"sudonters/zootler/internal/table"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/mido/vm"

	"github.com/etc-sudonters/substrate/skelly/graph"
)

var (
	ErrWorksetEmpty   = errors.New("workset empty")
	ErrNoProgressMade = errors.New("no progress made")
)

type Exploration struct {
	age     age
	nodes   *nodes
	world   *World
	objects objects.Table
}

func (this *Exploration) ExploreAccessible() (bitset.Bitset32, error) {
	var reached bitset.Bitset32
	var workset bitset.Bitset32

	if this.nodes.workset.IsEmpty() {
		return reached, ErrWorksetEmpty
	}

	vm := vm.VM{Objects: &this.objects}
	physical := this.world.Physical
	visited := &this.nodes.visited
	workset, this.nodes.workset = this.nodes.workset, bitset.Bitset32{}

	for visiting := range nodebiter(&workset).All {
		neighbors := successors(visiting, physical)
		for neighbor := range nodebiter(&neighbors).All {
			if bitset.IsSet32(visited, neighbor) {
				bitset.Unset32(&neighbors, neighbor)
				continue
			}

			source, sourceErr := this.getTransitRule(visiting, neighbor)
			if sourceErr != nil {
				return reached, sourceErr
			}

			obj, execErr := vm.Execute(source.ByteCode)
			if execErr != nil {
				return reached, execErr
			}
			if obj == nil {
				return reached, fmt.Errorf("%s: %s produced nil result", source.OriginatingRegion, source.Destination)
			}

			result, isBool := obj.(objects.Boolean)
			if !isBool {
				return reached, fmt.Errorf("%s: %s produced non boolean result %v", source.OriginatingRegion, source.Destination, obj)
			}

			if result {
				bitset.Unset32(&neighbors, neighbor)
				bitset.Set32(&this.nodes.workset, neighbor)
				bitset.Set32(&this.nodes.visited, neighbor)
				bitset.Set32(&reached, neighbor)
			}
		}

		if !neighbors.IsEmpty() {
			bitset.Set32(&this.nodes.workset, visiting)
		}
	}

	return reached, nil
}

func nodebiter(set *bitset.Bitset32) bitset.IterOf32[graph.Node] {
	return bitset.Iter64T[graph.Node](set)
}

func (this *Exploration) getTransitRule(origin, destination graph.Node) (mido.CompiledSource, error) {
	panic("not implemented")
}

func successors(node graph.Node, physical *graph.Directed) bitset.Bitset32 {
	neighbors, err := physical.Successors(node)
	if err != nil {
		panic(err)
	}
	return bitset.Upgrade(neighbors)

}

func (this *Exploration) init(world *World) {
	this.world = world
	this.nodes = &this.world.nodes[this.age]

	this.objects = objects.NewTable(
		objects.CloneTableFrom(world.objects),
		objects.TableWithBuiltIns(createBuiltins(this)),
	)
}

func createBuiltins(this *Exploration) objects.BuiltInFunctions {
	lookup := map[string]objects.BuiltInFunction{
		"has":                 {Name: "has", Params: 2, Fn: this.has},
		"has_anyof":           {Name: "has_anyof", Params: -1, Fn: this.has_anyof},
		"has_bottle":          {Name: "has_bottle", Params: 0, Fn: this.has_bottle},
		"has_dungeon_rewards": {Name: "has_dungeon_rewards", Params: 1, Fn: this.has_dungeon_rewards},
		"has_every":           {Name: "has_every", Params: -1, Fn: this.has_every},
		"has_hearts":          {Name: "has_hearts", Params: 1, Fn: this.has_hearts},
		"has_medallions":      {Name: "has_medallions", Params: 1, Fn: this.has_medallions},
		"has_stones":          {Name: "has_stones", Params: 1, Fn: this.has_stones},
		"is_adult":            {Name: "is_adult", Params: 0, Fn: constBuiltIn(this.age == EXPLORE_AS_ADULT).Call},
		"is_child":            {Name: "is_child", Params: 0, Fn: constBuiltIn(this.age == EXPLORE_AS_CHILD).Call},
		"is_starting_age":     {Name: "is_starting_age", Params: 0, Fn: constBuiltIn(this.age == this.world.startAge).Call},
	}

	table := make(objects.BuiltInFunctions, len(this.world.builtins))
	for name, i := range this.world.builtins {
		def, exists := lookup[name]
		if !exists {
			panic(fmt.Errorf("built-in %q declared but not provided", name))
		}
		table[i] = def
	}

	return table
}

func (this Exploration) QuantityFor(id CollectionId) uint8 {
	return this.world.collected[id]
}

func (this Exploration) HasBottle() bool {
	store := this.world.DataStore
	q := store.CreateQuery()
	q.Load(query.MustAsColumnId[CollectionId](store))
	q.Exists(query.MustAsColumnId[components.IsBottle](store))

	bottles, err := store.Retrieve(q)
	if err != nil {
		panic(err)
	}

	for _, bottle := range bottles.All {
		id := bottle.Values[0].(CollectionId)
		if this.world.collected[id] > 0 {
			return true
		}
	}

	return false
}

func (this Exploration) HasDungeonRewards(count int) bool {
	total := this.sumCollectionFor(query.MustAsColumnId[components.DungeonReward](this.world.DataStore))
	return total >= count
}

func (this Exploration) HasHearts(count int) bool {
	store := this.world.DataStore
	q := store.CreateQuery()
	q.Load(query.MustAsColumnId[CollectionId](store))
	q.Load(query.MustAsColumnId[components.PieceOfHeart](store))
	hearts, err := store.Retrieve(q)
	if err != nil {
		panic(err)
	}

	var total int
	for _, heart := range hearts.All {
		id := heart.Values[0].(CollectionId)
		mul := heart.Values[1].(components.PieceOfHeart)

		total += int(mul) * int(this.world.collected[id])
	}

	return total >= count
}

func (this Exploration) HasMedallions(count int) bool {
	total := this.sumCollectionFor(query.MustAsColumnId[components.Medallion](this.world.DataStore))
	return total >= count
}

func (this Exploration) HasStones(count int) bool {
	total := this.sumCollectionFor(query.MustAsColumnId[components.SpiritualStone](this.world.DataStore))
	return total >= count
}

func (this Exploration) sumCollectionFor(kind table.ColumnId) int {
	store := this.world.DataStore
	q := store.CreateQuery()
	q.Load(query.MustAsColumnId[CollectionId](store))
	q.Exists(kind)
	rows, err := store.Retrieve(q)
	if err != nil {
		panic(err)
	}

	var sum int
	for _, token := range rows.All {
		id := token.Values[0].(CollectionId)
		sum += int(this.world.collected[id])
	}
	return sum
}

func (this *Exploration) has(args []objects.Object) (objects.Object, error) {
	ptr, isPtr := args[0].(objects.Ptr)
	qty, isQty := args[1].(objects.Number)

	switch {
	case !isPtr:
		return nil, fmt.Errorf("expected arg 0 to be pointer, got %T", args[0])
	case !isQty:
		return nil, fmt.Errorf("expected arg 1 to be number, got %T", args[1])
	default:
		id := this.world.Translation.FromPtr(ptr)
		collected := this.QuantityFor(id)
		return objects.Boolean(float64(collected) >= float64(qty)), nil
	}
}

func (this *Exploration) has_anyof(args []objects.Object) (objects.Object, error) {
	for i, arg := range args {
		ptr, isPtr := arg.(objects.Ptr)
		if !isPtr {
			return nil, fmt.Errorf("expected arg %d to be pointer, got %T", i, arg)
		}
		id := this.world.Translation.FromPtr(ptr)
		if this.QuantityFor(id) >= 1 {
			return objects.Boolean(true), nil
		}
	}

	return objects.Boolean(false), nil
}

func (this *Exploration) has_bottle(args []objects.Object) (objects.Object, error) {
	return objects.Boolean(this.HasBottle()), nil
}

func (this *Exploration) has_dungeon_rewards(args []objects.Object) (objects.Object, error) {
	qty, isQty := args[0].(objects.Number)
	if !isQty {
		return nil, fmt.Errorf("expected arg 0 to be number, got %T", args[0])
	}
	return objects.Boolean(this.HasDungeonRewards(int(qty))), nil
}

func (this *Exploration) has_every(args []objects.Object) (objects.Object, error) {
	for i, arg := range args {
		ptr, isPtr := arg.(objects.Ptr)
		if !isPtr {
			return nil, fmt.Errorf("expected arg %d to be pointer, got %T", i, arg)
		}
		id := this.world.Translation.FromPtr(ptr)
		if this.QuantityFor(id) < 1 {
			return objects.Boolean(false), nil
		}
	}

	return objects.Boolean(true), nil
}

func (this *Exploration) has_hearts(args []objects.Object) (objects.Object, error) {
	qty, isQty := args[0].(objects.Number)
	if !isQty {
		return nil, fmt.Errorf("expected arg 0 to be number, got %T", args[0])
	}
	return objects.Boolean(this.HasHearts(int(qty))), nil
}

func (this *Exploration) has_medallions(args []objects.Object) (objects.Object, error) {
	qty, isQty := args[0].(objects.Number)
	if !isQty {
		return nil, fmt.Errorf("expected arg 0 to be number, got %T", args[0])
	}
	return objects.Boolean(this.HasMedallions(int(qty))), nil
}

func (this *Exploration) has_stones(args []objects.Object) (objects.Object, error) {
	qty, isQty := args[0].(objects.Number)
	if !isQty {
		return nil, fmt.Errorf("expected arg 0 to be number, got %T", args[0])
	}
	return objects.Boolean(this.HasStones(int(qty))), nil
}

type constBuiltIn bool

func (this constBuiltIn) Call([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
