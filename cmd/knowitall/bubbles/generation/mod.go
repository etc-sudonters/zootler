package generation

import (
	"sudonters/libzootr/cmd/knowitall/bubbles/spheres"
	"sudonters/libzootr/magicbean"
	"sudonters/libzootr/magicbean/tracking"

	tea "github.com/charmbracelet/bubbletea"
)

type searches = map[magicbean.Age]*magicbean.Search

func New(gen *magicbean.Generation, names tracking.NameTable, searches searches) Model {
	return Model{
		gen:      gen,
		spheres:  spheres.New(),
		search:   searches,
		names:    names,
		discache: make(discache, 32),
	}
}

type Model struct {
	gen      *magicbean.Generation
	names    tracking.NameTable
	spheres  spheres.Model
	search   map[magicbean.Age]*magicbean.Search
	discache discache
}

func (this Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spheres.WantSearch:
		return this, runSearch(this.gen, this.search, this.names)
	case SearchResult:
		syncCmd := this.spheres.PushSphere(spheres.Details(msg))
		return this, syncCmd
	case spheres.DisassemblyRequested:
		return this, disassemble(this.gen, msg.Id, this.discache)
	}

	var sphereCommand tea.Cmd
	this.spheres, sphereCommand = this.spheres.Update(msg)
	cmds = append(cmds, sphereCommand)
	return this, tea.Batch(cmds...)
}

func (this Model) Init() tea.Cmd {
	return this.spheres.Init()
}

func (this Model) View() string {
	return this.spheres.View()
}
