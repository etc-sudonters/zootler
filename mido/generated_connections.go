package mido

import (
	"errors"
	"fmt"
	"slices"
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/optimizer"
	"sudonters/zootler/mido/symbols"
)

const generatedConnectionsKey connectionsKey = "generated-connections"
const currentLocationKey currentKey = "current-key"
const generatorKey generatorKeyType = "connection-generator"

type generatorKeyType string

func CompilerWithConnectionGeneration(register func(*CompileEnv) func(*symbols.Sym)) ConfigureCompiler {
	return func(outer *CompileEnv) {
		outer.Symbols.Declare("at", symbols.COMPILER_FUNCTION)
		outer.Symbols.Declare("here", symbols.COMPILER_FUNCTION)
		connections := ConnectionGeneration(outer.Optimize.Context, outer.Symbols, register(outer))
		outer.Optimize.Context.Store(generatorKey, connections)
		outer.Optimize.AddOptimizer(func(inner *CompileEnv) ast.Rewriter {
			connections := inner.Optimize.Context.Retrieve(generatorKey).(ConnectionGenerator)
			return ast.Rewriter{Invoke: func(node ast.Invoke, rewriting ast.Rewriting) (ast.Node, error) {
				symbol := ast.LookUpNodeInTable(inner.Symbols, node.Target)
				if symbol == nil {
					return node, nil
				}

				switch symbol.Name {
				case "at":
					return connections.At(node.Args, rewriting)
				case "here":
					return connections.Here(node.Args, rewriting)
				default:
					return node, nil
				}
			}}
		})
	}
}

func SetCurrentLocation(ctx *optimizer.Context, where string) {
	ctx.Store(currentLocationKey, where)
}

func CompileGeneratedConnections(codegen *codegen) ([]CompiledSource, error) {
	var compiled []CompiledSource
	ctx := codegen.Optimize.Context
	connections := SwapGeneratedConnections(ctx)
	for size := connections.Size(); size > 0; size = connections.Size() {
		compiling := make([]CompiledSource, size)
		var currentlyCompiling int
		for source := range connections.All {
			var compileErr error
			this := &compiling[currentlyCompiling]
			this.Source = *source
			this.ByteCode, compileErr = codegen.CompileSource(&this.Source)
			if compileErr != nil {
				return compiled, compileErr
			}
			currentlyCompiling++
		}
		compiled = slices.Concat(compiled, compiling)
		connections = SwapGeneratedConnections(ctx)
	}

	return compiled, nil
}

func SwapGeneratedConnections(ctx *optimizer.Context) GeneratedConnections {
	var conns GeneratedConnections
	stored := ctx.Retrieve(generatedConnectionsKey)
	if stored != nil {
		conns = stored.(GeneratedConnections)
	} else {
		conns = GeneratedConnections{make(map[string][]generated), 0}
	}
	ctx.Store(generatedConnectionsKey, GeneratedConnections{make(map[string][]generated), conns.generation + 1})
	return conns
}

func ConnectionGeneration(ctx *optimizer.Context, symbolTable *symbols.Table, register func(*symbols.Sym)) ConnectionGenerator {
	var conns ConnectionGenerator
	conns.at, conns.here = symbolTable.Declare("at", symbols.COMPILER_FUNCTION), symbolTable.Declare("here", symbols.COMPILER_FUNCTION)
	conns.symbols = symbolTable
	conns.ctx = ctx
	conns.register = register
	ctx.Store(generatedConnectionsKey, GeneratedConnections{make(map[string][]generated), 0})
	return conns
}

type generated struct {
	Source
	rank int
}

type GeneratedConnections struct {
	connections map[string][]generated
	generation  int
}

func (this GeneratedConnections) All(f func(*Source) bool) {
	for _, connections := range this.connections {
		for _, connection := range connections {
			source := connection.Source
			if !f(&source) {
				return
			}
		}
	}
}

func (this GeneratedConnections) appendto(parent string) *generated {
	items := this.connections[parent]
	rank := len(items)
	items = append(items, generated{
		rank: rank,
	})
	this.connections[parent] = items
	return &this.connections[parent][rank]
}

func (this GeneratedConnections) Size() int {
	var size int
	if len(this.connections) == 0 {
		return size
	}

	for _, items := range this.connections {
		size += len(items)
	}

	return size
}

type connectionsKey string
type currentKey string

type ConnectionGenerator struct {
	at, here *symbols.Sym
	symbols  *symbols.Table
	ctx      *optimizer.Context
	register func(*symbols.Sym)
}

func (this ConnectionGenerator) At(args []ast.Node, _ ast.Rewriting) (ast.Node, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("expected 2 args for at, received %d", len(args))
	}
	origin, ok := args[0].(ast.String)
	if !ok {
		return nil, fmt.Errorf("expected first argument to 'at' to be string: %#v", args[0])
	}
	return this.generate(string(origin), args[1]), nil
}

func (this ConnectionGenerator) Here(args []ast.Node, _ ast.Rewriting) (ast.Node, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expected 1 arg for here, received %d", len(args))
	}
	current := this.ctx.Retrieve(currentLocationKey)
	if current == nil {
		return nil, errors.New("expected current location to be set")
	}
	return this.generate(current.(string), args[0]), nil
}

func (this ConnectionGenerator) generate(originationName string, rule ast.Node) ast.Node {
	origination := this.symbols.Declare(originationName, symbols.LOCATION)
	connections := this.ctx.Retrieve(generatedConnectionsKey).(GeneratedConnections)
	connection := connections.appendto(origination.Name)

	destinationName := fmt.Sprintf("Token$%04d#%04d@%s", connections.generation, connection.rank, origination.Name)
	destination := this.symbols.Declare(destinationName, symbols.EVENT)
	this.register(destination)

	connection.OriginatingRegion = origination.Name
	connection.Destination = destinationName
	connection.Ast = rule

	return ast.Invoke{
		Target: ast.IdentifierFrom(this.symbols.Declare("has", symbols.FUNCTION)),
		Args: []ast.Node{
			ast.IdentifierFrom(destination), ast.Number(1),
		},
	}
}
