package galoshes

import (
	"errors"
	"fmt"
	"iter"
	"sudonters/libzootr/internal/skelly"
	"sudonters/libzootr/internal/table"

	"github.com/etc-sudonters/substrate/skelly/bitset32"
)

type EntityVarSets struct {
	vars  map[Variable][]table.ColumnId
	attrs map[Attribute]table.ColumnId
}

func (this EntityVarSets) Track(v Variable, attr Attribute) error {
	col, exists := this.attrs[attr]
	if !exists {
		return fmt.Errorf("unknown attribute: %q", attr)
	}

	sets := this.vars[v]
	sets = append(sets, col)
	this.vars[v] = sets
	return nil
}

func (this EntityVarSets) Candidates(tbl *table.Table, err *error) iter.Seq2[Variable, bitset32.Bitset] {
	return func(yield func(Variable, bitset32.Bitset) bool) {
		for variable, cols := range this.vars {
			candidates, candidateErr := candidates(tbl, cols)
			if candidateErr != nil {
				*err = candidateErr
				return
			}
			if !yield(variable, candidates) {
				return
			}
		}
	}
}

func candidates(tbl *table.Table, cols []table.ColumnId) (bitset32.Bitset, error) {
	switch len(cols) {
	case 0:
		return bitset32.Bitset{}, errors.New("no columns provided")
	case 1:
		return tbl.MembersOfColumns(cols[0])
	}

	seed, err := tbl.MembersOfColumns(cols[0])
	if err != nil {
		return bitset32.Bitset{}, err
	}
	sets := make([]bitset32.Bitset, len(cols)-2)
	for i, id := range cols[1:] {
		sets[i], err = tbl.MembersOfColumns(id)
		if err != nil {
			return bitset32.Bitset{}, err
		}
	}

	return skelly.IntersectAll(seed, sets...), nil
}
