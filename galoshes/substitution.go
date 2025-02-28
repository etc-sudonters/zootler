package galoshes

import (
	"fmt"

	"github.com/etc-sudonters/substrate/skelly/bitset64"
)

type Substitutions map[TypeVar]Type

func (this Substitutions) Combine(other Substitutions) Substitutions {
	new := make(Substitutions, max(len(this), len(other)))
	for tv, ty := range this {
		new[tv] = ty
	}
	for tv, ty := range other {
		new[tv] = Substitute(ty, this)
	}

	return new
}

func Substitute(t Type, subs Substitutions) Type {
	switch t := t.(type) {
	case TypeVar:
		var terminal Type
		set := bitset64.Bitset{}
		seen := []TypeVar{}
		for {
			if !set.Set(uint64(t)) {
				// todo, return err instead
				panic(fmt.Errorf("already seen %v: %v\n%v", t, seen, subs))
			}
			seen = append(seen, t)
			typ, exists := subs[t]
			if !exists {
				terminal = t
				break
			}
			tv, isTv := typ.(TypeVar)
			if !isTv {
				terminal = typ
				break
			}
			t = tv
		}
		termVar, _ := terminal.(TypeVar)
		for _, tv := range seen {
			if tv == termVar {
				continue
			}
			subs[tv] = terminal
		}
		if tt, isTT := terminal.(TypeTuple); isTT {
			return Substitute(tt, subs)
		}
		return terminal
	case TypeTuple:
		tt := TypeTuple{make([]Type, len(t.Types))}
		for i := range t.Types {
			tt.Types[i] = Substitute(t.Types[i], subs)
		}
		return tt
	case TypeString, TypeNumber, TypeBool:
		return t
	default:
		panic(fmt.Errorf("unknown type %#v", t))
	}
}

type subber struct {
	subs Substitutions
}

func (this *subber) VisitFindNode(node *FindNode) {
	for i := range node.Finding {
		this.VisitVarNode(node.Finding[i])
	}
	for i := range node.Clauses {
		this.VisitClauseNode(node.Clauses[i])
	}
	for i := range node.Rules {
		this.VisitRuleDeclNode(node.Rules[i])
	}

	node.Type = Substitute(node.Type, this.subs)
}

func (this *subber) VisitInsertNode(node *InsertNode) {
	for i := range node.Inserting {
		this.VisitInsertTripletNode(node.Inserting[i])
	}
	for i := range node.Clauses {
		this.VisitClauseNode(node.Clauses[i])
	}
	for i := range node.Rules {
		this.VisitRuleDeclNode(node.Rules[i])
	}
	node.Type = Substitute(node.Type, this.subs)
}

func (this *subber) VisitInsertTripletNode(node *InsertTripletNode) {
	this.visitTriplet(&node.TripletNode)
}

func (this *subber) VisitRuleDeclNode(node *RuleDeclNode) {
	node.Type = Substitute(node.Type, this.subs)
	for i := range node.Args {
		this.VisitVarNode(node.Args[i])
	}
	for i := range node.Clauses {
		this.VisitClauseNode(node.Clauses[i])
	}
}

func (this *subber) VisitClauseNode(node ClauseNode) {
	switch node := node.(type) {
	case *TripletClauseNode:
		this.VisitTripletClauseNode(node)
	case *RuleClauseNode:
		this.VisitRuleClauseNode(node)
	}
}

func (this *subber) VisitTripletClauseNode(node *TripletClauseNode) {
	this.visitTriplet(&node.TripletNode)
}

func (this *subber) VisitRuleClauseNode(node *RuleClauseNode) {
	for i := range node.Args {
		this.VisitValueNode(node.Args[i])
	}
	node.Type = Substitute(node.Type, this.subs)
}

func (this *subber) VisitAttrNode(node *AttrNode) {
	node.Type = Substitute(node.Type, this.subs)
}

func (this *subber) VisitValueNode(node ValueNode) {
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

func (this *subber) VisitVarNode(node *VarNode) {
	node.Type = Substitute(node.Type, this.subs)
}

func (this *subber) VisitEntityNode(node *EntityNode) {
	if node.Var != nil {
		if node.Var.Type == nil {
			panic(fmt.Errorf("%s has var but no type", node.Var.Name))
		}
		node.Var.Type = Substitute(node.Var.Type, this.subs)
	}
}

func (this *subber) VisitNumber(node *NumberNode)     {}
func (this *subber) VisitBoolNode(node *BoolNode)     {}
func (this *subber) VisitStringNode(node *StringNode) {}

func (this *subber) visitTriplet(node *TripletNode) {
	this.VisitValueNode(node.Id)
	this.VisitValueNode(node.Attribute)
	this.VisitValueNode(node.Value)
	node.Type = Substitute(node.Type, this.subs)
}
