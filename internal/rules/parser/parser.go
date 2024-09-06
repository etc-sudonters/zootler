package parser

import (
	"fmt"
	"github.com/etc-sudonters/substrate/peruse"
	"github.com/etc-sudonters/substrate/slipup"
	"strings"
)

func ParseFunctionDecl(decl, body string) (FunctionDecl, error) {
	var f FunctionDecl

	funcDecl, funcDeclErr := Parse(decl)
	if funcDeclErr != nil {
		return f, slipup.Describe(funcDeclErr, "while parsing function decl")
	}

	switch d := funcDecl.(type) {
	case *Identifier:
		f.Identifier = d.Value
		break
	case *Call:
		ident, wasIdent := d.Callee.(*Identifier)
		if !wasIdent {
			return f, slipup.Createf("unsupported function decl identifier: %v", d)
		}

		f.Identifier = ident.Value
		f.Parameters = make([]string, len(d.Args))

		for i := range d.Args {
			switch a := d.Args[i].(type) {
			case *Identifier:
				f.Parameters[i] = a.Value
				break
			default:
				return f, slipup.Createf("unsupported function parameter identifier: %v", d)
			}
		}

		break
	default:
		return f, slipup.Createf("unsupported function decl identifier: %v", d)
	}

	funcBody, funcBodyErr := Parse(body)
	if funcBodyErr != nil {
		return f, slipup.Describe(funcBodyErr, "while parsing function body")
	}

	f.Body = funcBody
	return f, nil
}

func Parse(raw string) (Expression, error) {
	l := NewRulesLexer(raw)
	p := NewRulesParser(l)
	return p.Parse()
}

func MustParse(raw string) Expression {
	expr, err := Parse(raw)
	if err != nil {
		panic(err)
	}
	return expr
}

func NewRulesParser(tokens peruse.TokenStream) *peruse.Parser[Expression] {
	g := NewRulesGrammar()
	return peruse.NewParser(&g, tokens)
}

func UnaryOpFromTok(t peruse.Token) UnaryOpKind {
	switch t.Literal {
	case string(UnaryNot):
		return UnaryNot
	default:
		panic(fmt.Errorf("invalid unaryop %q", t))
	}
}

func BoolOpFromTok(t peruse.Token) BoolOpKind {
	switch s := strings.ToLower(t.Literal); s {
	case string(BoolOpAnd):
		return BoolOpAnd
	case string(BoolOpOr):
		return BoolOpOr
	default:
		panic(fmt.Errorf("invalid boolop %q", t))
	}
}

func BinOpFromTok(t peruse.Token) BinOpKind {
	switch t.Literal {
	case string(BinOpLt):
		return BinOpLt
	case string(BinOpEq):
		return BinOpEq
	case string(BinOpNotEq):
		return BinOpNotEq
	case string(BinOpContains):
		return BinOpContains
	default:
		panic(fmt.Errorf("invalid binop %q", t))
	}
}
