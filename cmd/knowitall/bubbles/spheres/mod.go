package spheres

import (
	"sudonters/libzootr/components"
	"sudonters/libzootr/mido/vm"
	"sudonters/libzootr/zecs"
	"sync/atomic"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var lastId uint64 = 1

func nextId() uint64 {
	return atomic.AddUint64(&lastId, 1)
}

type focus int

const (
	focusSidebar   = 1
	focusBreakdown = 2
)

type Model struct {
	id        uint64
	current   int
	spheres   []Details
	sidebar   list.Model
	breakdown breakdown
	focus     focus

	size tea.WindowSizeMsg
}

type Builder = func(*Model)

func listDefaults(l *list.Model) {
	l.SetShowFilter(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.Title = ""
	l.DisableQuitKeybindings()
	l.SetShowTitle(false)
}

func New(opts ...Builder) Model {
	var m Model
	width, height := 30, 20
	sidebar := list.New(nil, summaryDelegate{}, width, height)
	listDefaults(&sidebar)
	sidebar.SetShowTitle(true)
	sidebar.Title = "Spheres"
	sidebar.Styles.Title = currentStyle
	sidebar.Styles.NoItems = lipgloss.NewStyle().SetString("No Spheres Loaded")
	m.id = nextId()
	m.sidebar = sidebar
	for i := range opts {
		opts[i](&m)
	}
	m.breakdown = m.createBreakdown()
	m.focus = focusBreakdown
	return m
}

func (this *Model) PushSphere(sphere Details) tea.Cmd {
	sphere.I = len(this.spheres)
	this.spheres = append(this.spheres, sphere)
	this.current = sphere.I
	cmd := this.syncItems()
	this.createBreakdown()
	this.focus = focusBreakdown
	return cmd
}

func (this Model) Init() tea.Cmd {
	return tea.Batch(RequestSearch, this.breakdown.Init())
}

func (this *Model) syncItems() tea.Cmd {
	summaries := make([]list.Item, len(this.spheres))
	for i := range summaries {
		summaries[i] = summarize(this.spheres[i])
	}
	cmd := this.sidebar.SetItems(summaries)
	this.sidebar.Select(this.current)
	return cmd
}

func (this *Model) setSize(size tea.WindowSizeMsg) {
	this.size = size

	this.sidebar.SetHeight(this.size.Height - 10)
	this.sidebar.SetWidth(30)

	(&this.breakdown).setSize(this.sizeBreakdown())
}

func (this Model) sizeBreakdown() tea.WindowSizeMsg {
	return tea.WindowSizeMsg{
		Height: this.size.Height - 10,
		Width:  this.size.Width - 50,
	}
}

func (this *Model) createBreakdown() breakdown {
	var breakdown breakdown
	(&breakdown).setSize(this.sizeBreakdown())
	if len(this.spheres) > 0 {
		breakdown.sphere = &this.spheres[this.current]
		breakdown.summary = this.sidebar.SelectedItem().(Summary)
	}

	breakdown.nodes = newEdgesTab(&breakdown)
	breakdown.sum = newSummary(&breakdown)
	breakdown.tokens = newTokenTab(&breakdown)
	breakdown.dis = newDisTab(&breakdown)
	breakdown.search = newSearchTab(&breakdown)

	return breakdown
}

func (this Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var focusedCmd tea.Cmd
	addCmd := func(c tea.Cmd) {
		cmds = append(cmds, c)
	}
	if this.focus == 0 {
		this.focus = focusSidebar
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		(&this).setSize(msg)
		goto exit
	case Disassembly:
		var cmd tea.Cmd
		this.breakdown, cmd = this.breakdown.Update(msg)
		addCmd(cmd)
		goto sync
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab, tea.KeyShiftTab:
			if this.focus == focusSidebar {
				this.focus = focusBreakdown
			} else {
				this.focus = focusSidebar
			}
			goto sync
		case tea.KeyEnter:
			if this.focus == focusSidebar {
				this.current = this.sidebar.Index()
				this.focus = focusBreakdown
				this.breakdown = (&this).createBreakdown()
				goto sync
			}
		}
	}

	switch this.focus {
	case focusSidebar:
		this.sidebar, focusedCmd = this.sidebar.Update(msg)
	case focusBreakdown:
		this.breakdown, focusedCmd = this.breakdown.Update(msg)
	}
	cmds = append(cmds, focusedCmd)

sync:
	if len(this.spheres) != len(this.sidebar.Items()) {
		cmd := (&this).syncItems()
		cmds = append(cmds, cmd)
	}

exit:
	return this, tea.Batch(cmds...)
}

func (this Model) View() string {
	sidebar := this.sidebar.View()
	breakdown := this.breakdown.View()
	render := spheresBlock.Render
	return render(lipgloss.JoinHorizontal(
		lipgloss.Top,
		render(sidebar),
		render(breakdown),
	))
}

var spheresBlock = lipgloss.NewStyle().MarginLeft(3)

type WantSearch struct{}

func RequestSearch() tea.Msg {
	return WantSearch{}
}

type DisassemblyRequested struct {
	Id zecs.Entity
}

type Disassembly struct {
	Values any
	Id     zecs.Entity
	Name   components.Name
	Code   components.RuleCompiled
	Src    components.RuleSource
	Ast    components.RuleParsed
	Opt    components.RuleOptimized
	Dis    vm.Disassembly
	Err    error
}

func RequestDisassembly(id zecs.Entity) tea.Cmd {
	return func() tea.Msg {
		return DisassemblyRequested{id}
	}
}
