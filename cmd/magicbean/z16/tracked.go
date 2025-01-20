package z16

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

func NewTokens(ocm *zecs.Ocm) Tokens {
	return Tokens{named[magicbean.Token](ocm), ocm}
}

type Nodes struct {
	regions, placements namedents
	transit             zecs.Tracked[directed]
	parent              *zecs.Ocm
}

type Region struct {
	zecs.Proxy
	name   magicbean.Name
	parent Nodes
}

type Placement struct {
	zecs.Proxy
	name magicbean.Name
}

type Transit struct {
	Edge zecs.Proxy
	t    directed
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

func (this Region) Connects(other Region) Transit {
	directed := directed{From: this.Entity(), To: other.Entity()}
	transit := Transit{this.parent.transit.For(directed), directed}
	transit.Edge.Attach(namef("%s -> %s", this.name, other.name), magicbean.EdgeTransit)
	return transit
}

func (this Region) Has(node Placement) zecs.Proxy {
	edge := this.parent.transit.For(directed{From: this.Entity(), To: node.Entity()})
	edge.Attach(namef("%s -> %s", this.name, node.name), magicbean.EdgePlacement)
	return edge
}

func (this Placement) DefaultToken(token Token) {
	this.Attach(magicbean.DefaultPlacement(token.Entity()))
}

func (this Placement) Owns(token Token) {
	t := token.Entity()
	p := this.Entity()
	this.Attach(
		magicbean.HoldsToken(t),
		magicbean.Fixed{},
	)
	token.Attach(
		magicbean.HeldAt(p),
		magicbean.Fixed{},
	)
}

type Tokens struct {
	tokens namedents
	parent *zecs.Ocm
}

type Token struct {
	zecs.Proxy
	name name
}

func (this Tokens) Named(name name) Token {
	token := this.tokens.For(name)
	token.Attach(magicbean.Token{})
	return Token{token, name}
}
