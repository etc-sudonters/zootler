package parse

import (
	"fmt"
	"sudonters/libzootr/internal/skelly"
)

func NewAnnotator() *TypeAnnotator {
	return AnnotatorWith(make(Substitutions))
}

func AnnotatorWith(subs Substitutions) *TypeAnnotator {
	ta := new(TypeAnnotator)
	ta.envs = skelly.NewStack[*AstEnv](4)
	ta.Substitutions = subs
	ta.pushEnv()
	return ta
}

type TypeAnnotator struct {
	envs skelly.Stack[*AstEnv]

	Substitutions Substitutions
}

func (this *TypeAnnotator) pushEnv() {
	this.envs.Push(NewAstEnv())
}

func (this *TypeAnnotator) popEnv() *AstEnv {
	return this.envs.Pop()
}

func (this *TypeAnnotator) env() *AstEnv {
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

func (this *TypeAnnotator) getOrTV(name string) (Type, isTypeVarEnum) {
	if ty := this.env().GetNamed(name); ty != nil {
		return ty, isTv(ty)
	}
	tv := NextTypeVar()
	this.env().AddNamed(name, tv)
	return tv, isTypeVar
}

func (this *TypeAnnotator) addSubstitution(tv TypeVar, t Type) {
	this.Substitutions[tv] = t
}

func (this *TypeAnnotator) VisitFindNode(node *FindNode) {
	this.pushEnv()
	for i := range node.Finding {
		this.VisitVarNode(node.Finding[i])
	}
	for i := range node.Rules {
		this.VisitRuleDeclNode(node.Rules[i])
	}
	for _, clause := range node.Clauses {
		this.VisitClauseNode(clause)
	}

	tt := make([]Type, len(node.Finding))
	for i := range node.Finding {
		tt[i] = node.Finding[i].Type
	}
	tv := NextTypeVar()
	node.Type = tv
	this.addSubstitution(tv, TypeTuple{tt})
	subsitutions, err := Unify(node.Type, TypeTuple{tt}, this.Substitutions)
	if err != nil {
		panic(err)
	}

	this.Substitutions = this.Substitutions.Combine(subsitutions)
	node.Env = *(this.popEnv())
}

func (this *TypeAnnotator) VisitInsertNode(node *InsertNode) {
	this.pushEnv()
	for i := range node.Inserting {
		this.VisitInsertTripletNode(node.Inserting[i])
	}
	for i := range node.Rules {
		this.VisitRuleDeclNode(node.Rules[i])
	}
	for _, clause := range node.Clauses {
		this.VisitClauseNode(clause)
	}
	tv := NextTypeVar()
	node.Type = tv
	this.addSubstitution(tv, TypeVoid{})
	node.Env = *(this.popEnv())
}

func (this *TypeAnnotator) VisitClauseNode(node ClauseNode) {
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

func (this *TypeAnnotator) VisitRuleDeclNode(node *RuleDeclNode) {
	var isTv isTypeVarEnum
	node.Type, isTv = this.getOrTV(node.Name)
	tt := make([]Type, len(node.Args))
	this.pushEnv()
	for i := range node.Args {
		this.VisitVarNode(node.Args[i])
		tt[i] = node.Args[i].GetType()
	}
	for i := range node.Clauses {
		this.VisitClauseNode(node.Clauses[i])
	}

	if isTv == isTypeVar {
		this.addSubstitution(node.Type.(TypeVar), TypeTuple{tt})
	}

	node.Env = *(this.popEnv())
}

func (this *TypeAnnotator) VisitInsertTripletNode(node *InsertTripletNode) {
	this.visitTriplet(&node.TripletNode)
}

func (this *TypeAnnotator) VisitTripletClauseNode(node *TripletClauseNode) {
	this.visitTriplet(&node.TripletNode)
}

func (this *TypeAnnotator) VisitRuleClauseNode(node *RuleClauseNode) {
	node.Type, _ = this.getOrTV(node.Name)
	for i := range node.Args {
		this.VisitValueNode(node.Args[i])
	}

	args := make([]Type, len(node.Args))
	for i := range node.Args {
		args[i] = node.Args[i].GetType()
	}

	subsitutions, err := Unify(node.Type, TypeTuple{args}, this.Substitutions)
	if err != nil {
		panic(err)
	}

	this.Substitutions = this.Substitutions.Combine(subsitutions)
}

func (this *TypeAnnotator) VisitValueNode(node ValueNode) {
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
	case *AttrNode:
		this.VisitAttrNode(node)
	default:
		panic(fmt.Errorf("unknown value type %T", node))
	}
}

func (this *TypeAnnotator) VisitVarNode(node *VarNode) {
	node.Type, _ = this.getOrTV(node.Name)
}

func (this *TypeAnnotator) VisitEntityNode(node *EntityNode) {
	node.Type = TypeNumber{}
	if node.Var != nil {
		var isTv isTypeVarEnum
		node.Var.Type, isTv = this.getOrTV(node.Var.Name)
		if isTv == isTypeVar {
			tv := node.Var.Type.(TypeVar)
			this.addSubstitution(tv, TypeNumber{})
		}
	}
}

func (this *TypeAnnotator) VisitAttrNode(node *AttrNode) {
	node.Type, _ = this.getOrTV(node.Name)
}

func (this *TypeAnnotator) VisitNumber(node *NumberNode)     {}
func (this *TypeAnnotator) VisitBoolNode(node *BoolNode)     {}
func (this *TypeAnnotator) VisitStringNode(node *StringNode) {}

func (this *TypeAnnotator) visitTriplet(node *TripletNode) {
	tv := NextTypeVar()
	node.Type = tv
	this.VisitEntityNode(node.Id)
	this.VisitAttrNode(node.Attribute)
	this.VisitValueNode(node.Value)

	this.addSubstitution(tv, TypeTuple{[]Type{
		node.Id.Type,
		node.Attribute.Type,
		node.Value.GetType(),
	}})
}

func Unify(t1, t2 Type, subs Substitutions) (Substitutions, error) {
	if t1.StrictlyEq(t2) {
		return subs, nil
	}

	if tv, isTv := t1.(TypeVar); isTv {
		return unifyVar(tv, t2, subs)
	}

	if tv, isTv := t2.(TypeVar); isTv {
		return unifyVar(tv, t1, subs)
	}

	tt1, isTT1 := t1.(TypeTuple)
	tt2, isTT2 := t2.(TypeTuple)
	if isTT1 == isTT2 {
		return unifyTuples(tt1, tt2, subs)
	}

	return nil, CannotUnify{t1, t2}
}

type CannotUnify struct {
	T1, T2 Type
}

func (this CannotUnify) Error() string {
	return fmt.Sprintf("cannot unify %s and %s", this.T1, this.T2)
}

func unifyTuples(tt1, tt2 TypeTuple, subs Substitutions) (Substitutions, error) {
	if len(tt1.Types) != len(tt2.Types) {
		return nil, CannotUnify{tt1, tt2}
	}

	var err error
	for i := range tt1.Types {
		subs, err = Unify(tt1.Types[i], tt2.Types[i], subs)
		if err != nil {
			return subs, err
		}
	}

	return subs, nil
}

func unifyVar(tv TypeVar, ty Type, subs Substitutions) (Substitutions, error) {
	if ty2, exists := subs[tv]; exists {
		return Unify(ty2, ty, subs)
	}
	if tv2, isTv := ty.(TypeVar); isTv {
		if ty, exists := subs[tv2]; exists {
			return Unify(tv, ty, subs)
		}
	}

	if reoccurs(tv, ty, subs) {
		return nil, TypeReoccurs{tv}
	}

	subs[tv] = ty
	return subs, nil
}

func reoccurs(tv TypeVar, ty Type, subs Substitutions) bool {
	if tv.StrictlyEq(ty) {
		return true
	}

	if tv2, isTv := ty.(TypeVar); isTv {
		if ty, exists := subs[tv2]; exists {
			return reoccurs(tv, ty, subs)
		}
	}

	if tt, isTT := ty.(TypeTuple); isTT {
		for i := range tt.Types {
			if reoccurs(tv, tt.Types[i], subs) {
				return true
			}
		}
	}

	return false
}

type TypeReoccurs struct {
	Var TypeVar
}

func (this TypeReoccurs) Error() string {
	return fmt.Sprintf("%s reoccurs in itself", this.Var)
}
