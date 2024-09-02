package zasm

import (
	"errors"
	"sudonters/zootler/internal/intern"
	"sudonters/zootler/internal/rules/ast"

	"github.com/etc-sudonters/substrate/skelly/stack"
	"github.com/etc-sudonters/substrate/slipup"
)

func Scratch() Block {
	return Block{
		Label: "$<>SCRATCH",
		Instr: nil,
	}
}

type MacroExpander interface {
	Expand(*ast.Call, AssemblerScope) (Instructions, error)
}

func NewAssembly() Assembly {
	return Assembly{
		Data:   Data{},
		Blocks: make(map[string]Block, 16),
	}
}

func NewDataBuilder() DataBuilder {
	var db DataBuilder
	db.Strs = intern.NewStrHeaper()
	db.Consts = intern.NewInterner[Value]()
	db.Registers = intern.NewInterner[string]()
	return db
}

type Assembly struct {
	Data   Data
	Blocks map[string]Block
}

type Block struct {
	Label string
	Instr Instructions
}

type Data struct {
	Strs      intern.StrHeap
	Consts    []Value
	Registers map[string]int
}

type DataBuilder struct {
	Strs      intern.StrHeaper
	Consts    intern.HashIntern[Value]
	Registers intern.HashIntern[string]
}

func (a *Assembly) Block(label string) (Block, error) {
	var b Block
	if _, exists := a.Blocks[label]; exists {
		return b, slipup.Createf("block %s exists already", label)
	}
	b.Label = label
	return b, nil
}

func (a *Assembly) Load(db DataBuilder) error {
	a.Data.Strs = db.Strs.Heap()
	a.Data.Consts = make([]Value, db.Consts.Len())
	a.Data.Registers = make(map[string]int, db.Registers.Len())

	for handle, value := range db.Consts.All {
		a.Data.Consts[int(handle)] = value
	}

	for register, name := range db.Registers.All {
		a.Data.Registers[name] = int(register)
	}

	return nil
}

type AssemblerScope struct {
	Replacements replacements
	Assembler    *Assembler
}

func (s *AssemblerScope) InParentScope(f func() error) error {
	s.Assembler.scope = nil
	if len(s.Assembler.scopes) > 1 {
		s.Assembler.scope = s.Assembler.scopes[1]
	}
	defer func() {
		s.Assembler.scope = s.Replacements
	}()
	return f()
}

type replacements = map[string]Instructions

type Assembler struct {
	Data      DataBuilder
	Macros    MacroExpander
	Functions map[string]int

	scopes stack.S[replacements]
	scope  replacements
}

func (a *Assembler) startScope() (AssemblerScope, func()) {
	scope := AssemblerScope{
		Replacements: make(replacements),
		Assembler:    a,
	}
	a.scopes.Push(scope.Replacements)
	a.scope = scope.Replacements
	endScope := func() {
		a.scopes.Pop()
		a.scope = nil
		if a.scopes.Len() != 0 {
			a.scope = a.scopes[0]
		}
	}
	return scope, endScope
}

func (a *Assembler) AssembleInto(b *Block, node ast.Node) (err error) {
	b.Instr, err = ast.Transform(a, node)
	return err
}

func (a *Assembler) Comparison(node *ast.Comparison) (Instructions, error) {
	var instrs Instructions
	lhs, lhsErr := ast.Transform(a, node.LHS)
	rhs, rhsErr := ast.Transform(a, node.RHS)

	if joined := errors.Join(lhsErr, rhsErr); joined != nil {
		return nil, joined
	}

	sameLoad := lhs.MatchOne(INSTR_ANY_LOAD) && lhs.Match(rhs)
	switch node.Op {
	case ast.AST_CMP_EQ:
		if sameLoad {
			return instrs.WriteLoadBool(true), nil
		}
		return lhs.Concat(rhs, instrs.WriteOp(OP_CMP_EQ)), nil
	case ast.AST_CMP_NQ:
		if sameLoad {
			return instrs.WriteLoadBool(false), nil
		}
		return lhs.Concat(rhs, instrs.WriteOp(OP_CMP_NQ)), nil
	case ast.AST_CMP_LT:
		if sameLoad {
			return instrs.WriteLoadBool(false), nil
		}
		return lhs.Concat(rhs, instrs.WriteOp(OP_CMP_LT)), nil
	default:
		panic("unsupported comparison operator")
	}
}

