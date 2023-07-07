package entity

import "github.com/etc-sudonters/rando/set"

type Query func() (TagName, set.Operation[Id])

func Tagged(s TagName) Query {
	return func() (TagName, set.Operation[Id]) {
		return s, set.Intersection[Id]
	}
}

func NotTagged(s TagName) Query {
	return func() (TagName, set.Operation[Id]) {
		return s, set.Difference[Id]
	}
}
