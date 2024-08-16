package rules

import (
	"sudonters/zootler/internal"
	"sudonters/zootler/internal/components"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/internal/rules/parser"
	"sudonters/zootler/internal/slipup"
	"sudonters/zootler/internal/table"
)

type RuleAliases map[string]struct{}
type RandomizerSettings map[string]any

func (r RandomizerSettings) GetLiteral(name string) (*parser.Literal, bool) {
	if l, exists := r[name]; exists {
		if lit, casted := l.(*parser.Literal); casted {
			return lit, true
		}
	}

	return nil, false
}

func (r RandomizerSettings) GetNestedLiteral(name, subname string) (*parser.Literal, bool) {
	if nest, nestExists := r[name]; nestExists {
		if items, nestCast := nest.(map[string]any); nestCast {
			if l, exists := items[name]; exists {
				if lit, casted := l.(*parser.Literal); casted {
					return lit, true
				}
			}
		}
	}

	return nil, false
}

type AdvancementTokens struct {
	tokens map[string]uint64
}

func (a *AdvancementTokens) Retrieve(name string) *NamedAdvancementToken {
	normaled := internal.Normalize(name)
	if id, tracked := a.tokens[normaled]; tracked {
		return &NamedAdvancementToken{
			TokenId:  table.RowId(id),
			Name:     components.Name(name),
			Normaled: normaled,
		}
	}
	return nil
}

func IndexAdvancementTokens(eng query.Engine) (*AdvancementTokens, error) {
	at := AdvancementTokens{
		tokens: make(map[string]uint64, 256),
	}

	q := eng.CreateQuery()
	q.Load(query.MustAsColumnId[components.Name](eng))
	q.Exists(query.MustAsColumnId[components.Advancement](eng))

	tokens, tokenLoadErr := eng.Retrieve(q)
	if tokenLoadErr != nil {
		return nil, slipup.Describe(tokenLoadErr, "while indexing advancement tokens")
	}

	for tokens.MoveNext() {
		row := tokens.Current()
		normaled := internal.Normalize(row.Values[0].(components.Name))
		at.tokens[normaled] = uint64(row.Id)
	}

	return &at, nil
}

type NamedAdvancementToken struct {
	TokenId  table.RowId
	Name     components.Name
	Normaled string
}

func (n *NamedAdvancementToken) init(r *table.RowTuple) {
	n.Name = r.Values[0].(components.Name)
	n.TokenId = r.Id
}
