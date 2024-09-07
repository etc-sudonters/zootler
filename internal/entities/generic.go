package entities

import (
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"
)

// These are very common things to want to discuss so provide shared
// referencable convienences
type Entity interface {
	Name() components.Name
	Id() table.RowId
	AddComponents(vt table.Values) error
}

type Token struct {
	componenthaver
}

type Location struct {
	componenthaver
}

type Edge struct {
	componenthaver
	stash map[string]any
}

func (e *Edge) StashRawRule(rule components.RawLogic) {
	if _, exists := e.stash["rawrule"]; exists {
		return
	}
	e.stash["rawrule"] = rule
}

func (e Edge) GetRawRule() components.RawLogic {
	rule := e.stash["rawrule"]
	return rule.(components.RawLogic)
}

type componenthaver struct {
	rid  table.RowId
	name components.Name
	eng  query.Engine
}

func (c componenthaver) Name() components.Name {
	return c.name
}

func (c componenthaver) Id() table.RowId {
	return c.rid
}

func (c componenthaver) AddComponents(vt table.Values) error {
	return c.eng.SetValues(c.rid, vt)
}
