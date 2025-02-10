package magicbean

import (
	"fmt"
	"sudonters/libzootr/components"
	"sudonters/libzootr/eligibity"
	"sudonters/libzootr/internal"
	"sudonters/libzootr/internal/settings"
	"sudonters/libzootr/zecs"
)

type Criteria struct {
	settings *settings.Zootr
}

func (this Criteria) FlagToken(token zecs.Proxy, critera *eligibity.Critera) {
	membership, err := token.Membership()
	internal.PanicOnError(err)

	if zecs.IsMemberOf[components.Song](membership) {
		switch this.settings.Shuffling.Songs {
		case settings.ShuffleSongsOnSong:
			critera.Set(eligibity.Song)
		case settings.ShuffleSongsOnRewards:
			critera.Set(eligibity.DungeonReward)
		case settings.ShuffleSongsAnywhere:
			critera.Set(eligibity.Anywhere)
		}
	} else {
		// TODO
		critera.Set(eligibity.Junk)
	}

}

func (this Criteria) FlagLocation(location zecs.Proxy, critera *eligibity.Critera) {
	membership, err := location.Membership()
	internal.PanicOnError(err)

	if zecs.IsMemberOf[components.Fixed](membership) {
		if !zecs.IsMemberOf[components.HoldsToken](membership) {
			panic(fmt.Errorf("%v is fixed but has no token", location.Entity()))
		}
		return
	}

	if zecs.IsMemberOf[components.Song](membership) {
		switch this.settings.Shuffling.Songs {
		case settings.ShuffleSongsOnSong:
			critera.Set(eligibity.Song)
		case settings.ShuffleSongsOnRewards:
			critera.Set(eligibity.Song)
			critera.Set(eligibity.Junk)
		case settings.ShuffleSongsAnywhere:
			critera.Set(eligibity.Any)
		}
	} else {
		// TODO
		critera.Set(eligibity.Junk)
	}
}
