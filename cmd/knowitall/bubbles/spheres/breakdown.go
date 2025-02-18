package spheres

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type breakdown struct {
	sphere        *Details
	summary       Summary
	activeTab     int
	height, width int

	tabWidth, viewWidth, viewHeight int

	sum    summaryTab
	nodes  edgesTab
	tokens tokensTab
	dis    disassemblyTab
	search runSearchTab
}

func (this breakdown) Init() tea.Cmd {
	cmds := []tea.Cmd{}
	addCmd := func(cmd tea.Cmd) {
		cmds = append(cmds, cmd)
	}
	addCmd((&this.sum).init(&this))
	addCmd((&this.nodes).init(&this))
	addCmd((&this.tokens).init(&this))
	addCmd((&this.dis).init(&this))
	addCmd((&this.search).init(&this))

	return tea.Batch(cmds...)
}

func (this *breakdown) setSize(size tea.WindowSizeMsg) {
	this.height = size.Height
	this.width = size.Width
	this.tabWidth = sizeTab(this.width, len(breakDowntabs))
	tabs := renderBreakDownTabs(this.tabWidth, this.activeTab)
	this.viewWidth = lipgloss.Width(tabs) - windowStyle.GetHorizontalFrameSize()
	this.viewHeight = this.height - 4
}

func (this *breakdown) updateActiveTab(msg tea.Msg) []tea.Cmd {
	var cmds []tea.Cmd
	addCmd := func(cmd tea.Cmd) {
		cmds = append(cmds, cmd)
	}
	switch this.activeTab {
	case TAB_SUMMARY:
		addCmd(this.sum.update(msg, this))
	case TAB_EDGES:
		addCmd((&this.nodes).update(msg, this))
	case TAB_TOKENS:
		addCmd(this.tokens.update(msg, this))
	case TAB_DISASSEMBLY:
		addCmd(this.dis.update(msg, this))
	case TAB_SEARCH:
		addCmd(this.search.update(msg, this))
	}
	return cmds
}

func (this *breakdown) updateAllTabs(msg tea.Msg) []tea.Cmd {
	var cmds []tea.Cmd
	addCmd := func(cmd tea.Cmd) {
		cmds = append(cmds, cmd)
	}
	addCmd(this.sum.update(msg, this))
	addCmd((&this.nodes).update(msg, this))
	addCmd(this.tokens.update(msg, this))
	addCmd(this.dis.update(msg, this))
	addCmd(this.search.update(msg, this))
	return cmds
}

func (this breakdown) Update(msg tea.Msg) (breakdown, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		this.setSize(msg)
		produced := this.updateAllTabs(tea.WindowSizeMsg{
			Height: this.viewHeight,
			Width:  this.viewWidth,
		})
		cmds = append(cmds, produced...)
		goto exit
	case Disassembly:
		this.dis.disassembled = msg
		this.activeTab = TAB_DISASSEMBLY
		goto exit
	case tea.KeyMsg:
		switch str := msg.String(); str {
		case "h":
			(&this).tabLeft()
			goto exit
		case "l":
			(&this).tabRight()
			goto exit
		}
	}

	cmds = append(cmds, this.updateActiveTab(msg)...)

exit:
	return this, tea.Batch(cmds...)
}

func (this *breakdown) tabLeft() {
	this.activeTab = max(0, this.activeTab-1)
}

func (this *breakdown) tabRight() {
	this.activeTab = min(this.activeTab+1, len(breakDowntabs)-1)
}

func (this breakdown) View() string {
	var view strings.Builder
	fmt.Fprintln(&view, renderBreakDownTabs(this.tabWidth, this.activeTab))
	fmt.Fprint(&view, this.renderActiveTab(windowStyle.Width(this.viewWidth).Height(this.viewHeight)))
	return docStyle.Render(view.String())
}

func (this breakdown) renderActiveTab(style lipgloss.Style) string {
	var content strings.Builder
	switch this.activeTab {
	case TAB_SUMMARY:
		this.sum.view(&content, &this)
	case TAB_EDGES:
		(&this.nodes).view(&content, &this)
	case TAB_TOKENS:
		this.tokens.view(&content, &this)
	case TAB_DISASSEMBLY:
		this.dis.view(&content)
	case TAB_SEARCH:
		this.search.view(&content)
	}
	return style.Render(contentStyle.Render(content.String()))
}
