package components

import "sudonters/zootler/internal/entity"

type TokenArchetype struct{}
type EventArchetype struct{}
type LocationArchetype struct{}

func (e EventArchetype) Apply(entity entity.View) error {
	entity.Add(Token{})
	entity.Add(Event{})
	return nil
}
