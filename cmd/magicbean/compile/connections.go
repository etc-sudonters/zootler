package compile

import (
	"fmt"
	"hash/fnv"
	"sudonters/zootler/cmd/magicbean/z2"
	"sudonters/zootler/internal/query"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/mido/optimizer"
	"sudonters/zootler/mido/symbols"
)

type CompilerConnector struct {
	tokens  z2.Tokens
	regions z2.Regions
	symbols *symbols.Table
	objects *objects.TableBuilder
}

func (this CompilerConnector) AddConnectionTo(regionName string, rule ast.Node) (*symbols.Sym, error) {
	hash := fnv.New64()
	fmt.Fprint(hash, regionName)
	ast.Hash64(rule, hash)
	nodeId := hash.Sum64()

	region := this.regions.RegionNamed(z2.Name(regionName))
	dest := this.regions.RegionNamed(z2.Name(fmt.Sprintf("Region##%x", nodeId)))
	token := this.tokens.Entity(z2.Name(fmt.Sprintf("Token##%x", nodeId)))

	dest.Attach(z2.HoldsToken(token.Entity()), z2.FixedPlacement{}, z2.Generated{})
	token.Attach(z2.HeldAt(dest.Entity()), z2.FixedPlacement{}, z2.Generated{}, z2.Event{})

	edge := this.regions.Connect(&region, &dest)
	edge.Attach(z2.ParsedSource(rule), z2.Generated{})

	symbol := this.symbols.Declare(string(token.Name), symbols.TOKEN)
	this.objects.AddPointer(symbol.Name, objects.Pointer(objects.OpaquePointer(0xdead), objects.PtrToken))

	return symbol, nil
}

func WithGeneratedConnections(engine query.Engine) mido.ConfigureCompiler {
	return func(outer *mido.CompileEnv) {
		outer.Symbols.Declare("at", symbols.COMPILER_FUNCTION)
		outer.Symbols.Declare("here", symbols.COMPILER_FUNCTION)
		outer.Symbols.Declare("has", symbols.FUNCTION)

		outer.Optimize.AddOptimizer(func(env *mido.CompileEnv) ast.Rewriter {
			named := z2.Named(engine)
			return optimizer.NewConnectionGeneration(
				env.Optimize.Context, env.Symbols,
				CompilerConnector{
					tokens: z2.Tokens{Entities: named},
					regions: z2.Regions{
						Regions:     named,
						Connections: z2.Tracked[z2.Connection](engine),
					},
					symbols: env.Symbols,
					objects: env.Objects,
				},
			)
		})

	}
}
