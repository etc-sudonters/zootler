package logic

import "github.com/etc-sudonters/zootler/pkg/entity"

type hasQuantityOf struct {
	desired Name
	qty     uint
}

// fulfilled if the specified trick is enabled
type TrickEnabled struct {
	Trick Name
	// we only need to look up once
	cache *bool
}

// fulfilled if we've collected at least N of the desired token
func HasQuantityOf(desired Name, qty uint) hasQuantityOf {
	return hasQuantityOf{desired, qty}
}

func (h hasQuantityOf) Fulfill(q entity.Queryable) (bool, error) {
	acquiredTokens, err := q.Query(
		entity.With[Collected]{},
		entity.With[Token]{},
		entity.With[Name]{},
	)
	if err != nil {
		return false, err
	}

	count := 0

	var tokName Name
	for _, tok := range acquiredTokens {
		err := tok.Get(&tokName)
		if err != nil {
			return false, err
		}

		if tokName == h.desired {
			count++
		}
	}

	return count == int(h.qty), nil
}

func (t *TrickEnabled) Fulfill(q entity.Queryable) (bool, error) {
	if t.cache != nil {
		return *t.cache, nil
	}

	enabledTricks, err := q.Query(
		entity.With[Trick]{},
		entity.With[Enabled]{},
		entity.With[Name]{},
	)

	if err != nil {
		return false, err
	}

	var tokName Name
	for _, tok := range enabledTricks {
		if err := tok.Get(&tokName); err != nil {
			return false, err
		}

		if tokName == t.Trick {
			*t.cache = true
			return true, nil
		}
	}

	*t.cache = false
	return false, nil
}
