package analysis

import (
	"sudonters/zootler/icearrow/ast"
)

func analyze(node ast.Node, ctx *AnalysisContext) report {
	var a analyzer
	a.canExpand = true
	a.canPromote = true
	a.ctx = ctx
	ast.Visit(&a, node)
	return a.report
}

type report struct {
	expansions, promotions, compares, branches bool
}

type analyzer struct {
	canExpand, canPromote bool
	report                report
	ctx                   *AnalysisContext
}

func (a *analyzer) Comparison(node *ast.Comparison) error {
	a.report.compares = true
	if err := ast.Visit(a, node.LHS); err != nil {
		return err
	}
	if err := ast.Visit(a, node.RHS); err != nil {
		return err
	}
	return nil
}

func (a *analyzer) BooleanOp(node *ast.BooleanOp) error {
	a.report.branches = true
	if err := ast.Visit(a, node.LHS); err != nil {
		return err
	}
	if err := ast.Visit(a, node.RHS); err != nil {
		return err
	}
	return nil
}

func (a *analyzer) Call(node *ast.Call) error {
	isExpandable := a.ctx.isExpandable(node.Callee)
	a.report.expansions = a.report.expansions || isExpandable
	return nil
}

func (a *analyzer) Identifier(node *ast.Identifier) error {
	isExpandable := a.canExpand && a.ctx.isExpandable(node.Name)
	a.report.expansions = a.report.expansions || isExpandable
	return nil
}

func (a *analyzer) Literal(node *ast.Literal) error {
	return nil
}

func (a *analyzer) Empty(node *ast.Empty) error { return nil }

func copyfragment(node ast.Node) (ast.Node, error) {
	return ast.Transform(copier{}, node)
}

type copier struct{}

func (c copier) Comparison(node *ast.Comparison) (ast.Node, error) {
	op := new(ast.Comparison)
	op.Op = node.Op
	op.LHS, _ = ast.Transform(c, node.LHS)
	op.RHS, _ = ast.Transform(c, node.RHS)
	return op, nil
}

func (c copier) BooleanOp(node *ast.BooleanOp) (ast.Node, error) {
	op := new(ast.BooleanOp)
	op.Op = node.Op
	op.LHS, _ = ast.Transform(c, node.LHS)
	op.RHS, _ = ast.Transform(c, node.RHS)
	return op, nil
}

func (c copier) Call(node *ast.Call) (ast.Node, error) {
	call := new(ast.Call)
	call.Callee = node.Callee
	call.Args = make([]ast.Node, len(node.Args))
	for idx, arg := range node.Args {
		call.Args[idx], _ = ast.Transform(c, arg)
	}
	return call, nil
}

func (c copier) Identifier(node *ast.Identifier) (ast.Node, error) {
	ident := new(ast.Identifier)
	ident.Kind = node.Kind
	ident.Name = node.Name
	return ident, nil
}

func (c copier) Literal(node *ast.Literal) (ast.Node, error) {
	lit := new(ast.Literal)
	lit.Kind = node.Kind
	lit.Value = node.Value
	return lit, nil
}

func (c copier) Empty(node *ast.Empty) (ast.Node, error) {
	return node, nil
}
