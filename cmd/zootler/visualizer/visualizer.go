package visualizer

import (
	"context"
	"sudonters/zootler/cmd/zootler/visualizer/application"
	"sudonters/zootler/cmd/zootler/visualizer/foldertabs"
	"sudonters/zootler/cmd/zootler/visualizer/listpanel"
	"sudonters/zootler/cmd/zootler/visualizer/panels"
	"sudonters/zootler/internal/entity/bitpool"
	"sudonters/zootler/internal/reitertools"
	"sudonters/zootler/pkg/world"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func Visualize(w world.World) visualizer {
	return visualizer{w}
}

type visualizer struct {
	w world.World
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
				panels.WithClose(closeTabForEmptyPanel(m.Title())),
			)
			return false, foldertabs.AddTab(m.Title(), newPanels, true)
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

func (v visualizer) Run(ctx context.Context) error {
	tbl, err := bitpool.ExtractComponentTable(v.w.Entities.Pool)
	if err != nil {
		return err
	}

	rows := reitertools.ToSlice(tbl.Rows())
	componentPanel := createComponentOverview(rows, v.w.Entities.Pool)
	(&componentPanel).Focus()

	p := panels.New(
		panels.WithPanels(componentPanel),
		panels.WithUpdate(popout),
		panels.WithClose(closeTabForEmptyPanel("Components")),
	)

	tabs := foldertabs.New(
		foldertabs.WithTab("Components", p),
		foldertabs.WithStyle(foldertabs.Style{
			Active:    lipgloss.NewStyle().Padding(0, 0).Border(tabBorderWithBottom("┘", " ", "└"), true),
			Inactive:  lipgloss.NewStyle().Padding(0, 0).Border(tabBorderWithBottom("┴", "─", "┴"), true),
			ModelView: lipgloss.NewStyle().Padding(1).Border(lipgloss.NormalBorder()).UnsetBorderTop(),
			BorderDelegate: func(border *lipgloss.Border, active, first, last bool) {
				if first && active {
					border.BottomLeft = "│"
				} else if first && !active {
					border.BottomLeft = "├"
				} else if last && active {
					border.BottomRight = "│"
				} else if last && !active {
					border.BottomRight = "┤"
				}
			},
			HorizontalFill: "─",
		}),
		foldertabs.WithKeys(foldertabs.Keys{
			NextTab: key.NewBinding(key.WithKeys("[")),
			PrevTab: key.NewBinding(key.WithKeys("]")),
		}),
	)

	app := application.New(tabs, application.WithStyle(lipgloss.NewStyle().Padding(2)))
	_, err = tea.NewProgram(app, tea.WithAltScreen(), tea.WithContext(ctx)).Run()
	return err
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
