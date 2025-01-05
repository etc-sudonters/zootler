package optimizer

import (
	"fmt"
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/symbols"
)

/*
   can_use(Dins_Fire) => _is_magic_item(Dins_Fire) and Dins_Fire and Magic_Meter or False
   can_use(Dins_Fire) => (Dins_Fire == Dins_Fire or Dins_Fire == Farores_Wind or Dins_Fire == Nayrus_Love) and Dins_Fire and Magic_Meter
   can_use(Dins_Fire) => AnyOf(True, False, False) and Dins_Fire and Magic_Meter
   can_use => True and Dins_Fire and Magic_Meter
   can_use => has_every(Dins_Fire, Magic_Meter)


   can_use(item) => (_is_magic_item(item) and item and Magic_Meter)
           or (_is_adult_item(item) and is_adult and item)
           or (_is_magic_arrow(item) and is_adult and item and Bow and Magic_Meter)
           or (_is_child_item(item) and is_child and item)

   _is_magic_item(item) => item == Dins_Fire or item == Farores_Wind or item == Nayrus_Love or item == Lens_of_Truth

   _is_adult_item(item) => item == Bow or item == Megaton_Hammer or item == Iron_Boots or item == Hover_Boots or item == Hookshot or item == Longshot or item == Silver_Gauntlets or item == Golden_Gauntlets or item == Goron_Tunic or item == Zora_Tunic or item == Scarecrow or item == Distant_Scarecrow or item == Mirror_Shield

   _is_child_item(item) => item == Slingshot or item == Boomerang or item == Kokiri_Sword or item == Sticks or item == Deku_Shield,

   _is_magic_arrow(item) => item == Fire_Arrows or item == Light_Arrows or item == Blue_Fire_Arrows
*/

type funcTableKey string

const (
	FunctionTableKey funcTableKey = "function-table"
)

func InlineCalls(ctx *Context, syms *symbols.Table, funcs *ast.FunctionTable) ast.Rewriter {
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
	syms     *symbols.Table
	funcs    *ast.FunctionTable
	replacer *replacer
}

func (this *inliner) Invoke(node ast.Invoke, rewrite ast.Rewriting) (ast.Node, error) {
	target, casted := node.Target.(ast.Identifier)
	if !casted {
		return node, nil
	}

	fn, exists := this.funcs.Get(target)
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

	body, replaceErr := replacer.Rewriter().Rewrite(fn.Body)
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
