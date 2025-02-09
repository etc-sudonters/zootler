package optimizer

import (
	"errors"
	"fmt"
	"sudonters/libzootr/mido/ast"
	"sudonters/libzootr/mido/symbols"
)

type currentKey string

const currentLocationKey currentKey = "current-key"

func SetCurrentLocation(ctx *Context, where string) {
	ctx.Store(currentLocationKey, where)
}

func GetCurrentLocation(ctx *Context) string {
	answer := ctx.Retrieve(currentLocationKey)
	if answer != nil {
		return answer.(string)
	}
	return ""
}

type CompilerConnector interface {
	AddConnectionTo(regionName string, rule ast.Node) (*symbols.Sym, error)
}

func NewConnectionGeneration(ctx *Context, syms *symbols.Table, conns CompilerConnector) ast.Rewriter {
	plugin := connectionplugin{
		ctx,
		conns,
		syms,
		syms.Declare("at", symbols.FUNCTION),
		syms.Declare("here", symbols.FUNCTION),
		syms.Declare("has", symbols.FUNCTION),
	}

	return ast.Rewriter{
		Invoke: plugin.Invoke,
	}
}

type connectionplugin struct {
	ctx           *Context
	conns         CompilerConnector
	symbols       *symbols.Table
	at, here, has *symbols.Sym
}

func (this connectionplugin) Invoke(node ast.Invoke, _ ast.Rewriting) (ast.Node, error) {
	target := ast.LookUpNodeInTable(this.symbols, node.Target)

	switch {
	case target == nil:
	case target.Index == this.at.Index:
		return this.replaceAt(node.Args)
	case target.Index == this.here.Index:
		return this.replaceHere(node.Args)
	}

	return node, nil
}

func (this connectionplugin) replaceAt(args []ast.Node) (ast.Node, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("expected 2 args for at, received %d", len(args))
	}
	origin, ok := args[0].(ast.String)
	if !ok {
		return nil, fmt.Errorf("expected first argument to 'at' to be string: %#v", args[0])
	}
	return this.generate(string(origin), args[1])

}

func (this connectionplugin) replaceHere(args []ast.Node) (ast.Node, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expected 1 arg for here, received %d", len(args))
	}
	current := GetCurrentLocation(this.ctx)
	if current == "" {
		return nil, errors.New("expected current location to be set")
	}
	return this.generate(current, args[0])

}

func (this connectionplugin) generate(where string, rule ast.Node) (ast.Node, error) {
	symbol, err := this.conns.AddConnectionTo(where, rule)
	if err != nil {
		return nil, err
	}
	return this.createHas(symbol), nil
}

func (this connectionplugin) createHas(token *symbols.Sym) ast.Node {
	return ast.Invoke{
		Target: ast.IdentifierFrom(this.has),
		Args: []ast.Node{
			ast.IdentifierFrom(token),
			ast.Number(1),
		},
	}
}
