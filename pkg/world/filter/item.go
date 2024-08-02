package filter

import (
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/pkg/world/components"

	"github.com/etc-sudonters/substrate/mirrors"
)

func Collected(b entity.FilterBuilder) entity.FilterBuilder {
	return b.With(mirrors.TypeOf[components.Collected]())
}

func NotCollected(b entity.FilterBuilder) entity.FilterBuilder {
	return b.Without(mirrors.TypeOf[components.Collected]())
}

func Song(b entity.FilterBuilder) entity.FilterBuilder {
	return b.With(mirrors.TypeOf[components.OcarinaSong]())
}

func Placeable(b entity.FilterBuilder) entity.FilterBuilder {
	return b.With(mirrors.TypeOf[components.Placeable]())
}

func Item(b entity.FilterBuilder) entity.FilterBuilder {
	return b.With(mirrors.TypeOf[components.CollectableGameToken]())
}

func Location(b entity.FilterBuilder) entity.FilterBuilder {
	return b.With(mirrors.TypeOf[components.Location]())
}

func Inhabits(b entity.FilterBuilder) entity.FilterBuilder {
	return b.With(mirrors.TypeOf[components.Inhabits]())
}

func Inhabited(b entity.FilterBuilder) entity.FilterBuilder {
	return b.With(mirrors.TypeOf[components.Inhabited]())
}
