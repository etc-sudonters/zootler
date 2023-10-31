package visualizer

import (
	"fmt"
	"io"
	"strings"
	"sudonters/zootler/cmd/zootler/visualizer/listpanel"
	"sudonters/zootler/cmd/zootler/visualizer/panels"
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/internal/entity/componenttable"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func createComponentOverview(rows []componenttable.RowData, pool entity.Pool) listpanel.Model {
	listItems := make([]list.Item, len(rows))

	for i, r := range rows {
		listItems[i] = overviewItem{r: r}
	}

	l := list.New(
		listItems,
		overviewDelegate{pool},
		0, 0,
	)
	l.Title = "Components Available"
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()

	panel := listpanel.New(listpanel.WithList(l))
	return panel
}

var (
	titleStyle    = lipgloss.NewStyle().MarginLeft(2)
	itemStyle     = lipgloss.NewStyle().PaddingLeft(4)
	selectedStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	pagination    = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle     = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type overviewItem struct{ r componenttable.RowData }

func (r overviewItem) FilterValue() string { return r.r.Type().Name() }

type overviewDelegate struct {
	pool entity.Pool
}

func (d overviewDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	r, ok := item.(overviewItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s:%d (%d/%d)", r.r.Type().Name(), r.r.Id(), r.r.Len(), r.r.Capacity())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprintf(w, fn(str))
}

func (d overviewDelegate) Height() int {
	return 1
}

func (d overviewDelegate) Spacing() int {
	return 0
}

func (d overviewDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			return panels.CreateListPanel(func(s panels.Size) listpanel.Model {
				cur := m.SelectedItem().(overviewItem)
				return createComponentDrillIn(cur.r, d.pool, s)
			})
		}

	}

	return cmd
}
