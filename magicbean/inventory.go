package magicbean

import (
	"sudonters/zootler/zecs"
)

func NewInventory() Inventory {
	return Inventory{make(map[zecs.Entity]float64)}
}

type Inventory struct {
	onhand map[zecs.Entity]float64
}

func (this *Inventory) CollectOne(entity zecs.Entity) {
	this.Collect(entity, 1)
}

func (this *Inventory) Collect(entity zecs.Entity, n float64) {
	has := this.onhand[entity]
	this.onhand[entity] = has + n
}

func (this *Inventory) Remove(entity zecs.Entity, n float64) float64 {
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

func (this *Inventory) Count(entity zecs.Entity) float64 {
	return this.onhand[entity]
}

func (this *Inventory) Sum(entities []zecs.Entity) float64 {
	var total float64

	for _, entity := range entities {
		total += this.onhand[entity]
	}

	return total
}
