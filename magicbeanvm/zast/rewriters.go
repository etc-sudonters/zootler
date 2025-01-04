package zast

import (
	"errors"
	"fmt"
	"strings"
)

func IgnoreNode[T Ast]() RewriteFunc[T] {
	return func(t T, _ rewrite) (Ast, error) {
		return t, nil
	}
}

func TagIdentifiers(ctx *RewriteContext) Rewriter {
	tagIdents := tagIdents{ctx}
	return Rewriter{
		Identifier: tagIdents.Identifier,
		Value:      tagIdents.Value,
	}
}

func PullLateExpansions(ctx *RewriteContext) Rewriter {
	return Rewriter{Invoke: pullLateExpansions{ctx}.Invoke}
}

func EliminateConstOperations() Rewriter {
	var eco eliminateConstOps
	return Rewriter{
		Boolean:    eco.Boolean,
		Comparison: eco.Comparison,
	}
}

func LiftIntoFunctions(ctx *RewriteContext) Rewriter {
	lift := liftIntoFunctions{ctx}
	return Rewriter{
		Comparison: lift.Comparison,
		Identifier: lift.Identifier,
		Invoke:     IgnoreNode[Invoke](),
	}

}

type tagIdents struct {
	*RewriteContext
}

func (ti tagIdents) Identifier(node Identifier, rewrite rewrite) (Ast, error) {
	if node.Kind != IdentNotSet {
		return node, nil
	}

	if name, trimmed := strings.CutPrefix(node.Name, trickEnabledPrefix); trimmed {
		return Identifier{name, IdentTrick}, nil
	}

	kind, exists := ti.names[node.Name]
	if !exists {
		return node, nil
	}

	return Identifier{node.Name, kind}, nil
}

func (pn tagIdents) Value(node Value, rewrite rewrite) (Ast, error) {
	switch value := node.any.(type) {
	case string:
		if kind, exists := pn.names[value]; exists {
			return Identifier{value, kind}, nil
		}

		mangled := strings.ReplaceAll(value, " ", "_")

		if kind, exists := pn.names[mangled]; exists {
			return Identifier{mangled, kind}, nil
		}
	}
	return node, nil
}

type pullLateExpansions struct {
	*RewriteContext
}

func (ple pullLateExpansions) Invoke(node Invoke, rewrite rewrite) (Ast, error) {
	switch node.Target.Name {
	case "at":
		where, err := ple.extractLocation(node)
		if err != nil {
			return nil, err
		}
		return ple.attachAt(where, ple.extractRule(node, 1))
	case "here":
		return ple.attachHere(node)
	default:
		return node, nil
	}

}

func (ple pullLateExpansions) attachHere(node Invoke) (Ast, error) {
	raw, exists := ple.Retrieve(CurrentKey)
	if !exists || raw == nil {
		return nil, ErrCurrentLocationNotSet
	}

	current, ok := raw.(string)
	if !ok || current == "" {
		return nil, ErrCurrentLocationNotSet
	}

	return ple.attachAt(current, ple.extractRule(node, 0))
}

func (ple pullLateExpansions) attachAt(where string, rule Ast) (Ast, error) {
	var expansions LateExpansions
	raw, exists := ple.Retrieve(LateExpansionKey)
	if !exists || raw == nil {
		expansions = make(LateExpansions)
		ple.Store(LateExpansionKey, expansions)
	} else {
		expansions = raw.(LateExpansions)
	}

	identifer := expansions.Attach(where, rule)
	return InvokeHas(identifer, 1), nil
}

func (ple pullLateExpansions) extractRule(node Invoke, idx int) Ast {
	return node.Args[idx]
}

func (ple pullLateExpansions) extractLocation(node Invoke) (string, error) {
	value, castToValue := node.Args[0].(Value)
	if !castToValue {
		return "", ErrCouldNotExtractName
	}

	if value.Kind != ValueString {
		return "", ErrCouldNotExtractName
	}

	return value.any.(string), nil
}

type eliminateConstOps struct{}

func (eco eliminateConstOps) Comparison(node Comparison, rewrite rewrite) (Ast, error) {
	operands, operandErr := RewriteMany(rewrite, node.LHS, node.RHS)
	if operandErr != nil {
		return nil, errors.Join(operandErr, ErrCouldNotLower)
	}

	lhs, lhsIsIdent := operands[0].(Identifier)
	rhs, rhsIsIdent := operands[1].(Identifier)

	if !lhsIsIdent || !rhsIsIdent || lhs.Name != rhs.Name {
		return Comparison{operands[0], operands[1], node.Op}, nil
	}

	return LiteralBoolean(true), nil
}

