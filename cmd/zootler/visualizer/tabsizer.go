package visualizer

import tea "github.com/charmbracelet/bubbletea"
import "sudonters/zootler/cmd/zootler/visualizer/foldertabs"

type tabResizer struct {
	t foldertabs.Model
}

func (t tabResizer) Init() tea.Cmd {
	return t.t.Init()
}

func (t tabResizer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		(&t.t).SetWidth(msg.Width)
		return t, nil
	}

	m, cmd := t.t.Update(msg)
	t.t = m.(foldertabs.Model)
	return t, cmd
}

func (t tabResizer) View() string {
	return t.t.View()
}
