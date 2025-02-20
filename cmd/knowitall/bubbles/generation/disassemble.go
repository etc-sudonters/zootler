package generation

import (
	"sudonters/libzootr/cmd/knowitall/bubbles/explore"
	"sudonters/libzootr/components"
	"sudonters/libzootr/magicbean"
	"sudonters/libzootr/mido/compiler"
	"sudonters/libzootr/mido/vm"
	"sudonters/libzootr/zecs"

	tea "github.com/charmbracelet/bubbletea"
)

type discache map[zecs.Entity]explore.RuleDisassembled

func disassemble(gen *magicbean.Generation, edge zecs.Entity, cache discache) tea.Cmd {
	return func() tea.Msg {
		dis, exists := cache[edge]
		if exists {
			return dis
		}
		ocm := &gen.Ocm
		values, err := ocm.GetValues(edge,
			zecs.Get[components.Name], zecs.Get[components.RuleCompiled],
		)
		dis.Id = edge
		dis.Err = err
		if err == nil {
			name := values[0].(components.Name)
			code := values[1].(components.RuleCompiled)
			disassembled := vm.Disassemble(compiler.Bytecode(code), &gen.Objects)
			dis.Name = string(name)
			dis.Disassembly = disassembled
			cache[edge] = dis
		}
		return dis
	}
}
