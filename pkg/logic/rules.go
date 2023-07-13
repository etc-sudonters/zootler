package logic

import "github.com/etc-sudonters/zootler/pkg/entity"

type hasQuantityOf struct {
	desired Name
	qty     uint
}

type TrickEnabled struct {
	Trick Name
	// we only need to look up once
	cache *bool
}

func HasQuantityOf(desired Name, qty uint) hasQuantityOf {
	return hasQuantityOf{desired, qty}
}

func (h hasQuantityOf) Fulfill(q entity.Queryable) (bool, error) {
	acquiredTokens, err := q.Query(
		Collected{},
		entity.With[Advancement]{},
		entity.With[Token]{},
		entity.Load[Name]{},
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
		Trick(""),
		entity.With[Enabled]{},
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
