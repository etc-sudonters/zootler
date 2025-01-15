package z2

import (
	"sudonters/zootler/internal/bundle"
	"sudonters/zootler/internal/query"
)

func CreateQuery(eng query.Engine) Query {
	return Query{eng, eng.CreateQuery()}
}

type Query struct {
	eng query.Engine
	q   query.Query
}

func (this *Query) Build(builders ...BuildQuery) {
	for _, build := range builders {
		build(this)
	}
}

func (this Query) Execute() (bundle.Interface, error) {
	return this.eng.Retrieve(this.q)
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
