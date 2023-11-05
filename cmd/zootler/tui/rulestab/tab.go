package rulestab

import (
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/tui/application"
	"sudonters/zootler/internal/tui/listpanel"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewParseTab(pool entity.Pool) parsePanel {
	var p parsePanel
	p.list = createRuleList(pool, tea.WindowSizeMsg{Width: 0, Height: 0}) // its fine don't worry
	p.parsed = make(map[entity.Model]ruleView)
	return p
}

type parseFocus int

const (
	focusList parseFocus = iota
	focusView
)

type parsePanel struct {
	list   listpanel.Model
	cur    entity.Model
	parsed map[entity.Model]ruleView
	view   string
}

func (p parsePanel) Init() tea.Cmd {
	return application.WantWindowSize
}

func (p parsePanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd // always safe to overwrite

	updateList := func(msg tea.Msg) {
		p.list, cmd = p.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	changeViewContent := func() {
		p.view = p.parsed[p.cur].View()
	}

	switch msg := msg.(type) {
	case changeRuleView:
		p.cur = entity.Model(msg)
		changeViewContent()
		cmd = tea.Batch(cmds...)
		return p, cmd
	case addRuleView:
		p.parsed[msg.id] = msg.v
		if msg.focus {
			p.cur = msg.id
			changeViewContent()
		}
		cmd = tea.Batch(cmds...)
		return p, cmd
	}

	updateList(msg)
	cmd = tea.Batch(cmds...)
	return p, cmd
}

var viewStyle = lipgloss.NewStyle().Border(lipgloss.DoubleBorder())
var listStyle = lipgloss.NewStyle().Border(lipgloss.BlockBorder())

func (p parsePanel) View() string {
	var rendered []string
	rendered = append(rendered, listStyle.Render(p.list.View()))
	rendered = append(rendered, viewStyle.Render(p.view))
	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}

func (p parsePanel) current() *ruleView {
	if p.cur == entity.INVALID_ENTITY {
		return nil
	}

	e, ok := p.parsed[p.cur]
	if !ok {
		return nil
	}
	return &e
}
