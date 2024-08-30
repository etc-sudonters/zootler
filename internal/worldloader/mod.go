package worldloader

import (
	"fmt"
	"iter"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"

	"github.com/etc-sudonters/substrate/skelly/graph"
	"github.com/etc-sudonters/substrate/slipup"
)

type LogicLocation struct {
	BossRoom      bool              `json:"is_boss_room"`
	DoesTimePass  bool              `json:"time_passes"`
	Dungeon       string            `json:"dungeon"`
	HintRegion    string            `json:"hint"`     // hint zone
	HintRegionAlt string            `json:"alt_hint"` // more specific location -- restricted to gannon's chambers?
	Name          components.Name   `json:"region_name"`
	Savewarp      string            `json:"savewarp"` // effectively another exit
	Scene         string            `json:"scene"`    // similar to hint?
	Events        map[string]string `json:"events"`
	Exits         map[string]string `json:"exits"`
	Locations     map[string]string `json:"locations"`
}

func (l LogicLocation) Components() table.Values {
	vt := table.Values{
		components.Location{},
	}

	if l.HintRegion != "" {
		vt = append(vt, components.HintRegion{
			Name: l.HintRegion,
			Alt:  l.HintRegionAlt,
		})
	}

	if l.BossRoom {
		vt = append(vt, components.BossRoom{})
	}

	if l.Dungeon != "" {
		vt = append(vt, components.Dungeon(l.Dungeon))
	}

	if l.DoesTimePass {
		vt = append(vt, components.TimePasses{})
	}

	if l.Savewarp != "" {
		vt = append(vt, components.SavewarpName(l.Savewarp))
	}

	return vt
}

func (l LogicLocation) Edges(yield func(locedge) bool) {
	type pair struct {
		edges map[string]string
		kind  table.Value
	}

	var edge locedge
	edge.components = make(table.Values, 2)
	edge.components[0] = components.Location{}

	pairs := []pair{
		{l.Events, components.EventEdge{}},
		{l.Exits, components.ExitEdge{}},
		{l.Locations, components.CheckEdge{}},
	}

	forAll := func(p pair) bool {
		edge.components[1] = p.kind
		for name, rule := range p.edges {
			edge.name = components.Name(name)
			edge.rule = rule
			if !yield(edge) {
				return false
			}
		}
		return true
	}

	for _, pair := range pairs {
		if !forAll(pair) {
			return
		}
	}
}

type Locations struct {
	eng   query.Engine
	loc   map[internal.NormalizedStr]*LocationBuilder
	edges map[internal.NormalizedStr]*EdgeBuilder
	g     graph.Builder
}

func (l *Locations) Build(name components.Name) (*LocationBuilder, error) {
	normaled := internal.Normalize(name)
	if loc, exists := l.loc[normaled]; exists {
		return loc, nil
	}

	id, err := l.eng.InsertRow(name, components.Location{})
	if err != nil {
		return nil, slipup.Describef(err, "while creating location %s", name)
	}

	loc := &LocationBuilder{
		parent:   l,
		id:       id,
		name:     name,
		normaled: normaled,
	}
	l.loc[normaled] = loc

	return loc, nil
}

func (l *Locations) Connect(origin, dest *LocationBuilder, rule string) (*EdgeBuilder, error) {
	name := components.Name(fmt.Sprintf("%s -> %s", origin.name, dest.name))
	normaled := internal.Normalize(name)
	if eb, exists := l.edges[normaled]; exists {
		return eb, nil
	}

	edgeId, insertErr := l.eng.InsertRow(name, components.Edge{
		Origin: entity.Model(origin.id),
		Dest:   entity.Model(dest.id),
	})

	if insertErr != nil {
		return nil, slipup.Describef(insertErr, "while inserting edge %s", name)
	}

	l.g.AddEdge(graph.Origination(origin.id), graph.Destination(dest.id))
	return &EdgeBuilder{edgeId, name, l, rule, origin, dest}, nil
}

func (l *Locations) MustEachEdge(edges map[string]string) iter.Seq[*LocationBuilder] {
	return func(yield func(*LocationBuilder) bool) {
		for name := range edges {
			normaled := internal.Normalize(name)
			loc, mustExist := l.loc[normaled]
			if !mustExist {
				panic(slipup.Createf("expected location %s to exist", name))
			}

			if !yield(loc) {
				return
			}
		}
	}
}

type Tokens struct {
	eng  query.Engine
	item map[internal.NormalizedStr]table.RowId
}

func (t *Tokens) Attach(name string, comps ...table.Value) (table.RowId, error) {
	normaled := internal.Normalize(name)
	id, exists := t.item[normaled]
	if !exists {
		var insertErr error
		id, insertErr = t.eng.InsertRow(components.Name(name))
		if insertErr != nil {
			return table.INVALID_ROWID, slipup.Describef(insertErr, "while creating token %s", name)
		}
	}

	if setErr := t.eng.SetValues(id, table.Values(comps)); setErr != nil {
		return id, slipup.Describef(setErr, "while attaching components to %s", name)
	}

	return id, nil
}

type LocationBuilder struct {
	id       table.RowId
	name     components.Name
	normaled internal.NormalizedStr
	parent   *Locations
}

func (lb *LocationBuilder) Attach(vs table.Values) error {
	return lb.parent.eng.SetValues(lb.id, vs)
}

type EdgeBuilder struct {
	id           table.RowId
	name         components.Name
	parent       *Locations
	rule         string
	origin, dest *LocationBuilder
}

func (eb *EdgeBuilder) Attach(vs table.Values) error {
	return eb.parent.eng.SetValues(eb.id, vs)
}

type locedge struct {
	name       components.Name
	rule       string
	components []table.Value
}

func NewLocations(eng query.Engine) (*Locations, error) {
	q := eng.CreateQuery()
	q.Load(query.MustAsColumnId[components.Name](eng))
	q.Exists(query.MustAsColumnId[components.Location](eng))
	rows, qErr := eng.Retrieve(q)
	if qErr != nil {
		return nil, qErr
	}

	tbl := new(Locations)
	tbl.eng = eng
	tbl.loc = make(map[internal.NormalizedStr]*LocationBuilder, rows.Len())

	for id, tup := range rows.All {
		name, ok := tup.Values[0].(components.Name)
		if !ok {
			return nil, fmt.Errorf("could not cast row %d %+v to 'components.Name'", id, tup)
		}

		loc, locErr := tbl.Build(name)
		if locErr != nil {
			return nil, locErr
		}
		tbl.loc[internal.Normalize(name)] = loc
	}

	return tbl, nil
}

func NewTokens(eng query.Engine) (*Tokens, error) {
	q := eng.CreateQuery()
	q.Load(query.MustAsColumnId[components.Name](eng))
	q.Exists(query.MustAsColumnId[components.Advancement](eng))
	q.Exists(query.MustAsColumnId[components.Advancement](eng))

	rows, qErr := eng.Retrieve(q)
	if qErr != nil {
		return nil, qErr
	}

	tbl := new(Tokens)
	tbl.eng = eng
	tbl.item = make(map[internal.NormalizedStr]table.RowId, rows.Len())

	for id, tup := range rows.All {
		name, ok := tup.Values[0].(components.Name)
		if !ok {
			return nil, fmt.Errorf("could not cast row %d %+v to 'components.Name'", id, tup)
		}

		tbl.item[internal.Normalize(name)] = id
	}
	return tbl, nil
}
