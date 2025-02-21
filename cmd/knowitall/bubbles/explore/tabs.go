package explore

import (
	"sudonters/libzootr/magicbean"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func newTabs(sphere *NamedSphere) tabs {
	return tabs{
		tabs: []tab{
			{"SUMMARY", summary{sphere}},
			{"INVENTORY", newCollected(sphere)},
			{"ADULT", newEdges(sphere, magicbean.AgeAdult)},
			{"CHILD", newEdges(sphere, magicbean.AgeChild)},
			{"DISASSEMBLY", disassembly{}},
			{"SEARCH", search{}},
		},
	}
}

type tab struct {
	display string
	mount   tea.Model
}

type tabs struct {
	tabs []tab
	curr int

	tabWidth    int
	displaySize tea.WindowSizeMsg
}

func (this tabs) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(this.tabs))
	for i := range this.tabs {
		cmds[i] = this.tabs[i].mount.Init()
	}
	return tea.Batch(cmds...)
}

func (this tabs) Update(msg tea.Msg) (tabs, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmd := this.resize(msg)
		return this, cmd
	case RuleDisassembled:
		return this, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "l":
			this.next()
			return this, nil
		case "right", "h":
			this.prev()
			return this, nil
		case "I":
			this.focus(TAB_INVENTORY)
			return this, nil
		case "A":
			this.focus(TAB_ADULT)
			return this, nil
		case "C":
			this.focus(TAB_CHILD)
			return this, nil
		case "S":
			this.focus(TAB_SEARCH)
			return this, nil
		case "D":
			this.focus(TAB_DISASSEMBLY)
			return this, nil
		}
	}
	cmd := this.updateActiveTab(msg)
	return this, cmd
}

func (this tabs) View() string {
	tabs := this.renderTabs()
	body := this.renderActiveTab()
	return lipgloss.JoinVertical(lipgloss.Left, tabs, body)
}

func (this *tabs) resize(msg tea.WindowSizeMsg) tea.Cmd {
	this.tabWidth = sizeTab(msg.Width, len(tabTpl))
	tabs := this.renderTabs()
	this.displaySize = tea.WindowSizeMsg{
		Width:  lipgloss.Width(tabs) - windowStyle.GetHorizontalFrameSize(),
		Height: msg.Height - 2,
	}
	cmds := make([]tea.Cmd, len(this.tabs))
	for i := range this.tabs {
		this.tabs[i].mount, cmds[i] = this.tabs[i].mount.Update(tea.WindowSizeMsg{
			Width:  this.displaySize.Width - 10,
			Height: this.displaySize.Height - 10,
		})
	}

	return tea.Batch(cmds...)
}

func sizeTab(width int, tabCount int) int {
	return (width - tabCount) / tabCount
}

func (this *tabs) updateActiveTab(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	tab := &this.tabs[this.curr]
	tab.mount, cmd = tab.mount.Update(msg)
	return cmd
}

func (this *tabs) next() {
	this.curr = min(this.curr+1, TAB_SEARCH)
}

func (this *tabs) prev() {
	this.curr = max(this.curr-1, TAB_SUMMARY)
}

func (this *tabs) focus(tab int) {
	this.curr = tab
}

func (this tabs) renderTabs() string {
	tabs := make([]string, len(this.tabs))
	for i, tab := range this.tabs {
		first, last, active := i == 0, i == len(this.tabs)-1, i == this.curr

		style := inactiveStyle.Border(inactiveTabBorder)
		if active {
			style = currentStyle.Border(activeTabBorder)
		}

		style = tabStyle.Inherit(style).Width(this.tabWidth).Height(1).AlignHorizontal(lipgloss.Center)

		border, _, _, _, _ := style.GetBorder()
		if first && active {
			border.BottomLeft = "│"
		} else if first && !active {
			border.BottomLeft = "├"
		} else if last && active {
			border.BottomRight = "│"
		} else if last && !active {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		tabs[i] = style.Render(tab.display)
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	return row

}

func (this tabs) renderActiveTab() string {
	content := this.tabs[this.curr].mount.View()
	return windowStyle.Width(this.displaySize.Width).Height(this.displaySize.Height).Render(content)
}

const (
	TAB_SUMMARY     = 0
	TAB_INVENTORY   = 1
	TAB_ADULT       = 2
	TAB_CHILD       = 3
	TAB_DISASSEMBLY = 4
	TAB_SEARCH      = 5
)

var tabTpl []tea.Model = []tea.Model{
	(*summary)(nil),
	(*collected)(nil),
	(*edges)(nil), // adult
	(*edges)(nil), // child
	(*disassembly)(nil),
	(*search)(nil),
}

type summary struct {
	sphere *NamedSphere
}

func (_ summary) Init() tea.Cmd { return nil }

func (this summary) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return this, nil
}

func (this summary) View() string {
	if this.sphere == nil {
		return "NO SPHERE LOADED"
	}

	return "SPHERE LOADED"
}

type disassembly struct {
	dis RuleDisassembled
}

func (_ disassembly) Init() tea.Cmd { return nil }

func (this disassembly) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_ = msg
	return this, nil
}

func (this disassembly) View() string {
	return "DISASSEMBLY"
}

type search struct {
}

func (_ search) Init() tea.Cmd {
	return nil
}

func (this search) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			return this, RequestNextSphere
		}
	}
	return this, nil
}

func (this search) View() string {
	return "[ RUN SPHERE SEARCH ]"
}
