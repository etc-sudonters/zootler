package decorator

import tea "github.com/charmbracelet/bubbletea"

type Decorator interface {
	Update(model tea.Model, msg tea.Msg) (bool, tea.Model, tea.Cmd)
	Init(m tea.Model) (tea.Model, tea.Cmd)
	View(m tea.Model, view string) string
}

func Decorate(m tea.Model, i Decorator) tea.Model {
	return decorator{m, i}
}

type decorator struct {
	m tea.Model
	d Decorator
}

func (u decorator) Init() tea.Cmd {
	var cmds []tea.Cmd
	if u.m != nil {
		cmds = append(cmds, u.m.Init())
	}

	var cmd tea.Cmd
	u.m, cmd = u.d.Init(u.m)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (u decorator) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	cont, m, cmd := u.d.Update(u.m, msg)
	u.m = m
	cmds = append(cmds, cmd)
	if !cont {
		return u, tea.Batch(cmds...)
	}

	m, cmd = u.m.Update(msg)
	u.m = m
	cmds = append(cmds, cmd)
	return u, tea.Batch(cmds...)
}

func (u decorator) View() string {
	return u.d.View(u.m, u.m.View())
}
