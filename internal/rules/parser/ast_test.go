package parser

import "fmt"

type ValidatingVisitor struct {
	parsed Expression
}

func (v ValidatingVisitor) unexpected(expected Expression) error {
	return fmt.Errorf("\nExpected\n%+v\n\nActual\n%+v", expected, v.parsed)
}

func (v ValidatingVisitor) trace(err error) error {
	return fmt.Errorf("-> %+T %w", v.parsed, err)
}

func (v ValidatingVisitor) VisitBinOp(o *BinOp) error {
	e, ok := v.parsed.(*BinOp)
	if !ok || o.Op != e.Op {
		return v.unexpected(o)
	}

	err := Visit(ValidatingVisitor{e.Left}, o.Left)
	if err != nil {
		return v.trace(err)
	}
	err = Visit(ValidatingVisitor{e.Right}, o.Right)
	if err != nil {
		return v.trace(err)
	}

	return nil
}

func (v ValidatingVisitor) VisitBoolOp(o *BoolOp) error {
	e, ok := v.parsed.(*BoolOp)
	if !ok || o.Op != e.Op {
		return v.unexpected(o)
	}

	err := Visit(ValidatingVisitor{e.Left}, o.Left)
	if err != nil {
		return v.trace(err)
	}
	err = Visit(ValidatingVisitor{e.Right}, o.Right)
	if err != nil {
		return v.trace(err)
	}

	return nil
}

func (v ValidatingVisitor) VisitCall(o *Call) error {
	e, ok := v.parsed.(*Call)
	if !ok || len(o.Args) != len(e.Args) {
		return v.unexpected(o)
	}

	err := Visit(ValidatingVisitor{e.Callee}, o.Callee)
	if err != nil {
		return v.trace(err)
	}

	for i := range o.Args {
		err = Visit(ValidatingVisitor{e.Args[i]}, o.Args[i])
		if err != nil {
			return v.trace(err)
		}
	}

	return nil
}

func (v ValidatingVisitor) VisitIdentifier(o *Identifier) error {
	e, ok := v.parsed.(*Identifier)
	if !ok || o.Value != e.Value {
		return v.unexpected(o)
	}
	return nil
}

func (v ValidatingVisitor) VisitSubscript(o *Subscript) error {
	e, ok := v.parsed.(*Subscript)
	if !ok {
		return v.unexpected(o)
	}

	err := Visit(ValidatingVisitor{e.Target}, o.Target)
	if err != nil {
		return v.trace(err)
	}

	err = Visit(ValidatingVisitor{e.Index}, o.Index)
	if err != nil {
		return v.trace(err)
	}

	return nil
}

func (v ValidatingVisitor) VisitTuple(o *Tuple) error {
	e, ok := v.parsed.(*Tuple)
	if !ok || len(o.Elems) != len(e.Elems) {
		return v.unexpected(o)
	}

	for i := range o.Elems {
		err := Visit(ValidatingVisitor{e.Elems[i]}, o.Elems[i])
		if err != nil {
			return v.trace(err)
		}
	}

	return nil
}

func (v ValidatingVisitor) VisitUnary(o *UnaryOp) error {
	e, ok := v.parsed.(*UnaryOp)
	if !ok || o.Op != e.Op {
		return v.unexpected(o)
	}

	err := Visit(ValidatingVisitor{e.Target}, o.Target)
	if err != nil {
		err = v.trace(err)
	}
	return err
}

func (v ValidatingVisitor) VisitLiteral(o *Literal) error {
	e, ok := v.parsed.(*Literal)
	if !ok || o.Kind != e.Kind || o.Value != e.Value {
		return v.unexpected(o)
	}

	return nil
}
