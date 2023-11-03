package panels

import (
	"sudonters/zootler/internal/tui/listpanel"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/etc-sudonters/substrate/bag"
)

type Opt func(*Model)

type UpdateDelegate func(*Model, *listpanel.Model, int, tea.Msg) (bool, tea.Cmd)
type CloseDelegate func(*Model, *listpanel.Model, int) (bool, tea.Cmd)

type Size struct {
	Height, Width int
}
type createListPanel func(Size) listpanel.Model

func CreateListPanel(c createListPanel) tea.Cmd {
	return func() tea.Msg {
		return c
	}
}

type Keys struct {
	NextSibling, LastSibling, CloseCurrent key.Binding
}

type Styles struct {
	Focused, Blurred lipgloss.Style
}

func WithUpdate(u UpdateDelegate) Opt {
	return func(p *Model) {
		p.updating = u
	}
}

func WithPanels(p listpanel.Model, ps ...listpanel.Model) Opt {
	return func(panels *Model) {
		panels.p = append(panels.p, p)
		panels.p = append(panels.p, ps...)
	}
}

func WithKeys(k Keys) Opt {
	return func(p *Model) {
		p.keys = k
	}
}

func WithMaxDisplay(n int) Opt {
	return func(p *Model) {
		p.maxDisplay = n
	}
}

func WithClose(d CloseDelegate) Opt {
	return func(p *Model) {
		p.closing = d
	}
}

func defaultKeys() Keys {
	return Keys{
		NextSibling:  key.NewBinding(key.WithKeys("tab")),
		LastSibling:  key.NewBinding(key.WithKeys("shift+tab")),
		CloseCurrent: key.NewBinding(key.WithKeys("X")),
	}
}

func New(opts ...Opt) Model {
	var p Model
	p.keys = defaultKeys()
	p.maxDisplay = 3

	for _, o := range opts {
		o(&p)
	}

	return p
}

type Model struct {
	p []listpanel.Model

	keys Keys

	idx        int
	maxDisplay int

	height, width int

	Styles Styles

	updating UpdateDelegate
	closing  CloseDelegate
}

func (p Model) Init() tea.Cmd {
	return nil
}

func (p Model) current() *listpanel.Model {
	return &p.p[p.idx]
}

func (p Model) BlurCurrent() {
	p.current().Blur()
}

func (p Model) FocusCurrent() {
	p.current().Focus()
}

func (p Model) Index() int {
	return p.idx
}

func (p Model) Len() int {
	return len(p.p)
}

func (p *Model) Append(l listpanel.Model) {
	p.p = append(p.p, l)
	p.idx = len(p.p) - 1
}

func (p *Model) PopCurrent() listpanel.Model {
	panel := p.p[p.idx]
	p.p = append(p.p[:p.idx], p.p[p.idx+1:]...)
	p.idx = bag.Max(bag.Min(p.idx, len(p.p)-1), 0)
	p.FocusCurrent()
	return panel
}

func (p Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.height = msg.Height
		p.width = msg.Width

		cmds := make([]tea.Cmd, len(p.p))

		for i := range p.p {
			p.p[i], cmds[i] = p.p[i].Update(tea.WindowSizeMsg{
				Height: p.height,
				Width:  p.width / p.maxDisplay,
			})
		}

	case createListPanel:
		(&p).Append(msg(Size{
			Height: p.height,
			Width:  p.width / p.maxDisplay,
		}))
		p.FocusCurrent()
		return p, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, p.keys.NextSibling):
			if p.idx+1 == len(p.p) {
				return p, nil
			}

			p.BlurCurrent()
			p.idx++
			p.FocusCurrent()
			return p, nil
		case key.Matches(msg, p.keys.LastSibling):
			if p.idx == 0 {
				return p, nil
			}

			p.BlurCurrent()
			p.idx--
			p.FocusCurrent()
			return p, nil
		case key.Matches(msg, p.keys.CloseCurrent):
			// don't close the root if we're not the only one left
			if p.closing != nil {
				cont, cmd := p.closing(&p, p.current(), len(p.p)-1)
				cmds = append(cmds, cmd)
				if !cont {
					return p, tea.Batch(cmds...)
				}
			}

			if p.idx == 0 {
				return p, nil
			}

			(&p).PopCurrent()
			return p, nil
		}
	}

	if p.updating != nil {
		passMsg, delegateCmd := p.updating(&p, p.current(), p.idx, msg)

		if !passMsg {
			return p, delegateCmd
		}

		cmds = append(cmds, delegateCmd)
	}

	var cmd tea.Cmd
	panel := p.p[p.idx]
	panel, cmd = panel.Update(msg)
	cmds = append(cmds, cmd)
	p.p[p.idx] = panel

	return p, tea.Batch(cmds...)
}

func (p Model) View() string {
	if p.maxDisplay >= len(p.p) {
		return p.renderAll()
	}

	return p.renderWindow()
}

func (p Model) renderAll() string {
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

func (p Model) renderWindow() string {
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
