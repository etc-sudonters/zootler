package magicbeanvm

import (
	"errors"
	"fmt"
	"sudonters/zootler/internal/ruleparser"
	"sudonters/zootler/internal/settings"
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/compiler"
	"sudonters/zootler/magicbeanvm/objects"
	"sudonters/zootler/magicbeanvm/optimizer"
	"sudonters/zootler/magicbeanvm/symbols"
	"sudonters/zootler/magicbeanvm/vm"

	"github.com/etc-sudonters/substrate/peruse"
)

type SourceKind string
type SourceString string

func (this SourceKind) AsSymbolKind() symbols.Kind {
	switch this {
	case SourceCheck:
		return symbols.LOCATION
	case SourceEvent:
		return symbols.EVENT
	case SourceTransit:
		return symbols.TRANSIT
	default:
		panic("unreachable")
	}
}

const (
	_             SourceKind = ""
	SourceCheck              = "SourceCheck"
	SourceEvent              = "SourceEvent"
	SourceTransit            = "SourceTransit"
)

type CompiledSource struct {
	Source
	ByteCode compiler.ByteCode
}

func ptr[T any](what T) *T {
	return &what
}

type Source struct {
	Kind              SourceKind
	String            SourceString
	Ast               ast.Node
	OriginatingRegion string
	Destination       string
}

type ConfigureCompiler func(*CompilationEnvironment)

func WithCompilerFunctions(f func(*CompilationEnvironment) optimizer.CompilerFunctions) ConfigureCompiler {
	return func(env *CompilationEnvironment) {
		funcs := f(env)
		env.Optimization.CompilerFuncs = funcs
		env.Optimization.AddOptimizer(func(innerEnv *CompilationEnvironment) ast.Rewriter {
			return optimizer.RunCompilerFunctions(innerEnv.Symbols, innerEnv.Optimization.CompilerFuncs)
		})
	}
}

func WithTokens(names []string) ConfigureCompiler {
	return func(env *CompilationEnvironment) {
		env.Symbols.DeclareMany(symbols.TOKEN, names)
	}
}

func Defaults() ConfigureCompiler {
	return func(env *CompilationEnvironment) {
		env.Optimization.Passes = 10

		env.Symbols.DeclareMany(symbols.COMP_TIME, optimizer.CompileTimeNames())
		env.Symbols.DeclareMany(symbols.BUILT_IN, objects.BuiltInFunctionNames())
		env.Symbols.DeclareMany(symbols.GLOBAL, vm.GlobalNames())
		env.Symbols.DeclareMany(symbols.SETTING, settings.Names())

		env.OnSourceLoad(func(env *CompilationEnvironment, src *Source) {
			SetCurrentLocation(env.Optimization.Context, src.OriginatingRegion)
		})
		env.Optimization.AddOptimizer(func(env *CompilationEnvironment) ast.Rewriter {
			return optimizer.InlineCalls(env.Optimization.Context, env.Symbols, env.Functions)
		})
		env.Optimization.AddOptimizer(func(env *CompilationEnvironment) ast.Rewriter {
			return optimizer.FoldConstants(env.Symbols)
		})
		env.Optimization.AddOptimizer(func(env *CompilationEnvironment) ast.Rewriter {
			return optimizer.EnsureFuncs(env.Symbols, env.Functions)
		})
		env.Optimization.AddOptimizer(func(env *CompilationEnvironment) ast.Rewriter {
			return optimizer.CollapseHas(env.Symbols)
		})
		env.Optimization.AddOptimizer(func(env *CompilationEnvironment) ast.Rewriter {
			return optimizer.PromoteTokens(env.Symbols)
		})
	}
}

func NewCompileEnvironment(configure ...ConfigureCompiler) CompilationEnvironment {
	var env CompilationEnvironment
	env.Grammar = ruleparser.NewRulesGrammar()
	env.Symbols = ptr(symbols.NewTable())
	env.Objects = ptr(objects.NewTableBuilder())
	env.Optimization.Context = ptr(optimizer.NewCtx())

	for i := range configure {
		configure[i](&env)
	}

	return env
}

