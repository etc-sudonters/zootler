package preprocessor

import (
	"fmt"
	"strings"
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/rules/parser"

	"github.com/etc-sudonters/substrate/slipup"
)

type ValuesTable map[internal.NormalizedStr]*parser.Literal

func (c ValuesTable) cheat(target string) *parser.Literal {
	if strings.HasPrefix(target, "logic") {
		return parser.BoolLiteral(false)
	}

	return nil
}

func (c ValuesTable) Resolve(target string) (*parser.Literal, bool) {
	if cheat := c.cheat(target); cheat != nil {
		return cheat, true
	}
	v, ok := c[internal.Normalize(target)]
	return v, ok
}

func (c ValuesTable) ResolveNested(target, index string) (*parser.Literal, bool) {
	key := internal.Normalize(strings.Join([]string{target, index}, "__"))
	value, exists := c[key]
	return value, exists
}

func (c ValuesTable) MustResolve(target string) (*parser.Literal, error) {
	if cheat := c.cheat(target); cheat != nil {
		return cheat, nil
	}
	v, exists := c[internal.Normalize(target)]
	if !exists {
		return nil, slipup.Createf("could not resolve %s from compile time environment", target)
	}
	return v, nil
}

func (c ValuesTable) MustResolveNested(target, index string) (*parser.Literal, error) {
	v, exists := c.ResolveNested(target, index)
	if !exists {
		return nil, slipup.Createf("could not resolve %s %s from compile time environment", target, index)
	}
	return v, nil
}

func (c ValuesTable) ResolveAsToken(target string) (*parser.Literal, error) {
	expr, err := c.MustResolve(target)
	if err != nil {
		return nil, err
	}
	return parser.AssertAs[*parser.Literal](expr)
}

type FunctionTable map[string]parser.FunctionDecl

func (fns FunctionTable) BuiltIn(name string) error {
	if _, exists := fns[name]; exists {
		return slipup.Createf("function %s is already declared", name)
	}
	fns[name] = parser.FunctionDecl{Identifier: name}
	return nil
}

func (fns FunctionTable) Decl(f parser.FunctionDecl) error {
	if _, exists := fns[f.Identifier]; exists {
		return slipup.Createf("function %s is already declared", f.Identifier)
	}

	fns[f.Identifier] = f
	return nil
}

func (fns FunctionTable) Retrieve(name string) (*parser.FunctionDecl, error) {
	f, exists := fns[name]
	if !exists {
		return nil, slipup.Createf("could not resolve function '%s'", name)
	}

	return &f, nil
}

type DelayedRules map[string][]DelayedRule

type DelayedRule struct {
	Target, Name string
	Rule         parser.Expression
}

func (d DelayedRules) Add(target string, rule parser.Expression) string {
	rules, exist := d[target]
	if !exist {
		rules = make([]DelayedRule, 4)
	}

	tokenname := fmt.Sprintf("%s %d", target, len(rules)+1)

	rules = append(rules, DelayedRule{
		Target: target,
		Name:   tokenname,
		Rule:   rule,
	})
	d[target] = rules
	return tokenname
}
