package zast

import (
	"fmt"
	"sudonters/zootler/internal/ruleparser"

	"github.com/etc-sudonters/substrate/stageleft"
)

func LowerParseTree(pt ruleparser.Tree) (Ast, error) {
	return lower(pt)
}

type Ast interface {
	AstKind() AstKind
}

type AstKind uint8
type CompareOp uint8
type BoolOp uint8
type ValueKind uint8
type IdentKind uint8

const (
	_ CompareOp = iota
	CompareEqual
	CompareNotEqual
	CompareLessThan

	_ BoolOp = iota
	BoolAnd
	BoolOr
	BoolInvert

	_ ValueKind = iota
	ValueString
	ValueNumber
	ValueBoolean

	IdentNotSet   IdentKind = 0x00
	IdentToken              = 0x01
	IdentVariable           = 0x02
	IdentSetting            = 0x03
	IdentLocation           = 0x04
	IdentBuiltIn            = 0x05
	IdentFunc               = 0x06
	IdentTrick              = 0x07
	IdentUnknown            = 0xFF

	_ AstKind = iota
	AstComparison
	AstBoolean
	AstIdentifier
	AstInvoke
	AstValue
	AstHole
)

type Comparison struct {
	LHS, RHS Ast
	Op       CompareOp
}

type Boolean struct {
	LHS, RHS Ast
	Op       BoolOp
}

type Identifier struct {
	Name string
	Kind IdentKind
}

type Value struct {
	any
	Kind ValueKind
}

type Invoke struct {
	Target Identifier
	Args   []Ast
}

type Hole struct{}

func (ast Comparison) AstKind() AstKind {
	return AstComparison
}

func (ast Boolean) AstKind() AstKind {
	return AstBoolean
}

func (ast Identifier) AstKind() AstKind {
	return AstIdentifier
}

func (ast Identifier) Eq(other Identifier) bool {
	return ast.Name == other.Name
}

func (ast Value) AstKind() AstKind {
	return AstValue
}

func (ast Invoke) AstKind() AstKind {
	return AstInvoke
}

func (ast Hole) AstKind() AstKind {
	return AstHole
}

func lower(pt ruleparser.Tree) (Ast, error) {
	switch pt := pt.(type) {
	case *ruleparser.BinOp:
		if pt.Op == ruleparser.BinOpContains {
			return createInvoke(loadSetting2Ident, pt.Right, pt.Left)
		}
		lhs, lhsErr := lower(pt.Left)
		if lhsErr != nil {
			panic(fmt.Errorf("error handling not impled: %w", lhsErr))
		}
		rhs, rhsErr := lower(pt.Right)
		if rhsErr != nil {
			panic(fmt.Errorf("error handling not impled: %w", rhsErr))
		}

		op, opExists := compareOps[pt.Op]
		if !opExists {
			panic(stageleft.AttachExitCode(
				fmt.Errorf("unknown comparison operator %T", pt),
				stageleft.ExitCode(91),
			))
		}

		return Comparison{lhs, rhs, op}, nil
	case *ruleparser.BoolOp:
		lhs, lhsErr := lower(pt.Left)
		if lhsErr != nil {
			panic(fmt.Errorf("error handling not impled: %w", lhsErr))
		}
		rhs, rhsErr := lower(pt.Right)
		if rhsErr != nil {
			panic(fmt.Errorf("error handling not impled: %w", rhsErr))
		}

		op, opExists := booleanOps[pt.Op]
		if !opExists {
			panic(stageleft.AttachExitCode(
				fmt.Errorf("unknown comparison operator %T", pt),
				stageleft.ExitCode(91),
			))
		}

		return Boolean{lhs, rhs, op}, nil
	case *ruleparser.Call:
		target, targetErr := lower(pt.Callee)
		if targetErr != nil {
			panic(fmt.Errorf("error handling not impled: %w", targetErr))
		}
		identifer, isIdent := target.(Identifier)
		if !isIdent {
			panic(fmt.Errorf("expected identifer for function target, received: %T", target))
		}
		return createInvoke(FuncIdent(identifer.Name), pt.Args...)
	case *ruleparser.Identifier:
		return Identifier{Name: pt.Value}, nil
	case *ruleparser.Subscript:
		return createInvoke(loadSetting2Ident, pt.Target, pt.Index)
	case *ruleparser.Tuple:
		return createInvoke(hasIdent, pt.Elems...)
	case *ruleparser.UnaryOp:
		if pt.Op != ruleparser.UnaryNot {
			panic(fmt.Errorf("expected unary not, received: %T", pt))
		}
		target, targetErr := lower(pt.Target)
		if targetErr != nil {

			panic(fmt.Errorf("error handling not impled: %w", targetErr))
		}
		return Boolean{target, Hole{}, BoolInvert}, nil
	case *ruleparser.Literal:
		kind, exists := valueKinds[pt.Kind]
		if !exists {
			panic(fmt.Errorf("unknown literal kind: %T", pt))
		}

		return Value{pt.Value, kind}, nil
	default:
		panic(stageleft.AttachExitCode(
			fmt.Errorf("unknown node type %T", pt),
			stageleft.ExitCode(90),
		))
	}
}

