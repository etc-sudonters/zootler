package entities

import (
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/slipup"
)

type Tokens struct {
	*genericmap[Token]
}

type Edges struct {
	*genericmap[Edge]
}

type Locations struct {
	*genericmap[Location]
}

type Map[T Entity] interface {
	Entity(components.Name) (T, error)
	All(func(T) bool)
}

type genericmap[T Entity] struct {
	cache map[internal.NormalizedStr]T
	eng   query.Engine
	def   table.Values
	fact  func(table.RowId, components.Name, table.Values) T
}

func (e *genericmap[T]) CacheLen() int {
	return len(e.cache)
}

func (e *genericmap[T]) Entity(name components.Name) (T, error) {
	normaled := internal.Normalize(name)
	if ent, exists := e.cache[normaled]; exists {
		return ent, nil
	}

	ent, err := e.eng.InsertRow(name)
	if err != nil {
		var t T
		return t, err
	}

	if err := e.eng.SetValues(ent, e.def); err != nil {
		var t T
		return t, err
	}

	t := e.fact(ent, name, e.def)
	e.cache[normaled] = t
	return t, nil
}

func (e *genericmap[T]) All(yield func(T) bool) {
	for _, t := range e.cache {
		if !yield(t) {
			return
		}
	}
}

func (e *genericmap[T]) init(f func(query.Engine, query.Query)) error {
	q := e.eng.CreateQuery()
	f(e.eng, q)
	rows, err := e.eng.Retrieve(q)
	if err != nil {
		return slipup.Describe(err, "while initializing map")
	}

	for rid, row := range rows.All {
		m := row.ColumnMap()
		name, _ := table.FromColumnMap[components.Name](m)
		e.cache[internal.Normalize(name)] = e.fact(rid, name, e.def)
	}

	return nil
}

func LocationMap(eng query.Engine) (Locations, error) {
	locs, err := newmap(
		eng,
		func(eng query.Engine, q query.Query) {
			q.Exists(query.MustAsColumnId[components.Location](eng))
		},
		func(id table.RowId, name components.Name, _ table.Values) Location {
			var t Location
			t.rid = id
			t.name = name
			t.eng = eng
			return t
		},
		components.Location{},
	)

	return Locations{locs}, err
}

func TokenMap(eng query.Engine) (Tokens, error) {
	toks, err := newmap(
		eng,
		func(eng query.Engine, q query.Query) {
			q.Exists(query.MustAsColumnId[components.CollectableGameToken](eng))
		},
		func(id table.RowId, name components.Name, _ table.Values) Token {
			var t Token
			t.rid = id
			t.name = name
			t.eng = eng
			return t
		},
		components.CollectableGameToken{},
	)
	return Tokens{toks}, err
}

func EdgeMap(eng query.Engine) (Edges, error) {
	edges, err := newmap(
		eng,
		func(eng query.Engine, q query.Query) {
			q.Exists(query.MustAsColumnId[components.Edge](eng))
		},
		func(id table.RowId, name components.Name, _ table.Values) Edge {
			var t Edge
			t.rid = id
			t.name = name
			t.eng = eng
			t.stash = make(map[string]any, 8)
			return t
		},
		components.Edge{},
	)

	return Edges{edges}, err
}

func newmap[T Entity](
	eng query.Engine,
	q func(query.Engine, query.Query),
	fact func(table.RowId, components.Name, table.Values) T,
	def ...table.Value,
) (*genericmap[T], error) {
	var g genericmap[T]
	g.cache = make(map[internal.NormalizedStr]T, 256)
	g.def = def
	g.fact = fact
	g.eng = eng

	if err := g.init(q); err != nil {
		return &g, err
	}
	return &g, nil
}
