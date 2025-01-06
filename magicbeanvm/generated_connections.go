package magicbeanvm

import (
	"errors"
	"fmt"
	"slices"
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/optimizer"
	"sudonters/zootler/magicbeanvm/symbols"
)

const generatedConnectionsKey connectionsKey = "generated-connections"
const currentLocationKey currentKey = "current-key"

func SetCurrentLocation(ctx *optimizer.Context, where string) {
	ctx.Store(currentLocationKey, where)
}

func CompileGeneratedConnections(compiler *codegen) ([]CompiledSource, error) {
	var compiled []CompiledSource
	ctx := compiler.Optimize.Context
	connections := SwapGeneratedConnections(ctx)
	for size := connections.Size(); size > 0; size = connections.Size() {
		var offset int
		compiling := make([]CompiledSource, size)
		for source := range connections.All {
			var compileErr error
			compiling := &compiling[offset]
			compiling.Source = *source
			compiling.ByteCode, compileErr = compiler.CompileSource(&compiling.Source)
			if compileErr != nil {
				return compiled, compileErr
			}

			offset++
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

func ConnectionGeneration(ctx *optimizer.Context, tbl *symbols.Table) ConnectionGenerator {
	var conns ConnectionGenerator
	conns.at, conns.here = tbl.Declare("at", symbols.FUNCTION), tbl.Declare("here", symbols.FUNCTION)
	conns.tbl = tbl
	conns.ctx = ctx
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
	tbl      *symbols.Table
	ctx      *optimizer.Context
}

func (this ConnectionGenerator) At(args []ast.Node) (ast.Node, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("expected 2 args for at, received %d", len(args))
	}
	origin, ok := args[0].(ast.String)
	if !ok {
		return nil, fmt.Errorf("expected first argument to 'at' to be string: %#v", args[0])
	}
	return this.generate(string(origin), args[0]), nil
}

func (this ConnectionGenerator) Here(args []ast.Node) (ast.Node, error) {
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
	origination := this.tbl.Declare(originationName, symbols.LOCATION)
	connections := this.ctx.Retrieve(generatedConnectionsKey).(GeneratedConnections)
	connection := connections.appendto(origination.Name)

	destinationName := fmt.Sprintf("Token$%04d#%04d@%s", connections.generation, connection.rank, origination.Name)
	destination := this.tbl.Declare(destinationName, symbols.EVENT)

	connection.OriginatingRegion = origination.Name
	connection.Destination = destinationName
	connection.Ast = rule

	return ast.Invoke{
		Target: ast.IdentifierFrom(this.tbl.Declare("has", symbols.FUNCTION)),
		Args: []ast.Node{
			ast.IdentifierFrom(destination), ast.Number(1),
		},
	}
}