func (eco eliminateConstOps) Boolean(node Boolean, rewrite rewrite) (Ast, error) {
	operands, operandErr := RewriteMany(rewrite, node.LHS, node.RHS)
	if operandErr != nil {
		return nil, errors.Join(operandErr, ErrCouldNotLower)
	}

	lhs, rhs := operands[0], operands[1]

	lhsValue, lhsIsBool := lhs.(Value)
	rhsValue, rhsIsBool := rhs.(Value)

	lhsIsBool = lhsIsBool && lhsValue.Kind == ValueBoolean
	rhsIsBool = rhsIsBool && rhsValue.Kind == ValueBoolean

	switch {
	case !lhsIsBool && !rhsIsBool:
		return Boolean{operands[0], operands[1], node.Op}, nil
	case lhsIsBool && node.Op == BoolInvert:
		return LiteralBoolean(!(lhsValue.any.(bool))), nil
	case lhsIsBool && rhsIsBool && node.Op == BoolAnd:
		return LiteralBoolean(lhsValue.any.(bool) && rhsValue.any.(bool)), nil
	case lhsIsBool && rhsIsBool && node.Op == BoolOr:
		return LiteralBoolean(lhsValue.any.(bool) || rhsValue.any.(bool)), nil
	case lhsIsBool && node.Op == BoolAnd:
		switch lhsValue.any.(bool) {
		case true:
			return rhs, nil
		case false:
			return LiteralBoolean(false), nil
		}
	case lhsIsBool && node.Op == BoolOr:
		switch lhsValue.any.(bool) {
		case true:
			return LiteralBoolean(true), nil
		case false:
			return rhs, nil
		}
	case rhsIsBool && node.Op == BoolAnd:
		switch rhsValue.any.(bool) {
		case true:
			return lhs, nil
		case false:
			return LiteralBoolean(false), nil
		}
	case rhsIsBool && node.Op == BoolOr:
		switch rhsValue.any.(bool) {
		case true:
			return LiteralBoolean(true), nil
		case false:
			return lhs, nil
		}
	}

	panic("unknown combination")
}

type liftIntoFunctions struct {
	*RewriteContext
}

func (l liftIntoFunctions) Comparison(c Comparison, rewrite rewrite) (Ast, error) {
	operands, err := RewriteMany(rewrite, c.LHS, c.RHS)
	if err != nil {
		return nil, err
	}

	lhs, rhs := operands[0], operands[1]

	switch {
	case lhs.AstKind() != AstInvoke && rhs.AstKind() != AstInvoke,
		lhs.AstKind() == AstInvoke && rhs.AstKind() == AstInvoke:
	case lhs.AstKind() == AstInvoke && rhs.AstKind() != AstInvoke:
		invoke := lhs.(Invoke)
		switch {
		case invoke.Target.Eq(loadSettingIdent):
			return Invokes(compareToSettingIdent, invoke.Args[0], rhs, LiteralNumber(c.Op)), nil
		}
	case rhs.AstKind() == AstInvoke:
		invoke := rhs.(Invoke)
		switch {
		case invoke.Target.Eq(loadSettingIdent):
			return Invokes(compareToSettingIdent, invoke.Args[0], lhs, LiteralNumber(c.Op)), nil
		}

	}
	return Comparison{lhs, rhs, c.Op}, nil
}

func (l liftIntoFunctions) Identifier(i Identifier, rewrite rewrite) (Ast, error) {
	switch i.Kind {
	case IdentTrick:
		return InvokeIsTrickEnabled(i.Name), nil
	case IdentSetting:
		return InvokeLoadSetting(i.Name), nil
	case IdentToken:
		return InvokeHas(i, 1), nil
	case IdentFunc, IdentBuiltIn:
		return Invokes(i), nil
	default:
		return i, nil
	}
}

const (
	LateExpansionKey   lateExpansionKey = "lateExpansions"
	CurrentKey         currentKey       = "current"
	trickEnabledPrefix                  = "logic_"
)

var (
	ErrCouldNotLower         = errors.New("could not lower parse tree")
	ErrCurrentLocationNotSet = errors.New("current location not set")
	ErrCouldNotExtractName   = errors.New("could not extract name")
)

type currentKey string
type lateExpansionKey string
type LateExpansion struct {
	AttachedTo Identifier
	Token      Identifier
	Rule       Ast
}
type LateExpansions map[string][]LateExpansion

func (expansions LateExpansions) Attach(where string, rule Ast) Identifier {
	attached := expansions[where]
	rank := len(attached)
	name := fmt.Sprintf("Rule@%d@%s", rank, where)
	late := LateExpansion{LocationIdent(where), TokenIdent(name), rule}
	attached = append(attached, late)
	expansions[where] = attached
	return late.Token
}
