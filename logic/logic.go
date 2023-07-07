package logic

import "github.com/etc-sudonters/rando/entity"

type Queryable interface {
	Query(...entity.Query) []*entity.View
}

type Rule interface {
	Fulfill(Queryable) (bool, error)
}

type RuleFunc func(Queryable) (bool, error)

func (r RuleFunc) Fulfill(q Queryable) (bool, error) {
	return r(q)
}

func AllRules(rs ...Rule) RuleFunc {
	return func(q Queryable) (bool, error) {
		for _, r := range rs {
			pass, err := r.Fulfill(q)
			if err != nil {
				return false, err
			}
			if !pass {
				return false, nil
			}
		}

		return true, nil
	}
}

func AnyRule(rs ...Rule) RuleFunc {
	return func(q Queryable) (bool, error) {
		for _, r := range rs {
			pass, err := r.Fulfill(q)
			if err != nil {
				return false, err
			}

			if pass {
				return true, nil
			}
		}

		return false, nil
	}
}
