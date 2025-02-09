package tracking

import (
	"sudonters/libzootr/magicbean"
	"sudonters/libzootr/zecs"
)

type Placement struct {
	zecs.Proxy
	name magicbean.Name
}

func (this Placement) DefaultToken(token Token) {
	this.Attach(magicbean.DefaultPlacement(token.Entity()))
}

func (this Placement) Fixed(token Token) {
	this.Attach(magicbean.HoldsToken(token.Entity()), magicbean.Fixed{})
	token.Attach(magicbean.HeldAt(this.Entity()), magicbean.Fixed{})
}

func (this Placement) Holding(token Token) {
	this.Attach(magicbean.HoldsToken(token.Entity()))
	token.Attach(magicbean.HeldAt(this.Entity()))
}
