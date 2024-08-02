package components

import (
	"regexp"
	"strings"
	"sudonters/zootler/internal/entity"

	"github.com/etc-sudonters/substrate/mirrors"
)

type TokenArchetype struct {
	Strs mirrors.TypedStrings
}
type EventArchetype struct {
	T TokenArchetype
}
type LocationArchetype struct{}

func (t TokenArchetype) Apply(entity entity.View) error {
	var name Name
	if err := entity.Get(&name); err != nil {
		return err
	}

	escaped := EscapeName(name)
	comp := t.Strs.InstanceOf(escaped)
	if err := entity.Add(comp); err != nil {
		return err
	}

	if err := entity.Add(CollectableGameToken{}); err != nil {
		return err
	}

	return nil
}

func (e EventArchetype) Apply(entity entity.View) error {
	e.T.Apply(entity)
	entity.Add(Event{})
	return nil
}

var _nameEscapeRe *regexp.Regexp = regexp.MustCompile(`['()\[\]-]`)

func EscapeName(name Name) string {
	escaped := _nameEscapeRe.ReplaceAllLiteralString(string(name), "")
	return strings.ReplaceAll(escaped, " ", "_")
}
