package query

import (
	"fmt"
	"strings"
	"sudonters/libzootr/galoshes/parse"
	"sudonters/libzootr/internal/table"
)

type Engine struct {
	tbl *table.Table

	attrs map[string]table.ColumnId
}

type tripletvalue struct {
	name  string
	value any
	typ   parse.Type
}

type tripletId struct {
	name string
	id   table.RowId
}

type triplet struct {
	id    tripletId
	attr  string
	value tripletvalue
}

func (this triplet) String() string {
	view := &strings.Builder{}
	view.WriteString("[ ")
	if this.id.name != "" {
		fmt.Fprintf(view, "$%s ", this.id.name)
	} else {
		fmt.Fprintf(view, "%d ", this.id.id)
	}
	fmt.Fprintf(view, "%s ", this.attr)
	if this.value.name != "" {
		fmt.Fprintf(view, "$%s", this.value.name)
	} else {
		fmt.Fprintf(view, "%v", this.value.value)
	}
	view.WriteString(" ]")
	return view.String()
}

type QueryPlan struct {
	returning []string
	triplets  []triplet
}

func (this QueryPlan) String() string {
	view := &strings.Builder{}
	fmt.Fprintf(view, "returning: %v\n", this.returning)
	fmt.Fprintln(view, "triplets")
	for _, triplet := range this.triplets {
		fmt.Fprintln(view, triplet)
	}
	return view.String()
}

type planner struct {
	qp    *QueryPlan
	decls map[string]*parse.RuleDeclNode
}

func BuildQueryPlan(ast parse.QueryNode) QueryPlan {
	var qp QueryPlan
	planner := planner{&qp, nil}
	planner.planQuery(ast)
	return qp
}

func (this *planner) planQuery(node parse.QueryNode) {
	switch node := node.(type) {
	case *parse.FindNode:
		this.planFind(node)
	default:
		panic(fmt.Errorf("unknown query type: %#v", node))
	}
	this.finalizePlan()
}

func (this *planner) planFind(node *parse.FindNode) {
	this.qp.returning = make([]string, len(node.Finding))
	for i, finding := range node.Finding {
		this.qp.returning[i] = finding.Name
	}
	if len(node.Rules) > 0 {
		this.decls = make(map[string]*parse.RuleDeclNode, len(node.Rules))
		for _, decl := range node.Rules {
			this.decls[decl.Name] = decl
		}
	}

	for _, clause := range node.Clauses {
		this.applyClause(clause)
	}

}

func (this *planner) applyClause(clause parse.ClauseNode) {
	switch clause := clause.(type) {
	case *parse.RuleClauseNode:
		this.applyRuleClause(clause)
	case *parse.TripletClauseNode:
		this.applyTripletClause(clause)
	default:
		panic(fmt.Errorf("unknown clause type: %#v", clause))
	}
}

func (this *planner) applyRuleClause(clause *parse.RuleClauseNode) {
	triplets := this.produceTriplets(clause, make(map[string]parse.ValueNode))
	this.qp.triplets = append(this.qp.triplets, triplets...)
}

func (this *planner) produceTriplets(clause *parse.RuleClauseNode, env map[string]parse.ValueNode) []triplet {
	var triplets []triplet
	decl := this.decls[clause.Name]
	args := make(map[string]parse.ValueNode, len(decl.Args))
	for i, arg := range decl.Args {
		name := arg.Name
		if value, exists := env[name]; exists {
			args[arg.Name] = value
		} else {
			args[arg.Name] = clause.Args[i]
		}
	}
	for _, clause := range decl.Clauses {
		var produced []triplet
		switch clause := clause.(type) {
		case *parse.RuleClauseNode:
			produced = this.produceTriplets(clause, args)
		case *parse.TripletClauseNode:
			produced = append(produced, this.produceBoundTriplet((*clause).TripletNode, args))
		default:
			panic(fmt.Errorf("unknown clause type %#v", clause))
		}
		triplets = append(triplets, produced...)
	}
	return triplets
}

func (this *planner) produceBoundTriplet(clause parse.TripletNode, args map[string]parse.ValueNode) triplet {
	if clause.Id.Var != nil {
		arg := args[clause.Id.Var.Name]
		if arg == nil {
			name := clause.Id.Var.Name
			view := parse.AstRender{Sink: &strings.Builder{}}
			clause := parse.TripletClauseNode{TripletNode: clause}
			view.VisitTripletClauseNode(&clause)
			panic(fmt.Errorf("%s isn't bound in %v\n%v", name, view.Sink.String(), args))
		}
		switch arg := arg.(type) {
		case *parse.NumberNode:
			clause.Id = &parse.EntityNode{Value: uint32(arg.Value), Type: parse.TypeNumber{}}
		case *parse.VarNode:
			clause.Id = &parse.EntityNode{Var: arg, Type: parse.TypeNumber{}}
		default:
			panic(fmt.Errorf("unsupported id value %#v", arg))
		}
	}

	if variable, isVar := clause.Value.(*parse.VarNode); isVar {
		arg := args[variable.Name]
		if arg == nil {
			panic(fmt.Errorf("%s isn't bound", clause.Id.Var.Name))
		}
		clause.Value = arg
	}

	return this.produceTriplet(clause)
}

func (this *planner) produceTriplet(clause parse.TripletNode) triplet {
	var trip triplet
	if clause.Id.Var != nil {
		trip.id.name = clause.Id.Var.Name
	} else {
		trip.id.id = table.RowId(clause.Id.Value)
	}
	trip.attr = clause.Attribute.Name
	bindValue(clause.Value, &trip.value)
	return trip
}

func bindValue(value parse.ValueNode, trip *tripletvalue) {
	switch value := value.(type) {
	case *parse.StringNode:
		trip.value = value.Value
		trip.typ = value.GetType()
	case *parse.NumberNode:
		trip.value = value.Value
		trip.typ = value.GetType()
	case *parse.BoolNode:
		trip.value = value.Value
		trip.typ = value.GetType()
	case *parse.VarNode:
		trip.name = value.Name
		trip.value = value.GetType()
	default:
		panic(fmt.Errorf("unsupported triplet value type %#v", value))
	}
}

func (this *planner) applyTripletClause(clause *parse.TripletClauseNode) {
	this.qp.triplets = append(this.qp.triplets, this.produceTriplet(clause.TripletNode))
}

func (_ planner) finalizePlan() {}
