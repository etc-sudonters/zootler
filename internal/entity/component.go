package entity

import "github.com/etc-sudonters/substrate/mirrors"

type ComponentId uint64

const INVALID_COMPONENT ComponentId = 0

// arbitrary attachments to a Model
type Component interface{}

var ComponentType = mirrors.TypeOf[Component]()

func ComponentName(c Component) string {
	if c == nil {
		return "nil"
	}

	return PierceComponentType(c).Name()
}
