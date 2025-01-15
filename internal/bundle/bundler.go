package bundle

import (
	"errors"
	"sudonters/zootler/internal/skelly/bitset"
	"sudonters/zootler/internal/table"
)

var ErrExpectSingleRow = errors.New("expected exactly 1 row")

func Bundle(fill bitset.Bitset32, columns table.Columns) (Interface, error) {
	switch fill.Len() {
	case 0:
		return Empty{}, nil
	case 1:
		return Single(fill, columns)
	default:
		return Many(fill, columns), nil
	}
}
