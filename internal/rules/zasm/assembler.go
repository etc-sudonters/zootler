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
	Data        DataBuilder
	Macros      MacroExpander
	Functions   map[string]int
	DebugOutput func(string, ...any)

	scopes stack.S[replacements]
	scope  replacements
}

func (a *Assembler) debug(tpl string, vs ...any) {
	if a.DebugOutput != nil {
		a.DebugOutput(tpl, vs...)
	}
}

func (a *Assembler) startScope() (AssemblerScope, func()) {
	scope := AssemblerScope{
		Replacements: make(replacements),
		Assembler:    a,
	}
	a.scopes.Push(scope.Replacements)
	a.scope = scope.Replacements
	a.debug("starting scope")
	endScope := func() {
		a.debug("ending scope")
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
	instrs := Tape()
	lhs, lhsErr := ast.Transform(a, node.LHS)
	rhs, rhsErr := ast.Transform(a, node.RHS)

	if joined := errors.Join(lhsErr, rhsErr); joined != nil {
		return nil, joined
	}

	sameLoad := lhs.MatchOne(INSTR_ANY_LOAD) && lhs.Match(rhs)
	if sameLoad {
		a.debug("matched samed load (%s) on %+v", lhs, node)
	}
	switch node.Op {
	case ast.AST_CMP_EQ:
		if sameLoad {
			return instrs.WriteLoadBool(true).Instructions(), nil
		}
		return instrs.Concat(lhs, rhs).WriteOp(OP_CMP_EQ).Instructions(), nil
	case ast.AST_CMP_NQ:
		if sameLoad {
			return instrs.WriteLoadBool(false).Instructions(), nil
		}
		return instrs.Concat(lhs, rhs).WriteOp(OP_CMP_NQ).Instructions(), nil
	case ast.AST_CMP_LT:
		if sameLoad {
			return instrs.WriteLoadBool(false).Instructions(), nil
		}
		return instrs.Concat(lhs, rhs).WriteOp(OP_CMP_LT).Instructions(), nil
	default:
		panic("unsupported comparison operator")
	}
}

func (a *Assembler) BooleanOp(node *ast.BooleanOp) (Instructions, error) {

	runLhs := func() (Instructions, error) {
		a.debug("lhs %s", node.LHS)
		lhs, lhsErr := ast.Transform(a, node.LHS)
		if lhsErr != nil {
			return nil, lhsErr
		}

		if lhs == nil {
			a.debug("no output lhs")
		}

		if lhs.MatchOne(INSTR_LOAD_BOOL) {
			switch node.Op {
			case ast.AST_BOOL_NEGATE:
				a.debug("inverting compile time bool")
				return Tape().WriteLoadBool(lhs.MatchOne(INSTR_LOAD_BOOL_FALSE)).Instructions(), nil
			case ast.AST_BOOL_AND:
				if lhs.MatchOne(INSTR_LOAD_BOOL_FALSE) {
					a.debug("dropping rhs of AND because lhs is always FALSE %+v", node)
					return lhs, nil
				}
				return nil, nil
			case ast.AST_BOOL_OR:
				if lhs.MatchOne(INSTR_LOAD_BOOL_TRUE) {
					a.debug("dropping rhs of OR because lhs is always TRUE %+v", node)
					return lhs, nil
				}
				return nil, nil
			}
		}

		if node.Op == ast.AST_BOOL_NEGATE {
			a.debug("could not eliminate runtime bool flip")
			return Tape().Concat(lhs).WriteOp(OP_BOOL_NEGATE).Instructions(), nil
		}

		return lhs, nil
	}

	runRhs := func() (Instructions, error) {
		a.debug("rhs %s", node.RHS)
		rhs, rhsErr := ast.Transform(a, node.RHS)
		if rhsErr != nil {
			return nil, rhsErr
		}

		if rhs == nil {
			a.debug("no output rhs")
		}

		if rhs.MatchOne(INSTR_LOAD_BOOL) {
			switch node.Op {
			case ast.AST_BOOL_AND:
				if rhs.MatchOne(INSTR_LOAD_BOOL_FALSE) {
					a.debug("dropping lhs of AND because rhs is always FALSE %+v", node)
					return rhs, nil
				}
				return nil, nil
			case ast.AST_BOOL_OR:
				if rhs.MatchOne(INSTR_LOAD_BOOL_TRUE) {
					a.debug("dropping lhs of OR because rhs is always TRUE %+v", node)
					return rhs, nil
				}
				return nil, nil
			}
		}

		return rhs, nil
	}

	iw := Tape()
	lhs, lhsErr := runLhs()
	if lhsErr != nil || node.Op == ast.AST_BOOL_NEGATE {
		return lhs, lhsErr
	}

	rhs, rhsErr := runRhs()
	if rhsErr != nil {
		return nil, rhsErr
	}

	iw.Concat(lhs, rhs)
	if len(lhs) != 0 && len(rhs) != 0 {
		a.debug("could not elminate alternate %+v", node)
		switch node.Op {
		case ast.AST_BOOL_AND:
			iw.WriteOp(OP_BOOL_AND)
			break
		case ast.AST_BOOL_OR:
			iw.WriteOp(OP_BOOL_OR)
			break
		}
	}

	if iw.Len() == 0 {
		return nil, slipup.Createf("dropped all branches in boolop %+v", node)
	}

	return iw.Instructions(), nil
}

var fastCalls map[string]Op

func fastCallsInit() {
	if fastCalls == nil {
		fastCalls = map[string]Op{
			"has":            OP_CHK_QTY,
			"load_setting":   OP_CHK_SET_1,
			"load_setting_2": OP_CHK_SET_2,
			"load_trick":     OP_CHK_TRK,
		}
	}
}

func (a *Assembler) Call(node *ast.Call) (Instructions, error) {
	iw := Tape()
	if node.Macro {
		scope, stopScope := a.startScope()
		defer stopScope()
		a.debug("expanding macro %+v", node)
		expansion, err := a.Macros.Expand(node, scope)
		return expansion, err
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

	for idx, arg := range node.Args {
		load, loadErr := ast.Transform(a, arg)
		loadErrs = errors.Join(loadErrs, loadErr)
		if len(load) == 0 {
			a.debug("macro arg %d expanded to nothing: %v", idx, arg)

		}
		iw.Concat(load)
	}

	if loadErrs != nil {
		return nil, loadErrs
	}

	fastCallsInit()
	if fastOp, exists := fastCalls[node.Callee]; exists {
		return iw.WriteOp(fastOp).Instructions(), nil
	}

	_, err := iw.WriteCall(arity)
	return iw.Instructions(), err
}

func (a *Assembler) Identifier(node *ast.Identifier) (Instructions, error) {
	if a.scope != nil {
		replacement, exists := a.scope[node.Name]
		if exists {
			return replacement, nil
		}
	}

	iw := Tape()
	register := a.Data.Registers.Intern(node.Name)
	return iw.WriteLoadIdent(uint32(register)).Instructions(), nil
}

func (a *Assembler) Literal(node *ast.Literal) (Instructions, error) {
	switch node.Kind {
	case ast.AST_LIT_NUM:
		v := Pack(node.Value.(float64))
		c := a.Data.Consts.Intern(v)
		return Tape().WriteLoadConst(uint32(c)).Instructions(), nil
	case ast.AST_LIT_BOOL:
		return Tape().WriteLoadBool(node.Value.(bool)).Instructions(), nil
	case ast.AST_LIT_STR:
		s := node.Value.(string)
		str := a.Data.Strs.Intern(s)
		return Tape().WriteLoadStr(str).Instructions(), nil
	default:
		panic("invalid literal kind")
	}
}

func (a *Assembler) Empty(node *ast.Empty) (Instructions, error) {
	// not required to emit code
	// like it or not, this is what performance looks like
	a.debug("empty node")
	return nil, nil
}
