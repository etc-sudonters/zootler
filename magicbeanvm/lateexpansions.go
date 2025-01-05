package magicbeanvm

import (
	"errors"
	"fmt"
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/optimizer"
	"sudonters/zootler/magicbeanvm/symbols"
)

const LateExpansionKey lateKey = "late-expansions"
const CurrentLocationKey currentKey = "current-key"

func ExtractLateExpansions(ctx *optimizer.Context, tbl *symbols.Table) ast.Rewriter {
	var late late
	late.at, late.here = tbl.Declare("at", symbols.FUNCTION), tbl.Declare("here", symbols.FUNCTION)
	late.tbl = tbl
	captured := make(LateExpansions)
	late.ctx = ctx
	ctx.Store(LateExpansionKey, captured)
	return ast.Rewriter{
		Invoke: late.Invoke,
	}
}

type LateExpansion struct {
	Token, AttachedTo symbols.Index
	Rule              ast.Node
	Rank              int
}
type LateExpansions map[string][]LateExpansion

func (this LateExpansions) AppendTo(parent string) *LateExpansion {
	items := this[parent]
	rank := len(items)
	items = append(items, LateExpansion{
		Rank: rank,
	})
	this[parent] = items
	return &this[parent][rank]
}

func (this LateExpansions) Size() int {
	var size int
	if len(this) == 0 {
		return size
	}

	for _, items := range this {
		size += len(items)
	}

	return size
}

type lateKey string
type currentKey string

type late struct {
	at, here *symbols.Sym
	tbl      *symbols.Table
	ctx      *optimizer.Context
}

func (this late) Invoke(invoke ast.Invoke, rewrite ast.Rewriting) (ast.Node, error) {
	target := this.lookUpFromNode(invoke.Target)
	if target == nil || (!target.Eq(this.at) && !target.Eq(this.here)) {
		return invoke, nil
	}

	var attachedTo string
	var rule ast.Node
	switch target.Name {
	case "at":
		str, ok := invoke.Args[0].(ast.String)
		if !ok {
			return nil, fmt.Errorf("expected first argument to 'at' to be string: %#v", invoke.Args[0])
		}
		attachedTo = string(str)
		rule = invoke.Args[1]
	case "here":
		current := this.ctx.Retrieve(CurrentLocationKey)
		if current == nil {
			return nil, errors.New("expected current location to be set")
		}
		attachedTo = current.(string)
		rule = invoke.Args[0]
	default:
		panic("unreachable")
	}

	captured := this.ctx.Retrieve(LateExpansionKey).(LateExpansions)
	delayed := captured.AppendTo(attachedTo)
	token := this.tbl.Declare(fmt.Sprintf("Token#%04d@%s", delayed.Rank, attachedTo), symbols.TOKEN)
	delayed.AttachedTo = this.tbl.Declare(attachedTo, symbols.LOCATION).Index
	delayed.Token = token.Index
	delayed.Rule = rule

	return ast.Invoke{
		Target: ast.IdentifierFrom(this.tbl.Declare("has", symbols.FUNCTION)),
		Args: []ast.Node{
			ast.IdentifierFrom(token), ast.Number(1),
		},
	}, nil
}

func (this late) lookUpFromNode(node ast.Node) *symbols.Sym {
	switch node := node.(type) {
	case ast.Identifier:
		return this.tbl.LookUpByIndex(node.AsIndex())
	default:
		return nil
	}
}
