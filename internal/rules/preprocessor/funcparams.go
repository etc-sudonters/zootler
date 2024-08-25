package preprocessor

import (
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/rules/visitor"
)

type rewritefuncparams map[string]parser.Expression

func (r rewritefuncparams) TransformBinOp(op *parser.BinOp) (parser.Expression, error) {
	return visitor.TransformBinOp(r, op)
}

func (r rewritefuncparams) TransformBoolOp(op *parser.BoolOp) (parser.Expression, error) {
	return visitor.TransformBoolOp(r, op)
}

func (r rewritefuncparams) TransformCall(call *parser.Call) (parser.Expression, error) {
	return visitor.TransformCall(r, call)
}

func (r rewritefuncparams) TransformIdentifier(id *parser.Identifier) (parser.Expression, error) {
	replacement, exists := r[id.Value]
	if exists {
		return replacement, nil
	}
	return id, nil
}

func (r rewritefuncparams) TransformSubscript(subscript *parser.Subscript) (parser.Expression, error) {
	return visitor.Transform(r, subscript)
}

func (r rewritefuncparams) TransformTuple(tup *parser.Tuple) (parser.Expression, error) {
	return visitor.TransformTuple(r, tup)
}

func (r rewritefuncparams) TransformUnary(op *parser.UnaryOp) (parser.Expression, error) {
	return visitor.TransformUnary(r, op)
}

func (r rewritefuncparams) TransformLiteral(lit *parser.Literal) (parser.Expression, error) {
	return lit, nil
}
