package graph

import (
	"testing"

	"github.com/etc-sudonters/zootler/internal/queue"
	"github.com/etc-sudonters/zootler/internal/stack"
	"github.com/etc-sudonters/zootler/internal/testutils"
)

const (
	A Node = 1 << iota
	B Node = 1 << iota
	C Node = 1 << iota
	D Node = 1 << iota
	E Node = 1 << iota
	F Node = 1 << iota
)

func TestBFS(t *testing.T) {
	ctx, cc := testutils.CreateContextFrom(t)
	defer cc(testutils.ErrTestEnded)

	expectedTrip := queue.From([]Node{A, B, D, C})
	g := FromOriginationMap(OriginationMap{
		Origination(A): {Destination(B), Destination(D)},
		Origination(B): {Destination(D), Destination(C)},
		Origination(C): {Destination(A), Destination(D)},
	})

	visited := VisitQueue{}

	err := BreadthFirst[Destination]{
		Selector: DebugSelector[Destination]{t.Logf, Successors},
		Visitor:  DebugVisitor{t.Logf, &visited},
	}.Walk(ctx, g, A)

	testutils.Dump(t, g)
	testutils.Dump(t, expectedTrip)
	testutils.Dump(t, visited)

	if err != nil {
		t.Fatalf("error while traversing graph: %s", err)
	}

	testutils.ArrEq(expectedTrip, visited.Q, t)
}

func TestDFS(t *testing.T) {
	ctx, cc := testutils.CreateContextFrom(t)
	defer cc(testutils.ErrTestEnded)

	expectedTrip := stack.From([]Node{A, B, E, D, C, F})
	g := FromOriginationMap(OriginationMap{
		Origination(A): {Destination(D), Destination(B)},
		Origination(B): {Destination(E)},
		Origination(C): {Destination(F)},
		Origination(D): {Destination(C)},
	})

	visited := VisitStack{}

	err := DepthFirst[Destination]{
		Selector: DebugSelector[Destination]{t.Logf, Successors},
		Visitor:  DebugVisitor{t.Logf, &visited},
	}.Walk(ctx, g, A)

	testutils.Dump(t, g)
	testutils.Dump(t, expectedTrip)
	testutils.Dump(t, visited)

	if err != nil {
		t.Fatalf("error while traversing graph: %s", err)
	}

	testutils.ArrEq(expectedTrip, visited.S, t)
}
