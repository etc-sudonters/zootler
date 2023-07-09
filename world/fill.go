package world

import "context"

type Goal interface {
	Reachable(context.Context, World) (bool, error)
}

// modifies World but ensures Goal remains statisfied
type Filler interface {
	Fill(context.Context, World, Goal)
}

type RandomFill struct{}
type AssumedFill struct{}
type ForwardFill struct{}
