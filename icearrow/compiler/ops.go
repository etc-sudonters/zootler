package compiler

type IceArrowOp uint8

const (
	// 0x00 operational
	IA_NOP    IceArrowOp = 0x00
	IA_RETURN            = 0x0A // Top of stack is our return value
	// 0x10 load
	IA_LOAD_CONST  = 0x11 // next two bytes [lo, hi] u16 handle
	IA_LOAD_SYMBOL = 0x12 // next two bytes [lo, hi] u16 handle
	IA_LOAD_TRUE   = 0x13 // push packed true to top of stack
	IA_LOAD_FALSE  = 0x14 // push packed false to top of stack
	IA_LOAD_IMMED  = 0x15 // pack next byte as I32 and push to top of stack
	IA_LOAD_IMMED2 = 0x16 // pack [lo, hi] u16 as I32 and push to top of stack
	// 0x20 Reductions
	IA_REDUCE_ALL = 0x23 // top of stack is pop count
	IA_REDUCE_ANY = 0x24 // _MUST_ pop count but can short circuit
	// 0x30 Jumps
	IA_JUMP_UNCOND = 0x31 // top of stack is jump dest
	IA_JUMP_TRUE   = 0x32 // top of stack is dest, next is cond
	IA_JUMP_FALSE  = 0x33 // top of stack is dest, next is cond
	// 0x40 Calls
	IA_CALL_0 = 0x41 // stack: [ symbol ]
	IA_CALL_1 = 0x42 // stack: [ symbol arg1 ]
	IA_CALL_2 = 0x43 // stack: [ symbol arg2 arg1 ]
	// 0x60 Fast Calls
	IA_HAS_QTY    = 0x69 // next THREE bytes [ lo, hi, qty ] -> [ tok, qty ]
	IA_HAS_ALL    = 0x6A // top of stack is pop count
	IA_HAS_ANY    = 0x6B // _MUST_ pop all but may short circuit
	IA_IS_CHILD   = 0x6C
	IA_IS_ADULT   = 0x6D
	IA_HAS_BOTTLE = 0x6E
	// 0xF0 Temporary ops
	TEMP_IA_LOAD_STR = 0xF2
)
