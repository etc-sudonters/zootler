package table

import (
	"math"
)

type Predicate interface {
	ColumnId() ColumnId
	DecideIndex(cd ColumnData) IndexSelection
}

type IndexSelection struct {
	Index Index
	Len   int
}

type IndexSelections []IndexSelection

func (is IndexSelections) Len() int {
	return len(is)
}

func (is IndexSelections) Less(i int, j int) bool {
	return is[i].Len < is[j].Len
}

func (is IndexSelections) Swap(i int, j int) {
	is[i], is[j] = is[j], is[i]
}

type Selection struct {
	Column ColumnId
}

type Query struct {
	Load  []Selection // these are already inserted into the where
	Where []Predicate
}

type RowTuple struct {
	Row    RowId
	Values []Value
}

type LeastMatchesPredicate struct {
	Reference Value
	id        ColumnId
}

func (l LeastMatchesPredicate) ColumnId() ColumnId {
	return l.id
}

func (l LeastMatchesPredicate) DecideIndex(cd ColumnData) IndexSelection {
	var chosen Index
	matches := math.MaxInt64

	for _, idx := range cd.indexes {
		matched := idx.Matches(l.Reference)
		if matched < matches {
			matches = matched
			chosen = idx
		}
	}

	return IndexSelection{chosen, matches}
}
