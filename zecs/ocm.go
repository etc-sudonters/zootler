package zecs

import (
	"sudonters/zootler/internal/bundle"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"
)

type Entity = table.RowId
type DDL = table.DDL
type Value = table.Value
type Values = table.Values
type RowSet = bundle.Interface

func New() (Ocm, error) {
	eng, err := query.NewEngine()
	return Ocm{eng}, err
}

func Apply(ocm *Ocm, ddl []DDL) error {
	for i := range ddl {
		if _, err := ocm.eng.CreateColumn(ddl[i]()); err != nil {
			return err
		}
	}

	return nil
}

type Ocm struct {
	eng query.Engine
}

func (this *Ocm) Proxy(which Entity) Proxy {
	return Proxy{which, this}
}

func (this *Ocm) Engine() query.Engine {
	return this.eng
}

func (this *Ocm) Query() q {
	return q{this, this.eng.CreateQuery()}
}

type q struct {
	set *Ocm
	q   query.Query
}

func (this *q) Build(build BuildQuery, builds ...BuildQuery) *q {
	build(this)
	for _, b := range builds {
		b(this)
	}
	return this
}

type BuildQuery func(*q)

func Load[T Value](this *q) {
	this.q.Load(query.MustAsColumnId[T](this.set.eng))
}

func With[T Value](this *q) {
	this.q.Exists(query.MustAsColumnId[T](this.set.eng))
}

func WithOut[T Value](this *q) {
	this.q.NotExists(query.MustAsColumnId[T](this.set.eng))
}

func (this *q) Execute() (RowSet, error) {
	return this.set.eng.Retrieve(this.q)
}

func (this *q) Rows(yield bundle.RowIter) {
	rows, err := this.Execute()
	if err != nil {
		panic(err)
	}
	rows.All(yield)
}
