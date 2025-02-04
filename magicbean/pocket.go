package magicbean

import (
	"maps"
	"slices"
	"sudonters/zootler/zecs"
)

func NewPockets(inventory *Inventory, ocm *zecs.Ocm) Pocket {
	var pocket Pocket
	pocket.inventory = inventory
	pocket.heartPiece = zecs.FindOne(ocm, Name("Piece of Heart"), zecs.With[Token])
	pocket.scarecrowSong = zecs.FindOne(ocm, Name("Scarecrow Song"), zecs.With[Token])
	pocket.transcribe = zecs.IndexEntities[OcarinaNote](ocm)
	pocket.songs = zecs.IndexValue[SongNotes](ocm)
	pocket.bottles = zecs.EntitiesMatching(ocm, zecs.With[Bottle])
	pocket.stones = zecs.EntitiesMatching(ocm, zecs.With[Stone])
	pocket.meds = zecs.EntitiesMatching(ocm, zecs.With[Medallion])
	pocket.rewards = zecs.EntitiesMatching(ocm, zecs.With[DungeonReward])
	pocket.notes = slices.Collect(maps.Values(pocket.transcribe))
	return pocket
}

type Pocket struct {
	inventory  *Inventory
	transcribe map[OcarinaNote]zecs.Entity
	songs      map[zecs.Entity]SongNotes

	heartPiece, scarecrowSong             zecs.Entity
	bottles, stones, meds, rewards, notes []zecs.Entity
}

func (this Pocket) Has(entity zecs.Entity, n float64) bool {
	return this.inventory.Count(entity) >= n
}

func (this Pocket) HasEvery(entities []zecs.Entity) bool {
	for _, entity := range entities {
		if !this.Has(entity, 1) {
			return false
		}
	}
	return true
}

func (this Pocket) HasAny(entities []zecs.Entity) bool {
	for _, entity := range entities {
		if this.Has(entity, 1) {
			return true
		}
	}
	return false
}

func (this Pocket) HasBottle() bool {
	return this.HasAny(this.bottles)
}

func (this Pocket) HasStones(n float64) bool {
	return this.inventory.Sum(this.stones) >= n
}

func (this Pocket) HasMedallions(n float64) bool {
	return this.inventory.Sum(this.meds) >= n
}

func (this Pocket) HasDungeonRewards(n float64) bool {
	return this.inventory.Sum(this.rewards) >= n
}

func (this Pocket) HasHearts(n float64) bool {
	pieces := this.inventory.Count(this.heartPiece)
	hearts := pieces / 4
	return hearts >= n
}

func (this Pocket) HasAllNotes(entity zecs.Entity) bool {
	if entity == this.scarecrowSong {
		return this.inventory.Sum(this.notes) >= 2
	}

	song, exists := this.songs[entity]
	if !exists {
		panic("not a song!")
	}
	notes := []OcarinaNote(song)
	transcript := make([]zecs.Entity, len(notes))
	for i, note := range notes {
		entity, exists := this.transcribe[note]
		if !exists {
			panic("unknown note")
		}
		transcript[i] = entity
	}

	return this.HasEvery(transcript)
}
