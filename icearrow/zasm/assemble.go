package zasm

import (
	"errors"
	"sudonters/zootler/icearrow/ast"

	"github.com/etc-sudonters/substrate/slipup"
)

type ZasmMacroExpander interface {
	ExpandMacro(*Assembler, *ast.Call) Instructions
}

type Assembler struct {
	Data DataBuilder
}

type Assembly struct {
	I Instructions
}

func ToZasmOp[u ~uint8](op u) Op {
	switch uint8(op) {
	case uint8(ast.AST_CMP_EQ):
		return OP_CMP_EQ
	case uint8(ast.AST_CMP_NQ):
		return OP_CMP_NQ
	case uint8(ast.AST_CMP_LT):
		return OP_CMP_LT
	case uint8(ast.AST_BOOL_AND):
		return OP_BOOL_AND
	case uint8(ast.AST_BOOL_OR):
		return OP_BOOL_OR
	case uint8(ast.AST_BOOL_NEGATE):
		return OP_BOOL_NEGATE
	default:
		panic(slipup.Createf("unknown ast operation %v", op))
	}
}

func (a *Assembler) Assemble(tree ast.Node) (Assembly, error) {
	var ass Assembly
	i, err := ast.Transform(a, tree)
	ass.I = i
	return ass, err
}

func (a *Assembler) Comparison(node *ast.Comparison) (Instructions, error) {
	lhs, lhErr := ast.Transform(a, node.LHS)
	rhs, rhErr := ast.Transform(a, node.RHS)
	if joined := errors.Join(lhErr, rhErr); joined != nil {
		panic(joined)
	}
	if lhs == nil || rhs == nil {
		panic(slipup.Createf("expected both arms to exist"))
	}

	return IntoIW(lhs).Union(rhs).WriteOp(ToZasmOp(node.Op)).I, nil
}

func (a *Assembler) BooleanOp(node *ast.BooleanOp) (Instructions, error) {
	if node.Op == ast.AST_BOOL_NEGATE {
		return a.negate(node.LHS)
	}

	lhs, lhErr := ast.Transform(a, node.LHS)
	rhs, rhErr := ast.Transform(a, node.RHS)
	if joined := errors.Join(lhErr, rhErr); joined != nil {
		panic(joined)
	}
	if lhs == nil || rhs == nil {
		panic(slipup.Createf("expected both arms to exist"))
	}
	return IntoIW(lhs).Union(rhs).WriteOp(ToZasmOp(node.Op)).I, nil
}

func (a *Assembler) Call(node *ast.Call) (Instructions, error) {
	if instructions, didFastOps := a.tryFastCall(node); didFastOps {
		return instructions, nil
	}

	return a.slowCall(node)
}

func (a *Assembler) Identifier(node *ast.Identifier) (Instructions, error) {
	name := a.Data.Names.Intern(node.Name)
	return IW().WriteLoadIdent(name).I, nil
}

func (a *Assembler) Literal(node *ast.Literal) (Instructions, error) {
	switch v := node.Value.(type) {
	case bool:
		return IW().WriteLoadBool(v).I, nil
	case float64:
		if instructions, didFastConst := a.fastConst(v); didFastConst {
			return instructions, nil
		}
		c := a.Data.Consts.Intern(Pack(v))
		return IW().WriteLoadConst(c).I, nil
	case string:
		str := a.Data.Strs.Intern(v)
		return IW().WriteLoadStr(str).I, nil
	default:
		panic("unreachable")
	}
}

func (a *Assembler) Empty(node *ast.Empty) (Instructions, error) {
	return nil, nil
}

func (a *Assembler) negate(node ast.Node) (Instructions, error) {
	instructions, _ := ast.Transform(a, node)
	return IntoIW(instructions).WriteNegate().I, nil
}

func (a *Assembler) callArgs(args []ast.Node) Instructions {
	iw := IW()
	for _, arg := range args {
		instrs, _ := ast.Transform(a, arg)
		iw = iw.Union(instrs)
	}
	return iw.I
}

func (a *Assembler) slowCall(node *ast.Call) (Instructions, error) {
	iw := IntoIW(a.callArgs(node.Args))
	name := a.Data.Names.Intern(node.Callee)
	return iw.WriteCall(name, len(node.Args)).I, nil
}

func (a *Assembler) tryFastCall(call *ast.Call) (Instructions, bool) {
	if call.Callee == "has" && call.Args[1].Type() != ast.AST_NODE_LITERAL {
		return nil, false
	}

	fastOp, hasFast := map[string]Op{
		"has":            OP_CHK_QTY,
		"load_setting":   OP_CHK_SET,
		"load_setting_2": OP_CHK_SET2,
		"trick_enabled":  OP_CHK_TRK,
	}[call.Callee]

	if !hasFast {
		return nil, false
	}

	switch fastOp {
	case OP_CHK_QTY:
		name := ast.MustAssertAs[*ast.Identifier](call.Args[0]).Name
		qty := ast.MustAssertAs[*ast.Literal](call.Args[1]).Value.(float64)

		tok := a.Data.Names.Intern(name)
		u16Tok := uint16(tok)

		if uint32(u16Tok) != uint32(tok) {
			return nil, false
		}

		var payload [3]uint8
		payload[0] = uint8(0x00FF & u16Tok)
		payload[1] = uint8((0xFF00 & u16Tok) >> 8)
		payload[2] = uint8(qty)
		return IW().Write(Encode(fastOp, payload)).I, true
	case OP_CHK_TRK:
		name := ast.MustAssertAs[*ast.Identifier](call.Args[0]).Name
		trk := a.Data.Names.Intern(name)
		return IW().Write(EncodeOpAndU24(fastOp, uint32(trk))).I, true
	case OP_CHK_SET:
		name := ast.MustAssertAs[*ast.Identifier](call.Args[0]).Name
		setting := a.Data.Names.Intern(name)
		return IW().Write(EncodeOpAndU24(OP_CHK_SET, uint32(setting))).I, true
	case OP_CHK_SET2:
		name1 := ast.MustAssertAs[*ast.Identifier](call.Args[0]).Name
		var name2 string
		switch arg := call.Args[1].(type) {
		case *ast.Identifier:
			name2 = arg.Name
			break
			// TODO: back this out after AST drops string literal aliases
		case *ast.Literal:
			if arg.Kind != ast.AST_LIT_STR {
				panic("expected literal string")
			}
			name2 = arg.Value.(string)
		default:
			panic("expected literal string or identifier")

		}
		blk := uint32(a.Data.Names.Intern(name1))
		sub := uint32(a.Data.Names.Intern(name2))
		const u12Mask uint32 = 0x00000FFF

		if blk&u12Mask != blk || sub&u12Mask != sub {
			return nil, false
		}

		var payload uint32 = 0
		payload |= (blk & u12Mask) << 12
		payload |= (sub & u12Mask)
		return IW().Write(EncodeOpAndU24(fastOp, payload)).I, true
	}

	return nil, false
}

func (a *Assembler) fastConst(v float64) (Instructions, bool) {
	u32 := uint32(v)

	if float64(u32) != v || U24_MASK < u32 {
		return nil, false
	}
	return nil, false
}
