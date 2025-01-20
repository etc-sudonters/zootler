package main

import (
	"fmt"
	"sudonters/zootler/mido"
	"sudonters/zootler/mido/ast"
	"sudonters/zootler/mido/code"
	"sudonters/zootler/mido/objects"
	"sudonters/zootler/mido/vm"
)

func ExectuteAll(env *mido.CompileEnv, compiled []mido.CompiledSource) {
	objs := objects.TableFrom(env.Objects)
	engine := vm.VM{&objs, nil}

	for i := range compiled {
		compiled := compiled[i]
		result, err := engine.Execute(compiled.ByteCode)

		if err == nil && result != objects.Null {
			continue
		}
		fmt.Println()
		fmt.Printf("%s -> %s\n", compiled.OriginatingRegion, compiled.Destination)
		if compiled.String != "" {
			fmt.Println()
			fmt.Println("raw")
			fmt.Println(compiled.String)
		}
		if compiled.Ast != nil {
			fmt.Println()
			fmt.Println("ast")
			fmt.Println(ast.Render(compiled.Ast))
		}

		fmt.Println()
		fmt.Println(code.DisassembleToString(compiled.ByteCode.Tape))

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("produced nil result")
		}
		fmt.Println()
	}
}
