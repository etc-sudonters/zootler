package components

import "github.com/etc-sudonters/substrate/skelly/graph"

type BossRoom struct{}

type HintRegion struct {
	Name, Alt string
}

type Edge struct {
	Origin graph.Origination
	Dest   graph.Destination
}

type RawLogic struct {
	Rule string
}

type Helper struct{}
