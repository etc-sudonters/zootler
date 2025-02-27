package galoshes

import (
	"fmt"
	"sudonters/libzootr/internal/skelly"
)

type Type interface {
	isType()
}

type Unifies interface {
	Unify(Type, Type, Substitutions) (Type, Substitutions, error)
}

type Substitutions map[TypeVar]Type

type TypeVar uint64

func (_ TypeVar) isType() {}

type TypeVoid struct{}

func (_ TypeVoid) isType() {}

type TypeNumber struct{}

func (_ TypeNumber) isType() {}

type TypeString struct{}

func (_ TypeString) isType() {}

type TypeBool struct{}

func (_ TypeBool) isType() {}

type TypeTuple struct {
	Types []Type
}

func (_ TypeTuple) isType() {}

type TypeFunc struct {
	Params []Type
	Return Type
}

func (_ TypeFunc) isType() {}

type funcunifier func(Type, Type, Substitutions) (Type, Substitutions, error)

func (this funcunifier) Unify(t1, t2 Type, subs Substitutions) (Type, Substitutions, error) {
	return this(t1, t2, subs)
}

var DefaultUnifier Unifies = funcunifier(Unify)

func Unify(t1, t2 Type, subs Substitutions) (Type, Substitutions, error) {
	panic("unimplemented")
}

var _ AstVisitor = (*TypeAnnotator)(nil)

type TypeAnnotator struct {
	envs skelly.Stack[*AstEnv]
}

func (this TypeAnnotator) pushEnv() {
	this.envs.Push(NewAstEnv())
}

func (this TypeAnnotator) popEnv() *AstEnv {
	return this.envs.Pop()
}

func (this TypeAnnotator) env() *AstEnv {
	return this.envs.Top()
}

type isTypeVarEnum uint8

const (
	isTypeVar = 1
	isType    = 2
)

func IsTv(t Type) bool {
	return isTv(t) == isTypeVar
}

func isTv(t Type) isTypeVarEnum {
	_, istv := t.(TypeVar)
	if istv {
		return isTypeVar
	}
	return isType
}

func (this TypeAnnotator) getOrTV(name string) (Type, isTypeVarEnum) {
	if ty := this.env().GetNamed(name); ty != nil {
		return ty, isTv(ty)
	}
	tv := NextTypeVar()
	this.env().AddNamed(name, tv)
	return tv, isTypeVar
}

func (this TypeAnnotator) addSubstitution(tv TypeVar, t Type) {
	subs := this.env().InitialSubstitutions
	subs[tv] = t
}

func (this TypeAnnotator) VisitFindNode(node *FindNode) {
	this.pushEnv()
	tv := NextTypeVar()
	node.Type = tv
	for i := range node.Finding {
		this.VisitVarNode(&node.Finding[i])
	}
	for _, clause := range node.Clauses {
		this.VisitClauseNode(clause)
	}
	for i := range node.Rules {
		this.VisitRuleDeclNode(&node.Rules[i])
	}

	tt := make([]Type, len(node.Finding))
	for i := range node.Finding {
		tt[i] = node.Finding[i].Type
	}
	this.addSubstitution(tv, TypeTuple{tt})
	node.Env = *(this.popEnv())
}
func (this TypeAnnotator) VisitInsertNode(node *InsertNode) {
	this.pushEnv()
	tv := NextTypeVar()
	node.Type = tv
	this.addSubstitution(tv, TypeVoid{})
	for i := range node.Inserting {
		this.VisitInsertTripletNode(&node.Inserting[i])
	}
	for _, clause := range node.Clauses {
		this.VisitClauseNode(clause)
	}
	for i := range node.Rules {
		this.VisitRuleDeclNode(&node.Rules[i])
	}
	node.Env = *(this.popEnv())
}

func (this TypeAnnotator) VisitClauseNode(node ClauseNode) {
	switch node := node.(type) {
	case *TripletClauseNode:
		this.VisitTripletClauseNode(node)
		return
	case *RuleClauseNode:
		this.VisitRuleClauseNode(node)
		return
	default:
		panic(fmt.Errorf("unknown clause type %T", node))
	}
}

func (this TypeAnnotator) VisitRuleDeclNode(node *RuleDeclNode) {
	node.Type, _ = this.getOrTV(node.Name)
	this.pushEnv()
	for i := range node.Args {
		this.VisitVarNode(&node.Args[i])
	}

	for i := range node.Clauses {
		this.VisitClauseNode(node.Clauses[i])
	}
	node.Env = *(this.popEnv())
}

func (this TypeAnnotator) VisitInsertTripletNode(node *InsertTripletNode) {
	this.visitTriplet(&node.TripletNode)
}

func (this TypeAnnotator) VisitTripletClauseNode(node *TripletClauseNode) {
	this.visitTriplet(&node.TripletNode)
}

func (this TypeAnnotator) VisitRuleClauseNode(node *RuleClauseNode) {
	node.Type, _ = this.getOrTV(node.Name)
	for i := range node.Args {
		this.VisitValueNode(node.Args[i])
	}
}

func (this TypeAnnotator) VisitValueNode(node ValueNode) {
	switch node := node.(type) {
	case *VarNode:
		this.VisitVarNode(node)
	case *EntityNode:
		this.VisitEntityNode(node)
	case *NumberNode:
		this.VisitNumber(node)
	case *BoolNode:
		this.VisitBoolNode(node)
	case *StringNode:
		this.VisitStringNode(node)
	default:
		panic(fmt.Errorf("unknown value type %T", node))
	}
}

func (this TypeAnnotator) VisitVarNode(node *VarNode) {
	node.Type, _ = this.getOrTV(node.Name)
}
func (this TypeAnnotator) VisitEntityNode(node *EntityNode) {
	node.Type = TypeNumber{}
	if node.Var != nil {
		var isTv isTypeVarEnum
		node.Type, isTv = this.getOrTV(node.Var.Name)
		if isTv == isTypeVar {
			tv := node.Type.(TypeVar)
			this.addSubstitution(tv, TypeNumber{})
		}
	}
}

func (this TypeAnnotator) VisitAttrNode(node *AttrNode) {
	node.Type, _ = this.getOrTV(node.Name)
}
func (this TypeAnnotator) VisitNumber(node *NumberNode)     {}
func (this TypeAnnotator) VisitBoolNode(node *BoolNode)     {}
func (this TypeAnnotator) VisitStringNode(node *StringNode) {}

func (this TypeAnnotator) visitTriplet(node *TripletNode) {
	tv := NextTypeVar()
	node.Type = tv
	this.VisitEntityNode(&node.Id)
	this.VisitAttrNode(&node.Attribute)
	this.VisitValueNode(node.Value)

	this.addSubstitution(tv, TypeTuple{[]Type{
		node.Id.Type,
		node.Attribute.Type,
		node.Value.GetType(),
	}})
}
