package galoshes

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/etc-sudonters/substrate/peruse"
)

var grammar = NewGrammar()

func ParseString(script string) (Ast, error) {
	lexer := NewLexer(script)
	parser := peruse.NewParser(grammar, lexer)
	return parser.ParseAt(peruse.LOWEST)
}

func NewGrammar() peruse.Grammar[Ast] {
	g := peruse.NewGrammar[Ast]()
	mapping := map[TokenType]peruse.Parselet[Ast]{
		//TOKEN_ASSIGN:    parseAssign,
		TOKEN_FIND:    parseFind,
		TOKEN_INSERT:  parseInsert,
		TOKEN_COMMENT: parseComment,
	}

	for tok, fn := range mapping {
		g.Parse(tok, fn)
	}
	return g
}

func unexpectedToken(expected, actual TokenType, pos peruse.Pos) error {
	return fmt.Errorf("unexpected token at %d: expected %s but found %s", pos, TokenTypeString(expected), TokenTypeString(actual))
}

func parseBracketExpr[T Ast](p *peruse.Parser[Ast], expect TokenType, parselet peruse.Parselet[Ast]) ([]T, error) {
	var elems []T
	var empT T

	if !p.Expect(TOKEN_OPEN_BRACKET) {
		return nil, unexpectedToken(TOKEN_OPEN_BRACKET, p.Next.Type, p.Next.Pos)
	}

	for p.Expect(expect) {
		elem, err := parselet(p)
		if err != nil {
			return nil, fmt.Errorf("while parsing bracket-expr: %w", err)
		}

		t, isT := elem.(T)
		if !isT {
			return nil, fmt.Errorf("expected to parse %T but parsed %T instead", empT, elem)
		}
		elems = append(elems, t)
	}

	if !p.Expect(TOKEN_CLOSE_BRACKET) {
		return elems, unexpectedToken(TOKEN_CLOSE_BRACKET, p.Next.Type, p.Next.Pos)
	}

	return elems, nil
}

func parseVariables(p *peruse.Parser[Ast]) ([]Variable, error) {
	return parseBracketExpr[Variable](p, TOKEN_VARIABLE, parseVariable)
}

func parseTriplets(p *peruse.Parser[Ast]) ([]Triplet, error) {
	return parseBracketExpr[Triplet](p, TOKEN_OPEN_BRACKET, parseTriplet)
}

func parseConstraints(p *peruse.Parser[Ast]) ([]Constraint, error) {
	elms, err := parseBracketExpr[Ast](p, TOKEN_OPEN_BRACKET, parseConstraint)
	constraints := make([]Constraint, len(elms))
	for i := range elms {
		constraints[i] = elms[i].(Constraint)
	}
	return constraints, err
}

func parseFind(p *peruse.Parser[Ast]) (Ast, error) {
	var find Find
	start := p.Cur.Pos

	returning, returnErr := parseVariables(p)
	if returnErr != nil {
		return nil, fmt.Errorf("while parsing find at pos %d: %w", start, returnErr)
	}
	find.Finding = returning

	if !p.Expect(TOKEN_WHERE) {
		return nil, unexpectedToken(TOKEN_WHERE, p.Next.Type, p.Next.Pos)
	}

	constraints, constraintErr := parseConstraints(p)
	if constraintErr != nil {
		return nil, fmt.Errorf("failed to parse constraint: %w", constraintErr)
	}

	find.Constraints = constraints

	if p.Expect(TOKEN_RULES) {
		derivations, derivationsErr := parseBracketExpr[DerivationDecl](p, TOKEN_OPEN_BRACKET, parseDerivationDecl)
		if derivationsErr != nil {
			return nil, fmt.Errorf("failed to parse derivation: %w", derivationsErr)
		}
		find.Derivations = derivations
	}

	return find, nil
}

