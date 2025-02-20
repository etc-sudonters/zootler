package generation

import (
	"sudonters/libzootr/cmd/knowitall/bubbles/explore"
	"sudonters/libzootr/magicbean"
	"sudonters/libzootr/magicbean/tracking"
	"sudonters/libzootr/playthrough"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Searches struct {
	Adult *playthrough.Search
	Child *playthrough.Search
}

func New(gen *magicbean.Generation, names tracking.NameTable, searches Searches) Model {
	return Model{
		gen:      gen,
		search:   searches,
		names:    names,
		discache: make(discache, 32),
		explore:  explore.New(),
	}
}

type Model struct {
	gen        *magicbean.Generation
	names      tracking.NameTable
	search     Searches
	discache   discache
	explore    explore.Model
	statusLine string
}

func (this Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case explore.ExploreSphere:
		this.statusLine = "running search"
		return this, runSearch(msg, this.search, this.names)
	case explore.DisassembleRule:
		this.statusLine = "disassembling"
		return this, disassemble(this.gen, msg.Id, this.discache)
	case explore.RuleDisassembled:
		if msg.Name == "" && msg.Id != 0 {
			msg.Name = string(this.names[msg.Id])
		}
		break
	case explore.SphereExplored:
		this.statusLine = "search concluded"
		break
	}

	var xplrCmd tea.Cmd
	this.explore, xplrCmd = this.explore.Update(msg)
	cmds = append(cmds, xplrCmd)
	return this, tea.Batch(cmds...)
}

func (this Model) Init() tea.Cmd {
	return this.explore.Init()
}

func (this Model) View() string {
	xplr := this.explore.View()
	return lipgloss.JoinVertical(lipgloss.Left, xplr, this.statusLine)
}
