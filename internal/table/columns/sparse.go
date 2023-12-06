package columns

import "sudonters/zootler/internal/table"

func NewHashMap() HashMap {
	return HashMap{make(map[table.RowId]table.Value)}
}

type HashMap struct {
	entities map[table.RowId]table.Value
}

func (s HashMap) Get(e table.RowId) table.Value {
	return s.entities[e]
}

func (s HashMap) Set(e table.RowId, c table.Value) {
	s.entities[e] = c
}

func (s HashMap) Unset(e table.RowId) {
	delete(s.entities, e)
}