type SourceLoaded func(*CompilationEnvironment, *Source)
type Optimizer func(*CompilationEnvironment) ast.Rewriter
type Analyzer func(*CompilationEnvironment) ast.Visitor

type CompilationEnvironment struct {
	Grammar           peruse.Grammar[ruleparser.Tree]
	Symbols           *symbols.Table
	Functions         *ast.PartialFunctionTable
	Objects           *objects.TableBuilder
	CompilerFunctions optimizer.CompilerFunctions

	Optimization Optimization
	Analysis     Analysis
	onSourceLoad []SourceLoaded
}

func (this *CompilationEnvironment) OnSourceLoad(f SourceLoaded) {
	this.onSourceLoad = append(this.onSourceLoad, f)
}

type Analysis struct {
	pre  []Analyzer
	post []Analyzer
}

func (this *Analysis) PreOptimize(v Analyzer) {
	this.pre = append(this.pre, v)
}

func (this *Analysis) PostOptimize(v Analyzer) {
	this.post = append(this.post, v)
}

type Optimization struct {
	CompilerFuncs optimizer.CompilerFunctions
	Context       *optimizer.Context
	Optimiziers   []Optimizer
	Passes        int
}

func (this *Optimization) AddOptimizer(o Optimizer) {
	this.Optimiziers = append(this.Optimiziers, o)
}

func (this *CompilationEnvironment) BuildFunctionTable(declarations map[string]string) error {
	var err error
	var funcs ast.PartialFunctionTable
	funcs, err = ast.BuildCompilingFunctionTable(this.Symbols, this.Grammar, declarations)
	if err != nil {
		return err
	}
	this.Functions = &funcs
	return nil
}

func Compiler(env *CompilationEnvironment) compiling {
	optimizers := env.Optimization.Optimiziers
	analysis := env.Analysis
	compiling := compiling{
		CompilationEnvironment: env,
		rewriters:              make([]ast.Rewriter, len(optimizers)),
		preanalyzers:           make([]ast.Visitor, len(analysis.pre)),
		postanalyzers:          make([]ast.Visitor, len(analysis.post)),
	}
	for i := range compiling.rewriters {
		compiling.rewriters[i] = optimizers[i](env)
	}
	for i := range compiling.preanalyzers {
		compiling.preanalyzers[i] = analysis.pre[i](env)
	}
	for i := range compiling.postanalyzers {
		compiling.postanalyzers[i] = analysis.post[i](env)
	}
	return compiling
}

type compiling struct {
	*CompilationEnvironment
	rewriters     []ast.Rewriter
	preanalyzers  []ast.Visitor
	postanalyzers []ast.Visitor
}

func (this compiling) CompileSource(src *Source) (compiler.ByteCode, error) {
	var bytecode compiler.ByteCode
	for i := range this.onSourceLoad {
		this.onSourceLoad[i](this.CompilationEnvironment, src)
	}

	if src.Ast == nil {
		var astErr error
		src.Ast, astErr = ast.Parse(string(src.String), this.Symbols, this.Grammar)
		if astErr != nil {
			return bytecode, fmt.Errorf("%w: %w", ErrParse, astErr)
		}
	}

	for i := range this.Analysis.pre {
		this.preanalyzers[i].Visit(src.Ast)
	}

	for range this.Optimization.Passes {
		var rewriteErr error
		src.Ast, rewriteErr = ast.RewriteWithEvery(src.Ast, this.rewriters)
		if rewriteErr != nil {
			return bytecode, fmt.Errorf("%w: %w", ErrOptimization, rewriteErr)
		}
	}

	for i := range this.Analysis.post {
		this.postanalyzers[i].Visit(src.Ast)
	}

	var compileErr error
	bytecode, compileErr = compiler.Compile(src.Ast, this.Symbols, this.Objects)
	if compileErr != nil {
		compileErr = fmt.Errorf("%w: %w", ErrCompile, compileErr)
	}
	return bytecode, compileErr
}

var (
	ErrSourceLoad   = errors.New("source load")
	ErrParse        = errors.New("parsing")
	ErrOptimization = errors.New("optimization")
	ErrCompile      = errors.New("compile")
)
