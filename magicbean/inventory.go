package magicbean

import (
	"sudonters/libzootr/zecs"
)

func EmptyInventory() Inventory {
	return Inventory{make(map[zecs.Entity]uint32)}
}

type Inventory struct {
	onhand map[zecs.Entity]uint32
}

func (this *Inventory) CollectOne(entity zecs.Entity) {
	this.Collect(entity, 1)
}

func (this *Inventory) Collect(entity zecs.Entity, n uint32) {
	has := this.onhand[entity]
	this.onhand[entity] = has + n
}

func (this *Inventory) Remove(entity zecs.Entity, n uint32) uint32 {
	has := this.onhand[entity]

	switch {
	case has == 0:
		return 0
	case n == has:
		this.onhand[entity] = 0
		return n
	case n < has:
		this.onhand[entity] = has - n
		return n
	case n > has:
		this.onhand[entity] = 0
		return has
	default:
		panic("unreachable")
	}
}

func (this *Inventory) Count(entity zecs.Entity) uint32 {
	return this.onhand[entity]
}

func (this *Inventory) Sum(entities []zecs.Entity) uint32 {
	var total uint32

	for _, entity := range entities {
		total += this.onhand[entity]
	}

	return total
}