func parseVarOr[T TripletPart](p *peruse.Parser[Ast], expect TokenType, parse peruse.Parselet[Ast]) (MaybeVar[T], error) {
	var empT MaybeVar[T]
	if p.Cur.Type == expect {
		parsed, err := parse(p)
		if err != nil {
			return empT, fmt.Errorf("while parsing literal: %w", err)
		}
		t := parsed.(T)
		return MaybeVar[T]{Part: t}, nil
	}
	if p.Cur.Type == TOKEN_VARIABLE {
		variable, err := parseVariable(p)
		if err != nil {
			return empT, fmt.Errorf("while parsing variable: %w", err)
		}
		return MaybeVar[T]{Var: variable.(Variable)}, nil
	}

	return empT, fmt.Errorf("expected to parse %s or variable, but found %s", TokenTypeString(expect), TokenTypeString(p.Cur.Type))
}

func parseVarOrLiteral(p *peruse.Parser[Ast]) (MaybeVar[Literal], error) {
	if p.Cur.Type == TOKEN_VARIABLE {
		variable, _ := parseVariable(p)
		return MaybeVar[Literal]{Var: variable.(Variable)}, nil
	} else if p.Cur.Type == TOKEN_STRING {
		str, _ := parseString(p)
		return MaybeVar[Literal]{Part: str.(Literal)}, nil
	} else if p.Cur.Type == TOKEN_NUMBER {
		num, err := parseNumber(p)
		if err != nil {
			return MaybeVar[Literal]{}, err
		}
		return MaybeVar[Literal]{Part: num.(Literal)}, nil
	} else if p.Cur.Type == TOKEN_NIL {
		nil_, _ := parseNil(p)
		return MaybeVar[Literal]{Part: nil_.(Literal)}, nil
	} else if p.Cur.Type == TOKEN_TRUE || p.Cur.Type == TOKEN_FALSE {
		bool_, _ := parseBool(p)
		return MaybeVar[Literal]{Part: bool_.(Literal)}, nil
	}

	return MaybeVar[Literal]{}, unexpectedToken(TOKEN_VARIABLE, p.Cur.Type, p.Cur.Pos)
}

func parseConstraint(p *peruse.Parser[Ast]) (Ast, error) {
	if p.Expect(TOKEN_DERIVE) {
		derive, err := parseDerivationInvocation(p)
		if err != nil {
			err = fmt.Errorf("while parsing derivation constraint: %w", err)
		}
		return Constraint(derive.(DerivationInvoke)), err
	}

	if p.Expect(TOKEN_NUMBER) || p.Expect(TOKEN_VARIABLE) {
		triplet, err := parseTriplet(p)
		if err != nil {
			err = fmt.Errorf("while parsing triplet constraint: %w", err)
		}
		return Constraint(triplet.(Triplet)), err
	}

	return nil, fmt.Errorf("expected %s or %s but found %s at pos %d", "derivation", "triplet", TokenTypeString(p.Next.Type), p.Next.Pos)
}

func parseTriplet(p *peruse.Parser[Ast]) (Ast, error) {
	var triplet Triplet
	id, idErr := parseVarOr[Number](p, TOKEN_NUMBER, parseNumber)
	if idErr != nil {
		return triplet, fmt.Errorf("failed to parse triplet[0]: %w", idErr)
	}
	p.Consume()
	attr, attrErr := parseAttribute(p)
	if attrErr != nil {
		return triplet, fmt.Errorf("failed to parse triplet[1]: %w", attrErr)
	}
	p.Consume()
	value, valueErr := parseVarOrLiteral(p)
	if valueErr != nil {
		return triplet, fmt.Errorf("failed to parse triplet[2]: %w", valueErr)
	}

	if !p.Expect(TOKEN_CLOSE_BRACKET) {
		return triplet, unexpectedToken(TOKEN_CLOSE_BRACKET, p.Next.Type, p.Next.Pos)
	}
	triplet.Id = id
	triplet.Attr = attr.(Attribute)
	triplet.Value = value
	return triplet, nil
}

