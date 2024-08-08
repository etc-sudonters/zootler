package bundle

import (
	"errors"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/bitset"
)

var ErrExpectSingleRow = errors.New("expected exactly 1 row")

func ToMap[TKey comparable, TValue any](i Interface, f func(*table.RowTuple) (TKey, TValue, error)) (map[TKey]TValue, error) {
	m := make(map[TKey]TValue, i.Len())
	for i.MoveNext() {
		key, value, err := f(i.Current())
		if err != nil {
			return nil, err
		}

		m[key] = value
	}
	return m, nil
}

func Bundle(fill bitset.Bitset64, columns table.Columns) (Interface, error) {
	switch fill.Len() {
	case 0:
		return Empty{}, nil
	case 1:
		return Single(fill, columns)
	default:
		return Many(fill, columns), nil
	}
}