func createInvoke(what Identifier, args ...ruleparser.Tree) (Invoke, error) {
	invoke := Invoke{
		Target: what,
		Args:   make([]Ast, len(args)),
	}

	for i := range args {
		arg, err := lower(args[i])
		if err != nil {
			return invoke, err
		}
		invoke.Args[i] = arg
	}

	return invoke, nil
}

type Numbers interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8 |
		~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8 |
		~float64 | ~float32
}

func LiteralNumber[N Numbers](n N) Value {
	return Value{float64(n), ValueNumber}
}

func LiteralBoolean(b bool) Value {
	return Value{b, ValueBoolean}
}

func LiteralString(s string) Value {
	return Value{s, ValueString}
}

func BuiltInIdent(name string) Identifier {
	return Identifier{name, IdentBuiltIn}
}

func FuncIdent(name string) Identifier {
	return Identifier{name, IdentFunc}
}

func Invokes(what Identifier, args ...Ast) Invoke {
	return Invoke{what, args}
}

func TokenIdent(name string) Identifier {
	return Identifier{name, IdentToken}
}

func LocationIdent(name string) Identifier {
	return Identifier{name, IdentLocation}
}

func InvokeHas[N Numbers](what Identifier, qty N) Ast {
	return Invokes(hasIdent, what, LiteralNumber(qty))
}

func InvokeIsTrickEnabled(what string) Ast {
	return Invokes(trickEnabledIdent, LiteralString(what))
}

func InvokeLoadSetting(name string) Ast {
	return Invokes(loadSettingIdent, LiteralString(name))
}

var (
	compareOps = map[ruleparser.BinOpKind]CompareOp{
		ruleparser.BinOpEq:    CompareEqual,
		ruleparser.BinOpLt:    CompareLessThan,
		ruleparser.BinOpNotEq: CompareNotEqual,
	}

	booleanOps = map[ruleparser.BoolOpKind]BoolOp{
		ruleparser.BoolOpAnd: BoolAnd,
		ruleparser.BoolOpOr:  BoolOr,
	}

	valueKinds = map[ruleparser.LiteralKind]ValueKind{
		ruleparser.LiteralBool: ValueBoolean,
		ruleparser.LiteralNum:  ValueNumber,
		ruleparser.LiteralStr:  ValueString,
	}

	hasIdent               = BuiltInIdent("has")
	loadSetting2Ident      = BuiltInIdent("load_setting_2")
	loadSettingIdent       = BuiltInIdent("load_setting")
	trickEnabledIdent      = BuiltInIdent("is_trick_enabled")
	compareToSettingIdent  = BuiltInIdent("compare_to_setting")
	compareToVariableIdent = BuiltInIdent("compare_to_variable")
)
