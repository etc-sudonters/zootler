package panels

import (
	"sudonters/zootler/cmd/zootler/visualizer/listpanel"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/etc-sudonters/substrate/bag"
)

type createListPanel func() listpanel.Model

func CreateListPanel(c createListPanel) tea.Cmd {
	return func() tea.Msg {
		return c
	}
}

type PanelKeys struct {
	NextSibling, LastSibling, CloseCurrent key.Binding
}

type PanelStyles struct {
	Focused, Blurred lipgloss.Style
}

func defaultKeys() PanelKeys {
	return PanelKeys{
		NextSibling:  key.NewBinding(key.WithKeys("tab")),
		LastSibling:  key.NewBinding(key.WithKeys("shift+tab")),
		CloseCurrent: key.NewBinding(key.WithKeys("X")),
	}
}

func New(root listpanel.Model) panels {
	var p panels
	p.p = []listpanel.Model{root}
	p.keys = defaultKeys()
	p.maxDisplay = 3
	return p
}

type panels struct {
	p []listpanel.Model

	keys PanelKeys

	idx        int
	maxDisplay int

	Styles PanelStyles
}

func (p panels) Init() tea.Cmd {
	return nil
}

func (p panels) current() *listpanel.Model {
	return &p.p[p.idx]
}

func (p panels) blurCurrent() {
	p.current().Blur()
}

func (p panels) focusCurrent() {
	p.current().Focus()
}

func (p *panels) addPanelToEnd(l listpanel.Model) {
	p.p = append(p.p, l)
	p.idx = len(p.p) - 1
}

func (p panels) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case createListPanel:
		(&p).addPanelToEnd(msg())
		p.focusCurrent()
		return p, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, p.keys.NextSibling):
			if p.idx+1 == len(p.p) {
				return p, nil
			}

			p.blurCurrent()
			p.idx++
			p.focusCurrent()
			return p, nil
		case key.Matches(msg, p.keys.LastSibling):
			if p.idx == 0 {
				return p, nil
			}

			p.blurCurrent()
			p.idx--
			p.focusCurrent()
			return p, nil
		case key.Matches(msg, p.keys.CloseCurrent):
			if len(p.p) == 1 {
				return p, tea.Quit
			}

			// don't close the root if we're not the only one left
			if p.idx == 0 {
				return p, nil
			}

			(&p).removeAtCurrent()
			return p, nil
		}
	}

	// current() would mean a deref dance
	panel := p.p[p.idx]
	panel, cmd = panel.Update(msg)
	p.p[p.idx] = panel

	return p, cmd
}

func (p *panels) removeAtCurrent() {
	p.p = append(p.p[:p.idx], p.p[p.idx+1:]...)
	if p.idx == len(p.p) {
		p.idx--
	}
	p.focusCurrent()
}

func (p panels) View() string {
	if p.maxDisplay >= len(p.p) {
		return p.renderAll()
	}

	return p.renderWindow()
}

func (p panels) renderAll() string {
	views := make([]string, len(p.p))

	for i, v := range p.p {
		fn := p.Styles.Blurred.Render
		if i == p.idx {
			fn = p.Styles.Focused.Render
		}

		views[i] = fn(v.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, views...)
}

func (p panels) renderWindow() string {
	var views []string
	spacing := (p.maxDisplay - 1) / 2

	start := bag.Max(p.idx-spacing, 0)
	midway := start + spacing
	end := bag.Min(midway+spacing, len(p.p))
	// maybe try to add one at the end for an even display
	if p.maxDisplay&1 != 0 {
		end = bag.Min(end+1, len(p.p))
	}

	for i, v := range p.p[start:end] {
		fn := p.Styles.Blurred.Render
		if start+spacing == i {
			fn = p.Styles.Focused.Render
		}
		views = append(views, fn(v.View()))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, views...)

}
