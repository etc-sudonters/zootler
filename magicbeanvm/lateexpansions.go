package magicbeanvm

import (
	"errors"
	"fmt"
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/optimizer"
	"sudonters/zootler/magicbeanvm/symbols"
)

const lateExpansionKey lateKey = "late-expansions"
const currentLocationKey currentKey = "current-key"

func SetCurrentLocation(ctx *optimizer.Context, where string) {
	ctx.Store(currentLocationKey, where)
}

func SwapLateExpansions(ctx *optimizer.Context) LateExpansions {
	var late LateExpansions
	stored := ctx.Swap(lateExpansionKey, make(LateExpansions))
	if stored != nil {
		late = stored.(LateExpansions)
	} else {
		late = make(LateExpansions)
	}
	return late
}

func ExtractLateExpansions(ctx *optimizer.Context, tbl *symbols.Table) late {
	var late late
	late.at, late.here = tbl.Declare("at", symbols.FUNCTION), tbl.Declare("here", symbols.FUNCTION)
	late.tbl = tbl
	late.ctx = ctx
	ctx.Store(lateExpansionKey, make(LateExpansions))
	return late
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

func (this late) At(args []ast.Node) (ast.Node, error) {
	if len(args) != 2 {
		panic("wrong arg count")
	}
	str, ok := args[0].(ast.String)
	if !ok {
		return nil, fmt.Errorf("expected first argument to 'at' to be string: %#v", args[0])
	}
	replacement, _, err := this.rewrite(string(str), args[1])
	return replacement, err
}

func (this late) Here(args []ast.Node) (ast.Node, error) {
	if len(args) != 1 {
		panic("wrong arg count")
	}
	current := this.ctx.Retrieve(currentLocationKey)
	if current == nil {
		return nil, errors.New("expected current location to be set")
	}
	replacement, _, err := this.rewrite(current.(string), args[0])
	return replacement, err
}

func (this late) rewrite(attachedTo string, rule ast.Node) (ast.Node, *LateExpansion, error) {
	captured := this.ctx.Retrieve(lateExpansionKey).(LateExpansions)
	delayed := captured.AppendTo(attachedTo)
	token := this.tbl.Declare(fmt.Sprintf("Token#%04d@%s", delayed.Rank, attachedTo), symbols.EVENT)
	delayed.AttachedTo = this.tbl.Declare(attachedTo, symbols.LOCATION).Index
	delayed.Token = token.Index
	delayed.Rule = rule

	return ast.Invoke{
		Target: ast.IdentifierFrom(this.tbl.Declare("has", symbols.FUNCTION)),
		Args: []ast.Node{
			ast.IdentifierFrom(token), ast.Number(1),
		},
	}, delayed, nil
}
