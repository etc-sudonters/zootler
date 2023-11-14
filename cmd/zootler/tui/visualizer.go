package tui

import (
	"context"
	"sudonters/zootler/cmd/zootler/tui/rulestab"
	"sudonters/zootler/internal/entity/bitpool"
	"sudonters/zootler/internal/tui/application"
	"sudonters/zootler/internal/tui/foldertabs"
	"sudonters/zootler/internal/tui/listpanel"
	"sudonters/zootler/internal/tui/panels"
	"sudonters/zootler/pkg/world"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/etc-sudonters/substrate/reiterate"
)

func Tui(w world.World) tui {
	return tui{w}
}

type tui struct {
	w world.World
}

func (v tui) Run(ctx context.Context) error {
	tbl, err := bitpool.ExtractComponentTable(v.w.Entities.Pool)
	if err != nil {
		return err
	}

	rows := reiterate.ToSlice(tbl.Rows())
	componentPanel := createComponentOverview(rows, v.w.Entities.Pool)
	(&componentPanel).Focus()

	components := panels.New(
		panels.WithPanels(componentPanel),
		panels.WithUpdate(popout),
		panels.WithClose(closeTabForEmptyPanel("Components")),
	)

	rulesparsing := rulestab.NewParseTab(v.w.Entities.Pool)

	tabs := foldertabs.New(
		foldertabs.WithTab("Components", components),
		foldertabs.WithTab("Rules", rulesparsing),
		foldertabs.WithStyle(tabStyle),
		foldertabs.WithKeys(tabKeys),
	)

	app := application.New(
		tabs,
		application.WithStyle(lipgloss.NewStyle().Padding(2)),
		application.WithKeyDisplay(),
		application.WithMsgDisplay())

	_, err = tea.NewProgram(app, tea.WithAltScreen(), tea.WithContext(ctx)).Run()
	return err
}

func popout(p *panels.Model, m *listpanel.Model, idx int, msg tea.Msg) (bool, tea.Cmd) {
	// don't popout root
	if idx == 0 {
		return true, nil
	}
	promotePanelToTab := key.NewBinding(key.WithKeys("!"))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, promotePanelToTab):
			p.PopCurrent()
			newPanels := panels.New(
				panels.WithPanels(*m),
				panels.WithUpdate(popout),
				panels.WithClose(closeTabForEmptyPanel(m.L.Title)),
			)
			return false, foldertabs.AddTab(m.L.Title, newPanels, true)
		}
	}

	return true, nil
}

func closeTabForEmptyPanel(tabName string) panels.CloseDelegate {
	return func(p *panels.Model, m2 *listpanel.Model, lenAfterPop int) (bool, tea.Cmd) {
		if lenAfterPop != 0 {
			return true, nil
		}

		return false, foldertabs.CloseTab(tabName)
	}
}
