package optimizer

import (
	"fmt"
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/symbols"
)

type funcTableKey string

const (
	FunctionTableKey funcTableKey = "function-table"
)

func InlineCalls(ctx *Context, syms *symbols.Table, funcs *ast.PartialFunctionTable) ast.Rewriter {
	replacer := replacer{
		scopes: make([]map[ast.Identifier]ast.Node, 16),
		sp:     -1,
	}
	inliner := &inliner{ctx, syms, funcs, &replacer}
	return ast.Rewriter{
		Invoke: inliner.Invoke,
	}
}

type inliner struct {
	ctx      *Context
	symbols  *symbols.Table
	funcs    *ast.PartialFunctionTable
	replacer *replacer
}

func (this *inliner) Invoke(node ast.Invoke, rewrite ast.Rewriting) (ast.Node, error) {
	symbol := ast.LookUpNodeInTable(this.symbols, node.Target)
	if symbol == nil {
		return node, nil
	}

	if symbol.Kind != symbols.SCRIPTED_FUNC {
		return node, nil
	}

	fn, exists := this.funcs.Get(symbol.Name)
	if !exists {
		return node, nil
	}

	scope, buildErr := buildreplacements(fn.Params, node.Args)
	if buildErr != nil {
		return nil, buildErr
	}

	replacer := this.replacer
	replacer.PushScope(scope)
	defer replacer.PopScope()

	rewriter := replacer.Rewriter()
	body, replaceErr := rewriter.Rewrite(fn.Body)
	if replaceErr != nil {
		return nil, replaceErr
	}

	return rewrite(body)
}

type replacements = map[ast.Identifier]ast.Node

func buildreplacements(idents []ast.Identifier, values []ast.Node) (replacements, error) {
	if len(idents) != len(values) {
		return nil, fmt.Errorf("expected %d args but got %d", len(idents), len(values))
	}

	replacements := make(replacements, len(idents))
	for i := range idents {
		replacements[idents[i]] = values[i]
	}

	return replacements, nil
}

type replacer struct {
	scopes []replacements
	sp     int
}

func (r *replacer) Rewriter() ast.Rewriter {
	return ast.Rewriter{Identifier: r.Identifier}
}

func (r *replacer) PushScope(replacements replacements) {
	r.sp++
	r.scopes[r.sp] = replacements
}

func (r *replacer) PopScope() {
	r.scopes[r.sp] = nil
	r.sp--
}

func (r *replacer) LookInTop(node ast.Identifier) (ast.Node, bool) {
	replacement, exists := r.scopes[r.sp][node]
	return replacement, exists
}

func (r *replacer) Identifier(node ast.Identifier, _ ast.Rewriting) (ast.Node, error) {
	if replacement, exists := r.LookInTop(node); exists {
		return replacement, nil
	}

	return node, nil
}
