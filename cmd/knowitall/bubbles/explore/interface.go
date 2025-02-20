package explore

import (
	"sudonters/libzootr/components"
	"sudonters/libzootr/mido/vm"
	"sudonters/libzootr/playthrough"
	"sudonters/libzootr/zecs"

	tea "github.com/charmbracelet/bubbletea"
)

type sphereSelected int

type NamedEdge struct {
	Id   zecs.Entity
	Name components.Name
}

type NamedNode struct {
	Id   zecs.Entity
	Name components.Name
}

type NamedSphere struct {
	I     int
	Error error
	Adult playthrough.SearchSphere
	Child playthrough.SearchSphere
	Edges []NamedEdge
	Nodes []NamedNode
}

type SphereExplored struct {
	Err    error
	Sphere NamedSphere
}

type ExploreSphere struct{}

func RequestNextSphere() tea.Msg {
	return ExploreSphere{}
}

type DisassembleRule struct {
	Id zecs.Entity
}

type RuleDisassembled struct {
	Id          zecs.Entity
	Err         error
	Name        string
	Disassembly vm.Disassembly
}

func RequestDisassembly(edge zecs.Entity) tea.Cmd {
	return func() tea.Msg {
		return DisassembleRule{edge}
	}
}
