package ast

import "errors"

type Visiting func(Node) error

func (v Visiting) All(ast []Node) error {
	var err error
	for i := range ast {
		err = v(ast[i])
		if err != nil {
			err = errors.Join(err)
		}
	}
	return err
}

func DontVisit[N Node]() VisitFunc[N] {
	return func(n N, _ Visiting) error {
		return nil
	}
}

type VisitFunc[T Node] func(T, Visiting) error

type Visitor struct {
	AnyOf      VisitFunc[AnyOf]
	Bool       VisitFunc[Bool]
	Compare    VisitFunc[Compare]
	Every      VisitFunc[Every]
	Identifier VisitFunc[Identifier]
	Invert     VisitFunc[Invert]
	Invoke     VisitFunc[Invoke]
	Number     VisitFunc[Number]
	String     VisitFunc[String]
}

func (v Visitor) Visit(node Node) error {

	switch node := node.(type) {
	case AnyOf:
		if v.AnyOf == nil {
			return v.anyof(node, v.Visit)
		}
		return v.AnyOf(node, v.Visit)
	case Bool:
		if v.Bool == nil {
			return v.boolean(node)
		}
		return v.Bool(node, v.Visit)
	case Compare:
		if v.Compare == nil {
			return v.compare(node, v.Visit)
		}
		return v.Compare(node, v.Visit)
	case Every:
		if v.Every == nil {
			return v.every(node, v.Visit)
		}
		return v.Every(node, v.Visit)
	case Identifier:
		if v.Identifier == nil {
			return v.identifier(node)
		}
		return v.Identifier(node, v.Visit)
	case Invert:
		if v.Invert == nil {
			return v.invert(node)
		}
		return v.Invert(node, v.Visit)
	case Invoke:
		if v.Invoke == nil {
			return v.invoke(node)
		}
		return v.Invoke(node, v.Visit)
	case Number:
		if v.Number == nil {
			return v.number(node)
		}
		return v.Number(node, v.Visit)
	case String:
		if v.String == nil {
			return v.str(node)
		}
		return v.String(node, v.Visit)
	default:
		panic("not implemented")
	}
}

func (v Visitor) anyof(anyof AnyOf, visit Visiting) error {
	return visit.All(anyof)
}

func (v Visitor) boolean(_ Bool) error {
	return nil
}

func (v Visitor) compare(compare Compare, visit Visiting) error {
	return visit.All([]Node{compare.LHS, compare.RHS})
}

func (v Visitor) every(every Every, visit Visiting) error {
	return visit.All(every)
}

func (v Visitor) identifier(_ Identifier) error {
	return nil
}

func (v Visitor) invert(invert Invert) error {
	return v.Visit(invert.Inner)
}

func (v Visitor) invoke(invoke Invoke) error {
	err := v.Visit(invoke.Target)

	for i := range invoke.Args {
		argErr := v.Visit(invoke.Args[i])
		if argErr != nil {
			err = errors.Join(argErr)
		}
	}
	return err
}

func (v Visitor) number(_ Number) error {
	return nil
}

func (v Visitor) str(_ String) error {
	return nil
}
