package z2

import (
	"sudonters/zootler/internal/bundle"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/table"
)

func CreateQuery(eng query.Engine) Query {
	return Query{eng, eng.CreateQuery()}
}

type Query struct {
	eng query.Engine
	q   query.Query
}

func (this *Query) Build(builders ...BuildQuery) *Query {
	for _, build := range builders {
		build(this)
	}
	return this
}

func (this Query) Execute() (bundle.Interface, error) {
	return this.eng.Retrieve(this.q)
}

func (this Query) Rows(yield bundle.RowIter) {
	rows, err := this.Execute()
	if err != nil {
		panic(err)
	}
	rows.All(yield)
}

type BuildQuery func(*Query)

func QueryLoad[T any](q *Query) {
	q.q.Load(query.MustAsColumnId[T](q.eng))
}

func QueryWith[T any](q *Query) {
	q.q.Exists(query.MustAsColumnId[T](q.eng))

}

func QueryWithout[T any](q *Query) {
	q.q.NotExists(query.MustAsColumnId[T](q.eng))
}

func Fetch[T table.Value](eng query.Engine, entity Entity) T {
	tup, err := eng.GetValues(table.RowId(entity), table.ColumnIds{
		query.MustAsColumnId[T](eng),
	})
	if err != nil {
		panic(err)
	}
	return tup.Values[0].(T)
}
