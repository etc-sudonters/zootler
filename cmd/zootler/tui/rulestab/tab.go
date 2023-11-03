package rulestab

import (
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/tui/application"
	"sudonters/zootler/internal/tui/listpanel"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewParseTab(pool entity.Pool) parsePanel {
	var p parsePanel
	p.pool = pool
	p.list = createRuleList(pool, tea.WindowSizeMsg{Width: 0, Height: 0}) // its fine don't worry
	p.view = viewport.New(0, 0)
	p.ruleViews = make(map[entity.Model]ruleView)

	return p
}

type parsePanel struct {
	pool      entity.Pool
	list      listpanel.Model
	view      viewport.Model
	current   entity.Model
	ruleViews map[entity.Model]ruleView
}

func (p parsePanel) Init() tea.Cmd {
	return application.WantWindowSize
}

func (p parsePanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd // always safe to overwrite

	sync := func() {
		p.view.SetContent(p.ruleViews[p.current].View())
		cmds = append(cmds, viewport.Sync(p.view))
	}

	updateList := func(msg tea.Msg) {
		p.list, cmd = p.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	updateView := func(msg tea.Msg) {
		p.view, cmd = p.view.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		updateList(application.ResizeWindowMsg(msg, .25, 1))
		updateView(application.ResizeWindowMsg(msg, .75, 0))
	case changeRuleView:
		m := entity.Model(msg)
		if m != entity.INVALID_ENTITY && p.current != m {
			p.current = entity.Model(msg)
			sync()
		}
	case addRuleView:
		if msg.id != entity.INVALID_ENTITY {
			p.ruleViews[msg.id] = msg.v
			if msg.focus {
				p.current = msg.id
				sync()
			}
		}
	}

	cmd = tea.Batch(cmds...)
	return p, cmd
}

func (p parsePanel) View() string {
	var rendered []string

	rendered = append(rendered, p.list.View())

	if p.current != entity.INVALID_ENTITY {
		rendered = append(rendered, p.view.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}
