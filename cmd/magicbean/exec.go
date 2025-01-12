package main

import (
	"fmt"
	"sudonters/zootler/midologic"
	"sudonters/zootler/midologic/ast"
	"sudonters/zootler/midologic/code"
	"sudonters/zootler/midologic/objects"
	"sudonters/zootler/midologic/vm"
)

func ExectuteAll(env *midologic.CompileEnv, compiled []midologic.CompiledSource) {
	objTable := objects.NewTable(
		objects.TableFrom(*env.Objects),
	)
	engine := vm.New(&objTable)

	for i := range compiled {
		compiled := compiled[i]
		result, err := engine.Execute(compiled.ExecutionUnit())

		if err == nil && result != nil {
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

type constbuiltins bool

func (this constbuiltins) AtDampeTime([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) AtDay([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) AtNight([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) Has([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) HasAllNotesForSong([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) HasAnyOf([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) HasBottle([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) HasDungeonRewards([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) HasEvery([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) HasHearts([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) HasMedallions([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) HasStones([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) IsAdult([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) IsChild([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) IsStartingAge([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
func (this constbuiltins) RegionHasShortcuts([]objects.Object) (objects.Object, error) {
	return objects.Boolean(this), nil
}
