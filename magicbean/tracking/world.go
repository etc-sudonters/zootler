package tracking

import (
	"sudonters/zootler/magicbean"
	"sudonters/zootler/zecs"
)

type namedents = zecs.Tracked[magicbean.Name]
type directed = magicbean.Connection
type name = magicbean.Name

var namef = magicbean.NameF

func named[T zecs.Value](ocm *zecs.Ocm) namedents {
	return zecs.Tracking[name](ocm, zecs.With[T])
}

func NewNodes(ocm *zecs.Ocm) Nodes {
	var this Nodes
	this.regions = named[magicbean.Region](ocm)
	this.placements = named[magicbean.Placement](ocm)
	this.transit = zecs.Tracking[directed](ocm)
	this.parent = ocm
	return this
}

type Nodes struct {
	regions, placements namedents
	transit             zecs.Tracked[directed]
	parent              *zecs.Ocm
}

type Region struct {
	zecs.Proxy
	name   name
	parent Nodes
}

type Transit struct {
	zecs.Proxy
	name   name
	t      directed
	parent Nodes
}

func (this Nodes) Region(name name) Region {
	region := this.regions.For(name)
	region.Attach(magicbean.Region{})
	return Region{region, name, this}
}

func (this Nodes) Placement(name name) Placement {
	place := this.placements.For(name)
	place.Attach(magicbean.Placement{})
	return Placement{place, name}
}

func (this Region) ConnectsTo(other Region) Transit {
	return this.connect(other.Entity(), other.name, magicbean.EdgeTransit)
}

func (this Region) Has(node Placement) Transit {
	return this.connect(node.Entity(), node.name, magicbean.EdgePlacement)
}

func (this Region) connect(to zecs.Entity, toName name, kind magicbean.EdgeKind) Transit {
	name := namef("%s -> %s", this.name, toName)
	directed := directed{From: this.Entity(), To: to}
	transit := Transit{
		Proxy:  this.parent.transit.For(directed),
		name:   name,
		t:      directed,
		parent: this.parent,
	}
	transit.Proxy.Attach(name, kind)
	return transit
}
