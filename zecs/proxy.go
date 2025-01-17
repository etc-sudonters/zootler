package zecs

type Proxy struct {
	id     Entity
	parent *Ocm
}

func (this Proxy) Entity() Entity {
	return this.id
}

func (this *Proxy) Attach(values ...Value) error {
	return this.AttachAll(values)
}

func (this *Proxy) AttachAll(values Values) error {
	return this.parent.eng.SetValues(this.id, values)
}

func (this *Proxy) AttachFrom(from Attaching) error {
	return this.AttachAll(from.vs)
}

type Attaching struct {
	vs Values
}

func (this *Attaching) Add(v Value) {
	this.vs = append(this.vs, v)
}

type ComparableValue interface {
	Value
	comparable
}

func Tracking[Key ComparableValue](parent *Ocm, query ...BuildQuery) Tracked[Key] {
	q := parent.Query()
	rows, err := q.Build(Load[Key], query...).Execute()
	if err != nil {
		panic(err)
	}
	cache := make(map[Key]Entity, rows.Len())

	for row, tup := range rows.All {
		cache[tup.Values[0].(Key)] = Entity(row)
	}

	return Tracked[Key]{parent, cache}
}

type Tracked[Key ComparableValue] struct {
	parent *Ocm
	cache  map[Key]Entity
}

func (this *Tracked[Key]) For(key Key) Proxy {
	if entity, exists := this.cache[key]; exists {
		return Proxy{entity, this.parent}
	}

	row, err := this.parent.eng.InsertRow(key)
	if err != nil {
		panic(err)
	}
	entity := Entity(row)
	this.cache[key] = entity
	return Proxy{entity, this.parent}
}