func parseDerivationInvocation(p *peruse.Parser[Ast]) (Ast, error) {
	// [:the-name $maybe-var "maybe-literal"]
	var invoke DerivationInvoke

	// we're already here
	invoke.Name = p.Cur.Literal
	if p.Expect(TOKEN_CLOSE_BRACKET) {
		return nil, errors.New("expected at least one variable in derivation")
	}

	for !p.Expect(TOKEN_CLOSE_BRACKET) {
		p.Consume()
		val, valErr := parseVarOrLiteral(p)
		if valErr != nil {
			return nil, fmt.Errorf("while parsing derivation invocation: %w", valErr)
		}
		invoke.Accept = append(invoke.Accept, val)
	}

	return invoke, nil
}

func parseDerivationDecl(p *peruse.Parser[Ast]) (Ast, error) {
	/*
	   [ [:name $accept1 $accept2] [triplet, ...] ]
	*/
	var derive DerivationDecl

	if !p.Expect(TOKEN_DERIVE) {
		return nil, fmt.Errorf("while parsing derivation decl: %w", unexpectedToken(TOKEN_DERIVE, p.Next.Type, p.Next.Pos))

	}

	derive.Name = p.Cur.Literal
	for !p.Expect(TOKEN_CLOSE_BRACKET) {
		p.Consume()
		val, valErr := parseVariable(p)
		if valErr != nil {
			return nil, fmt.Errorf("while parsing derivation accepting: %w", valErr)
		}
		derive.Accepting = append(derive.Accepting, val.(Variable))
	}

	// TODO allow a bare constraint iff its the only thing that follows
	constraints, constraintErr := parseConstraints(p)
	if constraintErr != nil {
		return nil, fmt.Errorf("while parsing derivation constraints: %w", constraintErr)
	}
	derive.Constraints = constraints
	return derive, nil
}

func parseInsert(p *peruse.Parser[Ast]) (Ast, error) {
	var insert Insert

	triplets, tripletErr := parseTriplets(p)
	if tripletErr != nil {
		return nil, tripletErr
	}
	insert.Inserting = triplets

	if p.Expect(TOKEN_WHERE) {
		constraints, constraintErr := parseConstraints(p)
		if constraintErr != nil {
			return nil, fmt.Errorf("failed to parse constraints: %w", constraintErr)
		}

		insert.Constraints = constraints
	}

	if len(insert.Constraints) == 0 {
		return insert, nil
	}

	if p.Expect(TOKEN_RULES) {
		derivations, derivationsErr := parseBracketExpr[DerivationDecl](p, TOKEN_OPEN_BRACKET, parseDerivationDecl)
		if derivationsErr != nil {
			return nil, fmt.Errorf("failed to parse derivations: %w", derivationsErr)
		}
		insert.Derivations = derivations
	}

	return insert, nil

}

func parseBool(p *peruse.Parser[Ast]) (Ast, error) {
	lit := Literal{Value: p.Cur.Literal == trueWord, Kind: LiteralKindBool}
	return lit, nil
}

func parseNil(_ *peruse.Parser[Ast]) (Ast, error) {
	lit := Literal{Value: nil, Kind: LiteralKindNil}
	return lit, nil
}

func parseDiscard(_ *peruse.Parser[Ast]) (Ast, error) {
	return Discard, nil
}

func parseVariable(p *peruse.Parser[Ast]) (Ast, error) {
	variable := Variable(p.Cur.Literal)
	return variable, nil
}

func parseAttribute(p *peruse.Parser[Ast]) (Ast, error) {
	attr := Attribute(p.Cur.Literal)
	return attr, nil
}

func parseString(p *peruse.Parser[Ast]) (Ast, error) {
	lit := Literal{Value: p.Cur.Literal, Kind: LiteralKindString}
	return lit, nil
}

func parseNumber(p *peruse.Parser[Ast]) (Ast, error) {
	n, err := strconv.ParseFloat(p.Cur.Literal, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse %q as number", p.Cur)
	}
	return Literal{Value: n, Kind: LiteralKindNumber}, nil
}

func parseComment(p *peruse.Parser[Ast]) (Ast, error) {
	comment := Comment(p.Cur.Literal)
	return comment, nil
}
