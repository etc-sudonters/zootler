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

func (this *Ocm) Query() Q {
	return Q{this, this.eng.CreateQuery()}
}

type Q struct {
	set *Ocm
	q   query.Query
}

func (this *Q) Build(build BuildQuery, builds ...BuildQuery) *Q {
	build(this)
	for _, b := range builds {
		b(this)
	}
	return this
}

type BuildQuery func(*Q)

func Load[T Value](this *Q) {
	this.q.Load(query.MustAsColumnId[T](this.set.eng))
}

func With[T Value](this *Q) {
	this.q.Exists(query.MustAsColumnId[T](this.set.eng))
}

func WithOut[T Value](this *Q) {
	this.q.NotExists(query.MustAsColumnId[T](this.set.eng))
}

func (this *Q) Execute() (RowSet, error) {
	return this.set.eng.Retrieve(this.q)
}

func (this *Q) Rows(yield bundle.RowIter) {
	rows, err := this.Execute()
	if err != nil {
		panic(err)
	}
	rows.All(yield)
}
