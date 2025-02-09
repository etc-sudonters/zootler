package optimizer

import "context"

func NewCtx() Context {
	return Context{context.Background()}
}

type Context struct {
	ctx context.Context
}

func (this *Context) Store(key any, value any) {
	this.ctx = context.WithValue(this.ctx, key, value)
}

func (this *Context) Retrieve(key any) any {
	return this.ctx.Value(key)
}

func (this *Context) Swap(key any, value any) any {
	old := this.Retrieve(key)
	this.Store(key, value)
	return old
}
