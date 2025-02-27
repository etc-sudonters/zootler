package galoshes

import (
	"fmt"
)

type Ast interface {
}

type Constraint interface {
	Ast
	isConstraint() bool
}

type Find struct {
	Finding     []Variable
	Constraints []Constraint
	Derivations []DerivationDecl
}

type Insert struct {
	Inserting   []Triplet
	Constraints []Constraint
	Derivations []DerivationDecl
}

type DerivationDecl struct {
	Name        string
	Accepting   []Variable
	Constraints []Constraint
}

type DerivationInvoke struct {
	Name   string
	Accept []MaybeVar[Literal] // TODO accept attributes as well
}

func (this DerivationInvoke) isConstraint() bool { return true }

func Invoke(name string, accept ...MaybeVar[Literal]) DerivationInvoke {
	return DerivationInvoke{name, accept}
}

func (this DerivationInvoke) Eq(other DerivationInvoke) bool {
	if this.Name != other.Name {
		return false
	}

	if len(this.Accept) != len(other.Accept) {
		return false
	}

	for i, ours := range this.Accept {
		theirs := other.Accept[i]
		if !ours.Eq(theirs) {
			return false
		}
	}

	return true
}

type Triplet struct {
	Id    MaybeVar[Number]
	Attr  Attribute
	Value MaybeVar[Literal]
}

func (this Triplet) Eq(other Triplet) bool {
	return this.Id.Eq(other.Id) && this.Attr == other.Attr && this.Value.Eq(other.Value)
}

func (this Triplet) isConstraint() bool { return true }

type MaybeVar[T TripletPart] struct {
	Part T
	Var  Variable
}

func VarOf[T TripletPart](name string) MaybeVar[T] {
	return MaybeVar[T]{Var: Variable(name)}
}

func Part[T TripletPart](part T) MaybeVar[T] {
	return MaybeVar[T]{Part: part}
}

func (this MaybeVar[T]) Eq(other MaybeVar[T]) bool {
	if this.Var != other.Var {
		return false
	}
	return areEqualPart(this.Part, other.Part)
}

func areEqualPart(ours, theirs any) bool {
	switch ours := ours.(type) {
	case Literal:
		theirs, isLiteral := theirs.(Literal)
		return isLiteral && ours.Eq(theirs)
	case Attribute:
		theirs, isAttr := theirs.(Attribute)

		return isAttr && ours == theirs
	case Number:
		theirs, isNum := theirs.(Number)
		return isNum && ours == theirs
	default:
		panic(fmt.Errorf("unknown part type %#v", ours))
	}
}

type TripletPart interface {
	Number | Attribute | Literal
	Ast
}

type Variable string

type Attribute string

func (this Attribute) String() string {
	return string(this)
}

type LiteralKind string
type Literal struct {
	Kind  LiteralKind
	Value any
}

func (this Literal) String() string {
	return fmt.Sprintf("Literal{%v}", this.Value)
}

func (this Literal) Eq(other Literal) bool {
	return this == other
}

type literals interface {
	String | Number | Boolean | Nil
}

type String string

func (this String) String() string {
	return fmt.Sprintf("%q", string(this))
}

type Number float64

func (this Number) String() string {
	return fmt.Sprintf("%f", this)
}

type Boolean bool

func (this Boolean) String() string {
	return fmt.Sprintf("%t", this)
}

type Nil struct{}

func (this Nil) String() string {
	return "nil"
}

type Comment string

func (this Comment) String() string {
	return ";" + string(this)
}

const (
	LiteralKindBool   = "bool"
	LiteralKindString = "string"
	LiteralKindNumber = "number"
	LiteralKindNil    = "nil"

	Discard Variable = "_"
)

func LiteralString(str string) Literal {
	return Literal{Kind: LiteralKindString, Value: str}
}

func LiteralBool(b bool) Literal {
	return Literal{Kind: LiteralKindBool, Value: b}
}

func LiteralNumber(num float64) Literal {
	return Literal{Kind: LiteralKindNumber, Value: num}
}

func LiteralNil() Literal {
	return Literal{Kind: LiteralKindNil, Value: Nil{}} // distguinish from an _actual_ nil
}
