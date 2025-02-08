package tracking

import "sudonters/zootler/zecs"

type Set struct {
	Nodes  Nodes
	Tokens Tokens
}

func NewTrackingSet(ocm *zecs.Ocm) Set {
	return Set{
		Nodes:  NewNodes(ocm),
		Tokens: NewTokens(ocm),
	}
}