func (a *Assembler) BooleanOp(node *ast.BooleanOp) (Instructions, error) {
	var instrs Instructions
	lhs, lhsErr := ast.Transform(a, node.LHS)
	if lhsErr != nil {
		return nil, lhsErr
	}

	if lhs.MatchOne(INSTR_LOAD_BOOL) {
		switch node.Op {
		case ast.AST_BOOL_NEGATE:
			return instrs.WriteLoadBool(lhs.MatchOne(INSTR_LOAD_BOOL_FALSE)), nil
		case ast.AST_BOOL_AND:
			if lhs.MatchOne(INSTR_LOAD_BOOL_FALSE) {
				return instrs.Write(INSTR_LOAD_BOOL_FALSE), nil
			}
			break
		case ast.AST_BOOL_OR:
			if lhs.MatchOne(INSTR_LOAD_BOOL_TRUE) {
				return instrs.Write(INSTR_LOAD_BOOL_TRUE), nil
			}
			break
		}
	}

	if node.Op == ast.AST_BOOL_NEGATE {
		return instrs.Concat(lhs).WriteOp(OP_BOOL_NEGATE), nil
	}

	rhs, rhsErr := ast.Transform(a, node.RHS)
	if rhsErr != nil {
		return nil, rhsErr
	}

	if rhs.MatchOne(INSTR_LOAD_BOOL) {
		switch node.Op {
		case ast.AST_BOOL_AND:
			if rhs.MatchOne(INSTR_LOAD_BOOL_FALSE) {
				return instrs.Write(INSTR_LOAD_BOOL_FALSE), nil
			}
			break
		case ast.AST_BOOL_OR:
			if rhs.MatchOne(INSTR_LOAD_BOOL_TRUE) {
				return instrs.Write(INSTR_LOAD_BOOL_TRUE), nil
			}
			break
		}
	}

	instrs = instrs.Concat(lhs, rhs)
	if len(lhs) != 0 && len(rhs) != 0 {
		switch node.Op {
		case ast.AST_BOOL_AND:
			instrs = instrs.WriteOp(OP_BOOL_AND)
			break
		case ast.AST_BOOL_OR:
			instrs = instrs.WriteOp(OP_BOOL_OR)
			break
		}
	}

	if len(instrs) == 0 {
		return instrs, slipup.Createf("dropped all branches in boolop %+v", node)
	}

	return instrs, nil
}

func (a *Assembler) Call(node *ast.Call) (Instructions, error) {
	var instrs Instructions
	if node.Macro {
		scope, stopScope := a.startScope()
		defer stopScope()
		return a.Macros.Expand(node, scope)
	}

	arity, exists := a.Functions[node.Callee]
	if !exists {
		return nil, slipup.Createf("undeclared function %+v", node)
	}

	if arity != len(node.Args) {
		return nil, slipup.Createf(
			"'%s' expects %d args but got %d: %+v", node.Callee, arity, len(node.Args), node,
		)
	}

	var loadErrs error

	for _, arg := range node.Args {
		load, loadErr := ast.Transform(a, arg)
		loadErrs = errors.Join(loadErrs, loadErr)
		instrs = instrs.Concat(load)
	}

	if loadErrs != nil {
		return nil, loadErrs
	}

	return instrs.WriteCall(arity)
}

func (a *Assembler) Identifier(node *ast.Identifier) (Instructions, error) {
	if a.scope != nil {
		replacement, exists := a.scope[node.Name]
		if exists {
			return replacement, nil
		}
	}

	var instrs Instructions
	register := a.Data.Registers.Intern(node.Name)
	return instrs.WriteLoadIdent(uint32(register)), nil
}

func (a *Assembler) Literal(node *ast.Literal) (Instructions, error) {
	var instrs Instructions
	switch node.Kind {
	case ast.AST_LIT_NUM:
		v := Pack(node.Value.(float64))
		c := a.Data.Consts.Intern(v)
		return instrs.WriteLoadConst(uint32(c)), nil
	case ast.AST_LIT_BOOL:
		return instrs.WriteLoadBool(node.Value.(bool)), nil
	case ast.AST_LIT_STR:
		s := node.Value.(string)
		str := a.Data.Strs.Intern(s)
		return instrs.WriteLoadStr(str), nil
	default:
		panic("invalid literal kind")
	}
}

func (a *Assembler) Empty(node *ast.Empty) (Instructions, error) {
	// not required to emit code
	// like it or not, this is what performance looks like
	return nil, nil
}
